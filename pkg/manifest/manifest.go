package manifest

import (
	"encoding/json"
	"os"

	msg "github.com/aziontech/azion-cli/messages/manifest"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/utils"
	"github.com/davecgh/go-spew/spew"
)

// import (
// 	"context"
// 	"encoding/json"
// 	"fmt"
// 	"os"
// 	"strconv"

// 	msgcache "github.com/aziontech/azion-cli/messages/cache_setting"
// 	msgrule "github.com/aziontech/azion-cli/messages/delete/rules_engine"
//
// 	apiCache "github.com/aziontech/azion-cli/pkg/api/cache_setting"
// 	apiConnector "github.com/aziontech/azion-cli/pkg/api/edge_connector"
// 	apipurge "github.com/aziontech/azion-cli/pkg/api/realtime_purge"
// 	sdkedge "github.com/aziontech/azionapi-v4-go-sdk-dev/edge-api"
// 	edgesdk "github.com/aziontech/azionapi-v4-go-sdk/edge"
// 	"github.com/davecgh/go-spew/spew"

// 	apiEdgeApplications "github.com/aziontech/azion-cli/pkg/api/edge_applications"
// 	functionsApi "github.com/aziontech/azion-cli/pkg/api/edge_function"
// 	apiWorkloads "github.com/aziontech/azion-cli/pkg/api/workloads"
// 	"github.com/aziontech/azion-cli/pkg/cmdutil"
// 	"github.com/aziontech/azion-cli/pkg/contracts"
// 	"github.com/aziontech/azion-cli/pkg/logger"
// 	"github.com/aziontech/azion-cli/utils"
// 	"go.uber.org/zap"
// )

var (
	// CacheIds       map[string]int64
	// CacheIdsBackup map[string]int64
	// RuleIds        map[string]contracts.RuleIdsStruct
	// // OriginKeys       map[string]string
	// // OriginIds        map[string]int64
	// ConnectorIds     map[string]int64
	// FunctionIds      map[string]contracts.AzionJsonDataFunction
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
	spew.Dump(manifest.EdgeFunctions)

	return manifest, nil
}

// func (man *ManifestInterpreter) CreateResources(conf *contracts.AzionApplicationOptions, manifest *contracts.ManifestV4, functions map[string]contracts.AzionJsonDataFunction, f *cmdutil.Factory, projectConf string, msgs *[]string) error {

// 	logger.FInfoFlags(f.IOStreams.Out, msg.CreatingManifest, f.Format, f.Out)
// 	*msgs = append(*msgs, msg.CreatingManifest)

// 	client := apiEdgeApplications.NewClient(f.HttpClient, f.Config.GetString("api_v4_url"), f.Config.GetString("token"))
// 	clientCache := apiCache.NewClientV4(f.HttpClient, f.Config.GetString("api_v4_url"), f.Config.GetString("token"))
// 	clientWorkload := apiWorkloads.NewClient(f.HttpClient, f.Config.GetString("api_v4_url"), f.Config.GetString("token"))
// 	connectorClient := apiConnector.NewClient(f.HttpClient, f.Config.GetString("api_v4_url"), f.Config.GetString("token"))
// 	functionClient := functionsApi.NewClient(f.HttpClient, f.Config.GetString("api_v4_url"), f.Config.GetString("token"))
// 	ctx := context.Background()

// 	CacheIds = make(map[string]int64)
// 	CacheIdsBackup = make(map[string]int64)
// 	RuleIds = make(map[string]contracts.RuleIdsStruct)
// 	// OriginKeys = make(map[string]string)
// 	// OriginIds = make(map[string]int64)
// 	ConnectorIds = make(map[string]int64)
// 	FunctionIds = functions

// 	for _, cacheConf := range conf.CacheSettings {
// 		CacheIds[cacheConf.Name] = cacheConf.Id
// 	}

// 	for _, funcCong := range conf.Function {
// 		FunctionIds[funcCong.Name] = funcCong
// 	}

