package deploy

import (
	"context"
	"errors"
	"fmt"
	"path"
	"strconv"
	"strings"

	sdkv3 "github.com/aziontech/azionapi-go-sdk/edgeapplications"
	sdk "github.com/aziontech/azionapi-v4-go-sdk/edge"

	thoth "github.com/aziontech/go-thoth"
	"go.uber.org/zap"

	msg "github.com/aziontech/azion-cli/messages/deploy"
	apiapp "github.com/aziontech/azion-cli/pkg/api/edge_applications"
	api "github.com/aziontech/azion-cli/pkg/api/edge_function"
	apiori "github.com/aziontech/azion-cli/pkg/api/origin"
	apiworkload "github.com/aziontech/azion-cli/pkg/api/workloads"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/logger"
	vulcanPkg "github.com/aziontech/azion-cli/pkg/vulcan"
	"github.com/aziontech/azion-cli/utils"
)

// inject this code into worker.js
var injectIntoFunction = `
//---
//storages:
//   - name: assets
//     bucket: %s
//     prefix: %s
//---

`

var applicationUpdate = `
{
	"name": "%s",
	"active": true,
	"modules": {
		"edgeFunctionsEnabled": true
	}
}`

var connectorUpdate = `
{
	"name": "%s",
	"active": true,
	"type": "edge_storage",
	"typeProperties": {
		"bucket": "%s",
		"prefix": "%s"
	},
	"modules": {
		"loadBalancerEnabled": false,
		"originShieldEnabled": false
	}
}`

var functionUpdate = `
{
	"name": "%s",
	"path": ".edge/worker.js",
	"bindings": {
		"storage": {
		"bucket": "%s",
		"prefix": "%s"
		}
	}
}`

var storageUpdate = `
{
	"name": "%s",
	"edgeAccess": "read_only",
	"dir": ".edge/storage"
}`

var jsonTemplate = `{
	"scope": "global",
  	"preset": "%s",
	"edgeApplications": [
	{
		"name": "%s",
		"active": true,
		"modules": {
    		"edgeFunctionsEnabled": true
  		}
	}
	],
	"edgeConnectors": [
	{
		"name": "%s",
		"active": true,
		"type": "edge_storage",
		"typeProperties": {
			"bucket": "%s",
			"prefix": "%s"
		},
		"modules": {
			"loadBalancerEnabled": true,
			"originShieldEnabled": true
		}
	}
	],
	"edgeFunctions": [
	  {
		"name": "%s",
		"path": ".edge/worker.js",
		"bindings": {
		  "storage": {
			"bucket": "%s",
			"prefix": "%s"
		  }
		}
	  }
	],
	"edgeStorage": [
	  {
		"name": "%s",
		"edgeAccess": "read_only",
		"dir": ".edge/storage"
	  }
	]
  }`

func (cmd *DeployCmd) callBundlerInit(conf *contracts.AzionApplicationOptions) error {
	logger.Debug("Running bundler config update to update azion.config")
	// checking if vulcan major is correct
	vulcanVer, err := cmd.commandRunnerOutput(cmd.F, "npm show edge-functions version", []string{})
	if err != nil {
		return err
	}

	vul := vulcanPkg.NewVulcan()

	err = vul.CheckVulcanMajor(vulcanVer, cmd.F, vul)
	if err != nil {
		return err
	}

	// cmdVulcanInit := "config update"
	applicationUpdate = fmt.Sprintf(applicationUpdate, conf.Name)
	connectorUpdate = fmt.Sprintf(connectorUpdate, conf.Name, conf.Bucket, conf.Prefix)
	functionUpdate = fmt.Sprintf(functionUpdate, conf.Name, conf.Bucket, conf.Prefix)
	storageUpdate = fmt.Sprintf(storageUpdate, conf.Bucket)
	commands := []string{
		fmt.Sprintf("config update -k 'edgeApplications[0]' -v '%s'", applicationUpdate),
		fmt.Sprintf("config update -k 'edgeFunctions[0]' -v '%s'", functionUpdate),
		fmt.Sprintf("config update -k 'edgeConnectors[0]' -v '%s'", connectorUpdate),
		fmt.Sprintf("config update -k 'edgeStorage[0]' -v '%s'", storageUpdate),
	}

	for _, cmdStr := range commands {
		command := vul.Command("", cmdStr, cmd.F)
		logger.Debug("Running the following command", zap.Any("Command", command))

		err := cmd.commandRunInteractive(cmd.F, command)
		if err != nil {
			return err
		}
	}

	return nil
}

