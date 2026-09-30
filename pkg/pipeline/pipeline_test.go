package pipeline

import (
	"bytes"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/iostreams"
	"github.com/aziontech/azion-cli/pkg/logger"
	"go.uber.org/zap/zapcore"
)

// state is a stand-in for a generation's state type. ran records the order the
// steps executed in, so the tests assert on the sequence the runner produced.
type state struct {
	ran      []string
	declared map[string]int
}

func newState(declared map[string]int) *state {
	return &state{declared: declared}
}

// fake is a Pipeline over state. Pipelines are stateless: this is just a table.
type fake struct {
	name  string
	steps []Step[state]
}

func (f fake) Name() string         { return f.name }
func (f fake) Steps() []Step[state] { return f.steps }

// recorder builds a step that appends its own name to state.ran when it runs.
func recorder(name string, applies func(*state) bool, err error) Step[state] {
	return Step[state]{
		Name:    name,
		Applies: applies,
		Run: func(_ context.Context, s *state) error {
			s.ran = append(s.ran, name)
			return err
		},
	}
}

// declared is the predicate the real pipelines use: "the manifest declares at
// least one of these".
func declared(key string) func(*state) bool {
	return func(s *state) bool { return s.declared[key] > 0 }
}

func testFactory() (*cmdutil.Factory, *bytes.Buffer) {
	out := &bytes.Buffer{}
	return &cmdutil.Factory{
		IOStreams: &iostreams.IOStreams{Out: out, Err: &bytes.Buffer{}},
	}, out
}

func never(*state) bool { return false }

