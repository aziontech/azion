package manifest

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"
	"time"

	msg "github.com/aziontech/azion-cli/messages/manifest"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/utils"
	"github.com/briandowns/spinner"

	msgcache "github.com/aziontech/azion-cli/messages/cache_setting"
	msgrule "github.com/aziontech/azion-cli/messages/delete/rules_engine"
	apiApplications "github.com/aziontech/azion-cli/pkg/api/applications"
	apiCache "github.com/aziontech/azion-cli/pkg/api/cache_setting"
	apiConnector "github.com/aziontech/azion-cli/pkg/api/connector"
	functionsApi "github.com/aziontech/azion-cli/pkg/api/function"
	apiPurge "github.com/aziontech/azion-cli/pkg/api/realtime_purge"
	apiWorkloads "github.com/aziontech/azion-cli/pkg/api/workloads"
	"go.uber.org/zap"
)

var (
	CacheIds         map[string]int64
	CacheIdsBackup   map[string]int64
	RuleIds          map[string]contracts.RuleIdsStruct
	ConnectorIds     map[string]int64
	DeploymentIds    map[string]int64
	FunctionIds      map[string]contracts.AzionJsonDataFunction
	manifestFilePath = "/.edge/manifest.json"
)

type ManifestInterpreter struct {
	FileReader            func(path string) ([]byte, error)
	GetWorkDir            func() (string, error)
	WriteAzionJsonContent func(conf *contracts.AzionApplicationOptions, confPath string) error
}

func NewManifestInterpreter() *ManifestInterpreter {
	return &ManifestInterpreter{
		FileReader:            os.ReadFile,
		GetWorkDir:            utils.GetWorkingDir,
		WriteAzionJsonContent: utils.WriteAzionJsonContent,
	}
}

func (man *ManifestInterpreter) ManifestPath() (string, error) {
	pathWorkingDir, err := man.GetWorkDir()
	if err != nil {
		return "", err
	}
	return utils.Concat(pathWorkingDir, manifestFilePath), nil
}

func (man *ManifestInterpreter) ReadManifest(path string, f *cmdutil.Factory, msgs *[]string) (*contracts.ManifestV4, error) {
	logger.FInfoFlags(f.IOStreams.Out, msg.ReadingManifest, f.Format, f.Out)
	*msgs = append(*msgs, msg.ReadingManifest)
	manifest := &contracts.ManifestV4{}

	byteManifest, err := man.FileReader(path)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(byteManifest, &manifest)
	if err != nil {
		return nil, err
	}

	return manifest, nil
}