func (cmd *DeployCmd) doFunction(clients *Clients, ctx context.Context, conf *contracts.AzionApplicationOptions, msgs *[]string) error {
	FunctionIds = make(map[string]contracts.AzionJsonDataFunction)
	if len(conf.Function) == 0 {
		newFunc := contracts.AzionJsonDataFunction{
			Name: conf.Name,
			Args: "./azion/args.json",
			File: ".edge/worker.js",
		}
		var projName string
		functionId, err := cmd.createFunction(clients.EdgeFunction, ctx, conf, newFunc, msgs)
		if err != nil {
			for i := 0; i < 10; i++ {
				projName = fmt.Sprintf("%s-%s", conf.Name, utils.Timestamp())
				newFunc.Name = projName
				functionId, err := cmd.createFunction(clients.EdgeFunction, ctx, conf, newFunc, msgs)
				if err != nil {
					fmt.Println(err)
					if errors.Is(err, utils.ErrorNameInUse) && i < 9 {
						continue
					}
					return fmt.Errorf(msg.ErrorCreateFunction.Error(), err)
				}
				newFunc.ID = functionId
				break
			}
		} else {
			newFunc.ID = functionId
		}

		conf.Function = append(conf.Function, newFunc)

		err = cmd.WriteAzionJsonContent(conf, ProjectConf)
		if err != nil {
			logger.Debug("Error while writing azion.json file", zap.Error(err))
			return err
		}

		for {
			instance, err := cmd.createInstance(ctx, clients.EdgeApplication, conf, newFunc)
			if err != nil {
				// if the name is already in use, we ask for another one
				if errors.Is(err, utils.ErrorNameInUse) {
					logger.FInfoFlags(cmd.Io.Out, msg.FuncInstInUse, cmd.F.Format, cmd.F.Out)
					*msgs = append(*msgs, msg.FuncInstInUse)
					if Auto {
						projName = thoth.GenerateName()
					} else {
						projName, err = askForInput(msg.AskInputName, thoth.GenerateName())
						if err != nil {
							return err
						}
					}
					newFunc.InstanceName = projName
					continue
				}
				return err
			}
			newFunc.InstanceID = instance.GetId()
			break
		}

		conf.Function[0] = newFunc
		FunctionIds[newFunc.Name] = newFunc
		err = cmd.WriteAzionJsonContent(conf, ProjectConf)
		if err != nil {
			logger.Debug("Error while writing azion.json file", zap.Error(err))
			return err
		}

		return nil
	}

	for _, funcCong := range conf.Function {
		FunctionIds[funcCong.Name] = funcCong

		_, err := cmd.updateFunction(clients.EdgeFunction, ctx, conf, funcCong, msgs)
		if err != nil {
			return err
		}

		_, err = cmd.updateInstance(ctx, clients.EdgeApplication, conf, funcCong)
		if err != nil {
			return err
		}

	}

	return nil
}

