package manifest

import (
	"context"
	"fmt"
	"strconv"

	"encoding/json"

	msg "github.com/aziontech/azion-cli/messages/manifest"
	apiCache "github.com/aziontech/azion-cli/pkg/api/cache_setting"
	apiEdgeApplications "github.com/aziontech/azion-cli/pkg/api/edge_applications"
	apiConnector "github.com/aziontech/azion-cli/pkg/api/edge_connector"
	apiWorkloads "github.com/aziontech/azion-cli/pkg/api/workloads"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/logger"
	edgesdk "github.com/aziontech/azionapi-v4-go-sdk/edge-api"
	"go.uber.org/zap"
)

//TODO: FIX HERE

func transformEdgeConnectorRequest(connectorRequest edgesdk.EdgeConnectorPolymorphicRequest) *apiConnector.UpdateRequest {
	request := &apiConnector.UpdateRequest{}

	if connectorRequest.EdgeConnectorHTTPRequest != nil {
		bodyRequest := connectorRequest.EdgeConnectorHTTPRequest
		atts := bodyRequest.Attributes
		body := edgesdk.PatchedEdgeConnectorHTTPRequest{}
		if bodyRequest.Active != nil {
			body.SetActive(*bodyRequest.Active)
		}

		if bodyRequest.Name != "" {
			body.SetName(bodyRequest.Name)
		}

		body.SetType(bodyRequest.Type)

		body.SetAttributes(atts)
		request.PatchedEdgeConnectorHTTPRequest = &body
		return request
	}

	if connectorRequest.EdgeConnectorLiveIngestRequest != nil {
		body := edgesdk.PatchedEdgeConnectorLiveIngestRequest{}
		bodyRequest := connectorRequest.EdgeConnectorLiveIngestRequest
		if bodyRequest.Active != nil {
			body.SetActive(*bodyRequest.Active)
		}
		body.SetType(bodyRequest.Type)
		body.SetAttributes(bodyRequest.Attributes)

		if bodyRequest.Name != "" {
			body.SetName(bodyRequest.Name)
		}

		if bodyRequest.Type != "" {
			body.SetType(bodyRequest.Type)
		}
		request.PatchedEdgeConnectorLiveIngestRequest = &body
		return request
	}

	if connectorRequest.EdgeConnectorStorageRequest != nil {
		body := edgesdk.PatchedEdgeConnectorStorageRequest{}
		bodyRequest := connectorRequest.EdgeConnectorStorageRequest
		if bodyRequest.Active != nil {
			body.SetActive(*bodyRequest.Active)
		}

		if bodyRequest.Name != "" {
			body.SetName(bodyRequest.Name)
		}

		body.SetType(bodyRequest.Type)

		body.SetAttributes(bodyRequest.Attributes)

		request.PatchedEdgeConnectorStorageRequest = &body
		return request
	}

	return request
}

func transformWorkloadRequestUpdate(createRequest contracts.WorkloadManifest) *apiWorkloads.UpdateRequest {
	request := &apiWorkloads.UpdateRequest{}

	if createRequest.Name != "" {
		request.SetName(createRequest.Name)
	}
	if createRequest.Active != nil {
		request.SetActive(*createRequest.Active)
	}
	if len(createRequest.Domains) > 0 {
		request.SetDomains(createRequest.Domains)
	}
	if createRequest.Mtls != nil {
		request.SetMtls(*createRequest.Mtls)
	}
	if createRequest.Protocols != nil {
		request.SetProtocols(*createRequest.Protocols)
	}
	if createRequest.Tls != nil {
		request.SetTls(*createRequest.Tls)
	}

	return request
}

func transformWorkloadRequestCreate(createRequest contracts.WorkloadManifest, appid int64) *apiWorkloads.CreateRequest {
	request := &apiWorkloads.CreateRequest{}

	if createRequest.Name != "" {
		request.SetName(createRequest.Name)
	}
	if createRequest.Active != nil {
		request.SetActive(*createRequest.Active)
	}
	if len(createRequest.Domains) > 0 {
		request.SetDomains(createRequest.Domains)
	}
	if createRequest.Mtls != nil {
		request.SetMtls(*createRequest.Mtls)
	}
	if createRequest.Protocols != nil {
		request.SetProtocols(*createRequest.Protocols)
	}
	if createRequest.Tls != nil {
		request.SetTls(*createRequest.Tls)
	}

	return request
}

