package v6

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"runtime/debug"

	manifestmsg "github.com/aziontech/azion-cli/messages/manifest"
	deploy "github.com/aziontech/azion-cli/pkg/cmd/deploy_remote"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/logger"
	manifestv6 "github.com/aziontech/azion-cli/pkg/manifest/v6"
	"github.com/aziontech/azion-cli/pkg/pipeline"
	"go.uber.org/zap"
)

// manifestIsV6 reports whether the manifest on disk is one the v6 bundler
// wrote. A missing or unreadable one is not.
func manifestIsV6() bool {
	data, err := os.ReadFile(deploy.ManifestFile)
	if err != nil {
		return false
	}
	var head struct {
		Version int `json:"version"`
	}
	return json.Unmarshal(data, &head) == nil && head.Version == contracts.ManifestVersionV6
}

// V6DeployPipeline is the deploy flow for the v6 API. The step names are a
// contract with the timing summary (timer.go) and with the dry-run report.
type V6DeployPipeline struct{}

func (V6DeployPipeline) Name() string { return "deploy-v6" }

// Step names, in order. The build, read and storage names match the v4
// pipeline's, since those steps are the same.
const (
	StepGenerateManifest        = "GenerateManifest"
	StepBuild                   = "Build"
	StepReadManifest            = "ReadManifest"
	StepBucket                  = "Bucket"
	StepBundlerInit             = "BundlerInit"
	StepReadManifestAgain       = "ReadManifestAgain"
	StepCredentials             = "Credentials"
	StepUploadStaticFiles       = "UploadStaticFiles"
	StepManifestFunctions       = "ManifestFunctions"
	StepManifestNetworkLists    = "ManifestNetworkLists"
	StepManifestWAFs            = "ManifestWAFs"
	StepManifestConnectors      = "ManifestConnectors"
	StepVersionBuildDeps        = "VersionBuildDependencies"
	StepManifestApplication     = "ManifestApplication"
	StepManifestFirewalls       = "ManifestFirewalls"
	StepManifestCustomPages     = "ManifestCustomPages"
	StepVersionBuildPublishable = "VersionBuildPublishable"
	StepManifestEnvironments    = "ManifestEnvironments"
	StepManifestDeployments     = "ManifestDeployments"
	StepReleaseCreate           = "ReleaseCreate"
	StepReleaseBuild            = "ReleaseBuild"
	StepReleaseActivate         = "ReleaseActivate"
	StepManifestWorkloads       = "ManifestWorkloads"
	StepVerifyDeploy            = "VerifyDeploy"
	StepManifestPurge           = "ManifestPurge"
	StepOrphans                 = "ReportOrphans"
)

// declares reports whether the manifest has items of one kind. Before the
// manifest is read (a dry run without one) it answers yes, so the report
// lists the step as possible rather than hiding it.
func declares(s *DeployStateV6, count func(*contracts.ManifestV6) int) bool {
	if s.Manifest == nil {
		return true
	}
	return count(s.Manifest) > 0
}

func (V6DeployPipeline) Steps() []pipeline.Step[DeployStateV6] {
	steps := steps()
	for i := range steps {
		steps[i].Run = guarded(steps[i].Name, steps[i].Run)
	}
	return steps
}

// guarded turns a panic in a step into that step's error, so a bug (a nil
// pointer from a response nobody expected) fails the deploy with a message and
// still runs the deferred save of azion.json, instead of crashing the CLI.
func guarded(name string, run func(context.Context, *DeployStateV6) error) func(context.Context, *DeployStateV6) error {
	return func(ctx context.Context, s *DeployStateV6) (err error) {
		defer func() {
			if r := recover(); r != nil {
				logger.Debug("v6 deploy step panicked", zap.String("step", name), zap.Any("panic", r), zap.ByteString("stack", debug.Stack()))
				err = fmt.Errorf("%w: %s: %v", manifestmsg.ErrorV6Internal, stepLabels[name], r)
			}
		}()
		return run(ctx, s)
	}
}