// 	for _, ruleConf := range conf.RulesEngine.Rules {
// 		RuleIds[ruleConf.Name] = contracts.RuleIdsStruct{
// 			Id:    ruleConf.Id,
// 			Phase: ruleConf.Phase,
// 		}
// 	}

// 	// for _, originConf := range conf.Origin {
// 	// 	OriginKeys[originConf.Name] = originConf.OriginKey
// 	// 	OriginIds[originConf.Name] = originConf.OriginId
// 	// }

// 	for _, connectorConf := range conf.Connectors {
// 		ConnectorIds[connectorConf.Name] = connectorConf.Id
// 	}

// 	for _, funcMan := range manifest.EdgeFunctions {
// 		code, err := os.ReadFile(funcMan.Argument)
// 		if err != nil {
// 			return fmt.Errorf("Failed to read target code file: %w", err)
// 		}

// 		if funcConf := FunctionIds[funcMan.Name]; funcConf.ID > 0 {
// 			request := functionsApi.UpdateRequest{}
// 			request.SetActive(true)
// 			request.SetDefaultArgs(sdkedge.EdgeFunctionsDefaultArgs{Arg: funcMan.Args})
// 			request.SetName(funcMan.Name)
// 			request.SetCode(string(code))
// 			idString := strconv.FormatInt(funcConf.ID, 10)
// 			_, err := functionClient.Update(ctx, &request, idString)
// 			if err != nil {
// 				return err
// 			}
// 		} else {
// 			request := functionsApi.CreateRequest{}
// 			request.SetActive(true)
// 			request.SetDefaultArgs(sdkedge.EdgeFunctionsDefaultArgs{Arg: funcMan.Args})
// 			request.SetName(funcMan.Name)
// 			request.SetCode(string(code))
// 			resp, err := functionClient.Create(ctx, &request)
// 			if err != nil {
// 				return err
// 			}
// 			newFunc := contracts.AzionJsonDataFunction{
// 				ID:   resp.GetId(),
// 				Name: resp.GetName(),
// 				File: funcMan.Argument,
// 				Args: "./azion/args.json",
// 			}
// 			FunctionIds[funcMan.Name] = newFunc
// 			conf.Function = append(conf.Function, newFunc)
// 		}
// 	}

// 	err := man.WriteAzionJsonContent(conf, projectConf)
// 	if err != nil {
// 		logger.Debug("Error while writing azion.json file", zap.Error(err))
// 		return err
// 	}

// 	if len(manifest.EdgeApplications) > 0 {
// 		edgeappman := manifest.EdgeApplications[0]
// 		if conf.Application.ID > 0 {
// 			req := transformEdgeApplicationRequestUpdate(edgeappman)
// 			req.Id = conf.Application.ID
// 			_, err := client.Update(ctx, req)
// 			if err != nil {
// 				return err
// 			}
// 		} else {
// 			createreq := transformEdgeApplicationRequestCreate(edgeappman)
// 			resp, err := client.Create(ctx, createreq)
// 			if err != nil {
// 				return err
// 			}
// 			conf.Application.ID = resp.GetId()
// 		}

// 		cacheConf := []contracts.AzionJsonDataCacheSettings{}
// 		if len(edgeappman.Cache) > 0 {
// 			for _, cache := range edgeappman.Cache {
// 				if r := CacheIds[cache.Name]; r > 0 {
// 					request := transformCacheRequest(cache)
// 					updated, err := clientCache.Update(ctx, request, conf.Application.ID, r)
// 					if err != nil {
// 						return err
// 					}
// 					newCache := contracts.AzionJsonDataCacheSettings{
// 						Id:   updated.GetData().Id,
// 						Name: updated.GetData().Name,
// 					}
// 					cacheConf = append(cacheConf, newCache)
// 				} else {
// 					request := apiCache.Request{
// 						CacheSettingRequest: cache,
// 					}
// 					responseCache, err := clientCache.Create(ctx, request.CacheSettingRequest, conf.Application.ID)
// 					if err != nil {
// 						return err
// 					}
// 					newCache := contracts.AzionJsonDataCacheSettings{
// 						Id:   responseCache.GetId(),
// 						Name: responseCache.GetName(),
// 					}
// 					cacheConf = append(cacheConf, newCache)
// 					CacheIds[newCache.Name] = newCache.Id

