package deploy

import (
	"context"
	"os"
	"strconv"

	msg "github.com/aziontech/azion-cli/messages/deploy-remote"
	apiEdgeApplications "github.com/aziontech/azion-cli/pkg/api/v3/edge_applications"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/pipeline"
	manifestInt "github.com/aziontech/azion-cli/pkg/v3manifest"
	sdk "github.com/aziontech/azionapi-go-sdk/edgeapplications"
	"go.uber.org/zap"
)

// DeployState is what the v3 deploy steps read and write.
type DeployState struct {
	Cmd  *DeployCmd
	F    *cmdutil.Factory
	Ctx  context.Context
	Conf *contracts.AzionApplicationOptionsV3
	Msgs *[]string

	Clients     *Clients
	Interpreter *manifestInt.ManifestInterpreter

	// PathManifest is set by the manifest-path step and read by the read step.
	PathManifest string
	// Manifest is set by the read step and read by the apply and domain steps.
	Manifest *contracts.Manifest
}

// V3DeployPipeline is the deploy flow for the v3 API.
//
// It starts later than the v4 one. v3 builds before it reads azion.json, so the
// build stays in Run's setup alongside the config read — moving it after would
// change which of the two runs first. The pipeline begins where the state is
// complete, at the first resource operation.
//
// There is no timing summary here: v3 never had one.
type V3DeployPipeline struct{}

func (V3DeployPipeline) Name() string { return "deploy-v3" }

func (V3DeployPipeline) Steps() []pipeline.Step[DeployState] {
	return []pipeline.Step[DeployState]{
		{
			Name: "ManifestPath",
			Run: func(_ context.Context, s *DeployState) error {
				path, err := s.Interpreter.ManifestPath()
				if err != nil {
					return err
				}
				s.PathManifest = path
				return nil
			},
		},
		{
			Name: "Application",
			Run: func(_ context.Context, s *DeployState) error {
				return s.Cmd.doApplication(s.Clients.EdgeApplication, context.Background(), s.Conf, s.Msgs)
			},
		},
		{
			// On a first run the project needs its single origin and the
			// default rule pointed at it before anything else is applied.
			Name:    "FirstRunOriginAndRules",
			Applies: func(s *DeployState) bool { return !s.Conf.NotFirstRun },
			Run:     func(_ context.Context, s *DeployState) error { return s.firstRunOriginAndRules() },
		},
		{
			Name: "Bucket",
			Run: func(_ context.Context, s *DeployState) error {
				return s.Cmd.doBucket(s.Clients.Bucket, s.Ctx, s.Conf, s.Msgs)
			},
		},
		{
			Name: "UploadStaticFiles",
			Applies: func(s *DeployState) bool {
				// Kept here so the reason for skipping still reaches the debug
				// log, as it did when this was an if/else.
				if _, err := os.Stat(PathStatic); os.IsNotExist(err) {
					logger.Debug(msg.SkipUpload)
					return false
				}
				return true
			},
			Run: func(_ context.Context, s *DeployState) error {
				return s.Cmd.uploadFiles(s.F, s.Conf, s.Msgs)
			},
		},
		{
			Name: "Function",
			Run: func(_ context.Context, s *DeployState) error {
				s.Conf.Function.File = ".edge/worker.js"
				return s.Cmd.doFunction(s.Clients, s.Ctx, s.Conf, s.Msgs)
			},
		},
		{
			Name: "ReadManifest",
			Run: func(_ context.Context, s *DeployState) error {
				m, err := s.Interpreter.ReadManifest(s.PathManifest, s.F, s.Msgs)
				if err != nil {
					return err
				}
				s.Manifest = m
				return nil
			},
		},
		{
			Name: "Manifest",
			Run: func(_ context.Context, s *DeployState) error {
				return s.Interpreter.CreateResources(s.Conf, s.Manifest, s.F, s.Cmd.ProjectConf, s.Msgs)
			},
		},
		{
			// A manifest that declares no domain, or declares a named one,
			// still gets the project's own domain created here.
			Name: "Domain",
			Applies: func(s *DeployState) bool {
				return s.Manifest.Domain == nil || s.Manifest.Domain.Name != ""
			},
			Run: func(_ context.Context, s *DeployState) error {
				return s.Cmd.doDomain(s.Clients.Domain, s.Ctx, s.Conf, s.Msgs)
			},
		},
	}
}

// firstRunOriginAndRules creates the single origin, points the default rule at
// it, and deploys the project's rules.
//
// Every failure here tears the application down again with a cascade delete, so
// a half-created project is not left behind. That is why the block is one step:
// splitting it would put the cascade in three places.
func (s *DeployState) firstRunOriginAndRules() error {
	cmd, f, ctx, conf := s.Cmd, s.F, s.Ctx, s.Conf

	singleOriginId, err := cmd.doOriginSingle(s.Clients.Origin, ctx, conf, s.Msgs)
	if err != nil {
		return err
	}

	ruleDefaultID, err := s.Clients.EdgeApplication.GetRulesDefault(ctx, conf.Application.ID, "request")
	if err != nil {
		logger.Debug("Error while getting default rules engine", zap.Error(err))
		errCascade := callDeleteCascade(f)
		if errCascade != nil {
			return errCascade
		}
		return err
	}
	behaviors := make([]sdk.RulesEngineBehaviorEntry, 0)

	var behString sdk.RulesEngineBehaviorString
	behString.SetName("set_origin")

	behString.SetTarget(strconv.Itoa(int(singleOriginId)))

	behaviors = append(behaviors, sdk.RulesEngineBehaviorEntry{
		RulesEngineBehaviorString: &behString,
	})

	reqUpdateRulesEngine := apiEdgeApplications.UpdateRulesEngineRequest{
		IdApplication: conf.Application.ID,
		Phase:         "request",
		Id:            ruleDefaultID,
	}

	reqUpdateRulesEngine.SetBehaviors(behaviors)

	_, err = s.Clients.EdgeApplication.UpdateRulesEngine(ctx, &reqUpdateRulesEngine)
	if err != nil {
		logger.Debug("Error while updating default rules engine", zap.Error(err))
		errCascade := callDeleteCascade(f)
		if errCascade != nil {
			return errCascade
		}
		return err
	}

	if len(conf.RulesEngine.Rules) == 0 {
		err = cmd.doRulesDeploy(ctx, conf, s.Clients.EdgeApplication, s.Msgs)
		if err != nil {
			errCascade := callDeleteCascade(f)
			if errCascade != nil {
				return errCascade
			}
			return err
		}
	}

	return nil
}
