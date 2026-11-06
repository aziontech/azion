package v6

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	msg "github.com/aziontech/azion-cli/messages/deploy-remote"
	deploy "github.com/aziontech/azion-cli/pkg/cmd/deploy_remote"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/output"
	"github.com/aziontech/azion-cli/pkg/pipeline"
	"github.com/aziontech/azion-cli/pkg/poll"
	"github.com/aziontech/azion-cli/utils"
	sdk "github.com/aziontech/azionapi-v6-go-sdk/azion-api"
	"go.uber.org/zap"
)

// Runner holds what Run reads from outside, so tests can replace it.
type Runner struct {
	F       *cmdutil.Factory
	Cmd     *deploy.DeployCmd
	Clients *Clients
	Opts    Options

	ReadConf  func(dir string) (*contracts.AzionApplicationOptionsV6, error)
	WriteConf func(conf *contracts.AzionApplicationOptionsV6, dir string) error
	NewPrefix func() string
	WorkDir   func() (string, error)
	// Signals are the signals that interrupt the deploy.
	Signals []os.Signal
	// Poll tunes every wait; tests replace its clock.
	Poll poll.Options
}

// NewRunner builds a Runner over the real API, file system and signals. cmd
// carries the flags the v4 steps it shares read (config file, workers, bucket
// access and so on).
func NewRunner(f *cmdutil.Factory, cmd *deploy.DeployCmd, opts Options) *Runner {
	if opts.ConfigDir == "" {
		opts.ConfigDir = cmd.ProjectConf
	}
	// The v6 presets name the Environment with $ENVIRONMENT_NAME, which the
	// first run replaces in azion.config: --environment, or Production (D3).
	env := opts.Environment
	if env == "" {
		env = DefaultEnvironment
	}
	cmd.ExtraReplacements = map[string]string{"$ENVIRONMENT_NAME": env}
	return &Runner{
		F: f, Cmd: cmd.ForV6(), Opts: opts,
		Clients:   NewClients(f),
		ReadConf:  utils.GetAzionJsonContentV6,
		WriteConf: utils.WriteAzionJsonContentV6,
		NewPrefix: utils.Timestamp,
		WorkDir:   utils.GetWorkingDir,
		Signals:   []os.Signal{os.Interrupt, syscall.SIGTERM},
	}
}

// Run deploys the project to an API v6 account.
func (r *Runner) Run(parent context.Context) error {
	// With no signals, NotifyContext would relay every signal, including the
	// ones the Go runtime sends itself; only listen when there are some.
	ctx := parent
	if len(r.Signals) > 0 {
		var stop context.CancelFunc
		ctx, stop = signal.NotifyContext(parent, r.Signals...)
		defer stop()
	}

	start := time.Now()
	f := r.F
	msgs := []string{}

	conf, err := r.ReadConf(r.Opts.ConfigDir)
	if err != nil {
		return err
	}
	if conf.Name == "" {
		if wd, err := r.WorkDir(); err == nil {
			conf.Name = filepath.Base(wd)
		}
	}

	oldPrefix, newPrefix := conf.Prefix, conf.Prefix
	if conf.Prefix == "" || conf.RotatePrefix == nil || *conf.RotatePrefix {
		newPrefix = r.NewPrefix()
	}

	timing := &TimingSummaryV6{}
	s := &DeployStateV6{
		F: f, Opts: r.Opts, Cmd: r.Cmd, Clients: r.Clients, Conf: conf,
		Project: &deploy.Project{
			Name: conf.Name, Bucket: conf.Bucket, Prefix: conf.Prefix, Preset: conf.Preset,
			RotatePrefix: conf.RotatePrefix, NotFirstRun: conf.NotFirstRun,
		},
		OldPrefix: oldPrefix, NewPrefix: newPrefix,
		Poll: r.Poll, msgs: &msgs, timing: timing,
	}
	// A dry run changes nothing, azion.json included.
	if r.Opts.DryRun {
		return r.dryRun(ctx, s)
	}

	s.WriteConf = func() error { return r.WriteConf(conf, r.Opts.ConfigDir) }
	// Saved on every way out, including Ctrl-C and failures.
	defer s.writeConf()

	s.say(msg.V6RunningDeploy)
	_, err = pipeline.Run(ctx, V6DeployPipeline{}, s, &pipeline.Plan{Factory: f, Timing: timing.Record})
	if err != nil {
		return r.explain(ctx, s, err)
	}

	timing.Total = time.Since(start)
	if f.IOStreams != nil && !f.Silent {
		timing.Print(f.IOStreams.Out)
	}
	for _, w := range s.Warnings {
		s.say(w)
	}
	s.say(msg.V6DeploySuccessful)
	if s.Activated() {
		for _, w := range conf.Workloads {
			if w.Url != "" {
				s.say(fmt.Sprintf(msg.V6DeployOutputURL, w.Url))
			}
		}
	}

	out := output.SliceOutput{
		Messages:      msgs,
		GeneralOutput: output.GeneralOutput{Out: f.IOStreams.Out, Flags: f.Flags},
	}
	return output.Print(&out)
}

