package manifest

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	msg "github.com/aziontech/azion-cli/messages/manifest"
	apiApplications "github.com/aziontech/azion-cli/pkg/api/applications"
	apiCache "github.com/aziontech/azion-cli/pkg/api/cache_setting"
	apiConnector "github.com/aziontech/azion-cli/pkg/api/connector"
	apiWorkloads "github.com/aziontech/azion-cli/pkg/api/workloads"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/utils"
	edgesdk "github.com/aziontech/azionapi-v4-go-sdk-dev/edge-api"
	"go.uber.org/zap"
)

func transformEdgeConnectorRequest(connectorRequest edgesdk.ConnectorPolymorphicRequest) *apiConnector.UpdateRequest {
	request := &apiConnector.UpdateRequest{}

	if connectorRequest.ConnectorHTTPRequest != nil {
		bodyRequest := connectorRequest.ConnectorHTTPRequest
		atts := bodyRequest.Attributes
		body := edgesdk.PatchedConnectorHTTPRequest{}
		if bodyRequest.Active != nil {
			body.SetActive(*bodyRequest.Active)
		}

		if bodyRequest.Name != "" {
			body.SetName(bodyRequest.Name)
		}

		body.SetType(bodyRequest.Type)

		body.SetAttributes(atts)
		request.PatchedConnectorHTTPRequest = &body
		return request
	}

	if connectorRequest.ConnectorLiveIngestRequest != nil {
		body := edgesdk.PatchedConnectorLiveIngestRequest{}
		bodyRequest := connectorRequest.ConnectorLiveIngestRequest
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
		request.PatchedConnectorLiveIngestRequest = &body
		return request
	}

	if connectorRequest.ConnectorStorageRequest != nil {
		body := edgesdk.PatchedConnectorStorageRequest{}
		bodyRequest := connectorRequest.ConnectorStorageRequest
		if bodyRequest.Active != nil {
			body.SetActive(*bodyRequest.Active)
		}

		if bodyRequest.Name != "" {
			body.SetName(bodyRequest.Name)
		}

		body.SetType(bodyRequest.Type)

		body.SetAttributes(bodyRequest.Attributes)

		request.PatchedConnectorStorageRequest = &body
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

func transformWorkloadDeploymentRequestUpdate(updateRequest contracts.WorkloadDeployment, conf *contracts.AzionApplicationOptions) edgesdk.PatchedWorkloadDeploymentRequest {
	request := edgesdk.PatchedWorkloadDeploymentRequest{}

	if updateRequest.Name != "" {
		request.SetName(updateRequest.Name)
	}

	request.SetActive(updateRequest.Active)

	request.SetCurrent(updateRequest.Current)

	strategy := edgesdk.DeploymentStrategyDefaultDeploymentStrategyRequest{}
	attributes := edgesdk.DefaultDeploymentStrategyAttrsRequest{}

	if updateRequest.Strategy.Type != "" {
		strategy.SetType(updateRequest.Strategy.Type)
	}

	attributes.SetApplication(conf.Application.ID)
	strategy.SetAttributes(attributes)
	request.SetStrategy(strategy)

	return request
}

func transformWorkloadDeploymentRequestCreate(createRequest contracts.WorkloadDeployment, conf *contracts.AzionApplicationOptions) edgesdk.WorkloadDeploymentRequest {
	request := edgesdk.WorkloadDeploymentRequest{}

	if createRequest.Name != "" {
		request.SetName(createRequest.Name)
	}
	request.SetActive(createRequest.Active)

	request.SetCurrent(createRequest.Current)

	strategy := edgesdk.DeploymentStrategyDefaultDeploymentStrategyRequest{}
	attributes := edgesdk.DefaultDeploymentStrategyAttrsRequest{}

	if createRequest.Strategy.Type != "" {
		strategy.SetType(createRequest.Strategy.Type)
	}

	attributes.SetApplication(conf.Application.ID)
	strategy.SetAttributes(attributes)
	request.SetStrategy(strategy)

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

func transformEdgeApplicationRequestUpdate(edgeapprequest contracts.Applications) *apiApplications.UpdateRequest {
	request := &apiApplications.UpdateRequest{}

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

	if edgeapprequest.Name != "" {
		request.SetName(edgeapprequest.Name)
	}

	return request
}

func transformEdgeApplicationRequestCreate(edgeapprequest contracts.Applications) *apiApplications.CreateRequest {
	request := &apiApplications.CreateRequest{}

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

func transformRuleResponse(rule contracts.ManifestRule) *apiApplications.UpdateRulesEngineResponse {
	request := &apiApplications.UpdateRulesEngineResponse{}

	request.SetActive(rule.Active)
	// if rule.Behaviors != nil {
	// 	request.SetBehaviors(rule.Behaviors)
	// }
	if rule.Criteria != nil {
		request.SetCriteria(rule.Criteria)
	}

	request.SetDescription(rule.Description)
	request.SetName(rule.Name)

	return request
}

func transformRuleRequest(rule contracts.ManifestRule) *apiApplications.UpdateRulesEngineRequest {
	request := &apiApplications.UpdateRulesEngineRequest{}

	request.SetActive(rule.Active)
	if rule.Criteria != nil {
		request.SetCriteria(rule.Criteria)
	}

	request.SetDescription(rule.Description)
	request.SetName(rule.Name)

	return request
}

func getConnectorName(connector edgesdk.ConnectorPolymorphicRequest, defaultName string) (string, string) {
	if connector.ConnectorHTTPRequest != nil {
		return connector.ConnectorHTTPRequest.Name, "http"
	}

	if connector.ConnectorLiveIngestRequest != nil {
		return connector.ConnectorLiveIngestRequest.Name, "ingest"
	}

	if connector.ConnectorStorageRequest != nil {
		return connector.ConnectorStorageRequest.Name, "storage"
	}

	return defaultName, ""
}

func transformBehaviorsRequest(behaviors []contracts.ManifestRuleBehavior) ([]edgesdk.ApplicationRuleEngineRequestPhaseBehaviorsRequest, error) {
	behaviorsRequest := []edgesdk.ApplicationRuleEngineRequestPhaseBehaviorsRequest{}
	for _, behavior := range behaviors {
		var withArgs edgesdk.ApplicationRequestPhaseBehaviorWithArgsRequest
		var withoutArgs edgesdk.ApplicationRequestPhaseBehaviorWithoutArgsRequest
		var captureMatchGroups edgesdk.ApplicationRequestPhaseBehaviorCaptureMatchGroupsRequest
		var beh edgesdk.ApplicationRuleEngineRequestPhaseBehaviorsRequest
		switch behavior.Type {
		case "run_function":
			attributesJSON, err := json.Marshal(behavior.Attributes)
			if err != nil {
				return nil, err
			}
			var attributes edgesdk.ApplicationRuleEngineStringAttributes
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
				v := edgesdk.ApplicationRuleEngineStringAttributesValue{
					Int64: &funcToWorkWith.InstanceID,
				}
				attributes.SetValue(v)
				withArgs.SetAttributes(attributes)
			}
			beh.ApplicationRequestPhaseBehaviorWithArgsRequest = &withArgs
			behaviorsRequest = append(behaviorsRequest, beh)
		case "set_cache_policy":
			attributesJSON, err := json.Marshal(behavior.Attributes)
			if err != nil {
				return nil, err
			}
			var attributes edgesdk.ApplicationRuleEngineStringAttributes
			if err := json.Unmarshal(attributesJSON, &attributes); err != nil {
				return nil, err
			}
			withArgs.SetType("set_cache_policy")
			if attributes.Value.Int64 != nil {
				withArgs.SetAttributes(attributes)
			} else if attributes.Value.String != nil {
				cacheName := *attributes.Value.String
				if id := CacheIdsBackup[cacheName]; id > 0 {
					v := edgesdk.ApplicationRuleEngineStringAttributesValue{
						Int64: &id,
					}
					attributes.SetValue(v)
					withArgs.SetAttributes(attributes)
					delete(CacheIds, cacheName)
				} else {
					logger.Debug("Cache Setting not found", zap.Any("Target", *attributes.Value.String))
					return nil, msg.ErrorCacheNotFound
				}
			}
			beh.ApplicationRequestPhaseBehaviorWithArgsRequest = &withArgs
			behaviorsRequest = append(behaviorsRequest, beh)
		case "set_connector":
			attributesJSON, err := json.Marshal(behavior.Attributes)
			if err != nil {
				return nil, err
			}
			var attributes edgesdk.ApplicationRuleEngineStringAttributes
			if err := json.Unmarshal(attributesJSON, &attributes); err != nil {
				return nil, err
			}
			withArgs.SetType("set_connector")
			if attributes.Value.Int64 != nil {
				withArgs.SetAttributes(attributes)
			} else if attributes.Value.String != nil {
				connectorName := *attributes.Value.String
				if id := ConnectorIds[connectorName]; id > 0 {
					v := edgesdk.ApplicationRuleEngineStringAttributesValue{
						Int64: &id,
					}
					attributes.SetValue(v)
					withArgs.SetAttributes(attributes)
					// delete(ConnectorIds, connectorName)
				} else {
					logger.Debug("Connector not found", zap.Any("Target", connectorName))
					return nil, msg.ErrorConnectorNotFound
				}
			}
			beh.ApplicationRequestPhaseBehaviorWithArgsRequest = &withArgs
			behaviorsRequest = append(behaviorsRequest, beh)
		case "capture_match_groups":
			attributesJSON, err := json.Marshal(behavior.Attributes)
			if err != nil {
				return nil, err
			}
			var attributes edgesdk.ApplicationRuleEngineCaptureMatchGroupsAttributesRequest
			if err := json.Unmarshal(attributesJSON, &attributes); err != nil {
				return nil, err
			}
			captureMatchGroups.SetType("capture_match_groups")
			captureMatchGroups.SetAttributes(attributes)
			beh.ApplicationRequestPhaseBehaviorCaptureMatchGroupsRequest = &captureMatchGroups
			behaviorsRequest = append(behaviorsRequest, beh)
		case "redirect_to_301", "redirect_to_302", "filter_request_cookie", "rewrite_request", "add_request_header", "filter_request_header", "add_request_cookie":
			attributesJSON, err := json.Marshal(behavior.Attributes)
			if err != nil {
				return nil, err
			}
			var attributes edgesdk.ApplicationRuleEngineStringAttributes
			if err := json.Unmarshal(attributesJSON, &attributes); err != nil {
				return nil, err
			}
			withArgs.SetType(behavior.Type)
			withArgs.SetAttributes(attributes)
			beh.ApplicationRequestPhaseBehaviorWithArgsRequest = &withArgs
			behaviorsRequest = append(behaviorsRequest, beh)
		default:
			withoutArgs.SetType(behavior.Type)
			beh.ApplicationRequestPhaseBehaviorWithoutArgsRequest = &withoutArgs
			behaviorsRequest = append(behaviorsRequest, beh)
		}
	}

	return behaviorsRequest, nil
}

func transformBehaviorsResponse(behaviors []contracts.ManifestRuleBehavior) ([]edgesdk.ApplicationRuleEngineResponsePhaseBehaviorsRequest, error) {
	behaviorsResponse := []edgesdk.ApplicationRuleEngineResponsePhaseBehaviorsRequest{}

	for _, behavior := range behaviors {
		var withArgs edgesdk.ApplicationResponsePhaseBehaviorWithArgsRequest
		var withoutArgs edgesdk.ApplicationResponsePhaseBehaviorWithoutArgsRequest
		var captureMatchGroups edgesdk.ApplicationResponsePhaseBehaviorCaptureMatchGroupsRequest
		var beh edgesdk.ApplicationRuleEngineResponsePhaseBehaviorsRequest

		switch behavior.Type {
		case "capture_match_groups":
			attributesJSON, err := json.Marshal(behavior.Attributes)
			if err != nil {
				return nil, err
			}
			var attributes edgesdk.ApplicationRuleEngineCaptureMatchGroupsAttributes
			if err := json.Unmarshal(attributesJSON, &attributes); err != nil {
				return nil, err
			}
			captureMatchGroups.SetType("capture_match_groups")
			captureMatchGroups.SetAttributes(attributes)
			beh.ApplicationResponsePhaseBehaviorCaptureMatchGroupsRequest = &captureMatchGroups
			behaviorsResponse = append(behaviorsResponse, beh)
		case "enable_gzip", "deliver":
			withoutArgs.SetType(behavior.Type)
			beh.ApplicationResponsePhaseBehaviorWithoutArgsRequest = &withoutArgs
			behaviorsResponse = append(behaviorsResponse, beh)
		default:
			// Everything else is WithArgs string
			attributesJSON, err := json.Marshal(behavior.Attributes)
			if err != nil {
				return nil, err
			}
			var attributes edgesdk.ApplicationRuleEngineStringAttributes
			if err := json.Unmarshal(attributesJSON, &attributes); err != nil {
				return nil, err
			}
			withArgs.SetType(behavior.Type)
			withArgs.SetAttributes(attributes)
			beh.ApplicationResponsePhaseBehaviorWithArgsRequest = &withArgs
			behaviorsResponse = append(behaviorsResponse, beh)
		}
	}

	return behaviorsResponse, nil
}

func transformRuleRequestCreate(rule contracts.ManifestRule) edgesdk.ApplicationRequestPhaseRuleEngineRequest {
	request := edgesdk.ApplicationRequestPhaseRuleEngineRequest{}

	request.SetActive(rule.Active)
	if rule.Criteria != nil {
		request.SetCriteria(rule.Criteria)
	}
	request.SetDescription(rule.Description)
	request.SetName(rule.Name)

	return request
}

func transformRuleResponseCreate(rule contracts.ManifestRule) edgesdk.ApplicationResponsePhaseRuleEngineRequest {
	request := edgesdk.ApplicationResponsePhaseRuleEngineRequest{}

	request.SetActive(rule.Active)
	if rule.Criteria != nil {
		request.SetCriteria(rule.Criteria)
	}
	request.SetDescription(rule.Description)
	request.SetName(rule.Name)

	return request
}

func updateCache(f *cmdutil.Factory, cache contracts.ManifestCacheSetting, clientCache *apiCache.ClientV4, conf *contracts.AzionApplicationOptions, r int64, ctx context.Context, edgeappman contracts.Applications) (contracts.AzionJsonDataCacheSettings, error) {
	request := transformCacheRequest(cache)
	updated, err := clientCache.Update(ctx, request, conf.Application.ID, r)
	if errors.Is(err, utils.ErrorNotFound404) {
		logger.Debug("Cache Setting not found. Trying to create", zap.Any("Error", err))
		logger.FInfoFlags(f.IOStreams.Out, fmt.Sprintf(msg.MessageDeleteResource, "\n"), f.Format, f.Out)
		return createCache(cache, clientCache, conf, ctx, edgeappman)
	}
	if err != nil {
		return contracts.AzionJsonDataCacheSettings{}, err
	}
	newCache := contracts.AzionJsonDataCacheSettings{
		Id:   updated.GetData().Id,
		Name: updated.GetData().Name,
	}
	return newCache, nil
}

func createCache(cache contracts.ManifestCacheSetting, clientCache *apiCache.ClientV4, conf *contracts.AzionApplicationOptions, ctx context.Context, edgeappman contracts.Applications) (contracts.AzionJsonDataCacheSettings, error) {
	request := transformCacheRequestCreate(cache)
	responseCache, err := clientCache.Create(ctx, request.CacheSettingRequest, conf.Application.ID)
	if err != nil {
		return contracts.AzionJsonDataCacheSettings{}, err
	}
	newCache := contracts.AzionJsonDataCacheSettings{
		Id:   responseCache.GetId(),
		Name: responseCache.GetName(),
	}
	CacheIds[newCache.Name] = newCache.Id
	return newCache, nil
}
