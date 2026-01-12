package manifest

import (
	"context"
	"encoding/json"
	"errors"

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
	if connectorRequest.ConnectorHTTPRequest != nil {
		request := &apiConnector.UpdateRequest{}
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

	if connectorRequest.ConnectorRequest != nil {
		request := &apiConnector.UpdateRequest{}
		bodyRequest := connectorRequest.ConnectorRequest
		body := edgesdk.PatchedConnectorRequest{}
		if bodyRequest.Active != nil {
			body.SetActive(*bodyRequest.Active)
		}

		if bodyRequest.Name != "" {
			body.SetName(bodyRequest.Name)
		}

		body.SetType(bodyRequest.Type)

		body.SetAttributes(bodyRequest.Attributes)

		request.PatchedConnectorRequest = &body
		return request
	}

	return &apiConnector.UpdateRequest{}
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
	if edgeapprequest.Modules != nil {
		request.SetModules(*edgeapprequest.Modules)
	}

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

	if connector.ConnectorRequest != nil {
		return connector.ConnectorRequest.Name, "storage"
	}

	return defaultName, ""
}

func transformBehaviorsRequest(behaviors []contracts.ManifestRuleBehavior) ([]edgesdk.RequestPhaseBehaviorRequest, error) {
	behaviorsRequest := make([]edgesdk.RequestPhaseBehaviorRequest, 0, len(behaviors))
	for _, behavior := range behaviors {
		var argsString edgesdk.BehaviorString
		var argsInt edgesdk.BehaviorInteger
		var withoutArgs edgesdk.BehaviorNoArgs
		var captureMatchGroups edgesdk.BehaviorCapture
		var beh edgesdk.RequestPhaseBehaviorRequest
		switch behavior.Type {
		case "run_function":
			attributesJSON, err := json.Marshal(behavior.Attributes)
			if err != nil {
				return nil, err
			}
			var attributes edgesdk.BehaviorStringAttributes
			var attributesInt edgesdk.BehaviorIntegerAttributes
			if err := json.Unmarshal(attributesJSON, &attributes); err != nil {
				return nil, err
			}
			argsInt.SetType("run_function")
			if attributes.Value != "" {
				funcName := attributes.Value
				if _, ok := FunctionIds[funcName]; !ok {
					return nil, msg.ErrorFuncNotFound
				}
				funcToWorkWith := FunctionIds[funcName]
				attributesInt.SetValue(funcToWorkWith.InstanceID)
				argsInt.SetAttributes(attributesInt)
			}
			beh.BehaviorInteger = &argsInt
			behaviorsRequest = append(behaviorsRequest, beh)
		case "set_cache_policy":
			attributesJSON, err := json.Marshal(behavior.Attributes)
			if err != nil {
				return nil, err
			}
			var attributes edgesdk.BehaviorStringAttributes
			var attributesInt edgesdk.BehaviorIntegerAttributes
			if err := json.Unmarshal(attributesJSON, &attributes); err != nil {
				return nil, err
			}
			argsInt.SetType("set_cache_policy")
			if attributes.Value != "" {
				cacheName := attributes.Value
				if _, ok := CacheIdsBackup[cacheName]; !ok {
					return nil, msg.ErrorCacheNotFound
				}
				cacheToWorkWith := CacheIdsBackup[cacheName]
				attributesInt.SetValue(cacheToWorkWith)
				argsInt.SetAttributes(attributesInt)
				delete(CacheIds, cacheName)
			}
			beh.BehaviorInteger = &argsInt
			behaviorsRequest = append(behaviorsRequest, beh)
		case "set_connector":
			attributesJSON, err := json.Marshal(behavior.Attributes)
			if err != nil {
				return nil, err
			}
			var attributes edgesdk.BehaviorStringAttributes
			var attributesInt edgesdk.BehaviorIntegerAttributes
			if err := json.Unmarshal(attributesJSON, &attributes); err != nil {
				return nil, err
			}
			argsInt.SetType("set_connector")
			if attributes.Value != "" {
				connectorName := attributes.Value
				if _, ok := ConnectorIds[connectorName]; !ok {
					return nil, msg.ErrorConnectorNotFound
				}
				connectorToWorkWith := ConnectorIds[connectorName]
				attributesInt.SetValue(connectorToWorkWith)
				argsInt.SetAttributes(attributesInt)
			}
			beh.BehaviorInteger = &argsInt
			behaviorsRequest = append(behaviorsRequest, beh)
		case "capture_match_groups":
			attributesJSON, err := json.Marshal(behavior.Attributes)
			if err != nil {
				return nil, err
			}
			var attributes edgesdk.BehaviorCaptureMatchGroupsAttributes
			if err := json.Unmarshal(attributesJSON, &attributes); err != nil {
				return nil, err
			}
			captureMatchGroups.SetType("capture_match_groups")
			captureMatchGroups.SetAttributes(attributes)
			beh.BehaviorCapture = &captureMatchGroups
			behaviorsRequest = append(behaviorsRequest, beh)
		case "redirect_to_301", "redirect_to_302", "filter_request_cookie", "rewrite_request", "add_request_header", "filter_request_header", "add_request_cookie":
			attributesJSON, err := json.Marshal(behavior.Attributes)
			if err != nil {
				return nil, err
			}
			var attributes edgesdk.BehaviorStringAttributes
			if err := json.Unmarshal(attributesJSON, &attributes); err != nil {
				return nil, err
			}
			argsString.SetType(behavior.Type)
			argsString.SetAttributes(attributes)
			beh.BehaviorString = &argsString
			behaviorsRequest = append(behaviorsRequest, beh)
		default:
			withoutArgs.SetType(behavior.Type)
			beh.BehaviorNoArgs = &withoutArgs
			behaviorsRequest = append(behaviorsRequest, beh)
		}
	}

	return behaviorsRequest, nil
}

func transformBehaviorsResponse(behaviors []contracts.ManifestRuleBehavior) ([]edgesdk.ResponsePhaseBehaviorRequest, error) {
	behaviorsResponse := make([]edgesdk.ResponsePhaseBehaviorRequest, 0, len(behaviors))

	for _, behavior := range behaviors {
		var withoutArgs edgesdk.BehaviorNoArgs
		var captureMatchGroups edgesdk.BehaviorCapture
		var beh edgesdk.ResponsePhaseBehaviorRequest

		switch behavior.Type {
		case "capture_match_groups":
			attributesJSON, err := json.Marshal(behavior.Attributes)
			if err != nil {
				return nil, err
			}
			var attributes edgesdk.BehaviorCaptureMatchGroupsAttributes
			if err := json.Unmarshal(attributesJSON, &attributes); err != nil {
				return nil, err
			}
			captureMatchGroups.SetType("capture_match_groups")
			captureMatchGroups.SetAttributes(attributes)
			beh.BehaviorCapture = &captureMatchGroups
			behaviorsResponse = append(behaviorsResponse, beh)
		case "enable_gzip", "deliver":
			withoutArgs.SetType(behavior.Type)
			beh.BehaviorNoArgs = &withoutArgs
			behaviorsResponse = append(behaviorsResponse, beh)
		default:
			// Everything else is WithArgs - try both string and int
			attributesJSON, err := json.Marshal(behavior.Attributes)
			if err != nil {
				return nil, err
			}

			var argsString edgesdk.BehaviorString
			var argsInt edgesdk.BehaviorInteger

			var attributesString edgesdk.BehaviorStringAttributes
			errString := json.Unmarshal(attributesJSON, &attributesString)

			var attributesInt edgesdk.BehaviorIntegerAttributes
			errInt := json.Unmarshal(attributesJSON, &attributesInt)

			if errString == nil && attributesString.Value != "" {
				argsString.SetType(behavior.Type)
				argsString.SetAttributes(attributesString)
				beh.BehaviorString = &argsString
			} else if errInt == nil {
				argsInt.SetType(behavior.Type)
				argsInt.SetAttributes(attributesInt)
				beh.BehaviorInteger = &argsInt
			} else {
				return nil, errString
			}
			behaviorsResponse = append(behaviorsResponse, beh)
		}
	}

	return behaviorsResponse, nil
}

func transformRuleRequestCreate(rule contracts.ManifestRule) edgesdk.RequestPhaseRuleRequest {
	request := edgesdk.RequestPhaseRuleRequest{}

	request.SetActive(rule.Active)
	if rule.Criteria != nil {
		request.SetCriteria(rule.Criteria)
	}
	request.SetDescription(rule.Description)
	request.SetName(rule.Name)

	return request
}

func transformRuleResponseCreate(rule contracts.ManifestRule) edgesdk.ResponsePhaseRuleRequest {
	request := edgesdk.ResponsePhaseRuleRequest{}

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
		logger.FInfoFlags(f.IOStreams.Out, msg.MessageDeleteResource+"\n", f.Format, f.Out)
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