func TestRunExecutesStepsInOrderAndSkipsInapplicable(t *testing.T) {
	logger.New(zapcore.InfoLevel)
	p := fake{name: "test", steps: []Step[state]{
		recorder("first", nil, nil), // nil Applies means always
		recorder("skipped", never, nil),
		recorder("second", func(*state) bool { return true }, nil),
	}}
	s := newState(nil)

	res, err := Run(context.Background(), p, s, &Plan{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if want := []string{"first", "second"}; !equal(s.ran, want) {
		t.Errorf("ran %v, want %v", s.ran, want)
	}
	if want := []string{"first", "second"}; !equal(names(res.Steps), want) {
		t.Errorf("Steps = %v, want %v", names(res.Steps), want)
	}
	if want := []string{"skipped"}; !equal(res.Skipped, want) {
		t.Errorf("Skipped = %v, want %v", res.Skipped, want)
	}
	if res.Pipeline != "test" {
		t.Errorf("Pipeline = %q, want %q", res.Pipeline, "test")
	}
}

func TestStepsReceiveTheStateHandedToRun(t *testing.T) {
	logger.New(zapcore.InfoLevel)
	// The point of the type parameter: a step reads and writes the state the
	// caller supplied, rather than whatever a closure happened to capture.
	p := fake{steps: []Step[state]{
		{
			Name:    "reads-state",
			Applies: declared("functions"),
			Run: func(_ context.Context, s *state) error {
				s.ran = append(s.ran, "saw-functions")
				return nil
			},
		},
	}}

	withFunctions := newState(map[string]int{"functions": 1})
	if _, err := Run(context.Background(), p, withFunctions, &Plan{}); err != nil {
		t.Fatal(err)
	}
	if want := []string{"saw-functions"}; !equal(withFunctions.ran, want) {
		t.Errorf("ran %v, want %v", withFunctions.ran, want)
	}

	// Same step table, different state: the predicate now says no.
	without := newState(nil)
	if _, err := Run(context.Background(), p, without, &Plan{}); err != nil {
		t.Fatal(err)
	}
	if len(without.ran) != 0 {
		t.Errorf("ran %v, want nothing", without.ran)
	}
}

func TestRunStopsAtFirstError(t *testing.T) {
	logger.New(zapcore.InfoLevel)
	boom := errors.New("the server could not process the request")
	p := fake{steps: []Step[state]{
		recorder("ok", nil, nil),
		recorder("bad", nil, boom),
		recorder("never-reached", nil, nil),
	}}
	s := newState(nil)

	res, err := Run(context.Background(), p, s, &Plan{})
	if err == nil {
		t.Fatal("expected an error")
	}

	// Result.Steps holds only what completed before the failure: the failing
	// step returns before it is appended, and is named by StepError instead.
	if want := []string{"ok"}; !equal(names(res.Steps), want) {
		t.Errorf("Steps = %v, want %v — a failed step must not be listed as done", names(res.Steps), want)
	}
	if want := []string{"ok", "bad"}; !equal(s.ran, want) {
		t.Errorf("ran %v, want %v — a failed step must stop the run", s.ran, want)
	}

	// The message must be byte-identical to the underlying error: the code this
	// replaces returned step errors raw, so prefixing the step name here would
	// change the text of every deploy failure a user sees.
	if got, want := err.Error(), boom.Error(); got != want {
		t.Errorf("Error() = %q, want %q — step errors must not be reworded", got, want)
	}
	if !errors.Is(err, boom) {
		t.Error("errors.Is must see through StepError to the original")
	}
	var se *StepError
	if !errors.As(err, &se) {
		t.Fatal("errors.As must yield *StepError")
	}
	if se.Step != "bad" {
		t.Errorf("StepError.Step = %q, want %q", se.Step, "bad")
	}
}

func TestRunReportsTimingPerExecutedStep(t *testing.T) {
	logger.New(zapcore.InfoLevel)
	timed := map[string]int{}
	p := fake{steps: []Step[state]{
		recorder("ran", nil, nil),
		recorder("skipped", never, nil),
	}}

	_, err := Run(context.Background(), p, newState(nil), &Plan{
		Timing: func(name string, _ time.Duration) { timed[name]++ },
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if timed["ran"] != 1 {
		t.Errorf("timing for %q called %d times, want 1", "ran", timed["ran"])
	}
	if _, ok := timed["skipped"]; ok {
		t.Error("a skipped step must not be timed")
	}
}

func TestDryRunRunsNothingButReportsWhatWould(t *testing.T) {
	logger.New(zapcore.InfoLevel)
	p := fake{steps: []Step[state]{
		recorder("would-run", func(*state) bool { return true }, nil),
		recorder("would-skip", never, nil),
		recorder("would-explode", nil, errors.New("must not be called")),
	}}
	s := newState(nil)

	res, err := Run(context.Background(), p, s, &Plan{DryRun: true})
	if err != nil {
		t.Fatalf("a dry run must not fail: %v", err)
	}
	if len(s.ran) != 0 {
		t.Errorf("dry run executed %v, want nothing", s.ran)
	}
	if want := []string{"would-run", "would-explode"}; !equal(names(res.Steps), want) {
		t.Errorf("Steps = %v, want %v", names(res.Steps), want)
	}
	if want := []string{"would-skip"}; !equal(res.Skipped, want) {
		t.Errorf("Skipped = %v, want %v", res.Skipped, want)
	}
	if !res.DryRun {
		t.Error("Result.DryRun must be set")
	}
}

func TestAnnouncePrintsOnceAndIsRecorded(t *testing.T) {
	logger.New(zapcore.InfoLevel)
	f, out := testFactory()
	var msgs []string
	p := fake{steps: []Step[state]{recorder("noop", nil, nil)}}

	_, err := Run(context.Background(), p, newState(nil), &Plan{
		Factory: f, Msgs: &msgs, Announce: "creating resources",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !bytes.Contains(out.Bytes(), []byte("creating resources")) {
		t.Errorf("Announce was not printed; out = %q", out.String())
	}
	if want := []string{"creating resources"}; !equal(msgs, want) {
		t.Errorf("Msgs = %v, want %v", msgs, want)
	}
}

func TestRunToleratesNilPlanAndNilFactory(t *testing.T) {
	logger.New(zapcore.InfoLevel)
	p := fake{steps: []Step[state]{recorder("noop", nil, nil)}}
	if _, err := Run(context.Background(), p, newState(nil), nil); err != nil {
		t.Fatalf("nil plan: %v", err)
	}
	if _, err := Run(context.Background(), p, newState(nil), &Plan{Announce: "x", Spinner: "y"}); err != nil {
		t.Fatalf("nil factory: %v", err)
	}
}

// v4Steps is the step table of pkg/manifest/manifest.go, transcribed: the same
// names, the same order and the same predicates.
//
// Two things it pins down for the port in the next item. The names are the ones
// pkg/cmd/deploy_remote/timer.go switches on, so they are a contract. And cache
// settings, connectors and rules engine are nested inside the "there is an
// application" check in the current code — connectors are not applied for a
// manifest that declares connectors but no application — which a flat reading
// of the source would miss.
func v4Steps() []Step[state] {
	hasApp := declared("applications")
	and := func(a, b func(*state) bool) func(*state) bool {
		return func(s *state) bool { return a(s) && b(s) }
	}
	return []Step[state]{
		recorder("ManifestFunctions", declared("functions"), nil),
		recorder("ManifestFunctionInstances", and(hasApp, declared("instances")), nil),
		recorder("ManifestEdgeApplication", hasApp, nil),
		recorder("ManifestCacheSettings", and(hasApp, declared("caches")), nil),
		recorder("ManifestConnectors", and(hasApp, declared("connectors")), nil),
		recorder("ManifestRulesEngine", and(hasApp, declared("rules")), nil),
		recorder("ManifestWorkloads", declared("workloads"), nil),
		recorder("ManifestWorkloadDeployments", declared("deployments"), nil),
		recorder("ManifestFirewalls", declared("firewalls"), nil),
		recorder("ManifestPurge", declared("purge"), nil),
	}
}

// TestExpressesTheRealV4Pipeline is a design check, not a test of production
// code: it asserts the runner reproduces the order the hand-written sequence in
// pkg/manifest/manifest.go produces today.
func TestExpressesTheRealV4Pipeline(t *testing.T) {
	logger.New(zapcore.InfoLevel)
	p := fake{name: "v4", steps: v4Steps()}

	tests := []struct {
		name     string
		declared map[string]int
		want     []string
	}{
		{
			name: "everything declared",
			declared: map[string]int{
				"functions": 1, "applications": 1, "instances": 1, "caches": 1,
				"connectors": 1, "rules": 1, "workloads": 1, "deployments": 1,
				"firewalls": 1, "purge": 1,
			},
			want: []string{
				"ManifestFunctions", "ManifestFunctionInstances", "ManifestEdgeApplication",
				"ManifestCacheSettings", "ManifestConnectors", "ManifestRulesEngine",
				"ManifestWorkloads", "ManifestWorkloadDeployments", "ManifestFirewalls",
				"ManifestPurge",
			},
		},
		{
			name:     "connectors without an application are not applied",
			declared: map[string]int{"connectors": 1},
			want:     nil,
		},
		{
			name:     "cache settings and rules are also nested under the application",
			declared: map[string]int{"caches": 1, "rules": 1},
			want:     nil,
		},
		{
			name:     "functions and purge need no application",
			declared: map[string]int{"functions": 1, "purge": 1},
			want:     []string{"ManifestFunctions", "ManifestPurge"},
		},
		{
			name:     "workloads and firewalls need no application",
			declared: map[string]int{"workloads": 1, "firewalls": 1},
			want:     []string{"ManifestWorkloads", "ManifestFirewalls"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := newState(tt.declared)
			if _, err := Run(context.Background(), p, s, &Plan{}); err != nil {
				t.Fatal(err)
			}
			if !equal(s.ran, tt.want) {
				t.Errorf("ran %v, want %v", s.ran, tt.want)
			}
		})
	}
}

// TestStepNamesAreTheEstablishedOnes records the ten names that
// pkg/cmd/deploy_remote/timer.go switches on today, so the transcription above
// cannot drift while the port is being written.
//
// This is a transcription check, not a cross-file guard: both lists live in
// this file, and pkg/pipeline must not import a command package to read the
// real switch. The guard that ties the two together belongs in the package that
// owns both, once the v4 pipeline exists.
func TestStepNamesAreTheEstablishedOnes(t *testing.T) {
	want := []string{
		"ManifestFunctions", "ManifestFunctionInstances", "ManifestEdgeApplication",
		"ManifestCacheSettings", "ManifestConnectors", "ManifestRulesEngine",
		"ManifestWorkloads", "ManifestWorkloadDeployments", "ManifestFirewalls",
		"ManifestPurge",
	}
	got := make([]string, 0, len(want))
	for _, s := range v4Steps() {
		got = append(got, s.Name)
	}
	if !equal(got, want) {
		t.Errorf("step names = %v, want %v", got, want)
	}
}

func names(rs []StepResult) []string {
	out := make([]string, 0, len(rs))
	for _, r := range rs {
		out = append(out, r.Name)
	}
	return out
}

func equal(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