func transformEdgeApplicationRequestUpdate(edgeapprequest contracts.EdgeApplications) *apiEdgeApplications.UpdateRequest {
	request := &apiEdgeApplications.UpdateRequest{}

	if edgeapprequest.Active != nil {
		request.SetActive(*edgeapprequest.Active)
	}
	if edgeapprequest.Debug != nil {
		request.SetDebug(*edgeapprequest.Debug)
	}
	type Modules struct {
		EdgeCacheEnabled              bool `json:"edge_cache_enabled"`
		EdgeFunctionsEnabled          bool `json:"edge_functions_enabled"`
		ApplicationAcceleratorEnabled bool `json:"application_accelerator_enabled"`
		ImageProcessorEnabled         bool `json:"image_processor_enabled"`
		TieredCacheEnabled            bool `json:"tiered_cache_enabled"`
	}
	if edgeapprequest.Modules != nil {
		request.SetModules(*edgeapprequest.Modules)
	}

	//TODO: Fix Here
	if edgeapprequest.Name != "" {
		request.SetName(edgeapprequest.Name)
	}

	return request
}

func transformEdgeApplicationRequestCreate(edgeapprequest contracts.EdgeApplications) *apiEdgeApplications.CreateRequest {
	request := &apiEdgeApplications.CreateRequest{}

	if edgeapprequest.Active != nil {
		request.SetActive(*edgeapprequest.Active)
	}
	if edgeapprequest.Debug != nil {
		request.SetDebug(*edgeapprequest.Debug)
	}
	type Modules struct {
		EdgeCacheEnabled              bool `json:"edge_cache_enabled"`
		EdgeFunctionsEnabled          bool `json:"edge_functions_enabled"`
		ApplicationAcceleratorEnabled bool `json:"application_accelerator_enabled"`
		ImageProcessorEnabled         bool `json:"image_processor_enabled"`
		TieredCacheEnabled            bool `json:"tiered_cache_enabled"`
	}
	// if edgeapprequest.Modules != nil {
	// 	modules := edgesdk.EdgeApplicationModulesRequest{}
	// 	if edgeapprequest.Modules.ApplicationAcceleratorEnabled != nil {
	// 		modules.SetApplicationAcceleratorEnabled(*edgeapprequest.Modules.ApplicationAcceleratorEnabled)
	// 	}
	// 	if edgeapprequest.Modules.EdgeCacheEnabled != nil {
	// 		modules.SetEdgeCacheEnabled(*edgeapprequest.Modules.EdgeCacheEnabled)
	// 	}
	// 	if edgeapprequest.Modules.EdgeFunctionsEnabled != nil {
	// 		modules.SetEdgeFunctionsEnabled(*edgeapprequest.Modules.EdgeFunctionsEnabled)
	// 	}
	// 	if edgeapprequest.Modules.ImageProcessorEnabled != nil {
	// 		modules.SetImageProcessorEnabled(*edgeapprequest.Modules.ImageProcessorEnabled)
	// 	}
	// 	if edgeapprequest.Modules.TieredCacheEnabled != nil {
	// 		modules.SetTieredCacheEnabled(*edgeapprequest.Modules.TieredCacheEnabled)
	// 	}

	// 	request.SetModules(modules)
	// }

	//TODO: Fix Here
	if edgeapprequest.Name != "" {
		request.SetName(edgeapprequest.Name)
	}

	return request
}

func transformCacheRequest(cache contracts.ManifestCacheSetting) *apiCache.RequestUpdate {
	request := apiCache.RequestUpdate{}

	if cache.Name != "" {
		request.SetName(cache.Name)
	}

	if cache.BrowserCache != nil {
		request.SetBrowserCache(*cache.BrowserCache)
	}
	if cache.Modules != nil {
		request.SetModules(*cache.Modules)
	}

	return &request
}