func steps() []pipeline.Step[DeployStateV6] {
	return []pipeline.Step[DeployStateV6]{
		{
			Name: StepGenerateManifest,
			// A first run generates the manifest when there is none, and also
			// when the one there is not v6: a v4 build or init leaves one
			// behind, and reading it would fail.
			Applies: func(s *DeployStateV6) bool {
				return !s.Conf.NotFirstRun && !manifestIsV6()
			},
			Run: func(_ context.Context, s *DeployStateV6) error {
				return s.Cmd.GenerateManifest(s.Project, s.msgs)
			},
		},
		{
			Name: StepBuild,
			Applies: func(s *DeployStateV6) bool {
				return !s.Opts.SkipBuild && s.Conf.NotFirstRun
			},
			Run: func(_ context.Context, s *DeployStateV6) error {
				defer s.writeConf()
				return s.Cmd.Build(s.Project, s.OldPrefix, s.NewPrefix, s.msgs)
			},
		},
		{
			Name: StepReadManifest,
			Run: func(ctx context.Context, s *DeployStateV6) error {
				return s.readManifest(ctx, !s.runsBundlerInit())
			},
		},
		{
			Name: StepBucket,
			Applies: func(s *DeployStateV6) bool {
				return declares(s, func(m *contracts.ManifestV6) int { return len(m.Storage) })
			},
			Run: func(ctx context.Context, s *DeployStateV6) error {
				defer s.writeConf()
				return s.Cmd.EnsureBucket(ctx, deploy.StorageClient(s.Cmd), s.Project, s.msgs, s.Manifest.Storage)
			},
		},
		{
			Name:    StepBundlerInit,
			Applies: (*DeployStateV6).runsBundlerInit,
			Run: func(_ context.Context, s *DeployStateV6) error {
				defer s.writeConf()
				return s.Cmd.BundlerInit(s.Project, s.NewPrefix, s.msgs)
			},
		},
		{
			// The build above may have regenerated the manifest.
			Name: StepReadManifestAgain,
			Run: func(ctx context.Context, s *DeployStateV6) error {
				return s.readManifest(ctx, true)
			},
		},
		{
			Name:    StepCredentials,
			Applies: (*DeployStateV6).shouldUpload,
			Run: func(ctx context.Context, s *DeployStateV6) error {
				creds, err := s.Cmd.GetOrCreateCredentials(ctx, s.Project.Bucket, s.F.GetActiveProfile())
				if err != nil {
					return err
				}
				s.Creds = creds
				return nil
			},
		},
		{
			Name:    StepUploadStaticFiles,
			Applies: (*DeployStateV6).shouldUpload,
			Run: func(_ context.Context, s *DeployStateV6) error {
				return s.Cmd.UploadStatic(s.Project, s.msgs, s.Manifest.Storage, s.Creds)
			},
		},

		// Level 1: what the Application and Firewall depend on.
		applyStep(StepManifestFunctions, func(m *contracts.ManifestV6) int { return len(m.Functions) }, (*manifestv6.Applier).ApplyFunctions, false),
		applyStep(StepManifestNetworkLists, func(m *contracts.ManifestV6) int { return len(m.NetworkLists) }, (*manifestv6.Applier).ApplyNetworkLists, false),
		applyStep(StepManifestWAFs, func(m *contracts.ManifestV6) int { return len(m.WAFs) }, (*manifestv6.Applier).ApplyWAFs, false),
		applyStep(StepManifestConnectors, func(m *contracts.ManifestV6) int { return len(m.Connectors) }, (*manifestv6.Applier).ApplyConnectors, false),
		{
			Name:    StepVersionBuildDeps,
			Applies: func(s *DeployStateV6) bool { return declares(s, dependencyCount) },
			Run: func(ctx context.Context, s *DeployStateV6) error {
				return s.buildLevel(ctx, s.Dependencies)
			},
		},

		// Level 2: what a release serves.
		applyStep(StepManifestApplication, func(m *contracts.ManifestV6) int { return len(m.Applications) }, (*manifestv6.Applier).ApplyApplications, true),
		applyStep(StepManifestFirewalls, func(m *contracts.ManifestV6) int { return len(m.Firewalls) }, (*manifestv6.Applier).ApplyFirewalls, true),
		applyStep(StepManifestCustomPages, func(m *contracts.ManifestV6) int { return len(m.CustomPages) }, (*manifestv6.Applier).ApplyCustomPages, true),
		{
			Name:    StepVersionBuildPublishable,
			Applies: func(s *DeployStateV6) bool { return declares(s, publishableCount) },
			Run: func(ctx context.Context, s *DeployStateV6) error {
				return s.buildLevel(ctx, s.Publishable)
			},
		},

		{Name: StepManifestEnvironments, Run: method((*DeployStateV6).applyEnvironments)},
		{Name: StepManifestDeployments, Run: method((*DeployStateV6).applyDeployments)},
		{Name: StepReleaseCreate, Run: method((*DeployStateV6).releaseCreate)},
		{Name: StepReleaseBuild, Run: method((*DeployStateV6).releaseBuild)},
		{
			Name:    StepReleaseActivate,
			Applies: func(s *DeployStateV6) bool { return !s.Opts.NoActivate },
			Run:     method((*DeployStateV6).releaseActivate),
		},
		{
			// Without an activation the Workload is left alone: binding a
			// Deployment whose release is not live would serve nothing.
			Name: StepManifestWorkloads,
			Applies: func(s *DeployStateV6) bool {
				return !s.Opts.NoActivate && declares(s, func(m *contracts.ManifestV6) int { return len(m.Workloads) })
			},
			Run: method((*DeployStateV6).applyWorkloads),
		},
		{
			Name:    StepVerifyDeploy,
			Applies: func(s *DeployStateV6) bool { return !s.Opts.NoActivate },
			Run:     method((*DeployStateV6).verifyDeploy),
		},
		{
			Name: StepManifestPurge,
			Applies: func(s *DeployStateV6) bool {
				return !s.Opts.NoActivate && declares(s, func(m *contracts.ManifestV6) int { return len(m.Purge) })
			},
			Run: method((*DeployStateV6).purge),
		},
		{
			Name: StepOrphans,
			Run: func(_ context.Context, s *DeployStateV6) error {
				for _, w := range s.Applier.Orphans() {
					s.warn(w)
				}
				return nil
			},
		},
	}
}