// explain adds to a failure what the user needs to know about traffic: an
// interruption says whether it came before or after the activation, and a
// failure after the activation says the new release is live.
func (r *Runner) explain(ctx context.Context, s *DeployStateV6, err error) error {
	for _, w := range s.Warnings {
		s.say(w)
	}
	if ctx.Err() != nil || errors.Is(err, context.Canceled) {
		when := msg.V6InterruptedBefore
		if s.Activated() {
			when = msg.V6InterruptedAfter
		}
		s.say(fmt.Sprintf(msg.V6Interrupted, when))
		return err
	}
	var stepErr *pipeline.StepError
	if s.Activated() && errors.As(err, &stepErr) {
		s.say(fmt.Sprintf(msg.V6ActiveButIncomplete, stepLabels[stepErr.Step]))
	}
	logger.Debug("v6 deploy failed", zap.Error(err))
	return err
}

// dryRun reports the steps a deploy would run, from the same step list, with
// no request to the API. It reads the manifest when there is one, so the
// report can say which resource kinds are involved. A manifest it cannot use
// (one left by a v4 build, say) only loses the counts: the deploy builds the
// manifest again before reading it.
func (r *Runner) dryRun(ctx context.Context, s *DeployStateV6) error {
	if !deploy.ManifestMissing() {
		if err := s.readManifest(ctx, !s.runsBundlerInit()); err != nil {
			logger.Debug("dry run: the existing manifest is not usable", zap.Error(err))
			s.Manifest = nil
			s.say(fmt.Sprintf(msg.V6DryRunStaleManifest, err))
		}
	}
	res, err := pipeline.Run(ctx, V6DeployPipeline{}, s, &pipeline.Plan{Factory: r.F, DryRun: true})
	if err != nil {
		return err
	}
	s.say(msg.V6DryRunHeader)
	for _, step := range res.Steps {
		s.say(fmt.Sprintf(msg.V6DryRunStep, describeStep(s, step.Name)))
	}
	if len(res.Skipped) > 0 {
		labels := make([]string, len(res.Skipped))
		for i, name := range res.Skipped {
			labels[i] = stepLabels[name]
		}
		s.say(fmt.Sprintf(msg.V6DryRunSkipped, joinLabels(labels)))
	}
	return nil
}

// describeStep says what one step would do, with counts from the manifest.
func describeStep(s *DeployStateV6, name string) string {
	label := stepLabels[name]
	m := s.Manifest
	if m == nil {
		return label
	}
	count := map[string]int{
		StepManifestFunctions:       len(m.Functions),
		StepManifestNetworkLists:    len(m.NetworkLists),
		StepManifestWAFs:            len(m.WAFs),
		StepManifestConnectors:      len(m.Connectors),
		StepManifestApplication:     len(m.Applications),
		StepManifestFirewalls:       len(m.Firewalls),
		StepManifestCustomPages:     len(m.CustomPages),
		StepVersionBuildDeps:        dependencyCount(m),
		StepVersionBuildPublishable: publishableCount(m),
	}
	if n, ok := count[name]; ok {
		return fmt.Sprintf("%s: %d Draft(s)", label, n)
	}
	switch name {
	case StepManifestEnvironments:
		return fmt.Sprintf("%s: %s", label, joinLabels(namesIn(m.Environments, s.inScope)))
	case StepManifestDeployments:
		return fmt.Sprintf("%s: %s", label, joinLabels(s.targetDeployments()))
	case StepReleaseCreate:
		return fmt.Sprintf("%s: %d resource(s) on each Deployment, each at the version this run builds", label, publishableCount(m)+dependencyCount(m))
	}
	return label
}

// namesIn lists the Environments in scope by name.
func namesIn(envs []sdk.EnvironmentCreateRequest, inScope func(string) bool) []string {
	var out []string
	for _, e := range envs {
		if inScope(e.Name) {
			out = append(out, e.Name)
		}
	}
	return out
}

func joinLabels(labels []string) string {
	out := ""
	for i, l := range labels {
		if i > 0 {
			out += ", "
		}
		out += l
	}
	return out
}
