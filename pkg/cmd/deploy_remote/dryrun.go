package deploy

import (
	"context"
	"fmt"
	"os"

	msg "github.com/aziontech/azion-cli/messages/dryrun"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/logger"
	manifestInt "github.com/aziontech/azion-cli/pkg/manifest"
	"github.com/aziontech/azion-cli/pkg/output"
	"github.com/aziontech/azion-cli/pkg/pipeline"
	"go.uber.org/zap"
)

// DryRun describes what `azion deploy` would do, without doing any of it.
//
// It asks the real pipeline which steps apply — the same table Run executes —
// so the description cannot drift from the deploy. pipeline.Run in dry-run mode
// evaluates the predicates and never calls a step's Run, and nothing here
// issues an API request.
//
// The simulation this replaces kept its own copy of the conditions and had
// drifted: it announced a single origin and a default-rule update, which do not
// exist in the v4 flow, and gated the bucket on the project preset rather than
// on the manifest declaring storage.
func (cmd *DeployCmd) DryRun(f *cmdutil.Factory) error {
	msgs := []string{}
	report := func(line string) {
		msgs = append(msgs, line)
		logger.FInfoFlags(cmd.F.IOStreams.Out, line, f.Format, f.Out)
	}

	conf, err := cmd.GetAzionJsonContent(cmd.ProjectConf)
	if err != nil {
		logger.Debug("Failed to get Azion JSON content", zap.Error(err))
		return err
	}

	interpreter := cmd.Interpreter()
	state := &DeployState{
		Cmd: cmd, F: f, Ctx: context.Background(), Conf: conf, Msgs: &msgs,
		Clients:     NewClients(f),
		Interpreter: interpreter,
		Manifest:    &contracts.ManifestV4{},
	}

	report(msg.Header)

	// The manifest drives most of the predicates, so it is read up front —
	// which is what the step would do. A project that has not been built yet
	// has none, and the steps that depend on it are reported as not applying.
	manifestRead := false
	if path, err := interpreter.ManifestPath(); err == nil {
		if _, statErr := os.Stat(path); statErr == nil {
			if m, readErr := interpreter.ReadManifest(path, f, &msgs); readErr == nil {
				state.Manifest = m
				manifestRead = true
			}
		}
	}

	res, err := pipeline.Run(state.Ctx, V4DeployPipeline{}, state, &pipeline.Plan{DryRun: true})
	if err != nil {
		return err
	}

	for _, step := range res.Steps {
		switch step.Name {
		case "Build":
			report(msg.Build)
		case "Application":
			if conf.Application.ID == 0 {
				report(fmt.Sprintf(msg.CreateEdgeApp, conf.Name))
			} else {
				report(fmt.Sprintf(msg.UpdateEdgeApp, conf.Application.ID, conf.Name))
			}
		case "Bucket":
			report(fmt.Sprintf(msg.CreateBucket, conf.Name))
		case "BundlerInit":
			report(msg.BundlerInit)
		case "UploadStaticFiles":
			report(fmt.Sprintf(msg.UploadStaticFiles, conf.Bucket))
		case "RulesEngine":
			// v4 creates the preset's default rules. It does not offer a cache
			// setting — that prompt belongs to the v3 flow.
			report(fmt.Sprintf(msg.CreateRulesEngine, conf.Preset))
		case "Manifest":
			if !manifestRead {
				report(msg.SkipManifest)
				break
			}
			if err := reportManifest(f, conf, state.Manifest, cmd.ProjectConf, report); err != nil {
				return err
			}
		case "Workload":
			if conf.Workloads.Id == 0 {
				report(fmt.Sprintf(msg.CreateWorkload, conf.Name))
			} else {
				report(fmt.Sprintf(msg.UpdateWorkload, conf.Workloads.Id, conf.Name))
			}
		}
		// Steps with no case are plumbing — resolving the manifest path, reading
		// it, obtaining storage credentials — and have nothing to announce.
	}

	return output.Print(&output.SliceOutput{
		Messages: msgs,
		GeneralOutput: output.GeneralOutput{
			Out:   cmd.F.IOStreams.Out,
			Flags: cmd.F.Flags,
		},
	})
}

// reportManifest descends into the manifest pipeline and names the resources it
// would apply.
func reportManifest(
	f *cmdutil.Factory,
	conf *contracts.AzionApplicationOptions,
	manifest *contracts.ManifestV4,
	projectConf string,
	report func(string),
) error {
	steps, err := manifestInt.DryRunSteps(f, conf, manifest, projectConf)
	if err != nil {
		return err
	}
	report(msg.ManifestHeader)
	for _, name := range steps {
		label := manifestStepLabel(name)
		if n, counted := manifestStepCount(manifest, name); counted {
			report(fmt.Sprintf(msg.ManifestStepCount, label, n))
			continue
		}
		report(fmt.Sprintf(msg.ManifestStep, label))
	}
	return nil
}

// manifestStepCount reports how many of a resource the manifest declares, for
// the steps where a count is meaningful. It describes the manifest; it does not
// decide what runs — that is the pipeline's job.
func manifestStepCount(m *contracts.ManifestV4, step string) (int, bool) {
	app := contracts.Applications{}
	if len(m.Applications) > 0 {
		app = m.Applications[0]
	}
	switch step {
	case "ManifestFunctions":
		return len(m.Functions), true
	case "ManifestFunctionInstances":
		return len(app.FunctionsInstances), true
	case "ManifestCacheSettings":
		return len(app.CacheSettings), true
	case "ManifestConnectors":
		return len(m.Connectors), true
	case "ManifestRulesEngine":
		return len(app.Rules), true
	case "ManifestWorkloads":
		return len(m.Workloads), true
	case "ManifestWorkloadDeployments":
		return len(m.WorkloadDeployments), true
	case "ManifestFirewalls":
		return len(m.Firewalls), true
	case "ManifestPurge":
		return len(m.Purge), true
	default:
		// The application and the orphan cleanup are not counted.
		return 0, false
	}
}

// manifestStepLabel turns a manifest step name into something readable.
func manifestStepLabel(name string) string {
	switch name {
	case "ManifestFunctions":
		return "functions"
	case "ManifestFunctionInstances":
		return "function instances"
	case "ManifestEdgeApplication":
		return "application"
	case "ManifestCacheSettings":
		return "cache settings"
	case "ManifestConnectors":
		return "connectors"
	case "ManifestRulesEngine":
		return "rules engine"
	case "ManifestWorkloads":
		return "workloads"
	case "ManifestWorkloadDeployments":
		return "workload deployments"
	case "ManifestFirewalls":
		return "firewalls"
	case "ManifestPurge":
		return "purge"
	case "ManifestDeleteOrphanedResources":
		return "removing resources the manifest no longer declares"
	case "ManifestDomain":
		return "domain"
	case "ManifestOrigins":
		return "origins"
	default:
		return name
	}
}
