package manifest

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	msg "github.com/aziontech/azion-cli/messages/manifest"
	apiCache "github.com/aziontech/azion-cli/pkg/api/cache_setting"
	apiEdgeApplications "github.com/aziontech/azion-cli/pkg/api/edge_applications"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/utils"
	sdk "github.com/aziontech/azionapi-go-sdk/edgeapplications"
	"go.uber.org/zap"
)

var manifestFilePath = "/.edge/manifest.json"

type ManifestInterpreter struct {
	FileReader func(path string) ([]byte, error)
	GetWorkDir func() (string, error)
}

func NewManifestInterpreter() *ManifestInterpreter {
	return &ManifestInterpreter{
		FileReader: os.ReadFile,
		GetWorkDir: utils.GetWorkingDir,
	}
}

func (man *ManifestInterpreter) ManifestPath() (string, error) {
	pathWorkingDir, err := man.GetWorkDir()
	if err != nil {
		return "", err
	}

	return utils.Concat(pathWorkingDir, manifestFilePath), nil
}

func (man *ManifestInterpreter) ReadManifest(path string) (*contracts.Manifest, error) {
	manifest := &contracts.Manifest{}

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

func (man *ManifestInterpreter) CreateResources(path string, manifest *contracts.Manifest, f *cmdutil.Factory) error {

	fmt.Println("Create Resources entered")

	client := apiEdgeApplications.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))
	clientCache := apiCache.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))
	ctx := context.Background()

	conf := &contracts.AzionApplicationOptions{}
	byteAzionJson, err := os.ReadFile(path)
	if err != nil {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			logger.FInfo(f.IOStreams.Out, msg.CREATING)
		}
		return utils.ErrorUnmarshalAzionJsonFile
	} else {
		err = json.Unmarshal(byteAzionJson, &conf)
		fmt.Println("Unmarshal")
		if err != nil {
			logger.Debug("Error reading unmarshalling azion.json file", zap.Error(err))
			return msg.ErrorUnmarshalAzionJsonFile
		}
	}

	for i, cache := range manifest.CacheSettings {
		fmt.Println("cache ", +i)
		found := false
		for _, cacheConf := range conf.CacheSettings {
			fmt.Println("found cache")
			if cacheConf.Name == *cache.Name && conf.Application.ID != 0 && cacheConf.Id != 0 {
				found = true
				requestUpdate := makeCacheRequestUpdate(cache)
				clientCache.Update(ctx, requestUpdate, conf.Application.ID, cacheConf.Id)
				break
			}
		}
		if !found {
			fmt.Println("not found cache")
			found = true
			requestUpdate := makeCacheRequestCreate(cache)
			created, err := clientCache.Create(ctx, requestUpdate, conf.Application.ID)
			if err != nil {
				return err
			}
			newCache := contracts.AzionJsonDataCacheSettings{
				Id:   created.GetId(),
				Name: created.GetName(),
			}
			conf.CacheSettings = append(conf.CacheSettings, newCache)
		}
		found = false
	}

	for _, rule := range manifest.Rules {
		found := false
		fmt.Println("rules entered")
		for _, ruleConf := range conf.RulesEngine.Rules {
			if ruleConf.Name == *rule.Name && conf.Application.ID != 0 && ruleConf.Id != 0 {
				fmt.Println("rule found")
				found = true
				requestUpdate := makeRuleRequestUpdate(rule, conf)
				client.UpdateRulesEngine(ctx, requestUpdate)
				break
			}
		}
		if !found {
			found = true
			fmt.Println("Rule not found")
			requestCreate := makeRuleRequestCreate(rule, conf)
			created, err := client.CreateRulesEngine(ctx, conf.Application.ID, "request", requestCreate)
			if err != nil {
				return err
			}
			newRule := contracts.AzionJsonDataRules{
				Id:   created.GetId(),
				Name: created.GetName(),
			}
			conf.RulesEngine.Rules = append(conf.RulesEngine.Rules, newRule)
		}
		found = false
	}
	return nil
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

func makeRuleRequestUpdate(rule contracts.RuleEngine, conf *contracts.AzionApplicationOptions) *apiEdgeApplications.UpdateRulesEngineRequest {
	request := &apiEdgeApplications.UpdateRulesEngineRequest{}

	if rule.Name != nil {
		request.SetName(*rule.Name)
	}
	if rule.Description != nil {
		request.SetDescription(*rule.Description)
	}

	// HEHEHEHEHHEHE

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
					for _, c := range conf.CacheSettings {
						if c.Name == *rule.Name {
							str := strconv.FormatInt(c.Id, 10)
							behaviorString.SetTarget(str)
						}
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

	return request
}

func makeRuleRequestCreate(rule contracts.RuleEngine, conf *contracts.AzionApplicationOptions) *apiEdgeApplications.CreateRulesEngineRequest {
	request := &apiEdgeApplications.CreateRulesEngineRequest{}

	if rule.Name != nil {
		request.SetName(*rule.Name)
	}
	if rule.Description != nil {
		request.SetDescription(*rule.Description)
	}

	// HEHEHEHEHHEHE

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
					for _, c := range conf.CacheSettings {
						if c.Name == *rule.Name {
							str := strconv.FormatInt(c.Id, 10)
							behaviorString.SetTarget(str)
						}
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

	return request
}