func transformCacheRequestCreate(cache contracts.ManifestCacheSetting) *apiCache.Request {
	request := apiCache.Request{}

	if cache.Name != "" {
		request.SetName(cache.Name)
	}

	if cache.BrowserCache != nil {
		request.SetBrowserCache(*cache.BrowserCache)
	}
	if cache.Modules != nil {
		request.SetModules(*cache.Modules)
	}

	return &request
}

func transformRuleResponse(rule contracts.ManifestRule) *apiEdgeApplications.UpdateRulesEngineResponse {
	request := &apiEdgeApplications.UpdateRulesEngineResponse{}

	request.SetActive(rule.Active)
	// if rule.Behaviors != nil {
	// 	request.SetBehaviors(rule.Behaviors)
	// }
	if rule.Criteria != nil {
		request.SetCriteria(rule.Criteria)
	}

	//TODO: Fix Here
	request.SetDescription(rule.Description)
	request.SetName(rule.Name)

	return request
}

func transformRuleRequest(rule contracts.ManifestRule) *apiEdgeApplications.UpdateRulesEngineRequest {
	request := &apiEdgeApplications.UpdateRulesEngineRequest{}

	request.SetActive(rule.Active)
	if rule.Criteria != nil {
		request.SetCriteria(rule.Criteria)
	}

	request.SetDescription(rule.Description)
	request.SetName(rule.Name)

	return request
}

func doCacheForRule(ctx context.Context, client *apiEdgeApplications.Client, conf *contracts.AzionApplicationOptions, function contracts.AzionJsonDataFunction) (int64, error) {
	if function.CacheId > 0 {
		return function.CacheId, nil
	}
	var reqCache apiEdgeApplications.CreateCacheSettingsRequest
	reqCache.SetName("function-policy")
	// reqCache.BrowserCache = edgesdk.BrowserCacheModuleRequest{
	// 	Behavior: "honor",
	// }
	// reqCache.EdgeCache = edgesdk.EdgeCacheModuleRequest{
	// 	Behavior: "honor",
	// 	MaxAge:   0,
	// }
	// reqCache.ApplicationControls = edgesdk.ApplicationControlsModuleRequest{
	// 	CacheByQueryString: "all",
	// 	CacheByCookies:     "all",
	// }

	//TODO: Fix Here

	// create cache to function next
	str := strconv.FormatInt(conf.Application.ID, 10)
	cache, err := client.CreateCacheEdgeApplication(ctx, &reqCache, str)
	if err != nil {
		logger.Debug("Error while creating Cache Settings", zap.Error(err))
		return 0, err
	}

	// conf.Function.CacheId = cache.GetId()

	return cache.GetId(), nil
}

func getConnectorName(connector edgesdk.EdgeConnectorPolymorphicRequest, defaultName string) (string, string) {
	if connector.EdgeConnectorHTTPRequest != nil {
		return connector.EdgeConnectorHTTPRequest.Name, "http"
	}

	if connector.EdgeConnectorLiveIngestRequest != nil {
		return connector.EdgeConnectorLiveIngestRequest.Name, "ingest"
	}

	if connector.EdgeConnectorStorageRequest != nil {
		return connector.EdgeConnectorStorageRequest.Name, "storage"
	}

	return defaultName, ""
}

