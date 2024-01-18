package deploy

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	msg "github.com/aziontech/azion-cli/messages/deploy"
	apiEdgeApplications "github.com/aziontech/azion-cli/pkg/api/edge_applications"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/utils"
	sdk "github.com/aziontech/azionapi-go-sdk/edgeapplications"
	"go.uber.org/zap"

	thoth "github.com/aziontech/go-thoth"
)

type Manifest struct {
	Routes []Routes `json:"routes"`
	Fs     []any    `json:"fs"`
}

type Routes struct {
	From     string `json:"from"`
	To       string `json:"to"`
	Priority int    `json:"priority"`
	Type     string `json:"type"`
}

var manifestFilePath = "/.edge/manifest.json"

func readManifest(cmd *DeployCmd) (*Manifest, error) {
	logger.Debug("read manifest")
	pathWorkingDir, err := cmd.GetWorkDir()
	if err != nil {
		return nil, err
	}

	b, err := cmd.FileReader(utils.Concat(pathWorkingDir, manifestFilePath))
	if err != nil {
		return nil, err
	}

	manifest := Manifest{}
	err = json.Unmarshal(b, &manifest)
	if err != nil {
		return nil, err
	}

	return &manifest, err
}

// Interpreted TODO: better interpreted, removed flows edge application, domain
func (manifest *Manifest) Interpreted(f *cmdutil.Factory, cmd *DeployCmd, conf *contracts.AzionApplicationOptions, clients *Clients) error {
	logger.Debug("Execute manifest")
	ctx := context.Background()

	err := cmd.doApplication(clients.EdgeApplication, ctx, conf)
	if err != nil {
		return err
	}

	domainName, err := cmd.doDomain(clients.Domain, ctx, conf)
	if err != nil {
		return err
	}

	err = cmd.doBucket(clients.Bucket, ctx, conf)
	if err != nil {
		return err
	}

	// skip upload when type = javascript, typescript (storage folder does not exist in this case)
	if conf.Template != "javascript" && conf.Template != "typescript" {
		err = cmd.uploadFiles(f, conf)
		if err != nil {
			return err
		}
	}

	for _, route := range manifest.Routes {
		if route.Type == "compute" {
			conf.Function.File = ".edge/worker.js"
			err := cmd.doFunction(clients.EdgeFunction, ctx, conf)
			if err != nil {
				return err
			}

			reqIns := apiEdgeApplications.CreateInstanceRequest{}
			reqIns.SetEdgeFunctionId(conf.Function.ID)
			reqIns.SetName(conf.Name)
			reqIns.ApplicationId = conf.Application.ID
			instance, err := clients.EdgeApplication.CreateInstancePublish(ctx, &reqIns)
			if err != nil {
				logger.Debug("Error while creating edge function instance", zap.Error(err))
				return fmt.Errorf(msg.ErrorCreateInstance.Error(), err)
			}
			InstanceID = instance.GetId()
		}

		err = cmd.doOrigin(clients.EdgeApplication, clients.Origin, ctx, conf, route.Type)
		if err != nil {
			logger.Debug("Error while creating origin", zap.Error(err))
			return err
		}
		logger.FInfo(cmd.F.IOStreams.Out, msg.OriginsSuccessful)

		ruleDefaultID, err := clients.EdgeApplication.GetRulesDefault(ctx, conf.Application.ID, "request")
		if err != nil {
			logger.Debug("Error while getting default rules engine", zap.Error(err))
			return err
		}

		behaviors := make([]sdk.RulesEngineBehaviorEntry, 0)

		var behString sdk.RulesEngineBehaviorString
		behString.SetName("set_origin")

		if route.Type == "compute" {
			behString.SetTarget(strconv.Itoa(int(conf.Origin.SingleOriginID)))
		} else {
			behString.SetTarget(strconv.Itoa(int(conf.Origin.StorageOriginID)))
		}

		behaviors = append(behaviors, sdk.RulesEngineBehaviorEntry{
			RulesEngineBehaviorString: &behString,
		})

		reqUpdateRulesEngine := apiEdgeApplications.UpdateRulesEngineRequest{
			IdApplication: conf.Application.ID,
			Phase:         "request",
			Id:            ruleDefaultID,
		}

		reqUpdateRulesEngine.SetBehaviors(behaviors)

		_, err = clients.EdgeApplication.UpdateRulesEngine(ctx, &reqUpdateRulesEngine)
		if err != nil {
			logger.Debug("Error while updating default rules engine", zap.Error(err))
			return err
		}

		requestRules, err := requestRulesEngineManifest(conf.Origin.StorageOriginID, InstanceID, route)
		if err != nil {
			return err
		}

		_, err = clients.EdgeApplication.CreateRulesEngine(ctx, conf.Application.ID, "request", &requestRules)
		if err != nil {
			return err
		}

		switch route.From {
		case "/_next/static/":
			logger.Debug("Create Rules to route /_next/static/ rewrite")

			reqNextStatic := apiEdgeApplications.CreateRulesEngineRequest{}
			reqNextStatic.SetName("rule_rewrite_next_static")

			behaviors := make([]sdk.RulesEngineBehaviorEntry, 0)

			// ---------------------------------
			// capture match groups
			// ---------------------------------

			var behCaptureMatchGroups sdk.RulesEngineBehaviorObject

			behCaptureMatchGroups.SetName("capture_match_groups")

			behTarget := sdk.RulesEngineBehaviorObjectTarget{}
			behTarget.SetCapturedArray("capture")
			behTarget.SetSubject("${uri}")
			behTarget.SetRegex("/_next/static/(.*)")

			behCaptureMatchGroups.SetTarget(behTarget)

			behaviors = append(behaviors, sdk.RulesEngineBehaviorEntry{
				RulesEngineBehaviorObject: &behCaptureMatchGroups,
			})

			// ---------------------------------
			// rewrite request 
			// ---------------------------------

			var behRewriteRequest sdk.RulesEngineBehaviorString

			behRewriteRequest.SetName("rewrite_request")
			behRewriteRequest.SetTarget("/.next/static/%{capture[1]}")

			behaviors = append(behaviors, sdk.RulesEngineBehaviorEntry{
				RulesEngineBehaviorString: &behRewriteRequest,
			})

			// ---------------------------------
			// deliver 
			// ---------------------------------

			var behDeliver sdk.RulesEngineBehaviorString

			behDeliver.SetName("deliver")

			behaviors = append(behaviors, sdk.RulesEngineBehaviorEntry{
				RulesEngineBehaviorString: &behDeliver,
			})

			reqNextStatic.SetBehaviors(behaviors)
			
			criteria := make([][]sdk.RulesEngineCriteria, 1)
			for i := 0; i < 1; i++ {
				criteria[i] = make([]sdk.RulesEngineCriteria, 1)
			}
					
			criteria[0][0].SetConditional("if")
			criteria[0][0].SetVariable("${uri}")
			criteria[0][0].SetOperator("starts_with")
			criteria[0][0].SetInputValue(route.From)
			reqNextStatic.SetCriteria(criteria)

			_, err = clients.EdgeApplication.CreateRulesEngine(ctx, conf.Application.ID, "request", &reqNextStatic)
			if err != nil {
				return err
			}
		case "\\.(css|js|ttf|woff|woff2|pdf|svg|jpg|jpeg|gif|bmp|png|ico|mp4)$":
			logger.Debug("Create Rules to route \\.(css|js|ttf|woff|woff2|pdf|svg|jpg|jpeg|gif|bmp|png|ico|mp4)$ rewrite")

			reqAssets := apiEdgeApplications.CreateRulesEngineRequest{}
			reqAssets.SetName("rule_rewrite_assets")

			behaviors := make([]sdk.RulesEngineBehaviorEntry, 0)

			// ---------------------------------
			// capture match groups
			// ---------------------------------

			var behCaptureMatchGroups sdk.RulesEngineBehaviorObject

			behCaptureMatchGroups.SetName("capture_match_groups")

			behTarget := sdk.RulesEngineBehaviorObjectTarget{}
			behTarget.SetCapturedArray("capture")
			behTarget.SetSubject("${uri}")
			behTarget.SetRegex("/(.*)")

			behCaptureMatchGroups.SetTarget(behTarget)

			behaviors = append(behaviors, sdk.RulesEngineBehaviorEntry{
				RulesEngineBehaviorObject: &behCaptureMatchGroups,
			})

			// ---------------------------------
			// rewrite request 
			// ---------------------------------

			var behRewriteRequest sdk.RulesEngineBehaviorString

			behRewriteRequest.SetName("rewrite_request")
			behRewriteRequest.SetTarget("/public/%{capture[1]}")

			behaviors = append(behaviors, sdk.RulesEngineBehaviorEntry{
				RulesEngineBehaviorString: &behRewriteRequest,
			})

			// ---------------------------------
			// deliver 
			// ---------------------------------

			var behDeliver sdk.RulesEngineBehaviorString

			behDeliver.SetName("deliver")

			behaviors = append(behaviors, sdk.RulesEngineBehaviorEntry{
				RulesEngineBehaviorString: &behDeliver,
			})

			reqAssets.SetBehaviors(behaviors)
			
			criteria := make([][]sdk.RulesEngineCriteria, 1)
			for i := 0; i < 1; i++ {
				criteria[i] = make([]sdk.RulesEngineCriteria, 1)
			}
					
			criteria[0][0].SetConditional("if")
			criteria[0][0].SetVariable("${uri}")
			criteria[0][0].SetOperator("matches")
			criteria[0][0].SetInputValue(route.From)
			reqAssets.SetCriteria(criteria)

			_, err = clients.EdgeApplication.CreateRulesEngine(ctx, conf.Application.ID, "request", &reqAssets)
			if err != nil {
				return err
			}
		}

		err = cmd.WriteAzionJsonContent(conf)
		if err != nil {
			logger.Debug("Error while writing azion.json file", zap.Error(err))
			return err
		}
	}

	logger.FInfo(cmd.F.IOStreams.Out, msg.DeploySuccessful)
	logger.FInfo(cmd.F.IOStreams.Out, fmt.Sprintf(msg.DeployOutputDomainSuccess, utils.Concat("https://", domainName)))
	logger.FInfo(cmd.F.IOStreams.Out, msg.DeployPropagation)
	return nil
}

