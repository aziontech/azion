package deploy

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	msg "github.com/aziontech/azion-cli/messages/deploy"
	apiEdgeApplications "github.com/aziontech/azion-cli/pkg/api/edge_applications"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/utils"
	sdk "github.com/aziontech/azionapi-go-sdk/edgeapplications"
	"go.uber.org/zap"
	"strings"

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

	err := cmd.uploadFiles(f, conf.Prefix)
	if err != nil {
		return err
	}

	err = cmd.doApplication(clients.EdgeApplication, ctx, conf)
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

			// TODO: Review what to do when user updates Function ID directly in azion.json
			err = cmd.updateRulesEngine(clients.EdgeApplication, ctx, conf)
			if err != nil {
				logger.Debug("Error while updating rules engine", zap.Error(err))
				return err
			}
		}

		err = cmd.doOrigin(clients.EdgeApplication, clients.Origin, ctx, conf)
		if err != nil {
			logger.Debug("Error while creating origin", zap.Error(err))
			return err
		}
		logger.FInfo(cmd.F.IOStreams.Out, msg.OriginsSuccessful)

		requestRules, err := requestRulesEngineManifest(conf.Origin.ID, InstanceID, route)
		if err != nil {
			return err
		}

		_, err = clients.EdgeApplication.CreateRulesEngine(ctx, conf.Application.ID, "request", &requestRules)
		if err != nil {
			return err
		}

		err = cmd.WriteAzionJsonContent(conf)
		if err != nil {
			logger.Debug("Error while writing azion.json file", zap.Error(err))
			return err
		}
	}

	logger.FInfo(cmd.F.IOStreams.Out, msg.DeploySuccessful)
	logger.FInfo(cmd.F.IOStreams.Out, fmt.Sprintf(msg.DeployOutputDomainSuccess, "https://"+domainName))
	logger.FInfo(cmd.F.IOStreams.Out, msg.DeployPropagation)
	return nil
}

func requestRulesEngineManifest(originID, functionID int64, routes Routes) (apiEdgeApplications.CreateRulesEngineRequest, error) {
	logger.Debug("Create Rules Engine set origin")

	req := apiEdgeApplications.CreateRulesEngineRequest{}
	req.SetName(fmt.Sprintf("rules_manifest_%s", thoth.GenerateName()))

	behaviors := make([]sdk.RulesEngineBehaviorEntry, 0)
	var behStringCache sdk.RulesEngineBehaviorString

	if routes.Type == "compute" {
		behStringCache.SetName("run_function")
		behStringCache.SetTarget(fmt.Sprintf("%d", functionID))
	} else {
		behStringCache.SetName("set_origin")
		behStringCache.SetTarget(fmt.Sprintf("%d", originID))
	}

	behaviors = append(behaviors, sdk.RulesEngineBehaviorEntry{
		RulesEngineBehaviorString: &behStringCache,
	})

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
		doesNotMatch := "does_not_match"
		return doesNotMatch, nil
	}

	return "", errors.New("the value of 'from' not recognized")
}