func (cmd *DeployCmd) doApplication(client *apiapp.Client, ctx context.Context, conf *contracts.AzionApplicationOptions, msgs *[]string) error {
	if conf.Application.ID == 0 {
		var projName string
		for {
			applicationId, err := cmd.createApplication(client, ctx, conf, msgs)
			if err != nil {
				// if the name is already in use, we ask for another one
				if strings.Contains(err.Error(), utils.ErrorNameInUse.Error()) {
					if NoPrompt {
						return err
					}
					logger.FInfoFlags(cmd.Io.Out, msg.AppInUse, cmd.F.Format, cmd.F.Out)
					*msgs = append(*msgs, msg.AppInUse)
					if Auto {
						projName = fmt.Sprintf("%s-%s", conf.Name, utils.Timestamp())
						msgf := fmt.Sprintf(msg.NameInUseApplication, projName)
						logger.FInfoFlags(cmd.Io.Out, msgf, cmd.F.Format, cmd.F.Out)
						*msgs = append(*msgs, msgf)
					} else {
						projName, err = askForInput(msg.AskInputName, thoth.GenerateName())
						if err != nil {
							return err
						}
					}
					conf.Name = projName
					continue
				}
				return err
			}
			conf.Application.ID = applicationId
			break
		}

		err := cmd.WriteAzionJsonContent(conf, ProjectConf)
		if err != nil {
			logger.Debug("Error while writing azion.json file", zap.Error(err))
			return err
		}
	} else {
		err := cmd.updateApplication(client, ctx, conf, msgs)
		if err != nil {
			logger.Debug("Error while updating Edge Application", zap.Error(err))
			return err
		}
	}
	return nil
}

func (cmd *DeployCmd) doWorkload(client *apiworkload.Client, ctx context.Context, conf *contracts.AzionApplicationOptions, msgs *[]string) error {
	var workload apiworkload.WorkloadResponse
	var err error

	newWorkload := false
	if conf.Workloads.Id == 0 {
		var projName string
		for {
			workload, err = cmd.createWorkload(client, ctx, conf, msgs)
			if err != nil {
				// if the name is already in use, we ask for another one
				if strings.Contains(err.Error(), utils.ErrorNameInUse.Error()) {
					if NoPrompt {
						return err
					}
					logger.FInfoFlags(cmd.Io.Out, msg.DomainInUse, cmd.F.Format, cmd.F.Out)
					*msgs = append(*msgs, msg.DomainInUse)
					if Auto {
						projName = fmt.Sprintf("%s-%s", conf.Name, utils.Timestamp())
						msgf := fmt.Sprintf(msg.NameInUseApplication, projName)
						logger.FInfoFlags(cmd.Io.Out, msgf, cmd.F.Format, cmd.F.Out)
						*msgs = append(*msgs, msgf)
						projName = thoth.GenerateName()
					} else {
						projName, err = askForInput(msg.AskInputName, thoth.GenerateName())
						if err != nil {
							return err
						}
					}
					conf.Domain.Name = projName
					continue
				}
				return err
			}
			conf.Workloads.Id = workload.GetId()
			conf.Workloads.Name = workload.GetName()
			conf.Workloads.Domains = workload.GetDomains()
			conf.Workloads.Url = utils.Concat("https://", workload.GetDomains()[0].GetDomain())
			newWorkload = true
			break
		}

		err = cmd.WriteAzionJsonContent(conf, ProjectConf)
		if err != nil {
			logger.Debug("Error while writing azion.json file", zap.Error(err))
			return err
		}

	} else {
		workload, err = cmd.updateWorkload(client, ctx, conf, msgs)
		if err != nil {
			logger.Debug("Error while updating workload", zap.Error(err))
			return err
		}
	}

	if conf.RtPurge.PurgeOnPublish && !newWorkload {
		err = PurgeForUpdatedFiles(cmd, workload, ProjectConf, msgs)
		if err != nil {
			logger.Debug("Error while purging workload", zap.Error(err))
			return err
		}
	}

	return nil
}