func requestRulesEngineManifest(originID, functionID int64, routes Routes) (apiEdgeApplications.CreateRulesEngineRequest, error) {
	logger.Debug("Create Rules Engine set origin")

	req := apiEdgeApplications.CreateRulesEngineRequest{}
	req.SetName(fmt.Sprintf("rules_manifest_%s", thoth.GenerateName()))

	behaviors := make([]sdk.RulesEngineBehaviorEntry, 0)

	if routes.Type == "compute" {
		var behFunction sdk.RulesEngineBehaviorString
		behFunction.SetName("run_function")
		behFunction.SetTarget(fmt.Sprintf("%d", functionID))

		behaviors = append(behaviors, sdk.RulesEngineBehaviorEntry{
			RulesEngineBehaviorString: &behFunction,
		})
	} else {
		var behOrigin sdk.RulesEngineBehaviorString
		behOrigin.SetName("set_origin")
		behOrigin.SetTarget(fmt.Sprintf("%d", originID))

		behaviors = append(behaviors, sdk.RulesEngineBehaviorEntry{
			RulesEngineBehaviorString: &behOrigin,
		})
	}

	req.SetBehaviors(behaviors)

	criteria := make([][]sdk.RulesEngineCriteria, 1)
	for i := 0; i < 1; i++ {
		criteria[i] = make([]sdk.RulesEngineCriteria, 1)
	}

	operator, err := checkFieldFrom(routes.From)
	if err != nil {
		return req, err
	}

	criteria[0][0].SetConditional("if")
	criteria[0][0].SetVariable("${uri}")
	criteria[0][0].SetOperator(operator)
	criteria[0][0].SetInputValue(routes.From)
	req.SetCriteria(criteria)

	return req, nil
}

func checkFieldFrom(from string) (string, error) {
	if strings.HasPrefix(from, "/") {
		startWith := "starts_with"
		return startWith, nil
	} else if strings.HasPrefix(from, "\\.") && strings.HasSuffix(from, "$") {
		doesNotMatch := "matches"
		return doesNotMatch, nil
	}

	return "", errors.New("the value of 'from' could not be recognized")
}