// 				}
// 			}
// 		}

// 		//backup cache ids
// 		for k, v := range CacheIds {
// 			CacheIdsBackup[k] = v
// 		}

// 		conf.CacheSettings = cacheConf
// 		err = man.WriteAzionJsonContent(conf, projectConf)
// 		if err != nil {
// 			logger.Debug("Error while writing azion.json file", zap.Error(err))
// 			return err
// 		}

// 		ruleConf := []contracts.AzionJsonDataRules{}
// 		if len(edgeappman.Rules) > 0 {
// 			for _, rule := range edgeappman.Rules {
// 				if r := RuleIds[rule.GetName()]; r.Id > 0 {
// 					req := transformRuleRequest(rule)
// 					strid := strconv.FormatInt(r.Id, 10)
// 					req.Id = strid
// 					behaviorsRequest := []sdkedge.EdgeApplicationRuleEngineRequestPhaseBehaviorsRequest{}
// 					for _, behavior := range rule.Behaviors {
// 						switch behavior.Name {
// 						case "run_function":
// 							if _, ok := FunctionIds[*behavior.GetArgument().String]; !ok {
// 								return msg.ErrorFuncNotFound
// 							}
// 							funcToWorkWith := FunctionIds[*behavior.GetArgument().String]
// 							var behString sdkedge.EdgeApplicationRuleEngineRequestPhaseBehaviorsRequest
// 							var behSet sdkedge.EdgeApplicationRuleEngineRequestPhaseBehaviorsEdgeApplicationRuleEngineStringRequest
// 							var att sdkedge.EdgeApplicationRuleEngineStringAttributesRequest

// 							// var beh sdkedge.EdgeApplicationBehaviorFieldRequest
// 							cacheId, err := doCacheForRule(ctx, client, conf, funcToWorkWith)
// 							if err != nil {
// 								return err
// 							}

// 							behSet.SetType("set_cache_policy")
// 							att.SetValue(strconv.FormatInt(cacheId, 10))
// 							behSet.SetAttributes(att)
// 							behString.EdgeApplicationRuleEngineRequestPhaseBehaviorsEdgeApplicationRuleEngineStringRequest = &behSet
// 							behaviorsRequest = append(behaviorsRequest, behString)

// 							var behFunc sdkedge.EdgeApplicationRuleEngineRequestPhaseBehaviorsEdgeApplicationRuleEngineRunFunctionRequest

