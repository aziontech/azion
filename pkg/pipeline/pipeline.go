// Package pipeline runs a deploy as a declared list of steps instead of a
// hand-written sequence.
//
// Here a generation supplies a Pipeline: a name and an ordered list of steps
// over its own state type. The runner owns everything that was copied around
// them — the spinner, per-step timing, debug logging, dry-run reporting and
// attaching the failing step's name to an error. Adding a generation means
// supplying a different step list, not forking the runner.
//
// The state type is a parameter, not a shared schema. A v4 step reads a
// ResourceContext and a ManifestV4; a v3 step reads its own clients and its own
// manifest type; a v6 step will read whatever v6 turns out to need. The runner
// never looks inside it, so a pipeline whose state has nothing in common with
// another's still fits. Pipelines themselves are stateless: Steps() is a table,
// and the state is handed to Run.
package pipeline

import (
	"context"
	"time"

	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/briandowns/spinner"
	"go.uber.org/zap"
)

// Step is one unit of work in a deploy, over a generation's state type T.
type Step[T any] struct {
	// Name identifies the step in timing reports, debug logs and dry-run
	// output.
	//
	// It is a contract, not a label: pkg/cmd/deploy_remote/timer.go switches on
	// these names to fill its timing summary, so renaming one silently stops
	// that measurement rather than failing to build. The established names are
	// the Manifest-prefixed ones — "ManifestFunctions", "ManifestCacheSettings"
	// and so on — and a new step should follow that convention.
	Name string

	// Applies reports whether this run needs the step. A nil Applies means the
	// step always runs. This replaces the inline len(...) > 0 checks.
	Applies func(*T) bool

	// Run performs the step.
	Run func(context.Context, *T) error
}

// needed reports whether the step applies to this run.
func (s Step[T]) needed(state *T) bool {
	return s.Applies == nil || s.Applies(state)
}

// Pipeline is one generation's deploy, expressed as an ordered list of steps.
//
// Orphan cleanup is a step like any other, and belongs at the end of Steps().
// It is not a separate interface method because the current behaviour is that a
// failed step returns immediately and cleanup does not run — which is exactly
// what a trailing step gives. Making it a deferred hook would change that.
type Pipeline[T any] interface {
	// Name identifies the generation, for logs and dry-run output.
	Name() string
	// Steps returns the steps in dependency order.
	Steps() []Step[T]
}

// TimingCallback receives how long each executed step took.
type TimingCallback func(name string, duration time.Duration)

// Plan carries what the runner needs and a step does not: where to write, and
// how this invocation wants to be told what is happening. It is deliberately
// not generic — nothing in it varies by generation.
type Plan struct {
	// Factory supplies the output streams and the debug flag.
	Factory *cmdutil.Factory

	// Msgs accumulates the lines shown to the user, for --format json. The
	// runner appends only what it prints itself.
	Msgs *[]string

	// DryRun reports which steps would run without running any of them.
	DryRun bool

	// Timing, when set, is called once per executed step.
	Timing TimingCallback

	// Spinner is the suffix of the spinner shown while the pipeline runs. Empty
	// means no spinner. It is suppressed under --debug and in a dry run, as the
	// hand-written version already did.
	Spinner string

	// Announce is printed once before the run. Empty means nothing is printed.
	//
	// Spinner and Announce both exist because the two generations present this
	// differently today: v4 shows a spinner whose suffix is the "creating
	// resources" message, v3 prints that same message as a line and shows no
	// spinner. Naming both keeps the port behaviour-preserving; converging them
	// is a decision for the deploy commands, not for the runner.
	Announce string
}

// Result reports what a run did. It is what gives --dry-run an implementation
// derived from the step list rather than a second, separately maintained
// simulation.
type Result struct {
	// Pipeline is the name of the pipeline that ran.
	Pipeline string
	// Steps lists, in order, the steps this run acted on — those whose Applies
	// returned true.
	//
	// The name deliberately does not claim execution, because it would be wrong
	// in two of the three cases. On a successful run these ran to completion and
	// Duration is real. In a dry run none of them ran and Duration is zero. When
	// a step fails, the run stops and the failing step is not listed here at
	// all — it is named by StepError.Step — so this holds only what completed
	// before the failure.
	Steps []StepResult
	// Skipped lists the steps whose Applies returned false, in order.
	Skipped []string
	// DryRun reports whether this was a dry run.
	DryRun bool
}

// StepResult is one executed step.
type StepResult struct {
	Name     string
	Duration time.Duration
}

// StepError reports which step failed.
//
// Error deliberately returns the underlying message unchanged. The steps this
// replaces returned their errors raw, so prefixing the step name here would
// change the text of every deploy failure a user sees. Callers that want the
// step name ask for it with errors.As; callers that print the error get exactly
// what they got before.
type StepError struct {
	Step string
	Err  error
}

func (e *StepError) Error() string { return e.Err.Error() }
func (e *StepError) Unwrap() error { return e.Err }

// Run executes the pipeline's steps against state, in order, skipping those
// that do not apply, and stops at the first error.
//
// A step's error is returned wrapped in *StepError, which prints identically to
// the original.
func Run[T any](ctx context.Context, p Pipeline[T], state *T, plan *Plan) (*Result, error) {
	if plan == nil {
		plan = &Plan{}
	}

	res := &Result{Pipeline: p.Name(), DryRun: plan.DryRun}

	if plan.Announce != "" {
		plan.print(plan.Announce)
	}

	stop := plan.startSpinner()
	defer stop()

	for _, step := range p.Steps() {
		if !step.needed(state) {
			logger.Debug("Skipping pipeline step", zap.String("step", step.Name))
			res.Skipped = append(res.Skipped, step.Name)
			continue
		}

		if plan.DryRun {
			res.Steps = append(res.Steps, StepResult{Name: step.Name})
			continue
		}

		logger.Debug("Applying pipeline step", zap.String("step", step.Name))
		start := time.Now()
		if err := step.Run(ctx, state); err != nil {
			return res, &StepError{Step: step.Name, Err: err}
		}
		elapsed := time.Since(start)

		res.Steps = append(res.Steps, StepResult{Name: step.Name, Duration: elapsed})
		if plan.Timing != nil {
			plan.Timing(step.Name, elapsed)
		}
	}

	return res, nil
}

// startSpinner shows the spinner if this run wants one, and returns the
// function that stops it. The returned function is always safe to call.
func (plan *Plan) startSpinner() func() {
	if plan.Spinner == "" || plan.DryRun || plan.debug() {
		return func() {}
	}
	s := spinner.New(spinner.CharSets[7], 100*time.Millisecond)
	s.Suffix = " " + plan.Spinner
	s.FinalMSG = "\n"
	s.Start()
	return s.Stop
}

func (plan *Plan) debug() bool {
	return plan.Factory != nil && plan.Factory.Debug
}

// print writes a line to the user and records it for --format json.
func (plan *Plan) print(line string) {
	if plan.Factory != nil && plan.Factory.IOStreams != nil {
		logger.FInfoFlags(plan.Factory.IOStreams.Out, line, plan.Factory.Format, plan.Factory.Out)
	}
	if plan.Msgs != nil {
		*plan.Msgs = append(*plan.Msgs, line)
	}
}