func (man *ManifestInterpreter) CreateResources(conf *contracts.AzionApplicationOptions, manifest *contracts.ManifestV4, functions map[string]contracts.AzionJsonDataFunction, f *cmdutil.Factory, projectConf string, msgs *[]string) error {
	s := spinner.New(spinner.CharSets[7], 100*time.Millisecond)
	s.Suffix = " " + msg.CreatingManifest
	s.FinalMSG = "\n"
	if !f.Debug {
		s.Start() // Start the spinner
	}
	defer s.Stop()

	client := apiApplications.NewClient(f.HttpClient, f.Config.GetString("api_v4_url"), f.Config.GetString("token"))
	clientCache := apiCache.NewClientV4(f.HttpClient, f.Config.GetString("api_v4_url"), f.Config.GetString("token"))
	clientWorkload := apiWorkloads.NewClient(f.HttpClient, f.Config.GetString("api_v4_url"), f.Config.GetString("token"))
	connectorClient := apiConnector.NewClient(f.HttpClient, f.Config.GetString("api_v4_url"), f.Config.GetString("token"))
	functionClient := functionsApi.NewClient(f.HttpClient, f.Config.GetString("api_v4_url"), f.Config.GetString("token"))
	ctx := context.Background()

	CacheIds = make(map[string]int64)
	CacheIdsBackup = make(map[string]int64)
	RuleIds = make(map[string]contracts.RuleIdsStruct)
	ConnectorIds = make(map[string]int64)
	DeploymentIds = make(map[string]int64)
	FunctionIds = make(map[string]contracts.AzionJsonDataFunction)

	for _, cacheConf := range conf.CacheSettings {
		CacheIds[cacheConf.Name] = cacheConf.Id
	}

	for _, funcCong := range conf.Function {
		FunctionIds[funcCong.Name] = funcCong
	}

	for _, deploymentConf := range conf.Workloads.Deployments {
		DeploymentIds[deploymentConf.Name] = deploymentConf.Id
	}

	for _, ruleConf := range conf.RulesEngine.Rules {
		RuleIds[ruleConf.Name] = contracts.RuleIdsStruct{
			Id:    ruleConf.Id,
			Phase: ruleConf.Phase,
		}
	}

	for _, connectorConf := range conf.Connectors {
		ConnectorIds[connectorConf.Name] = connectorConf.Id
	}

	for _, funcMan := range manifest.Functions {
		code, err := os.ReadFile(path.Join(".edge", funcMan.Path))
		if err != nil {
			return fmt.Errorf(msg.ErrorReadCodeFile.Error(), err)
		}

		if funcConf := FunctionIds[funcMan.Name]; funcConf.ID > 0 {
			request := functionsApi.UpdateRequest{}
			request.SetActive(true)
			request.SetDefaultArgs(funcMan.DefaultArgs)
			request.SetName(funcMan.Name)
			request.SetCode(string(code))
			_, err := functionClient.Update(ctx, &request, funcConf.ID)
			if err != nil {
				return err
			}
		} else {
			request := functionsApi.CreateRequest{}
			request.SetActive(true)
			request.SetDefaultArgs(funcMan.DefaultArgs)
			request.SetName(funcMan.Name)
			request.SetCode(string(code))
			resp, err := functionClient.Create(ctx, &request)
			if err != nil {
				return err
			}
			newFunc := contracts.AzionJsonDataFunction{
				ID:   resp.GetId(),
				Name: resp.GetName(),
				File: funcMan.Argument,
				Args: "./azion/args.json",
			}
			FunctionIds[funcMan.Name] = newFunc
			conf.Function = append(conf.Function, newFunc)
		}
	}

	err := man.WriteAzionJsonContent(conf, projectConf)
	if err != nil {
		logger.Debug("Error while writing azion.json file", zap.Error(err))
		return err
	}

	for _, funcMan := range manifest.Applications[0].FunctionsInstances {
		if funcConf := FunctionIds[funcMan.Function]; funcConf.InstanceID > 0 {
			request := apiApplications.UpdateInstanceRequest{}
			request.SetActive(funcMan.Active)
			request.SetFunction(funcConf.ID)
			if len(funcMan.Args) > 0 {
				request.SetArgs(funcMan.Args)
			} else {
				args, err := unmarshalJsonArgs(funcConf.Args)
				if err != nil {
					return err
				}
				request.SetArgs(args)
			}
			request.SetName(funcMan.Name)
			_, err := client.UpdateInstance(ctx, &request, conf.Application.ID, funcConf.InstanceID)
			if err != nil {
				return err
			}
		} else {
			request := apiApplications.CreateInstanceRequest{}
			request.SetActive(true)
			if len(funcMan.Args) > 0 {
				request.SetArgs(funcMan.Args)
			} else {
				args, err := unmarshalJsonArgs(funcConf.Args)
				if err != nil {
					return err
				}
				request.SetArgs(args)
			}
			request.SetName(funcMan.Name)
			request.SetFunction(funcConf.ID)
			resp, err := client.CreateFuncInstances(ctx, &request, conf.Application.ID)
			if err != nil {
				return err
			}
			newFunc := contracts.AzionJsonDataFunction{
				ID:         funcConf.ID,
				CacheId:    funcConf.CacheId,
				Name:       funcConf.Name,
				File:       funcConf.File,
				Args:       funcConf.Args,
				InstanceID: resp.GetId(),
			}
			FunctionIds[funcConf.Name] = newFunc
		}
	}

	funcsToWrite := []contracts.AzionJsonDataFunction{}
	for _, funcs := range FunctionIds {
		funcsToWrite = append(funcsToWrite, funcs)
	}

	conf.Function = funcsToWrite

	err = man.WriteAzionJsonContent(conf, projectConf)
	if err != nil {
		logger.Debug("Error while writing azion.json file", zap.Error(err))
		return err
	}

	if len(manifest.Applications) > 0 {
		edgeappman := manifest.Applications[0]
		if conf.Application.ID > 0 {
			req := transformEdgeApplicationRequestUpdate(edgeappman)
			req.Id = conf.Application.ID
			_, err := client.Update(ctx, req)
			if err != nil {
				return err
			}
		} else {
			createreq := transformEdgeApplicationRequestCreate(edgeappman)
			resp, err := client.Create(ctx, createreq)
			if err != nil {
				return err
			}
			conf.Application.ID = resp.GetId()
		}

		CacheConf := []contracts.AzionJsonDataCacheSettings{}
		if len(edgeappman.CacheSettings) > 0 {
			for _, cache := range edgeappman.CacheSettings {
				if r := CacheIds[cache.Name]; r > 0 {
					updated, err := updateCache(f, cache, clientCache, conf, r, ctx, edgeappman)
					if err != nil {
						return err
					}
					CacheConf = append(CacheConf, updated)
				} else {
					newCache, err := createCache(cache, clientCache, conf, ctx, edgeappman)
					if err != nil {
						return err
					}
					CacheConf = append(CacheConf, newCache)
					CacheIds[newCache.Name] = newCache.Id
				}
			}
		}

		//backup cache ids
		for k, v := range CacheIds {
			CacheIdsBackup[k] = v
		}

		conf.CacheSettings = CacheConf
		err = man.WriteAzionJsonContent(conf, projectConf)
		if err != nil {
			logger.Debug("Error while writing azion.json file", zap.Error(err))
			return err
		}

		connectorConf := []contracts.AzionJsonDataConnectors{}
		if len(manifest.Connectors) > 0 {
			connector := manifest.Connectors[0]
			connName, connType := getConnectorName(connector, conf.Name)
			if id := ConnectorIds[connName]; id > 0 {
				request := transformEdgeConnectorRequest(connector)
				connectorResp, err := connectorClient.Update(ctx, request, id)
				if err != nil {
					return err
				}
				conn := contracts.AzionJsonDataConnectors{}
				switch connType {
				case "http":
					http := connectorResp.ConnectorHTTP
					conn.Id = http.GetId()
					conn.Name = http.GetName()
					conn.Address = http.Attributes.Addresses
				case "ingest":
					liveIngest := connectorResp.ConnectorLiveIngest
					conn.Id = liveIngest.GetId()
					conn.Name = liveIngest.GetName()
				case "storage":
					storage := connectorResp.ConnectorStorage
					conn.Id = storage.GetId()
					conn.Name = storage.GetName()
				default:
					return errors.New("failed to get Connector type")
				}
				connectorConf = append(connectorConf, conn)
			} else {
				request := apiConnector.CreateRequest{
					ConnectorPolymorphicRequest: connector,
				}
				connectorResp, err := connectorClient.Create(ctx, &request)
				if err != nil {
					return err
				}
				conn := contracts.AzionJsonDataConnectors{}
				switch connType {
				case "http":
					http := connectorResp.ConnectorHTTP
					conn.Id = http.GetId()
					conn.Name = http.GetName()
				case "ingest":
					liveIngest := connectorResp.ConnectorLiveIngest
					conn.Id = liveIngest.GetId()
					conn.Name = liveIngest.GetName()
				case "storage":
					storage := connectorResp.ConnectorStorage
					conn.Id = storage.GetId()
					conn.Name = storage.GetName()
				default:
					return errors.New("failed to get Connector type")
				}
				ConnectorIds[conn.Name] = conn.Id
				connectorConf = append(connectorConf, conn)
			}
		}

		conf.Connectors = connectorConf
		err = man.WriteAzionJsonContent(conf, projectConf)
		if err != nil {
			logger.Debug("Error while writing azion.json file", zap.Error(err))
			return err
		}

		ruleConf := []contracts.AzionJsonDataRules{}
		if len(edgeappman.Rules) > 0 {
			for _, rule := range edgeappman.Rules {
				if r := RuleIds[rule.Rule.Name]; r.Id > 0 {
					switch rule.Phase {
					case "request":
						req := transformRuleRequest(rule.Rule)
						req.IdApplication = conf.Application.ID
						req.Id = r.Id
						behs, err := transformBehaviorsRequest(rule.Rule.Behaviors)
						if errors.Is(err, utils.ErrorNotFound404) {
							logger.Debug("Rule not found. Skipping update", zap.Any("Error", err))
							edgeappman.Rules = append(edgeappman.Rules, rule)
							delete(RuleIds, rule.Rule.Name)
							continue
						}
						if err != nil {
							return err
						}
						req.Behaviors = behs
						updated, err := client.UpdateRulesEngineRequest(ctx, req)
						if err != nil {
							return err
						}
						newRule := contracts.AzionJsonDataRules{
							Id:    updated.GetId(),
							Name:  updated.GetName(),
							Phase: rule.Phase,
						}
						ruleConf = append(ruleConf, newRule)
						delete(RuleIds, updated.GetName())
					case "response":
						req := transformRuleResponse(rule.Rule)
						req.IdApplication = conf.Application.ID
						req.Id = r.Id
						behs, err := transformBehaviorsResponse(rule.Rule.Behaviors)
						if errors.Is(err, utils.ErrorNotFound404) {
							logger.Debug("Rule not found. Skipping update", zap.Any("Error", err))
							edgeappman.Rules = append(edgeappman.Rules, rule)
							delete(RuleIds, rule.Rule.Name)
							continue
						}
						if err != nil {
							return err
						}
						req.Behaviors = behs
						updated, err := client.UpdateRulesEngineResponse(ctx, req)
						if err != nil {
							return err
						}
						newRule := contracts.AzionJsonDataRules{
							Id:    updated.GetId(),
							Name:  updated.GetName(),
							Phase: rule.Phase,
						}
						ruleConf = append(ruleConf, newRule)
						delete(RuleIds, updated.GetName())
					default:
						return msg.ErrorInvalidPhase
					}

				} else {
					switch rule.Phase {
					case "request":
						req := &apiApplications.CreateRulesEngineRequest{}
						createRequest := transformRuleRequestCreate(rule.Rule)
						bh, err := transformBehaviorsRequest(rule.Rule.Behaviors)
						if err != nil {
							return err
						}
						req.ApplicationRequestPhaseRuleEngineRequest = createRequest
						req.Behaviors = bh
						created, err := client.CreateRulesEngineRequest(ctx, conf.Application.ID, rule.Phase, req)
						if err != nil {
							return err
						}
						newRule := contracts.AzionJsonDataRules{
							Id:    created.GetId(),
							Name:  created.GetName(),
							Phase: rule.Phase,
						}
						ruleConf = append(ruleConf, newRule)
					case "response":
						req := &apiApplications.CreateRulesEngineResponse{}
						createRequest := transformRuleResponseCreate(rule.Rule)
						bh, err := transformBehaviorsResponse(rule.Rule.Behaviors)
						if err != nil {
							return err
						}
						req.ApplicationResponsePhaseRuleEngineRequest = createRequest
						req.Behaviors = bh
						created, err := client.CreateRulesEngineResponse(ctx, conf.Application.ID, rule.Phase, req)
						if err != nil {
							return err
						}
						newRule := contracts.AzionJsonDataRules{
							Id:    created.GetId(),
							Name:  created.GetName(),
							Phase: rule.Phase,
						}
						ruleConf = append(ruleConf, newRule)
					default:
						return msg.ErrorInvalidPhase
					}
				}
			}

			conf.RulesEngine.Rules = ruleConf
			err = man.WriteAzionJsonContent(conf, projectConf)
			if err != nil {
				logger.Debug("Error while writing azion.json file", zap.Error(err))
				return err
			}
		}
	}

	err = man.WriteAzionJsonContent(conf, projectConf)
	if err != nil {
		logger.Debug("Error while writing azion.json file", zap.Error(err))
		return err
	}

	if len(manifest.Workloads) > 0 {
		workloadMan := manifest.Workloads[0]
		if conf.Workloads.Id > 0 {
			request := transformWorkloadRequestUpdate(workloadMan)
			request.Id = conf.Workloads.Id
			updated, err := clientWorkload.Update(ctx, request)
			if err != nil {
				return err
			}
			conf.Workloads.Domains = updated.GetDomains()
			conf.Workloads.Url = utils.Concat("https://", updated.GetWorkloadDomain())
		} else {
			request := transformWorkloadRequestCreate(workloadMan, conf.Application.ID)
			resp, err := clientWorkload.Create(ctx, request)
			if err != nil {
				return err
			}
			conf.Workloads.Id = resp.GetId()
			conf.Workloads.Name = resp.GetName()
			conf.Workloads.Domains = resp.GetDomains()
			conf.Workloads.Url = utils.Concat("https://", resp.GetWorkloadDomain())
		}
	}

	err = man.WriteAzionJsonContent(conf, projectConf)
	if err != nil {
		logger.Debug("Error while writing azion.json file", zap.Error(err))
		return err
	}

	if len(manifest.WorkloadDeployments) > 0 {
		for _, deployment := range manifest.WorkloadDeployments {
			if id := DeploymentIds[deployment.Name]; id > 0 {
				request := transformWorkloadDeploymentRequestUpdate(deployment, conf)
				_, err := clientWorkload.UpdateDeployment(ctx, request, conf.Workloads.Id, id)
				if err != nil {
					return err
				}
			} else {
				request := transformWorkloadDeploymentRequestCreate(deployment, conf)
				resp, err := clientWorkload.CreateDeployment(ctx, request, conf.Workloads.Id)
				if err != nil {
					return err
				}
				conf.Workloads.Deployments = append(conf.Workloads.Deployments, contracts.Deployments{
					Id:   resp.GetId(),
					Name: resp.GetName(),
				})
			}
		}
	}

	err = man.WriteAzionJsonContent(conf, projectConf)
	if err != nil {
		logger.Debug("Error while writing azion.json file", zap.Error(err))
		return err
	}

	clipurge := apiPurge.NewClient(f.HttpClient, f.Config.GetString("api_v4_url"), f.Config.GetString("token"))
	for _, purgeObj := range manifest.Purge {
		err = clipurge.PurgeCache(ctx, purgeObj.Items, purgeObj.Type, *purgeObj.Layer)
		if err != nil {
			logger.Debug("Error while purging domains", zap.Error(err))
			return err
		}
	}

	err = deleteResources(ctx, f, conf, msgs)
	if err != nil {
		return err
	}

	return nil
}