// 							behFunc.SetType("run_function")
// 							var attFunc sdkedge.EdgeApplicationRuleEngineRunFunctionAttributesRequest
// 							attFunc.SetValue(funcToWorkWith.InstanceID)
// 							behFunc.SetAttributes(attFunc)
// 							behString = sdkedge.EdgeApplicationRuleEngineRequestPhaseBehaviorsRequest{}
// 							behString.EdgeApplicationRuleEngineRequestPhaseBehaviorsEdgeApplicationRuleEngineRunFunctionRequest = &behFunc
// 							behaviorsRequest = append(behaviorsRequest, behString)
// 						case "set_cache_policy":
// 							if id := CacheIdsBackup[*behavior.GetArgument().String]; id > 0 {
// 								var beh edgesdk.EdgeApplicationBehaviorFieldRequest
// 								beh.SetName("set_cache_policy")
// 								beh.SetArgument(edgesdk.EdgeApplicationBehaviorPolymorphicArgumentRequest{
// 									Int64: &id,
// 								})
// 								delete(CacheIds, *behavior.GetArgument().String)
// 							} else {
// 								logger.Debug("Cache Setting not found", zap.Any("Target", *behavior.GetArgument().String))
// 								return msg.ErrorCacheNotFound
// 							}
// 						case "set_edge_connector":
// 							if id := ConnectorIds[*behavior.GetArgument().String]; id > 0 {
// 								var beh edgesdk.EdgeApplicationBehaviorFieldRequest
// 								beh.SetName("set_edge_connector")
// 								beh.SetArgument(edgesdk.EdgeApplicationBehaviorPolymorphicArgumentRequest{
// 									Int64: &id,
// 								})
// 								delete(ConnectorIds, *behavior.GetArgument().String)
// 							} else {
// 								logger.Debug("Edge Connector not found", zap.Any("Target", id))
// 								return msg.ErrorConnectorNotFound
// 							}
// 						default:
// 							behaviorsRequest = append(behaviorsRequest, behavior)
// 						}
// 					}
// 					req.Behaviors = behaviorsRequest
// 					updated, err := client.UpdateRulesEngine(ctx, req)
// 					if err != nil {
// 						return err
// 					}
// 					newRule := contracts.AzionJsonDataRules{
// 						Id:    updated.GetId(),
// 						Name:  updated.GetName(),
// 						Phase: updated.GetPhase(),
// 					}
// 					ruleConf = append(ruleConf, newRule)
// 				} else {
// 					req := &apiEdgeApplications.CreateRulesEngineRequest{
// 						EdgeApplicationRuleEngineRequest: rule,
// 					}
// 					appstring := strconv.FormatInt(conf.Application.ID, 10)
// 					behaviorsRequest := []edgesdk.EdgeApplicationBehaviorFieldRequest{}
// 					for _, behavior := range rule.Behaviors {
// 						switch behavior.Name {
// 						case "run_function":
// 							if _, ok := FunctionIds[*behavior.GetArgument().String]; !ok {
// 								return msg.ErrorFuncNotFound
// 							}
// 							funcToWorkWith := FunctionIds[*behavior.GetArgument().String]
// 							var beh edgesdk.EdgeApplicationBehaviorFieldRequest
// 							cacheId, err := doCacheForRule(ctx, client, conf, funcToWorkWith)
// 							if err != nil {
// 								return err
// 							}
// 							beh.SetName("set_cache_policy")
// 							beh.SetArgument(edgesdk.EdgeApplicationBehaviorPolymorphicArgumentRequest{
// 								Int64: &cacheId,
// 							})
// 							behaviorsRequest = append(behaviorsRequest, beh)
// 							str := strconv.FormatInt(funcToWorkWith.InstanceID, 10)
// 							beh.SetName("run_function")
// 							beh.SetArgument(edgesdk.EdgeApplicationBehaviorPolymorphicArgumentRequest{
// 								String: &str,
// 							})
// 							behaviorsRequest = append(behaviorsRequest, beh)
// 						case "set_cache_policy":
// 							if id := CacheIdsBackup[*behavior.GetArgument().String]; id > 0 {
// 								var beh edgesdk.EdgeApplicationBehaviorFieldRequest
// 								beh.SetName("set_cache_policy")
// 								beh.SetArgument(edgesdk.EdgeApplicationBehaviorPolymorphicArgumentRequest{
// 									Int64: &id,
// 								})
// 								delete(CacheIds, *behavior.GetArgument().String)
// 							} else {
// 								logger.Debug("Cache Setting not found", zap.Any("Target", *behavior.GetArgument().String))
// 								return msg.ErrorCacheNotFound
// 							}
// 						case "set_edge_connector":
// 							if id := ConnectorIds[*behavior.GetArgument().String]; id > 0 {
// 								var beh edgesdk.EdgeApplicationBehaviorFieldRequest
// 								beh.SetName("set_edge_connector")
// 								beh.SetArgument(edgesdk.EdgeApplicationBehaviorPolymorphicArgumentRequest{
// 									Int64: &id,
// 								})
// 								delete(ConnectorIds, *behavior.GetArgument().String)
// 							} else {
// 								logger.Debug("Edge Connector not found", zap.Any("Target", id))
// 								return msg.ErrorConnectorNotFound
// 							}
// 						default:
// 							behaviorsRequest = append(behaviorsRequest, behavior)
// 						}
// 					}
// 					req.Behaviors = behaviorsRequest
// 					created, err := client.CreateRulesEngine(ctx, appstring, rule.Phase, req)
// 					if err != nil {
// 						return err
// 					}
// 					newRule := contracts.AzionJsonDataRules{
// 						Id:    created.GetId(),
// 						Name:  created.GetName(),
// 						Phase: created.GetPhase(),
// 					}
// 					ruleConf = append(ruleConf, newRule)
// 				}
// 			}