func (cmd *DeployCmd) doRulesDeploy(ctx context.Context, conf *contracts.AzionApplicationOptions, client *apiapp.Client, msgs *[]string) error {
	if conf.NotFirstRun {
		return nil
	}
	var cacheId int64
	var authorize bool
	if Auto || NoPrompt {
		authorize = false
	} else {
		authorize = utils.Confirm(cmd.F.GlobalFlagAll, msg.AskCreateCacheSettings, false)
	}

	if authorize {
		var reqCache apiapp.CreateCacheSettingsRequest
		reqCache.SetName(conf.Name)

		// create Cache Settings
		strApp := strconv.FormatInt(conf.Application.ID, 10)
		cache, err := client.CreateCacheSettingsNextApplication(ctx, &reqCache, strApp)
		if err != nil {
			logger.Debug("Error while creating Cache Settings", zap.Error(err))
			return err
		}
		logger.FInfoFlags(cmd.F.IOStreams.Out, msg.CacheSettingsSuccessful, cmd.F.Format, cmd.F.Out)
		*msgs = append(*msgs, msg.CacheSettingsSuccessful)
		cacheId = cache.GetId()
	}

	appId := fmt.Sprintf("%d", conf.Application.ID)

	// creates gzip and cache rules
	err := client.CreateRulesEngineNextApplication(ctx, appId, cacheId, conf.Preset, authorize)
	if err != nil {
		logger.Debug("Error while creating rules engine", zap.Error(err))
		return err
	}

	conf.NotFirstRun = true
	return nil
}

func (cmd *DeployCmd) doOriginSingle(clientOrigin *apiori.Client, ctx context.Context, conf *contracts.AzionApplicationOptions, msgs *[]string) (int64, error) {
	var DefaultOrigin = [1]string{"api.azion.net"}

	if conf.NotFirstRun {
		return 0, nil
	}

	reqSingleOrigin := apiori.CreateRequest{}
	addresses := prepareAddresses(DefaultOrigin[:])
	reqSingleOrigin.SetAddresses(addresses)

	reqSingleOrigin.SetName(utils.Concat(conf.Name, "_single"))
	reqSingleOrigin.SetHostHeader("${host}")

	origin, err := clientOrigin.Create(ctx, conf.Application.ID, &reqSingleOrigin)
	if err != nil {
		logger.Debug("Error while creating default origin ", zap.Any("Error", err))
		return 0, err
	}
	logger.FInfoFlags(cmd.F.IOStreams.Out, msg.OriginsSuccessful, cmd.F.Format, cmd.F.Out)
	*msgs = append(*msgs, msg.OriginsSuccessful)
	newOrigin := contracts.AzionJsonDataOrigin{
		OriginId:  origin.GetOriginId(),
		OriginKey: origin.GetOriginKey(),
		Name:      origin.GetName(),
	}
	conf.Origin = append(conf.Origin, newOrigin)

	return newOrigin.OriginId, nil
}

func (cmd *DeployCmd) createFunction(client *api.Client, ctx context.Context, conf *contracts.AzionApplicationOptions, funcToCreate contracts.AzionJsonDataFunction, msgs *[]string) (int64, error) {
	reqCre := api.CreateRequest{}

	code, err := cmd.FileReader(funcToCreate.File)
	if err != nil {
		logger.Debug("Error while reading Edge Function file <"+funcToCreate.File+">", zap.Error(err))
		return 0, fmt.Errorf("%s: %w", msg.ErrorCodeFlag, err)
	}

	prependText := fmt.Sprintf(injectIntoFunction, conf.Bucket, conf.Prefix)
	newCode := append([]byte(prependText), code...)

	reqCre.SetCode(string(newCode))

	reqCre.SetActive(true)
	if funcToCreate.Name == "__DEFAULT__" {
		reqCre.SetName(conf.Name)
	} else {
		reqCre.SetName(funcToCreate.Name)
	}

	//Read args
	marshalledArgs, err := cmd.FileReader(funcToCreate.Args)
	if err != nil {
		logger.Debug("Error while reding args.json file <"+funcToCreate.Args+">", zap.Error(err))
		return 0, fmt.Errorf("%s: %w", msg.ErrorArgsFlag, err)
	}
	args := make(map[string]interface{})
	if err := cmd.Unmarshal(marshalledArgs, &args); err != nil {
		logger.Debug("Error while unmarshling args.json file <"+funcToCreate.Args+">", zap.Error(err))
		return 0, fmt.Errorf("%s: %w", msg.ErrorParseArgs, err)
	}

	reqCre.SetJsonArgs(args)
	response, err := client.Create(ctx, &reqCre)
	if err != nil {
		logger.Debug("Error while creating Edge Function", zap.Error(err), zap.Any("Name", reqCre.Name))
		return 0, err
	}
	msgf := fmt.Sprintf(msg.DeployOutputEdgeFunctionCreate, response.GetName(), response.GetId())
	logger.FInfoFlags(cmd.F.IOStreams.Out, msgf, cmd.F.Format, cmd.F.Out)
	*msgs = append(*msgs, msgf)
	return response.GetId(), nil
}