func transformBehaviorsRequest(behaviors []contracts.ManifestRuleBehavior) ([]edgesdk.EdgeApplicationRuleEngineRequestPhaseBehaviorsRequest, error) {
	behaviorsRequest := []edgesdk.EdgeApplicationRuleEngineRequestPhaseBehaviorsRequest{}
	for _, behavior := range behaviors {
		var withArgs edgesdk.EdgeApplicationRequestPhaseBehaviorWithArgsRequest
		var withoutArgs edgesdk.EdgeApplicationRequestPhaseBehaviorWithoutArgsRequest
		var captureMatchGroups edgesdk.EdgeApplicationRequestPhaseBehaviorCaptureMatchGroupsRequest
		var beh edgesdk.EdgeApplicationRuleEngineRequestPhaseBehaviorsRequest
		switch behavior.Type {
		case "run_function":
			attributesJSON, err := json.Marshal(behavior.Attributes)
			if err != nil {
				return nil, err
			}
			var attributes edgesdk.EdgeApplicationRuleEngineStringAttributes
			if err := json.Unmarshal(attributesJSON, &attributes); err != nil {
				return nil, err
			}
			withArgs.SetType("run_function")
			if attributes.Value.Int64 != nil {
				withArgs.SetAttributes(attributes)
			} else if attributes.Value.String != nil {
				funcName := *attributes.Value.String
				if _, ok := FunctionIds[funcName]; !ok {
					return nil, msg.ErrorFuncNotFound
				}
				funcToWorkWith := FunctionIds[funcName]
				v := edgesdk.EdgeApplicationRuleEngineStringAttributesValue{
					Int64: &funcToWorkWith.InstanceID,
				}
				attributes.SetValue(v)
				withArgs.SetAttributes(attributes)
			}
			beh.EdgeApplicationRequestPhaseBehaviorWithArgsRequest = &withArgs
			behaviorsRequest = append(behaviorsRequest, beh)
		case "set_cache_policy":
			attributesJSON, err := json.Marshal(behavior.Attributes)
			if err != nil {
				return nil, err
			}
			var attributes edgesdk.EdgeApplicationRuleEngineStringAttributes
			if err := json.Unmarshal(attributesJSON, &attributes); err != nil {
				return nil, err
			}
			withArgs.SetType("set_cache_policy")
			if attributes.Value.Int64 != nil {
				withArgs.SetAttributes(attributes)
			} else if attributes.Value.String != nil {
				cacheName := *attributes.Value.String
				if id := CacheIdsBackup[cacheName]; id > 0 {
					fmt.Println("ID", id)
					v := edgesdk.EdgeApplicationRuleEngineStringAttributesValue{
						Int64: &id,
					}
					attributes.SetValue(v)
					withArgs.SetAttributes(attributes)
					fmt.Println(cacheName)
					delete(CacheIds, cacheName)
				} else {
					logger.Debug("Cache Setting not found", zap.Any("Target", *attributes.Value.String))
					return nil, msg.ErrorCacheNotFound
				}
			}
			beh.EdgeApplicationRequestPhaseBehaviorWithArgsRequest = &withArgs
			behaviorsRequest = append(behaviorsRequest, beh)
		case "set_edge_connector":
			attributesJSON, err := json.Marshal(behavior.Attributes)
			if err != nil {
				return nil, err
			}
			var attributes edgesdk.EdgeApplicationRuleEngineStringAttributes
			if err := json.Unmarshal(attributesJSON, &attributes); err != nil {
				return nil, err
			}
			withArgs.SetType("set_edge_connector")
			if attributes.Value.Int64 != nil {
				withArgs.SetAttributes(attributes)
			} else if attributes.Value.String != nil {
				connectorName := *attributes.Value.String
				if id := ConnectorIds[connectorName]; id > 0 {
					v := edgesdk.EdgeApplicationRuleEngineStringAttributesValue{
						Int64: &id,
					}
					attributes.SetValue(v)
					withArgs.SetAttributes(attributes)
					delete(ConnectorIds, connectorName)
				} else {
					logger.Debug("Edge Connector not found", zap.Any("Target", connectorName))
					return nil, msg.ErrorConnectorNotFound
				}
			}
			beh.EdgeApplicationRequestPhaseBehaviorWithArgsRequest = &withArgs
			behaviorsRequest = append(behaviorsRequest, beh)
		case "capture_match_groups":
			attributesJSON, err := json.Marshal(behavior.Attributes)
			if err != nil {
				return nil, err
			}
			var attributes edgesdk.EdgeApplicationRuleEngineCaptureMatchGroupsAttributesRequest
			if err := json.Unmarshal(attributesJSON, &attributes); err != nil {
				return nil, err
			}
			captureMatchGroups.SetType("capture_match_groups")
			captureMatchGroups.SetAttributes(attributes)
			beh.EdgeApplicationRequestPhaseBehaviorCaptureMatchGroupsRequest = &captureMatchGroups
			behaviorsRequest = append(behaviorsRequest, beh)
		case "redirect_to_301", "redirect_to_302", "filter_request_cookie", "rewrite_request", "add_request_header", "filter_request_header", "add_request_cookie":
			attributesJSON, err := json.Marshal(behavior.Attributes)
			if err != nil {
				return nil, err
			}
			var attributes edgesdk.EdgeApplicationRuleEngineStringAttributes
			if err := json.Unmarshal(attributesJSON, &attributes); err != nil {
				return nil, err
			}
			withArgs.SetType(behavior.Type)
			withArgs.SetAttributes(attributes)
			beh.EdgeApplicationRequestPhaseBehaviorWithArgsRequest = &withArgs
			behaviorsRequest = append(behaviorsRequest, beh)
		default:
			withoutArgs.SetType(behavior.Type)
			beh.EdgeApplicationRequestPhaseBehaviorWithoutArgsRequest = &withoutArgs
			behaviorsRequest = append(behaviorsRequest, beh)
		}
	}

	return behaviorsRequest, nil
}