// 			conf.RulesEngine.Rules = ruleConf
// 			err = man.WriteAzionJsonContent(conf, projectConf)
// 			if err != nil {
// 				logger.Debug("Error while writing azion.json file", zap.Error(err))
// 				return err
// 			}
// 		}

// 	}

// 	err = man.WriteAzionJsonContent(conf, projectConf)
// 	if err != nil {
// 		logger.Debug("Error while writing azion.json file", zap.Error(err))
// 		return err
// 	}

// 	// if len(manifest.EdgeStorage) > 0 {
// 	// 	storageman := manifest.EdgeStorage[0]

// 	// }

// 	connectorConf := []contracts.AzionJsonDataConnectors{}
// 	if len(manifest.EdgeConnectors) > 0 {
// 		connector := manifest.EdgeConnectors[0]
// 		connName := getConnectorName(connector, conf.Name)
// 		if id := ConnectorIds[connName]; id > 0 {
// 			request := transformEdgeConnectorRequest(connector)
// 			idstring := strconv.FormatInt(id, 10)
// 			connectorResp, err := connectorClient.Update(ctx, request, idstring)
// 			if err != nil {
// 				return err
// 			}
// 			conn := contracts.AzionJsonDataConnectors{}
// 			// conn.Address = connectorResp.GetAddresses()
// 			conn.Id = connectorResp.GetId()
// 			conn.Name = connectorResp.GetName()
// 			connectorConf = append(connectorConf, conn)
// 		} else {
// 			request := apiConnector.CreateRequest{
// 				EdgeConnectorPolymorphicRequest: connector,
// 			}
// 			connectorResp, err := connectorClient.Create(ctx, &request)
// 			if err != nil {
// 				return err
// 			}
// 			conn := contracts.AzionJsonDataConnectors{}
// 			// conn.Address = connectorResp.GetAddresses()
// 			conn.Id = connectorResp.GetId()
// 			conn.Name = connectorResp.GetName()
// 			connectorConf = append(connectorConf, conn)
// 		}
// 	}

// 	conf.Connectors = connectorConf
// 	err = man.WriteAzionJsonContent(conf, projectConf)
// 	if err != nil {
// 		logger.Debug("Error while writing azion.json file", zap.Error(err))
// 		return err
// 	}

// 	if len(manifest.Workloads) > 0 {
// 		workloadMan := manifest.Workloads[0]
// 		if conf.Workloads.Id > 0 {
// 			request := transformWorkloadRequestUpdate(workloadMan)
// 			request.Id = conf.Workloads.Id
// 			updated, err := clientWorkload.Update(ctx, request)
// 			if err != nil {
// 				return err
// 			}
// 			conf.Workloads.Domains = updated.GetDomains()
// 			conf.Workloads.Url = utils.Concat("https://", conf.Workloads.Domains[0])
// 		} else {
// 			request := transformWorkloadRequestCreate(workloadMan, conf.Application.ID)
// 			resp, err := clientWorkload.Create(ctx, request)
// 			if err != nil {
// 				return err
// 			}
// 			conf.Workloads.Id = resp.GetId()
// 			conf.Workloads.Name = resp.GetName()
// 			conf.Workloads.Domains = resp.GetDomains()
// 			conf.Workloads.Url = utils.Concat("https://", resp.GetDomains()[0])
// 		}
// 	}