func (cmd *DeployCmd) updateFunction(client *api.Client, ctx context.Context, conf *contracts.AzionApplicationOptions, funcToUpdate contracts.AzionJsonDataFunction, msgs *[]string) (int64, error) {
	reqUpd := api.UpdateRequest{}

	code, err := cmd.FileReader(funcToUpdate.File)
	if err != nil {
		logger.Debug("Error while reading Edge Function file <"+funcToUpdate.File+">", zap.Error(err))
		return 0, fmt.Errorf("%s: %w", msg.ErrorCodeFlag, err)
	}

	prependText := fmt.Sprintf(injectIntoFunction, conf.Bucket, conf.Prefix)
	newCode := append([]byte(prependText), code...)

	reqUpd.SetCode(string(newCode))

	reqUpd.SetActive(true)
	if funcToUpdate.Name == "__DEFAULT__" {
		reqUpd.SetName(conf.Name)
	} else {
		reqUpd.SetName(funcToUpdate.Name)
	}

	//Read args
	marshalledArgs, err := cmd.FileReader(funcToUpdate.Args)
	if err != nil {
		logger.Debug("Error while reading args.json file <"+funcToUpdate.Args+">", zap.Error(err))
		return 0, fmt.Errorf("%s: %w", msg.ErrorArgsFlag, err)
	}
	args := make(map[string]interface{})
	if err := cmd.Unmarshal(marshalledArgs, &args); err != nil {
		logger.Debug("Error while unmarshling args.json file <"+funcToUpdate.Args+">", zap.Error(err))
		return 0, fmt.Errorf("%s: %w", msg.ErrorParseArgs, err)
	}

	reqUpd.SetJsonArgs(args)
	funcId := strconv.FormatInt(funcToUpdate.ID, 10)
	response, err := client.Update(ctx, &reqUpd, funcId)
	if err != nil {
		logger.Debug("Error while updating Edge Function", zap.Error(err), zap.Any("Name", reqUpd.Name))
		return 0, fmt.Errorf(msg.ErrorUpdateFunction.Error(), err)
	}

	msgf := fmt.Sprintf(msg.DeployOutputEdgeFunctionUpdate, response.GetName(), funcToUpdate.ID)
	logger.FInfoFlags(cmd.F.IOStreams.Out, msgf, cmd.F.Format, cmd.F.Out)
	*msgs = append(*msgs, msgf)
	return response.GetId(), nil
}

func (cmd *DeployCmd) createApplication(client *apiapp.Client, ctx context.Context, conf *contracts.AzionApplicationOptions, msgs *[]string) (int64, error) {
	reqApp := apiapp.CreateRequest{}
	if conf.Application.Name == "__DEFAULT__" {
		reqApp.SetName(conf.Name)
	} else {
		reqApp.SetName(conf.Application.Name)
	}

	application, err := client.Create(ctx, &reqApp)
	if err != nil {
		return 0, fmt.Errorf(msg.ErrorCreateApplication.Error(), err)
	}

	msgf := fmt.Sprintf(
		msg.DeployOutputEdgeApplicationCreate, application.GetName(), application.GetId())
	logger.FInfoFlags(cmd.F.IOStreams.Out, msgf, cmd.F.Format, cmd.F.Out)
	*msgs = append(*msgs, msgf)

	reqUpApp := apiapp.UpdateRequest{}
	mods := sdk.EdgeApplicationModulesRequest{}
	mods.SetEdgeFunctionsEnabled(true)
	mods.SetApplicationAcceleratorEnabled(true)
	reqUpApp.SetModules(mods)
	reqUpApp.Id = application.GetId()

	application, err = client.Update(ctx, &reqUpApp)
	if err != nil {
		return 0, fmt.Errorf(msg.ErrorUpdateApplication.Error(), err)
	}

	return application.GetId(), nil
}

