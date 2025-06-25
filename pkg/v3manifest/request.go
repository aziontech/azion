package manifest

import (
	"context"
	"fmt"
	"strconv"

	msg "github.com/aziontech/azion-cli/messages/manifest"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/logger"
	apiCache "github.com/aziontech/azion-cli/pkg/v3api/cache_setting"
	apiDomain "github.com/aziontech/azion-cli/pkg/v3api/domain"
	apiEdgeApplications "github.com/aziontech/azion-cli/pkg/v3api/edge_applications"
	apiOrigin "github.com/aziontech/azion-cli/pkg/v3api/origin"
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

func makeRuleRequestUpdate(rule contracts.RuleEngine, conf *contracts.AzionApplicationOptionsV3) (*apiEdgeApplications.UpdateRulesEngineRequest, error) {
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

	request.Behaviors = behaviors

	return request, nil
}

func makeRuleRequestCreate(rule contracts.RuleEngine, conf *contracts.AzionApplicationOptionsV3, client *apiEdgeApplications.Client, ctx context.Context) (*apiEdgeApplications.CreateRulesEngineRequest, error) {
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

	request.Behaviors = behaviors

	return request, nil
}

func makeOriginCreateRequest(origin contracts.Origin, conf *contracts.AzionApplicationOptionsV3) *apiOrigin.CreateRequest {
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

func makeOriginUpdateRequest(origin contracts.Origin, conf *contracts.AzionApplicationOptionsV3) *apiOrigin.UpdateRequest {
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

func doCacheForRule(ctx context.Context, client *apiEdgeApplications.Client, conf *contracts.AzionApplicationOptionsV3) (int64, error) {
	if conf.Function.CacheId > 0 {
		return conf.Function.CacheId, nil
	}
	var reqCache apiEdgeApplications.CreateCacheSettingsRequest
	reqCache.SetName("function-policy")
	reqCache.SetBrowserCacheSettings("honor")
	reqCache.SetCdnCacheSettings("honor")
	reqCache.SetCdnCacheSettingsMaximumTtl(0)
	reqCache.SetCacheByQueryString("all")
	reqCache.SetCacheByCookies("all")

	// create cache to function next
	cache, err := client.CreateCacheEdgeApplication(ctx, &reqCache, conf.Application.ID)
	if err != nil {
		logger.Debug("Error while creating Cache Settings", zap.Error(err))
		return 0, err
	}

	conf.Function.CacheId = cache.GetId()

	return cache.GetId(), nil
}

func makeDomainUpdateRequest(domain *contracts.Domains, conf *contracts.AzionApplicationOptionsV3) *apiDomain.UpdateRequest {
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

func makeDomainCreateRequest(domain *contracts.Domains, conf *contracts.AzionApplicationOptionsV3) *apiDomain.CreateRequest {
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