// 	err = man.WriteAzionJsonContent(conf, projectConf)
// 	if err != nil {
// 		logger.Debug("Error while writing azion.json file", zap.Error(err))
// 		return err
// 	}

// 	clipurge := apipurge.NewClient(f.HttpClient, f.Config.GetString("api_v4_url"), f.Config.GetString("token"))
// 	for _, purgeObj := range manifest.Purge {
// 		err = clipurge.PurgeCache(ctx, purgeObj.Items, purgeObj.Type, *purgeObj.Layer)
// 		if err != nil {
// 			logger.Debug("Error while purging domains", zap.Error(err))
// 			return err
// 		}
// 	}

// 	err = deleteResources(ctx, f, conf, msgs)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// // this is called to delete resources no longer present in manifest.json
// func deleteResources(ctx context.Context, f *cmdutil.Factory, conf *contracts.AzionApplicationOptions, msgs *[]string) error {
// 	client := apiEdgeApplications.NewClient(f.HttpClient, f.Config.GetString("api_v4_url"), f.Config.GetString("token"))
// 	clientCache := apiCache.NewClientV4(f.HttpClient, f.Config.GetString("api_v4_url"), f.Config.GetString("token"))
// 	// clientOrigin := apiOrigin.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))

// 	for _, value := range RuleIds {
// 		//since until [UXE-3599] was carried out we'd only cared about "request" phase, this check guarantees that if Phase is empty
// 		// we are probably dealing with a rule engine from a previous version
// 		phase := "request"
// 		if value.Phase != "" {
// 			phase = value.Phase
// 		}
// 		str := strconv.FormatInt(conf.Application.ID, 10)
// 		strRule := strconv.FormatInt(value.Id, 10)
// 		var statusInt int
// 		var err error
// 		switch phase {
// 		case "request":
// 			statusInt, err = client.DeleteRulesEngineRequest(ctx, str, phase, strRule)
// 		case "response":
// 			statusInt, err = client.DeleteRulesEngineResponse(ctx, str, phase, strRule)
// 		default:
// 			return msgrule.ErrorInvalidPhase
// 		}

// 		if statusInt == 404 {
// 			logger.Debug("Rule Engine not found. Skipping delete")
// 			continue
// 		}
// 		if err != nil {
// 			return err
// 		}
// 		msgf := fmt.Sprintf(msgrule.DeleteOutputSuccess+"\n", value.Id)
// 		logger.FInfoFlags(f.IOStreams.Out, msgf, f.Format, f.Out)
// 		*msgs = append(*msgs, msgf)
// 	}

// 	// for i, value := range OriginKeys {
// 	// 	if strings.Contains(i, "_single") {
// 	// 		continue
// 	// 	}
// 	// 	err := clientOrigin.DeleteOrigins(ctx, conf.Application.ID, value)
// 	// 	if err != nil {
// 	// 		return err
// 	// 	}
// 	// 	msgf := fmt.Sprintf(msgorigin.DeleteOutputSuccess+"\n", value)
// 	// 	logger.FInfoFlags(f.IOStreams.Out, msgf, f.Format, f.Out)
// 	// 	*msgs = append(*msgs, msgf)
// 	// }

// 	for _, value := range CacheIds {
// 		status, err := clientCache.Delete(ctx, conf.Application.ID, value)
// 		if status == 404 {
// 			logger.Debug("Cache Setting not found. Skipping delete")
// 			continue
// 		}
// 		if err != nil {
// 			return err
// 		}
// 		msgf := fmt.Sprintf(msgcache.DeleteOutputSuccess+"\n", value)
// 		logger.FInfoFlags(f.IOStreams.Out, msgf, f.Format, f.Out)
// 		*msgs = append(*msgs, msgf)
// 	}

// 	return nil
// }
