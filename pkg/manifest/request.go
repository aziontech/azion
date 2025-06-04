package manifest

import (
	"context"
	"fmt"
	"strconv"

	msg "github.com/aziontech/azion-cli/messages/manifest"
	apiCache "github.com/aziontech/azion-cli/pkg/api/cache_setting"
	apiDomain "github.com/aziontech/azion-cli/pkg/api/domain"
	apiEdgeApplications "github.com/aziontech/azion-cli/pkg/api/edge_applications"
	apiConnector "github.com/aziontech/azion-cli/pkg/api/edge_connector"
	apiOrigin "github.com/aziontech/azion-cli/pkg/api/origin"
	apiWorkloads "github.com/aziontech/azion-cli/pkg/api/workloads"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/logger"
	sdk "github.com/aziontech/azionapi-go-sdk/edgeapplications"
	"github.com/aziontech/azionapi-v4-go-sdk/edge"
	edgesdk "github.com/aziontech/azionapi-v4-go-sdk/edge"
	"go.uber.org/zap"
)

func transformEdgeConnectorRequest(connectorRequest edgesdk.EdgeConnectorPolymorphicRequest) *apiConnector.UpdateRequest {
	request := &apiConnector.UpdateRequest{}

	if connectorRequest.EdgeConnectorHTTPTypedRequest != nil {
		bodyRequest := connectorRequest.EdgeConnectorHTTPTypedRequest
		body := edgesdk.PatchedEdgeConnectorHTTPTypedRequest{}
		if bodyRequest.Active != nil {
			body.SetActive(*bodyRequest.Active)
		}
		if len(bodyRequest.Addresses) > 0 {
			body.SetAddresses(bodyRequest.Addresses)
		}
		if len(bodyRequest.ConnectionPreference) > 0 {
			body.SetConnectionPreference(bodyRequest.ConnectionPreference)
		}
		if bodyRequest.ConnectionTimeout != nil {
			body.SetConnectionTimeout(*bodyRequest.ConnectionTimeout)
		}
		if bodyRequest.LoadBalanceMethod != nil {
			body.SetLoadBalanceMethod(*bodyRequest.LoadBalanceMethod)
		}
		if bodyRequest.MaxRetries != nil {
			body.SetMaxRetries(*bodyRequest.MaxRetries)
		}

		body.SetModules(bodyRequest.Modules)

		if bodyRequest.Name != "" {
			body.SetName(bodyRequest.Name)
		}
		if bodyRequest.ReadWriteTimeout != nil {
			body.SetReadWriteTimeout(*bodyRequest.ReadWriteTimeout)
		}
		if bodyRequest.Tls != nil {
			body.SetTls(*bodyRequest.Tls)
		}
		if bodyRequest.Type != "" {
			body.SetType(bodyRequest.Type)
		}
		body.SetTypeProperties(bodyRequest.TypeProperties)
		request.PatchedEdgeConnectorHTTPTypedRequest = &body
		return request
	}

	if connectorRequest.EdgeConnectorLiveIngestTypedRequest != nil {
		body := edgesdk.PatchedEdgeConnectorLiveIngestTypedRequest{}
		bodyRequest := connectorRequest.EdgeConnectorLiveIngestTypedRequest
		if bodyRequest.Active != nil {
			body.SetActive(*bodyRequest.Active)
		}
		if len(bodyRequest.Addresses) > 0 {
			body.SetAddresses(bodyRequest.Addresses)
		}
		if len(bodyRequest.ConnectionPreference) > 0 {
			body.SetConnectionPreference(bodyRequest.ConnectionPreference)
		}
		if bodyRequest.ConnectionTimeout != nil {
			body.SetConnectionTimeout(*bodyRequest.ConnectionTimeout)
		}
		if bodyRequest.LoadBalanceMethod != nil {
			body.SetLoadBalanceMethod(*bodyRequest.LoadBalanceMethod)
		}
		if bodyRequest.MaxRetries != nil {
			body.SetMaxRetries(*bodyRequest.MaxRetries)
		}

		body.SetModules(bodyRequest.Modules)

		if bodyRequest.Name != "" {
			body.SetName(bodyRequest.Name)
		}
		if bodyRequest.ReadWriteTimeout != nil {
			body.SetReadWriteTimeout(*bodyRequest.ReadWriteTimeout)
		}
		if bodyRequest.Tls != nil {
			body.SetTls(*bodyRequest.Tls)
		}
		if bodyRequest.Type != "" {
			body.SetType(bodyRequest.Type)
		}
		body.SetTypeProperties(bodyRequest.TypeProperties)
		request.PatchedEdgeConnectorLiveIngestTypedRequest = &body
		return request
	}

	if connectorRequest.EdgeConnectorS3TypedRequest != nil {
		body := edgesdk.PatchedEdgeConnectorS3TypedRequest{}
		bodyRequest := edgesdk.PatchedEdgeConnectorS3TypedRequest{}
		if bodyRequest.Active != nil {
			body.SetActive(*bodyRequest.Active)
		}
		if len(bodyRequest.Addresses) > 0 {
			body.SetAddresses(bodyRequest.Addresses)
		}
		if len(bodyRequest.ConnectionPreference) > 0 {
			body.SetConnectionPreference(bodyRequest.ConnectionPreference)
		}
		if bodyRequest.ConnectionTimeout != nil {
			body.SetConnectionTimeout(*bodyRequest.ConnectionTimeout)
		}
		if bodyRequest.LoadBalanceMethod != nil {
			body.SetLoadBalanceMethod(*bodyRequest.LoadBalanceMethod)
		}
		if bodyRequest.MaxRetries != nil {
			body.SetMaxRetries(*bodyRequest.MaxRetries)
		}
		if bodyRequest.Modules != nil {
			body.SetModules(*bodyRequest.Modules)
		}
		if bodyRequest.Name != nil {
			body.SetName(*bodyRequest.Name)
		}
		if bodyRequest.ReadWriteTimeout != nil {
			body.SetReadWriteTimeout(*bodyRequest.ReadWriteTimeout)
		}
		if bodyRequest.Tls != nil {
			body.SetTls(*bodyRequest.Tls)
		}
		if bodyRequest.Type != nil {
			body.SetType(*bodyRequest.Type)
		}
		if bodyRequest.TypeProperties != nil {
			body.SetTypeProperties(*bodyRequest.TypeProperties)
		}

		request.PatchedEdgeConnectorS3TypedRequest = &body
		return request
	}

	if connectorRequest.EdgeConnectorStorageTypedRequest != nil {
		body := edgesdk.PatchedEdgeConnectorStorageTypedRequest{}
		bodyRequest := edgesdk.PatchedEdgeConnectorStorageTypedRequest{}
		if bodyRequest.Active != nil {
			body.SetActive(*bodyRequest.Active)
		}
		if len(bodyRequest.Addresses) > 0 {
			body.SetAddresses(bodyRequest.Addresses)
		}
		if len(bodyRequest.ConnectionPreference) > 0 {
			body.SetConnectionPreference(bodyRequest.ConnectionPreference)
		}
		if bodyRequest.ConnectionTimeout != nil {
			body.SetConnectionTimeout(*bodyRequest.ConnectionTimeout)
		}
		if bodyRequest.LoadBalanceMethod != nil {
			body.SetLoadBalanceMethod(*bodyRequest.LoadBalanceMethod)
		}
		if bodyRequest.MaxRetries != nil {
			body.SetMaxRetries(*bodyRequest.MaxRetries)
		}
		if bodyRequest.Modules != nil {
			body.SetModules(*bodyRequest.Modules)
		}
		if bodyRequest.Name != nil {
			body.SetName(*bodyRequest.Name)
		}
		if bodyRequest.ReadWriteTimeout != nil {
			body.SetReadWriteTimeout(*bodyRequest.ReadWriteTimeout)
		}
		if bodyRequest.Tls != nil {
			body.SetTls(*bodyRequest.Tls)
		}
		if bodyRequest.Type != nil {
			body.SetType(*bodyRequest.Type)
		}
		if bodyRequest.TypeProperties != nil {
			body.SetTypeProperties(*bodyRequest.TypeProperties)
		}

		request.PatchedEdgeConnectorStorageTypedRequest = &body
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
	if len(createRequest.AlternateDomains) > 0 {
		request.SetAlternateDomains(createRequest.AlternateDomains)
	}
	if len(createRequest.Domains) > 0 {
		request.SetDomains(createRequest.Domains)
	}
	if createRequest.Mtls != nil {
		request.SetMtls(*createRequest.Mtls)
	}
	if createRequest.NetworkMap != nil {
		request.SetNetworkMap(*createRequest.NetworkMap)
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

	request.SetEdgeApplication(appid)

	if createRequest.Name != "" {
		request.SetName(createRequest.Name)
	}
	if createRequest.Active != nil {
		request.SetActive(*createRequest.Active)
	}
	if len(createRequest.AlternateDomains) > 0 {
		request.SetAlternateDomains(createRequest.AlternateDomains)
	}
	if len(createRequest.Domains) > 0 {
		request.SetDomains(createRequest.Domains)
	}
	if createRequest.Mtls != nil {
		request.SetMtls(*createRequest.Mtls)
	}
	if createRequest.NetworkMap != nil {
		request.SetNetworkMap(*createRequest.NetworkMap)
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
		modules := edgesdk.EdgeApplicationModulesRequest{}
		if edgeapprequest.Modules.ApplicationAcceleratorEnabled != nil {
			modules.SetApplicationAcceleratorEnabled(*edgeapprequest.Modules.ApplicationAcceleratorEnabled)
		}
		if edgeapprequest.Modules.EdgeCacheEnabled != nil {
			modules.SetEdgeCacheEnabled(*edgeapprequest.Modules.EdgeCacheEnabled)
		}
		if edgeapprequest.Modules.EdgeFunctionsEnabled != nil {
			modules.SetEdgeFunctionsEnabled(*edgeapprequest.Modules.EdgeFunctionsEnabled)
		}
		if edgeapprequest.Modules.ImageProcessorEnabled != nil {
			modules.SetImageProcessorEnabled(*edgeapprequest.Modules.ImageProcessorEnabled)
		}
		if edgeapprequest.Modules.TieredCacheEnabled != nil {
			modules.SetTieredCacheEnabled(*edgeapprequest.Modules.TieredCacheEnabled)
		}

		request.SetModules(modules)
	}
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
	if edgeapprequest.Modules != nil {
		modules := edgesdk.EdgeApplicationModulesRequest{}
		if edgeapprequest.Modules.ApplicationAcceleratorEnabled != nil {
			modules.SetApplicationAcceleratorEnabled(*edgeapprequest.Modules.ApplicationAcceleratorEnabled)
		}
		if edgeapprequest.Modules.EdgeCacheEnabled != nil {
			modules.SetEdgeCacheEnabled(*edgeapprequest.Modules.EdgeCacheEnabled)
		}
		if edgeapprequest.Modules.EdgeFunctionsEnabled != nil {
			modules.SetEdgeFunctionsEnabled(*edgeapprequest.Modules.EdgeFunctionsEnabled)
		}
		if edgeapprequest.Modules.ImageProcessorEnabled != nil {
			modules.SetImageProcessorEnabled(*edgeapprequest.Modules.ImageProcessorEnabled)
		}
		if edgeapprequest.Modules.TieredCacheEnabled != nil {
			modules.SetTieredCacheEnabled(*edgeapprequest.Modules.TieredCacheEnabled)
		}

		request.SetModules(modules)
	}
	if edgeapprequest.Name != "" {
		request.SetName(edgeapprequest.Name)
	}

	return request
}

func transformCacheRequest(cache edgesdk.CacheSettingRequest) *apiCache.RequestUpdate {
	request := apiCache.RequestUpdate{}

	if cache.Name != "" {
		request.SetName(cache.Name)
	}
	request.SetApplicationControls(cache.ApplicationControls)
	request.SetBrowserCache(cache.BrowserCache)
	request.SetEdgeCache(cache.EdgeCache)
	request.SetSliceControls(cache.SliceControls)

	return &request
}

func transformRuleRequest(rule edgesdk.EdgeApplicationRuleEngineRequest) *apiEdgeApplications.UpdateRulesEngineRequest {
	request := &apiEdgeApplications.UpdateRulesEngineRequest{}

	if rule.Active != nil {
		request.SetActive(*rule.Active)
	}
	if rule.Behaviors != nil {
		request.SetBehaviors(rule.Behaviors)
	}
	if rule.Criteria != nil {
		request.SetCriteria(rule.Criteria)
	}
	if rule.Description != nil {
		request.SetDescription(*rule.Description)
	}
	if rule.Name != "" {
		request.SetName(rule.Name)
	}
	if rule.Phase != "" {
		request.SetPhase(rule.Phase)
	}

	return request
}

func makeCacheRequestUpdate(cache contracts.CacheSetting) *apiCache.UpdateRequest {
	request := &apiCache.UpdateRequest{}
	if cache.AdaptiveDeliveryAction != nil {
		request.SetAdaptiveDeliveryAction(*cache.AdaptiveDeliveryAction)
	}
	if cache.BrowserCacheSettings != nil {
		request.SetBrowserCacheSettings(*cache.BrowserCacheSettings)
	}
	if cache.BrowserCacheSettingsMaximumTtl != nil {
		request.SetBrowserCacheSettingsMaximumTtl(*cache.BrowserCacheSettingsMaximumTtl)
	}
	if cache.CacheByCookies != nil {
		request.SetCacheByCookies(*cache.CacheByCookies)
	}
	if cache.CacheByQueryString != nil {
		request.SetCacheByQueryString(*cache.CacheByQueryString)
	}
	if cache.CdnCacheSettings != nil {
		request.SetCdnCacheSettings(*cache.CdnCacheSettings)
	}
	if cache.CdnCacheSettingsMaximumTtl != nil {
		request.SetCdnCacheSettingsMaximumTtl(*cache.CdnCacheSettingsMaximumTtl)
	}
	if cache.CookieNames != nil {
		var cookieNames []string
		for _, name := range cache.CookieNames {
			if name != "" {
				cookieNames = append(cookieNames, name)
			}
		}
		request.SetCookieNames(cookieNames)
	}
	if cache.EnableCachingForOptions != nil {
		request.SetEnableCachingForOptions(*cache.EnableCachingForOptions)
	}
	if cache.EnableCachingForPost != nil {
		request.SetEnableCachingForPost(*cache.EnableCachingForPost)
	}
	if cache.EnableQueryStringSort != nil {
		request.SetEnableQueryStringSort(*cache.EnableQueryStringSort)
	}
	if cache.IsSliceConfigurationEnabled != nil {
		request.SetIsSliceConfigurationEnabled(*cache.IsSliceConfigurationEnabled)
	}
	if cache.IsSliceEdgeCachingEnabled != nil {
		request.SetIsSliceEdgeCachingEnabled(*cache.IsSliceEdgeCachingEnabled)
	}
	if cache.IsSliceL2CachingEnabled != nil {
		request.SetIsSliceL2CachingEnabled(*cache.IsSliceL2CachingEnabled)
	}
	if cache.L2CachingEnabled != nil {
		request.SetL2CachingEnabled(*cache.L2CachingEnabled)
	}
	if cache.Name != nil {
		request.SetName(*cache.Name)
	}
	if cache.QueryStringFields != nil {
		request.SetQueryStringFields(cache.QueryStringFields)
	}
	if cache.SliceConfigurationRange != nil {
		request.SetSliceConfigurationRange(*cache.SliceConfigurationRange)
	}

	return request

}

func makeCacheRequestCreate(cache contracts.CacheSetting) *apiCache.CreateRequest {
	request := &apiCache.CreateRequest{}
	if cache.AdaptiveDeliveryAction != nil {
		request.SetAdaptiveDeliveryAction(*cache.AdaptiveDeliveryAction)
	}
	if cache.BrowserCacheSettings != nil {
		request.SetBrowserCacheSettings(*cache.BrowserCacheSettings)
	}
	if cache.BrowserCacheSettingsMaximumTtl != nil {
		request.SetBrowserCacheSettingsMaximumTtl(*cache.BrowserCacheSettingsMaximumTtl)
	}
	if cache.CacheByCookies != nil {
		request.SetCacheByCookies(*cache.CacheByCookies)
	}
	if cache.CacheByQueryString != nil {
		request.SetCacheByQueryString(*cache.CacheByQueryString)
	}
	if cache.CdnCacheSettings != nil {
		request.SetCdnCacheSettings(*cache.CdnCacheSettings)
	}
	if cache.CdnCacheSettingsMaximumTtl != nil {
		request.SetCdnCacheSettingsMaximumTtl(*cache.CdnCacheSettingsMaximumTtl)
	}
	if cache.CookieNames != nil {
		var cookieNames []string
		for _, name := range cache.CookieNames {
			if name != "" {
				cookieNames = append(cookieNames, name)
			}
		}
		request.SetCookieNames(cookieNames)
	}
	if cache.DeviceGroup != nil {
		request.SetDeviceGroup(cache.DeviceGroup)
	}
	if cache.EnableCachingForOptions != nil {
		request.SetEnableCachingForOptions(*cache.EnableCachingForOptions)
	}
	if cache.EnableCachingForPost != nil {
		request.SetEnableCachingForPost(*cache.EnableCachingForPost)
	}
	if cache.EnableQueryStringSort != nil {
		request.SetEnableQueryStringSort(*cache.EnableQueryStringSort)
	}
	if cache.EnableStaleCache != nil {
		request.SetEnableStaleCache(*cache.EnableStaleCache)
	}
	if cache.IsSliceConfigurationEnabled != nil {
		request.SetIsSliceConfigurationEnabled(*cache.IsSliceConfigurationEnabled)
	}
	if cache.IsSliceEdgeCachingEnabled != nil {
		request.SetIsSliceEdgeCachingEnabled(*cache.IsSliceEdgeCachingEnabled)
	}
	if cache.IsSliceL2CachingEnabled != nil {
		request.SetIsSliceL2CachingEnabled(*cache.IsSliceL2CachingEnabled)
	}
	if cache.L2CachingEnabled != nil {
		request.SetL2CachingEnabled(*cache.L2CachingEnabled)
	}
	if cache.L2Region != nil {
		request.SetL2Region(*cache.L2Region)
	}
	if cache.Name != nil {
		request.SetName(*cache.Name)
	}
	if cache.QueryStringFields != nil {
		request.SetQueryStringFields(cache.QueryStringFields)
	}
	if cache.SliceConfigurationRange != nil {
		request.SetSliceConfigurationRange(*cache.SliceConfigurationRange)
	}

	return request
}

func makeRuleRequestUpdate(rule contracts.RuleEngine, conf *contracts.AzionApplicationOptions) (*apiEdgeApplications.UpdateRulesEngineRequest, error) {
	request := &apiEdgeApplications.UpdateRulesEngineRequest{}

	if rule.Description != nil {
		request.SetDescription(*rule.Description)
	}

	var rulesEngineCriteria [][]sdk.RulesEngineCriteria
	for _, itemCriterias := range rule.Criteria {
		var criterias []sdk.RulesEngineCriteria
		for _, itemCriteria := range itemCriterias {
			var criteria sdk.RulesEngineCriteria

			criteria.Conditional = itemCriteria.Conditional
			criteria.Variable = itemCriteria.Variable
			criteria.Operator = itemCriteria.Operator
			criteria.InputValue = itemCriteria.InputValue

			criterias = append(criterias, criteria)
		}
		rulesEngineCriteria = append(rulesEngineCriteria, criterias)
	}

	// request.Criteria = rulesEngineCriteria //TODO: correct criteria
	var behaviors []sdk.RulesEngineBehaviorEntry
	for _, v := range rule.Behaviors {
		if v.RulesEngineBehaviorObject != nil {
			if v.RulesEngineBehaviorObject.Target.CapturedArray != nil && v.RulesEngineBehaviorObject.Target.Regex != nil && v.RulesEngineBehaviorObject.Target.Subject != nil {
				var behaviorObject sdk.RulesEngineBehaviorObject
				behaviorObject.SetName(v.RulesEngineBehaviorObject.Name)
				behaviorObject.SetTarget(v.RulesEngineBehaviorObject.Target)
				behaviors = append(behaviors, sdk.RulesEngineBehaviorEntry{
					RulesEngineBehaviorObject: &behaviorObject,
				})
			} else {
				var behaviorString sdk.RulesEngineBehaviorString
				behaviorString.SetName(v.RulesEngineBehaviorObject.Name)
				behaviors = append(behaviors, sdk.RulesEngineBehaviorEntry{
					RulesEngineBehaviorString: &behaviorString,
				})
			}
		} else {
			if v.RulesEngineBehaviorString != nil {
				var behaviorString sdk.RulesEngineBehaviorString
				if v.RulesEngineBehaviorString.Name == "set_cache_policy" {
					if id := CacheIdsBackup[v.RulesEngineBehaviorString.Target]; id > 0 {
						str := strconv.FormatInt(id, 10)
						behaviorString.SetTarget(str)
						delete(CacheIds, v.RulesEngineBehaviorString.Target)
					} else {
						logger.Debug("Cache Setting not found", zap.Any("Target", v.RulesEngineBehaviorString.Target))
						return nil, msg.ErrorCacheNotFound
					}
				} else if v.RulesEngineBehaviorString.Name == "run_function" {
					str := strconv.FormatInt(conf.Function.InstanceID, 10)
					behaviorString.SetTarget(str)
				} else if v.RulesEngineBehaviorString.Name == "set_origin" {
					if id := OriginIds[v.RulesEngineBehaviorString.Target]; id > 0 {
						str := strconv.FormatInt(id, 10)
						behaviorString.SetTarget(str)
						delete(OriginKeys, v.RulesEngineBehaviorString.Target)
					} else {
						logger.Debug("Origin not found", zap.Any("Target", v.RulesEngineBehaviorString.Target))
						return nil, msg.ErrorOriginNotFound
					}
				} else {
					behaviorString.SetTarget(v.RulesEngineBehaviorString.Target)
				}
				behaviorString.SetName(v.RulesEngineBehaviorString.Name)
				behaviors = append(behaviors, sdk.RulesEngineBehaviorEntry{
					RulesEngineBehaviorString: &behaviorString,
				})
			}
		}

	}

	// request.Behaviors = behaviors //TODO: correct behaviors

	return request, nil
}

func makeRuleRequestCreate(rule contracts.RuleEngine, conf *contracts.AzionApplicationOptions, client *apiEdgeApplications.Client, ctx context.Context) (*apiEdgeApplications.CreateRulesEngineRequest, error) {
	request := &apiEdgeApplications.CreateRulesEngineRequest{}

	if rule.Description != nil {
		request.SetDescription(*rule.Description)
	}

	var rulesEngineCriteria [][]sdk.RulesEngineCriteria
	for _, itemCriterias := range rule.Criteria {
		var criterias []sdk.RulesEngineCriteria
		for _, itemCriteria := range itemCriterias {
			var criteria sdk.RulesEngineCriteria

			criteria.Conditional = itemCriteria.Conditional
			criteria.Variable = itemCriteria.Variable
			criteria.Operator = itemCriteria.Operator
			criteria.InputValue = itemCriteria.InputValue

			criterias = append(criterias, criteria)
		}
		rulesEngineCriteria = append(rulesEngineCriteria, criterias)
	}

	// request.Criteria = rulesEngineCriteria //TODO: correct criteria
	var behaviors []sdk.RulesEngineBehaviorEntry

	for _, v := range rule.Behaviors {
		if v.RulesEngineBehaviorObject != nil {
			if v.RulesEngineBehaviorObject.Target.CapturedArray != nil && v.RulesEngineBehaviorObject.Target.Regex != nil && v.RulesEngineBehaviorObject.Target.Subject != nil {
				var behaviorObject sdk.RulesEngineBehaviorObject
				behaviorObject.SetName(v.RulesEngineBehaviorObject.Name)
				behaviorObject.SetTarget(v.RulesEngineBehaviorObject.Target)
				behaviors = append(behaviors, sdk.RulesEngineBehaviorEntry{
					RulesEngineBehaviorObject: &behaviorObject,
				})
			} else {
				var behaviorString sdk.RulesEngineBehaviorString
				behaviorString.SetName(v.RulesEngineBehaviorObject.Name)
				behaviors = append(behaviors, sdk.RulesEngineBehaviorEntry{
					RulesEngineBehaviorString: &behaviorString,
				})
			}
		} else {
			if v.RulesEngineBehaviorString != nil {
				var behaviorString sdk.RulesEngineBehaviorString
				if v.RulesEngineBehaviorString.Name == "set_cache_policy" {
					if id := CacheIdsBackup[v.RulesEngineBehaviorString.Target]; id > 0 {
						str := strconv.FormatInt(id, 10)
						behaviorString.SetTarget(str)
						delete(CacheIds, v.RulesEngineBehaviorString.Target)
					} else {
						logger.Debug("Cache Setting not found", zap.Any("Target", v.RulesEngineBehaviorString.Target))
						return nil, msg.ErrorCacheNotFound
					}
				} else if v.RulesEngineBehaviorString.Name == "run_function" {
					var beh sdk.RulesEngineBehaviorString
					cacheId, err := doCacheForRule(ctx, client, conf)
					if err != nil {
						return nil, err
					}
					beh.SetName("set_cache_policy")
					beh.SetTarget(fmt.Sprintf("%d", cacheId))
					behaviors = append(behaviors, sdk.RulesEngineBehaviorEntry{
						RulesEngineBehaviorString: &beh,
					})
					str := strconv.FormatInt(conf.Function.InstanceID, 10)
					behaviorString.SetTarget(str)
				} else if v.RulesEngineBehaviorString.Name == "set_origin" {
					if id := OriginIds[v.RulesEngineBehaviorString.Target]; id > 0 {
						str := strconv.FormatInt(id, 10)
						behaviorString.SetTarget(str)
						delete(OriginKeys, v.RulesEngineBehaviorString.Target)
					} else {
						logger.Debug("Origin not found", zap.Any("Target", v.RulesEngineBehaviorString.Target))
						return nil, msg.ErrorOriginNotFound
					}
				} else {
					behaviorString.SetTarget(v.RulesEngineBehaviorString.Target)
				}
				behaviorString.SetName(v.RulesEngineBehaviorString.Name)
				behaviors = append(behaviors, sdk.RulesEngineBehaviorEntry{
					RulesEngineBehaviorString: &behaviorString,
				})
			}
		}

	}

	// request.Behaviors = behaviors //TODO: correct behaviors

	return request, nil
}

func makeOriginCreateRequest(origin contracts.Origin, conf *contracts.AzionApplicationOptions) *apiOrigin.CreateRequest {
	request := &apiOrigin.CreateRequest{}

	switch origin.OriginType {
	case "single_origin":
		if origin.HostHeader != "" {
			request.SetHostHeader(origin.HostHeader)
		}
		if len(origin.Addresses) > 0 {
			addresses := make([]sdk.CreateOriginsRequestAddresses, len(origin.Addresses))
			for i, item := range origin.Addresses {
				addresses[i].Address = item.Address
			}
			request.SetAddresses(addresses)
		}
		if origin.HmacAccessKey != nil {
			request.SetHmacAccessKey(*origin.HmacSecretKey)
		}
		if origin.HmacAuthentication != nil {
			request.SetHmacAuthentication(*origin.HmacAuthentication)
		}
		if origin.HmacRegionName != nil {
			request.SetHmacRegionName(*origin.HmacRegionName)
		}
		if origin.HmacSecretKey != nil {
			request.SetHmacSecretKey(*origin.HmacSecretKey)
		}
		if origin.OriginPath != nil {
			request.SetOriginPath(*origin.OriginPath)
		}
		if origin.OriginProtocolPolicy != nil {
			request.SetOriginProtocolPolicy(*origin.OriginProtocolPolicy)
		}

	case "object_storage":
		request.SetBucket(conf.Bucket)

		if origin.Prefix != "" {
			request.SetPrefix(origin.Prefix)
		} else {
			request.SetPrefix(conf.Prefix)
		}
	}

	if origin.OriginType != "" {
		request.SetOriginType(origin.OriginType)
	}

	return request
}

func makeOriginUpdateRequest(origin contracts.Origin, conf *contracts.AzionApplicationOptions) *apiOrigin.UpdateRequest {
	request := &apiOrigin.UpdateRequest{}

	switch origin.OriginType {
	case "single_origin":
		if origin.HostHeader != "" {
			request.SetHostHeader(origin.HostHeader)
		}
		if len(origin.Addresses) > 0 {
			addresses := make([]sdk.CreateOriginsRequestAddresses, len(origin.Addresses))
			for i, item := range origin.Addresses {
				addresses[i].Address = item.Address
			}
			request.SetAddresses(addresses)
		}
		if origin.HmacAccessKey != nil {
			request.SetHmacAccessKey(*origin.HmacSecretKey)
		}
		if origin.HmacAuthentication != nil {
			request.SetHmacAuthentication(*origin.HmacAuthentication)
		}
		if origin.HmacRegionName != nil {
			request.SetHmacRegionName(*origin.HmacRegionName)
		}
		if origin.HmacSecretKey != nil {
			request.SetHmacSecretKey(*origin.HmacSecretKey)
		}
		if origin.OriginPath != nil {
			request.SetOriginPath(*origin.OriginPath)
		}
		if origin.OriginProtocolPolicy != nil {
			request.SetOriginProtocolPolicy(*origin.OriginProtocolPolicy)
		}

	case "object_storage":
		request.SetBucket(conf.Bucket)

		if origin.Prefix != "" {
			request.SetPrefix(origin.Prefix)
		} else {
			request.SetPrefix(conf.Prefix)
		}
	}

	if origin.OriginType != "" {
		request.SetOriginType(origin.OriginType)
	}

	return request
}

func doCacheForRule(ctx context.Context, client *apiEdgeApplications.Client, conf *contracts.AzionApplicationOptions) (int64, error) {
	if conf.Function.CacheId > 0 {
		return conf.Function.CacheId, nil
	}
	var reqCache apiEdgeApplications.CreateCacheSettingsRequest
	reqCache.SetName("function-policy")
	// reqCache.SetBrowserCacheSettings("honor")
	// reqCache.SetCdnCacheSettings("honor")
	// reqCache.SetCdnCacheSettingsMaximumTtl(0)
	// reqCache.SetCacheByQueryString("all")
	// reqCache.SetCacheByCookies("all")

	// create cache to function next
	str := strconv.FormatInt(conf.Application.ID, 10)
	cache, err := client.CreateCacheEdgeApplication(ctx, &reqCache, str)
	if err != nil {
		logger.Debug("Error while creating Cache Settings", zap.Error(err))
		return 0, err
	}

	conf.Function.CacheId = cache.GetId()

	return cache.GetId(), nil
}

func makeDomainUpdateRequest(domain *contracts.Domains, conf *contracts.AzionApplicationOptions) *apiDomain.UpdateRequest {
	request := &apiDomain.UpdateRequest{}

	if domain.CnameAccessOnly != nil {
		request.SetCnameAccessOnly(*domain.CnameAccessOnly)
	}
	if len(domain.CrlList) > 0 {
		request.SetCrlList(domain.CrlList)
	}
	if domain.DigitalCertificateId != nil {
		request.SetDigitalCertificateId(*domain.DigitalCertificateId)
	}
	if domain.EdgeApplicationId > 0 {
		request.SetEdgeApplicationId(domain.EdgeApplicationId)
	} else {
		request.SetEdgeApplicationId(conf.Application.ID)
	}
	if domain.EdgeFirewallId > 0 {
		request.SetEdgeFirewallId(domain.EdgeApplicationId)
	}
	if domain.IsActive != nil {
		request.SetIsActive(*domain.IsActive)
	}
	if domain.IsMtlsEnabled != nil {
		request.SetIsMtlsEnabled(*domain.IsMtlsEnabled)
	}
	if domain.MtlsTrustedCaCertificateId > 0 {
		request.SetMtlsTrustedCaCertificateId(domain.MtlsTrustedCaCertificateId)
	}
	if domain.MtlsVerification != nil {
		request.SetMtlsVerification(*domain.MtlsVerification)
	}

	request.Name = &domain.Name
	request.Id = conf.Domain.Id
	request.SetCnames(domain.Cnames)

	return request
}

func makeDomainCreateRequest(domain *contracts.Domains, conf *contracts.AzionApplicationOptions) *apiDomain.CreateRequest {
	request := &apiDomain.CreateRequest{}

	if domain.CnameAccessOnly != nil {
		request.SetCnameAccessOnly(*domain.CnameAccessOnly)
	}
	if len(domain.CrlList) > 0 {
		request.SetCrlList(domain.CrlList)
	}
	if domain.DigitalCertificateId != nil {
		request.SetDigitalCertificateId(*domain.DigitalCertificateId)
	}
	if domain.EdgeApplicationId > 0 {
		request.SetEdgeApplicationId(domain.EdgeApplicationId)
	} else {
		request.SetEdgeApplicationId(conf.Application.ID)
	}
	if domain.EdgeFirewallId > 0 {
		request.SetEdgeFirewallId(domain.EdgeApplicationId)
	}
	if domain.IsActive != nil {
		request.SetIsActive(*domain.IsActive)
	}
	if domain.IsMtlsEnabled != nil {
		request.SetIsMtlsEnabled(*domain.IsMtlsEnabled)
	}
	if domain.MtlsTrustedCaCertificateId > 0 {
		request.SetMtlsTrustedCaCertificateId(domain.MtlsTrustedCaCertificateId)
	}
	if domain.MtlsVerification != nil {
		request.SetMtlsVerification(*domain.MtlsVerification)
	}

	request.Name = domain.Name
	request.SetCnames(domain.Cnames)

	return request
}

func getConnectorName(connector edge.EdgeConnectorPolymorphicRequest, defaultName string) string {
	if connector.EdgeConnectorHTTPTypedRequest != nil {
		return connector.EdgeConnectorHTTPTypedRequest.Name
	}

	if connector.EdgeConnectorLiveIngestTypedRequest != nil {
		return connector.EdgeConnectorLiveIngestTypedRequest.Name
	}

	if connector.EdgeConnectorS3TypedRequest != nil {
		return connector.EdgeConnectorS3TypedRequest.Name
	}

	if connector.EdgeConnectorStorageTypedRequest != nil {
		return connector.EdgeConnectorStorageTypedRequest.Name
	}
	return defaultName
}
