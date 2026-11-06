// Package v6 is the deploy for API v6 accounts, as planned in
// doc/V6-DEPLOY-PLAN.md.
//
// It runs on the same pipeline.Run as the v4 deploy, with its own steps and its
// own state. The difference is inside the steps: configuration goes into
// Drafts, each Draft is built and waited on, and only then is a release made,
// built, activated and checked. Nothing reaches traffic before the activation,
// so any failure before it leaves the site as it was.
package v6

import (
	"context"
	"sync"
	"time"

	apiPurge "github.com/aziontech/azion-cli/pkg/api/realtime_purge"
	apiDeployments "github.com/aziontech/azion-cli/pkg/api/v6/deployments"
	apiEnvironments "github.com/aziontech/azion-cli/pkg/api/v6/environments"
	apiWorkloads "github.com/aziontech/azion-cli/pkg/api/v6/workloads"
	deploy "github.com/aziontech/azion-cli/pkg/cmd/deploy_remote"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/logger"
	manifestv6 "github.com/aziontech/azion-cli/pkg/manifest/v6"
	"github.com/aziontech/azion-cli/pkg/poll"
	"github.com/aziontech/azion-cli/pkg/registry"
	"github.com/aziontech/azion-cli/pkg/token"
	"go.uber.org/zap"
)

// DefaultEnvironment is where a deploy goes when nothing names another
// Environment (doc/V6-DEPLOY-PLAN.md, D3).
const DefaultEnvironment = "Production"

// Options are the v6 deploy's flags.
type Options struct {
	ConfigDir     string
	Environment   string
	NoActivate    bool
	Timeout       time.Duration
	DryRun        bool
	SkipBuild     bool
	SkipFramework bool
}

// Clients are every API client the v6 deploy uses.
type Clients struct {
	Resources    *manifestv6.Clients
	Environments *apiEnvironments.Client
	Deployments  *apiDeployments.Client
	Workloads    *apiWorkloads.Client
	Purge        *apiPurge.Client
}

func NewClients(f *cmdutil.Factory) *Clients {
	return &Clients{
		Resources: &manifestv6.Clients{
			Applications: registry.ApplicationsV6(f),
			Functions:    registry.FunctionsV6(f),
			Connectors:   registry.ConnectorsV6(f),
			Firewalls:    registry.FirewallsV6(f),
			CustomPages:  registry.CustomPagesV6(f),
			NetworkLists: registry.NetworkListsV6(f),
			WAFs:         registry.WAFsV6(f),
		},
		Environments: registry.EnvironmentsV6(f),
		Deployments:  registry.DeploymentsV6(f),
		Workloads:    registry.WorkloadsV6(f),
		Purge:        registry.RealtimePurge(f),
	}
}

// DeployStateV6 is what the v6 deploy steps read and write. Each field below
// the clients is written by one step and read by later ones, which is the
// dependency the step order encodes.
type DeployStateV6 struct {
	F    *cmdutil.Factory
	Opts Options
	// Cmd is the v4 deploy command, for the steps v6 shares with it: build,
	// bucket and upload.
	Cmd     *deploy.DeployCmd
	Clients *Clients

	Conf    *contracts.AzionApplicationOptionsV6
	Project *deploy.Project
	// WriteConf saves azion.json; it runs after each step that changes it.
	WriteConf func() error

	OldPrefix, NewPrefix string

	Manifest *contracts.ManifestV6
	Applier  *manifestv6.Applier
	Creds    token.S3Credentials

	// Dependencies and Publishable are the Drafts the two build levels wait
	// on (doc/V6-DEPLOY-PLAN.md, section 3).
	Dependencies []manifestv6.Pending
	Publishable  []manifestv6.Pending

	// Environments are the manifest's, by manifest name, resolved to the
	// account's. mustExist names those the CLI may not create (filled in
	// when the manifest declares none).
	Environments map[string]apiEnvironments.Environment
	mustExist    map[string]bool
	// Targets are the Deployments this run releases to.
	Targets []*releaseTarget
	// Workloads are the Workloads this run saved.
	Workloads []appliedWorkload

	// Warnings are shown after the deploy; they do not fail it.
	Warnings []string

	// Poll tunes every wait; tests replace its clock.
	Poll poll.Options

	mu     sync.Mutex
	msgs   *[]string
	timing *TimingSummaryV6

	uploadDecided, upload bool
}

// say prints a line and records it for --format json. Builds of one level
// run concurrently, so it is safe to call from several goroutines.
func (s *DeployStateV6) say(line string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.F != nil && s.F.IOStreams != nil {
		logger.FInfoFlags(s.F.IOStreams.Out, line, s.F.Format, s.F.Out)
	}
	if s.msgs != nil {
		*s.msgs = append(*s.msgs, line)
	}
}

func (s *DeployStateV6) warn(line string) {
	s.Warnings = append(s.Warnings, line)
}

func (s *DeployStateV6) writeConf() {
	if s.WriteConf == nil {
		return
	}
	s.syncProject()
	if err := s.WriteConf(); err != nil {
		logger.Debug("Error while writing azion.json", zap.Error(err))
	}
}

// syncProject copies what the shared v4 steps changed back into azion.json.
func (s *DeployStateV6) syncProject() {
	if s.Project == nil {
		return
	}
	s.Conf.Bucket, s.Conf.Prefix = s.Project.Bucket, s.Project.Prefix
}

func (s *DeployStateV6) pollOptions() poll.Options {
	o := s.Poll
	if s.Opts.Timeout > 0 {
		o.Timeout = s.Opts.Timeout
	}
	return o
}

// shouldUpload is asked by two steps and answered once.
func (s *DeployStateV6) shouldUpload() bool {
	if !s.uploadDecided {
		s.uploadDecided = true
		s.upload = s.Cmd.ShouldUpload()
	}
	return s.upload
}

// cleanupContext gives best-effort cleanup (canceling a build, deleting a
// failed release) a short deadline of its own, so it still runs when ctx was
// canceled by Ctrl-C.
func cleanupContext(ctx context.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.WithoutCancel(ctx), 10*time.Second)
}