// this is called to delete resources no longer present in manifest.json
func deleteResources(ctx context.Context, f *cmdutil.Factory, conf *contracts.AzionApplicationOptions, msgs *[]string) error {
	client := apiApplications.NewClient(f.HttpClient, f.Config.GetString("api_v4_url"), f.Config.GetString("token"))
	clientCache := apiCache.NewClientV4(f.HttpClient, f.Config.GetString("api_v4_url"), f.Config.GetString("token"))
	// clientOrigin := apiOrigin.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))

	if conf.SkipDeletion != nil && *conf.SkipDeletion {
		logger.FInfoFlags(f.IOStreams.Out, msg.SkipDeletion, f.Format, f.Out)
		*msgs = append(*msgs, msg.SkipDeletion)
		return nil
	}

	for _, value := range RuleIds {
		//since until [UXE-3599] was carried out we'd only cared about "request" phase, this check guarantees that if Phase is empty
		// we are probably dealing with a rule engine from a previous version
		phase := "request"
		if value.Phase != "" {
			phase = value.Phase
		}
		var statusInt int
		var err error
		switch phase {
		case "request":
			statusInt, err = client.DeleteRulesEngineRequest(ctx, conf.Application.ID, phase, value.Id)
		case "response":
			statusInt, err = client.DeleteRulesEngineResponse(ctx, conf.Application.ID, phase, value.Id)
		default:
			return msgrule.ErrorInvalidPhase
		}

		if statusInt == 404 {
			logger.Debug("Rule Engine not found. Skipping delete")
			continue
		}
		if err != nil {
			return err
		}
		msgf := fmt.Sprintf(msgrule.DeleteOutputSuccess+"\n", value.Id)
		logger.FInfoFlags(f.IOStreams.Out, msgf, f.Format, f.Out)
		*msgs = append(*msgs, msgf)
	}

	for _, value := range CacheIds {
		status, err := clientCache.Delete(ctx, conf.Application.ID, value)
		if status == 404 {
			logger.Debug("Cache Setting not found. Skipping delete")
			continue
		}
		if err != nil {
			return err
		}
		msgf := fmt.Sprintf(msgcache.DeleteOutputSuccess+"\n", value)
		logger.FInfoFlags(f.IOStreams.Out, msgf, f.Format, f.Out)
		*msgs = append(*msgs, msgf)
	}

	return nil
}

func unmarshalJsonArgs(argsPath string) (map[string]interface{}, error) {
	marshalledArgs, err := os.ReadFile(argsPath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", msg.ErrorReadArgsFile, err)
	}
	args := make(map[string]interface{})
	if err := json.Unmarshal(marshalledArgs, &args); err != nil {
		return nil, fmt.Errorf("%s: %w", msg.ErrorUnmarshalArgsFile, err)
	}
	return args, nil
}