func (cmd *DeployCmd) updateApplication(client *apiapp.Client, ctx context.Context, conf *contracts.AzionApplicationOptions, msgs *[]string) error {
	reqApp := apiapp.UpdateRequest{}
	if conf.Application.Name == "__DEFAULT__" {
		reqApp.SetName(conf.Name)
	} else {
		reqApp.SetName(conf.Application.Name)
	}
	reqApp.Id = conf.Application.ID
	application, err := client.Update(ctx, &reqApp)
	if err != nil {
		return err
	}
	msgf := fmt.Sprintf(
		msg.DeployOutputEdgeApplicationUpdate, application.GetName(), application.GetId())
	logger.FInfoFlags(cmd.F.IOStreams.Out, msgf, cmd.F.Format, cmd.F.Out)
	*msgs = append(*msgs, msgf)
	return nil
}

func (cmd *DeployCmd) createWorkload(client *apiworkload.Client, ctx context.Context, conf *contracts.AzionApplicationOptions, msgs *[]string) (apiworkload.WorkloadResponse, error) {
	reqWork := apiworkload.CreateRequest{}
	if conf.Workloads.Name == "__DEFAULT__" || conf.Workloads.Name == "" {
		reqWork.SetName(conf.Name)
	} else {
		reqWork.SetName(conf.Workloads.Name)
	}
	reqWork.SetAlternateDomains([]string{})
	reqWork.SetActive(true)
	reqWork.SetEdgeApplication(conf.Application.ID)
	workload, err := client.Create(ctx, &reqWork)
	if err != nil {
		return nil, fmt.Errorf(msg.ErrorCreateDomain.Error(), err)
	}
	msgf := fmt.Sprintf(msg.DeployOutputWorkloadCreate, conf.Name, workload.GetId())
	logger.FInfoFlags(cmd.F.IOStreams.Out, msgf, cmd.F.Format, cmd.F.Out)
	*msgs = append(*msgs, msgf)
	return workload, nil
}

func (cmd *DeployCmd) updateWorkload(client *apiworkload.Client, ctx context.Context, conf *contracts.AzionApplicationOptions, msgs *[]string) (apiworkload.WorkloadResponse, error) {
	reqWork := apiworkload.UpdateRequest{}
	if conf.Workloads.Name == "__DEFAULT__" {
		reqWork.SetName(conf.Name)
	} else {
		reqWork.SetName(conf.Workloads.Name)
	}
	reqWork.Id = conf.Workloads.Id
	workload, err := client.Update(ctx, &reqWork)
	if err != nil {
		return nil, fmt.Errorf(msg.ErrorUpdateDomain.Error(), err)
	}
	msgf := fmt.Sprintf(msg.DeployOutputWorkloadUpdate, conf.Name, workload.GetId())
	logger.FInfoFlags(cmd.F.IOStreams.Out, msgf, cmd.F.Format, cmd.F.Out)
	*msgs = append(*msgs, msgf)
	return workload, nil
}

func prepareAddresses(addrs []string) (addresses []sdkv3.CreateOriginsRequestAddresses) {
	var addr sdkv3.CreateOriginsRequestAddresses
	for _, v := range addrs {
		addr.Address = v
		addresses = append(addresses, addr)
	}
	return
}

