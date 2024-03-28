package manifest

import (
	"strconv"

	msg "github.com/aziontech/azion-cli/messages/manifest"
	apiCache "github.com/aziontech/azion-cli/pkg/api/cache_setting"
	apiEdgeApplications "github.com/aziontech/azion-cli/pkg/api/edge_applications"
	apiOrigin "github.com/aziontech/azion-cli/pkg/api/origin"
	"github.com/aziontech/azion-cli/pkg/contracts"
	sdk "github.com/aziontech/azionapi-go-sdk/edgeapplications"
	"github.com/davecgh/go-spew/spew"
)

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
		request.SetCookieNames(cache.CookieNames)
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
		request.SetCookieNames(cache.CookieNames)
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

func makeRuleRequestUpdate(rule contracts.RuleEngine, cacheIds map[string]int64, conf *contracts.AzionApplicationOptions) (*apiEdgeApplications.UpdateRulesEngineRequest, error) {
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

	request.Criteria = rulesEngineCriteria
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
					if id := cacheIds[v.RulesEngineBehaviorString.Target]; id > 0 {
						str := strconv.FormatInt(id, 10)
						behaviorString.SetTarget(str)
					} else {
						return nil, msg.ErrorCacheNotFound
					}
				} else if v.RulesEngineBehaviorString.Name == "run_function" {
					str := strconv.FormatInt(conf.Function.ID, 10)
					behaviorString.SetTarget(str)
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

	request.Behaviors = behaviors

	return request, nil
}

func makeRuleRequestCreate(rule contracts.RuleEngine, cacheIds map[string]int64, conf *contracts.AzionApplicationOptions) (*apiEdgeApplications.CreateRulesEngineRequest, error) {
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

	request.Criteria = rulesEngineCriteria
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
					if id := cacheIds[v.RulesEngineBehaviorString.Target]; id > 0 {
						str := strconv.FormatInt(id, 10)
						behaviorString.SetTarget(str)
					} else {
						return nil, msg.ErrorCacheNotFound
					}
				} else if v.RulesEngineBehaviorString.Name == "run_function" {
					str := strconv.FormatInt(conf.Function.ID, 10)
					behaviorString.SetTarget(str)
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

	request.Behaviors = behaviors

	return request, nil
}

func makeOriginCreateRequest(origin contracts.Origin) *apiOrigin.CreateRequest {
	request := &apiOrigin.CreateRequest{}

	spew.Dump(origin)

	if origin.Bucket != "" {
		request.SetBucket(origin.Bucket)
	}
	if origin.Prefix != "" {
		request.SetPrefix(origin.Prefix)
	}
	if origin.OriginType != "" {
		request.SetOriginType(origin.OriginType)
	}

	return request
}

func makeOriginUpdateRequest(origin contracts.Origin) *apiOrigin.UpdateRequest {
	request := &apiOrigin.UpdateRequest{}

	if origin.Bucket != "" {
		request.SetBucket(origin.Bucket)
	}
	if origin.Prefix != "" {
		request.SetPrefix(origin.Prefix)
	}
	if origin.OriginType != "" {
		request.SetOriginType(origin.OriginType)
	}

	return request
}