func transformBehaviorsResponse(behaviors []contracts.ManifestRuleBehavior) ([]edgesdk.EdgeApplicationRuleEngineResponsePhaseBehaviorsRequest, error) {
	behaviorsResponse := []edgesdk.EdgeApplicationRuleEngineResponsePhaseBehaviorsRequest{}

	for _, behavior := range behaviors {
		var withArgs edgesdk.EdgeApplicationResponsePhaseBehaviorWithArgsRequest
		var withoutArgs edgesdk.EdgeApplicationResponsePhaseBehaviorWithoutArgsRequest
		var captureMatchGroups edgesdk.EdgeApplicationResponsePhaseBehaviorCaptureMatchGroupsRequest
		var beh edgesdk.EdgeApplicationRuleEngineResponsePhaseBehaviorsRequest

		switch behavior.Type {
		case "capture_match_groups":
			attributesJSON, err := json.Marshal(behavior.Attributes)
			if err != nil {
				return nil, err
			}
			var attributes edgesdk.EdgeApplicationRuleEngineCaptureMatchGroupsAttributes
			if err := json.Unmarshal(attributesJSON, &attributes); err != nil {
				return nil, err
			}
			captureMatchGroups.SetType("capture_match_groups")
			captureMatchGroups.SetAttributes(attributes)
			beh.EdgeApplicationResponsePhaseBehaviorCaptureMatchGroupsRequest = &captureMatchGroups
			behaviorsResponse = append(behaviorsResponse, beh)
		case "enable_gzip", "deliver":
			withoutArgs.SetType(behavior.Type)
			beh.EdgeApplicationResponsePhaseBehaviorWithoutArgsRequest = &withoutArgs
			behaviorsResponse = append(behaviorsResponse, beh)
		default:
			// Everything else is WithArgs string
			attributesJSON, err := json.Marshal(behavior.Attributes)
			if err != nil {
				return nil, err
			}
			var attributes edgesdk.EdgeApplicationRuleEngineStringAttributes
			if err := json.Unmarshal(attributesJSON, &attributes); err != nil {
				return nil, err
			}
			withArgs.SetType(behavior.Type)
			withArgs.SetAttributes(attributes)
			beh.EdgeApplicationResponsePhaseBehaviorWithArgsRequest = &withArgs
			behaviorsResponse = append(behaviorsResponse, beh)
		}
	}

	return behaviorsResponse, nil
}

func transformRuleRequestCreate(rule contracts.ManifestRule) edgesdk.EdgeApplicationRequestPhaseRuleEngineRequest {
	request := edgesdk.EdgeApplicationRequestPhaseRuleEngineRequest{}

	request.SetActive(rule.Active)
	if rule.Criteria != nil {
		request.SetCriteria(rule.Criteria)
	}
	request.SetDescription(rule.Description)
	request.SetName(rule.Name)

	return request
}

func transformRuleResponseCreate(rule contracts.ManifestRule) edgesdk.EdgeApplicationResponsePhaseRuleEngineRequest {
	request := edgesdk.EdgeApplicationResponsePhaseRuleEngineRequest{}

	request.SetActive(rule.Active)
	if rule.Criteria != nil {
		request.SetCriteria(rule.Criteria)
	}
	request.SetDescription(rule.Description)
	request.SetName(rule.Name)

	return request
}
