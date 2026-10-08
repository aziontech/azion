package deploy

import (
	"context"
	"fmt"
	"os"

	msg "github.com/aziontech/azion-cli/messages/dryrun"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/output"
	"github.com/aziontech/azion-cli/pkg/pipeline"
	manifestInt "github.com/aziontech/azion-cli/pkg/v3manifest"
	"github.com/aziontech/azion-cli/utils"
	"go.uber.org/zap"
)

// DryRun describes what `azion deploy` would do, without doing any of it.
//
// It asks the real v3 pipeline which steps apply — the same table Run executes —
// so the description cannot drift from the deploy. Nothing here issues an API
// request: dry-run mode evaluates the predicates and never calls a step's Run.
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
		Manifest:    &contracts.Manifest{},
	}

	report(msg.Header)

	manifestRead := false
	if path, err := interpreter.ManifestPath(); err == nil {
		if _, statErr := os.Stat(path); statErr == nil {
			if m, readErr := interpreter.ReadManifest(path, f, &msgs); readErr == nil {
				state.Manifest = m
				manifestRead = true
			}
		}
	}

	res, err := pipeline.Run(state.Ctx, V3DeployPipeline{}, state, &pipeline.Plan{DryRun: true})
	if err != nil {
		return err
	}

	for _, step := range res.Steps {
		switch step.Name {
		case "Application":
			if conf.Application.ID == 0 {
				report(fmt.Sprintf(msg.CreateEdgeApp, conf.Name))
			} else {
				report(fmt.Sprintf(msg.UpdateEdgeApp, conf.Application.ID, conf.Name))
			}
		case "FirstRunOriginAndRules":
			// One step, three announcements: it creates the single origin,
			// points the default rule at it and deploys the project's rules.
			single := utils.Concat(conf.Name, "_single")
			report(fmt.Sprintf(msg.CreateOriginSingle, single))
			report(fmt.Sprintf(msg.UpdateDefaultRule, single))
			if len(conf.RulesEngine.Rules) == 0 {
				report(msg.CreateRulesCache)
				logger.Debug("", zap.Any("Cache Setting information", msg.AskCreateCacheSettings))
			}
		case "Bucket":
			report(fmt.Sprintf(msg.CreateBucket, conf.Name))
		case "UploadStaticFiles":
			report(fmt.Sprintf(msg.UploadStaticFiles, conf.Bucket))
		case "Function":
			if conf.Function.ID == 0 {
				report(fmt.Sprintf(msg.CreateFunction, conf.Name))
			} else {
				report(fmt.Sprintf(msg.UpdateFunction, conf.Function.ID, conf.Name))
			}
		case "Manifest":
			if !manifestRead {
				report(msg.SkipManifest)
				break
			}
			steps, err := manifestInt.DryRunSteps(f, conf, state.Manifest, cmd.ProjectConf)
			if err != nil {
				return err
			}
			report(msg.ManifestHeader)
			for _, name := range steps {
				label := manifestStepLabel(name)
				n, counted := manifestStepCount(state.Manifest, name)
				switch {
				case !counted:
					report(fmt.Sprintf(msg.ManifestStep, label))
				case n == 0:
					// The v3 steps run even when the manifest declares none,
					// and that is the point: they empty the corresponding
					// section of azion.json. Saying "origins" alone would read
					// as if the manifest declared some.
					report(fmt.Sprintf(msg.ManifestStepCleared, label))
				default:
					report(fmt.Sprintf(msg.ManifestStepCount, label, n))
				}
			}
		case "Domain":
			if conf.Domain.Id == 0 {
				report(fmt.Sprintf(msg.CreateDomain, conf.Name))
			} else {
				report(fmt.Sprintf(msg.UpdateDomain, conf.Domain.Id, conf.Name))
			}
		}
		// Steps with no case are plumbing and have nothing to announce.
	}

	return output.Print(&output.SliceOutput{
		Messages: msgs,
		GeneralOutput: output.GeneralOutput{
			Out:   cmd.F.IOStreams.Out,
			Flags: cmd.F.Flags,
		},
	})
}

// manifestStepCount reports how many of a resource the manifest declares, for
// the steps where a count is meaningful. It describes the manifest; it does not
// decide what runs — that is the pipeline's job.
func manifestStepCount(m *contracts.Manifest, step string) (int, bool) {
	switch step {
	case "ManifestOrigins":
		return len(m.Origins), true
	case "ManifestCacheSettings":
		return len(m.CacheSettings), true
	case "ManifestRulesEngine":
		return len(m.Rules), true
	case "ManifestPurge":
		return len(m.Purge), true
	default:
		// The domain and the orphan cleanup are not counted.
		return 0, false
	}
}

// manifestStepLabel turns a manifest step name into something readable.
func manifestStepLabel(name string) string {
	switch name {
	case "ManifestDomain":
		return "domain"
	case "ManifestOrigins":
		return "origins"
	case "ManifestCacheSettings":
		return "cache settings"
	case "ManifestRulesEngine":
		return "rules engine"
	case "ManifestPurge":
		return "purge"
	case "ManifestDeleteOrphanedResources":
		return "removing resources the manifest no longer declares"
	default:
		return name
	}
}
