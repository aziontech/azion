package manifest

import (
	"context"
	"fmt"
	"strconv"

	msg "github.com/aziontech/azion-cli/messages/manifest"
	apiCache "github.com/aziontech/azion-cli/pkg/api/cache_setting"
	apiEdgeApplications "github.com/aziontech/azion-cli/pkg/api/edge_applications"
	apiOrigin "github.com/aziontech/azion-cli/pkg/api/origin"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/logger"
	sdk "github.com/aziontech/azionapi-go-sdk/edgeapplications"
	"go.uber.org/zap"
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

func makeRuleRequestUpdate(rule contracts.RuleEngine, cacheIds map[string]int64, conf *contracts.AzionApplicationOptions, originKeys map[string]int64) (*apiEdgeApplications.UpdateRulesEngineRequest, error) {
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
					str := strconv.FormatInt(conf.Function.InstanceID, 10)
					behaviorString.SetTarget(str)
				} else if v.RulesEngineBehaviorString.Name == "set_origin" {
					if id := originKeys[v.RulesEngineBehaviorString.Target]; id > 0 {
						str := strconv.FormatInt(id, 10)
						behaviorString.SetTarget(str)
					} else {
						fmt.Println(v.RulesEngineBehaviorString.Target)
						return nil, msg.ErrorCacheNotFound
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

	request.Behaviors = behaviors

	return request, nil
}

func makeRuleRequestCreate(rule contracts.RuleEngine, cacheIds map[string]int64, conf *contracts.AzionApplicationOptions, originKeys map[string]int64, client *apiEdgeApplications.Client, ctx context.Context) (*apiEdgeApplications.CreateRulesEngineRequest, error) {
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
					if id := originKeys[v.RulesEngineBehaviorString.Target]; id > 0 {
						str := strconv.FormatInt(id, 10)
						behaviorString.SetTarget(str)
					} else {
						fmt.Println(v.RulesEngineBehaviorString.Target)
						return nil, msg.ErrorCacheNotFound
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

	request.Behaviors = behaviors

	return request, nil
}

func makeOriginCreateRequest(origin contracts.Origin, conf *contracts.AzionApplicationOptions) *apiOrigin.CreateRequest {
	request := &apiOrigin.CreateRequest{}

	request.SetBucket(conf.Bucket)

	if origin.Prefix != "" {
		request.SetPrefix(origin.Prefix)
	} else {
		request.SetPrefix(conf.Prefix)
	}
	if origin.OriginType != "" {
		request.SetOriginType(origin.OriginType)
	}

	return request
}

func doCacheForRule(ctx context.Context, client *apiEdgeApplications.Client, conf *contracts.AzionApplicationOptions) (int64, error) {
	var reqCache apiEdgeApplications.CreateCacheSettingsRequest
	reqCache.SetName("function policy")
	reqCache.SetBrowserCacheSettings("honor")
	reqCache.SetCdnCacheSettings("honor")
	reqCache.SetCdnCacheSettingsMaximumTtl(0)
	reqCache.SetCacheByQueryString("all")
	reqCache.SetCacheByCookies("all")

	// create cache to function next
	cache, err := client.CreateCacheEdgeApplication(ctx, &reqCache, conf.Application.ID)
	if err != nil {
		logger.Debug("Error while creating cache settings", zap.Error(err))
		return 0, err
	}
	return cache.GetId(), nil
}

// func makeOriginUpdateRequest(origin contracts.Origin) *apiOrigin.UpdateRequest {
// 	request := &apiOrigin.UpdateRequest{}

// 	if origin.Bucket != "" {
// 		request.SetBucket(origin.Bucket)
// 	}
// 	if origin.Prefix != "" {
// 		request.SetPrefix(origin.Prefix)
// 	}
// 	if origin.OriginType != "" {
// 		request.SetOriginType(origin.OriginType)
// 	}

// 	return request
// }