func (cmd *DeployCmd) createInstance(ctx context.Context, client *apiapp.Client, conf *contracts.AzionApplicationOptions, funcToCreate contracts.AzionJsonDataFunction) (sdk.EdgeApplicationFunctionInstance, error) {
	logger.Debug("Create Instance")

	// create instance function
	reqIns := apiapp.CreateInstanceRequest{}
	reqIns.SetEdgeFunction(funcToCreate.ID)

	if funcToCreate.InstanceName == "__DEFAULT__" || funcToCreate.InstanceName == "" {
		reqIns.SetName(conf.Name)
	} else {
		reqIns.SetName(funcToCreate.InstanceName)
	}
	reqIns.ApplicationId = conf.Application.ID

	//Read args
	marshalledArgs, err := cmd.FileReader(funcToCreate.Args)
	if err != nil {
		logger.Debug("Error while reding args.json file <"+funcToCreate.Args+">", zap.Error(err))
		return sdk.EdgeApplicationFunctionInstance{}, fmt.Errorf("%s: %w", msg.ErrorArgsFlag, err)
	}
	args := make(map[string]interface{})
	if err := cmd.Unmarshal(marshalledArgs, &args); err != nil {
		logger.Debug("Error while unmarshling args.json file <"+funcToCreate.Args+">", zap.Error(err))
		return sdk.EdgeApplicationFunctionInstance{}, fmt.Errorf("%s: %w", msg.ErrorParseArgs, err)
	}
	reqIns.SetJsonArgs(args)

	appId := fmt.Sprintf("%d", conf.Application.ID)
	resp, err := client.CreateFuncInstances(ctx, &reqIns, appId)
	if err != nil {
		return sdk.EdgeApplicationFunctionInstance{}, err
	}

	return resp, nil
}

func (cmd *DeployCmd) updateInstance(ctx context.Context, client *apiapp.Client, conf *contracts.AzionApplicationOptions, funcToUpdate contracts.AzionJsonDataFunction) (sdk.EdgeApplicationFunctionInstance, error) {
	logger.Debug("Update Instance")

	// create instance function
	reqIns := apiapp.UpdateInstanceRequest{}
	reqIns.SetEdgeFunction(funcToUpdate.ID)

	if funcToUpdate.InstanceName == "__DEFAULT__" {
		reqIns.SetName(conf.Name)
	} else {
		reqIns.SetName(funcToUpdate.Name)
	}

	//Read args
	marshalledArgs, err := cmd.FileReader(funcToUpdate.Args)
	if err != nil {
		logger.Debug("Error while reding args.json file <"+funcToUpdate.Args+">", zap.Error(err))
		return sdk.EdgeApplicationFunctionInstance{}, fmt.Errorf("%s: %w", msg.ErrorArgsFlag, err)
	}
	args := make(map[string]interface{})
	if err := cmd.Unmarshal(marshalledArgs, &args); err != nil {
		logger.Debug("Error while unmarshling args.json file <"+funcToUpdate.Args+">", zap.Error(err))
		return sdk.EdgeApplicationFunctionInstance{}, fmt.Errorf("%s: %w", msg.ErrorParseArgs, err)
	}
	reqIns.SetJsonArgs(args)

	instID := strconv.FormatInt(funcToUpdate.InstanceID, 10)
	appID := strconv.FormatInt(conf.Application.ID, 10)
	resp, err := client.UpdateInstance(ctx, &reqIns, appID, instID)
	if err != nil {
		return sdk.EdgeApplicationFunctionInstance{}, err
	}

	return resp, nil
}

func checkArgsJson(cmd *DeployCmd, projectPath string) error {
	workingDir, err := cmd.GetWorkDir()
	if err != nil {
		return err
	}

	workingDirPath := path.Join(workingDir, projectPath, "args.json")

	_, err = cmd.FileReader(workingDirPath)
	if err != nil {
		if err := cmd.WriteFile(workingDirPath, []byte("{}"), 0644); err != nil {
			logger.Debug("Error while trying to create args.json file", zap.Error(err))
			return fmt.Errorf(utils.ErrorCreateFile.Error(), workingDirPath)
		}
	}

	return nil
}