// method adapts a state method to a step's Run.
func method(f func(*DeployStateV6, context.Context) error) func(context.Context, *DeployStateV6) error {
	return func(ctx context.Context, s *DeployStateV6) error { return f(s, ctx) }
}

func dependencyCount(m *contracts.ManifestV6) int {
	return len(m.Functions) + len(m.NetworkLists) + len(m.WAFs) + len(m.Connectors)
}

func publishableCount(m *contracts.ManifestV6) int {
	return len(m.Applications) + len(m.Firewalls) + len(m.CustomPages)
}

// applyStep is a step that writes one resource kind into Drafts and queues
// them for the build of their level.
func applyStep(name string, count func(*contracts.ManifestV6) int,
	apply func(*manifestv6.Applier, context.Context) ([]manifestv6.Pending, error), publishable bool,
) pipeline.Step[DeployStateV6] {
	return pipeline.Step[DeployStateV6]{
		Name:    name,
		Applies: func(s *DeployStateV6) bool { return declares(s, count) },
		Run: func(ctx context.Context, s *DeployStateV6) error {
			// Reaching the resources means this is no longer a first run.
			s.Conf.NotFirstRun = true
			pending, err := apply(s.Applier, ctx)
			if err != nil {
				return err
			}
			if publishable {
				s.Publishable = append(s.Publishable, pending...)
			} else {
				s.Dependencies = append(s.Dependencies, pending...)
			}
			return nil
		},
	}
}

// runsBundlerInit reports whether this run initialises the bundler: a first
// run, unless both the build and the framework setup are skipped. That step
// replaces the template's placeholders and builds the manifest again.
func (s *DeployStateV6) runsBundlerInit() bool {
	return !s.Conf.NotFirstRun && (!s.Opts.SkipBuild || !s.Opts.SkipFramework)
}

// readManifest reads the manifest, adds the Environment and Deployment it
// leaves to the CLI (withDefaults), then validates it, so a binding to the
// default Deployment is valid.
//
// validate is false only for the read before BundlerInit on a first run: the
// manifest there still has the template's placeholders, which BundlerInit
// replaces before the second read validates it.
func (s *DeployStateV6) readManifest(_ context.Context, validate bool) error {
	m, err := manifestv6.Parse(deploy.ManifestFile)
	if err != nil {
		logger.Debug("Error while reading the v6 manifest", zap.Error(err))
		return err
	}
	envName := s.Opts.Environment
	if envName == "" {
		envName = DefaultEnvironment
	}
	s.mustExist = withDefaults(m, envName, s.Conf.Name)
	if validate {
		if err := manifestv6.Validate(m); err != nil {
			logger.Debug("Error while validating the v6 manifest", zap.Error(err))
			return err
		}
	}
	s.Manifest = m
	s.Applier = &manifestv6.Applier{
		Clients:     s.Clients.Resources,
		Conf:        s.Conf,
		Manifest:    m,
		Say:         s.say,
		WriteConfig: func() error { s.writeConf(); return nil },
	}
	return nil
}

func (s *DeployStateV6) purge(ctx context.Context) error {
	for _, p := range s.Manifest.Purge {
		layer := "cache"
		if p.Layer != nil {
			layer = *p.Layer
		}
		if err := s.Clients.Purge.PurgeCache(ctx, p.Items, p.Type, layer); err != nil {
			return err
		}
		s.say(fmt.Sprintf(manifestmsg.ManifestPurgeSuccess, p.Type))
	}
	return nil
}
