package sync

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strconv"

	msg "github.com/aziontech/azion-cli/messages/sync"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/logger"
	edgeApp "github.com/aziontech/azion-cli/pkg/v3api/edge_applications"
	"github.com/aziontech/azion-cli/pkg/v3api/origin"
	varApi "github.com/aziontech/azion-cli/pkg/v3api/variables"
	vulcanPkg "github.com/aziontech/azion-cli/pkg/vulcan"
	"github.com/aziontech/azion-cli/utils"
	"go.uber.org/zap"
)

var (
	opts  *contracts.ListOptions
	ctx   context.Context = context.Background()
	words                 = []string{"PASSWORD", "PWD", "SECRET", "HASH", "ENCRYPTED", "PASSCODE", "AUTH", "TOKEN", "SECRET"}
)

func SyncLocalResources(f *cmdutil.Factory, info contracts.SyncOptsV3, synch *SyncCmd) error {
	opts = &contracts.ListOptions{
		PageSize: 1000,
		Page:     1,
	}

	if info.Conf.Application.ID <= 0 {
		return msg.ERRORNOTDEPLOYED
	}

	var err error
	manifest := &contracts.Manifest{}

	remoteCaches, err := synch.syncCache(info, f, manifest)
	if err != nil {
		return fmt.Errorf(msg.ERRORSYNC, err.Error())
	}

	remoteOrigins, err := synch.syncOrigin(info, f, manifest)
	if err != nil {
		return fmt.Errorf(msg.ERRORSYNC, err.Error())
	}

	err = synch.syncRules(info, f, manifest, remoteCaches, remoteOrigins)
	if err != nil {
		return fmt.Errorf(msg.ERRORSYNC, err.Error())
	}

	err = synch.syncEnv(f)
	if err != nil {
		return fmt.Errorf(msg.ERRORSYNC, err.Error())
	}

	if IaC {
		if IaCFormat != "mjs" && IaCFormat != "cjs" && IaCFormat != "js" && IaCFormat != "ts" {
			return msg.INVALIDFORMAT
		}

		err = synch.WriteManifest(manifest, "")
		if err != nil {
			return err
		}
		defer os.Remove("manifesttoconvert.json")
		fileName := fmt.Sprintf("azion.config.%s", IaCFormat)

		vul := vulcanPkg.NewVulcan()
		command := vul.Command("", "manifest -o %s transform %s", f)
		err = synch.CommandRunInteractive(f, fmt.Sprintf(command, fileName, "manifesttoconvert.json"))
		if err != nil {
			return err
		}
	}

	return nil
}

func (synch *SyncCmd) syncOrigin(info contracts.SyncOptsV3, f *cmdutil.Factory, manifest *contracts.Manifest) (map[string]contracts.AzionJsonDataOrigin, error) {
	remoteOriginIds := make(map[string]contracts.AzionJsonDataOrigin)
	client := origin.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))
	resp, err := client.ListOrigins(ctx, opts, info.Conf.Application.ID)
	if err != nil {
		return remoteOriginIds, err
	}

	originsAzion := []contracts.AzionJsonDataOrigin{}
	info.Conf.Origin = originsAzion

	for _, origin := range resp.Results {
		remoteOriginIds[strconv.FormatInt(*origin.OriginId, 10)] = contracts.AzionJsonDataOrigin{
			OriginId:  *origin.OriginId,
			OriginKey: *origin.OriginKey,
			Name:      origin.Name,
		}
		oEntry := contracts.Origin{
			Name:       origin.GetName(),
			OriginType: origin.GetOriginType(),
		}
		if oEntry.OriginType == "single_origin" {
			for _, o := range origin.Addresses {
				a := contracts.Address{
					Address: o.Address,
				}
				oEntry.Addresses = append(oEntry.Addresses, a)
			}
			oEntry.HmacAccessKey = origin.HmacAccessKey
			oEntry.HmacAuthentication = origin.HmacAuthentication
			oEntry.HmacRegionName = origin.HmacRegionName
			oEntry.HmacSecretKey = origin.HmacSecretKey
			oEntry.HostHeader = *origin.HostHeader
			oEntry.OriginPath = origin.OriginPath
			oEntry.OriginProtocolPolicy = origin.OriginProtocolPolicy
			oEntry.OriginType = *origin.OriginType
		}
		manifest.Origins = append(manifest.Origins, oEntry)
		newOrigin := contracts.AzionJsonDataOrigin{
			OriginId:  origin.GetOriginId(),
			OriginKey: origin.GetOriginKey(),
			Name:      origin.GetName(),
		}
		originsAzion = append(originsAzion, newOrigin)
		info.Conf.Origin = originsAzion
	}
	err = synch.WriteAzionJsonContent(info.Conf, ProjectConf)
	if err != nil {
		logger.Debug("Error while writing azion.json file", zap.Error(err))
		return remoteOriginIds, err
	}
	return remoteOriginIds, nil
}

func (synch *SyncCmd) syncCache(info contracts.SyncOptsV3, f *cmdutil.Factory, manifest *contracts.Manifest) (map[string]contracts.AzionJsonDataCacheSettings, error) {
	remoteCacheIds := make(map[string]contracts.AzionJsonDataCacheSettings)
	client := edgeApp.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))
	resp, err := client.ListCacheEdgeApp(context.Background(), info.Conf.Application.ID)
	if err != nil {
		return remoteCacheIds, err
	}

	cacheAzion := []contracts.AzionJsonDataCacheSettings{}
	info.Conf.CacheSettings = cacheAzion
	for _, cache := range resp {
		remoteCacheIds[strconv.FormatInt(cache.Id, 10)] = contracts.AzionJsonDataCacheSettings{
			Id:   cache.Id,
			Name: cache.Name,
		}
		var cookieNames []string
		for _, cookieName := range cache.CookieNames {
			cookieNames = append(cookieNames, *cookieName)
		}
		cEntry := contracts.CacheSetting{
			Name:                           &cache.Name,
			BrowserCacheSettings:           &cache.BrowserCacheSettings,
			BrowserCacheSettingsMaximumTtl: &cache.BrowserCacheSettingsMaximumTtl,
			CdnCacheSettings:               &cache.CdnCacheSettings,
			CdnCacheSettingsMaximumTtl:     &cache.CdnCacheSettingsMaximumTtl,
			CacheByQueryString:             &cache.CacheByQueryString,
			QueryStringFields:              cache.QueryStringFields,
			EnableQueryStringSort:          &cache.EnableCachingForOptions,
			CacheByCookies:                 &cache.CacheByCookies,
			CookieNames:                    cookieNames,
			AdaptiveDeliveryAction:         &cache.AdaptiveDeliveryAction,
			DeviceGroup:                    cache.DeviceGroup,
			EnableCachingForPost:           &cache.EnableCachingForPost,
			L2CachingEnabled:               &cache.L2CachingEnabled,
			IsSliceConfigurationEnabled:    cache.IsSliceConfigurationEnabled,
			IsSliceEdgeCachingEnabled:      cache.IsSliceEdgeCachingEnabled,
			IsSliceL2CachingEnabled:        cache.IsSliceL2CachingEnabled,
			SliceConfigurationRange:        cache.SliceConfigurationRange,
			EnableCachingForOptions:        &cache.EnableCachingForOptions,
			EnableStaleCache:               &cache.EnableStaleCache,
			L2Region:                       cache.L2Region.Get(),
		}
		manifest.CacheSettings = append(manifest.CacheSettings, cEntry)

		newCache := contracts.AzionJsonDataCacheSettings{
			Id:   cache.GetId(),
			Name: cache.GetName(),
		}
		cacheAzion = append(cacheAzion, newCache)
		info.Conf.CacheSettings = cacheAzion
	}
	err = synch.WriteAzionJsonContent(info.Conf, ProjectConf)
	if err != nil {
		logger.Debug("Error while writing azion.json file", zap.Error(err))
		return remoteCacheIds, err
	}
	return remoteCacheIds, nil
}

func (synch *SyncCmd) syncRules(info contracts.SyncOptsV3, f *cmdutil.Factory, manifest *contracts.Manifest,
	remoteCacheIds map[string]contracts.AzionJsonDataCacheSettings,
	remoteOriginIds map[string]contracts.AzionJsonDataOrigin) error {
	// Get request rules first
	client := edgeApp.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))
	resp, err := client.ListRulesEngine(context.Background(), opts, info.Conf.Application.ID, "request")
	if err != nil {
		return err
	}
	rulesAzion := []contracts.AzionJsonDataRules{}
	info.Conf.RulesEngine.Rules = rulesAzion
	for _, rule := range resp.Results {
		if rule.Name == "Default Rule" {
			//default rule is not added to azion.json or azion.config
			continue
		}
		mEntry := contracts.RuleEngine{
			Name:        rule.GetName(),
			IsActive:    rule.GetIsActive(),
			Order:       rule.GetOrder(),
			Phase:       rule.GetPhase(),
			Description: rule.Description,
			Behaviors:   rule.GetBehaviors(),
			Criteria:    rule.GetCriteria(),
		}

		for i, beh := range mEntry.Behaviors {
			if beh.RulesEngineBehaviorString != nil && beh.RulesEngineBehaviorString.Name == "set_origin" {
				mEntry.Behaviors[i].RulesEngineBehaviorString.Target = remoteOriginIds[beh.RulesEngineBehaviorString.Target].Name
			} else if beh.RulesEngineBehaviorString != nil && beh.RulesEngineBehaviorString.Name == "set_cache_policy" {
				mEntry.Behaviors[i].RulesEngineBehaviorString.Target = remoteCacheIds[beh.RulesEngineBehaviorString.Target].Name
			} else if beh.RulesEngineBehaviorString != nil && beh.RulesEngineBehaviorString.Name == "run_function" {
				mEntry.Behaviors[i].RulesEngineBehaviorString.Target = info.Conf.Function.Name
			}
		}
		manifest.Rules = append(manifest.Rules, mEntry)

		newRule := contracts.AzionJsonDataRules{
			Id:    rule.GetId(),
			Name:  rule.GetName(),
			Phase: rule.GetPhase(),
		}
		rulesAzion = append(rulesAzion, newRule)
		info.Conf.RulesEngine.Rules = rulesAzion
	}
	err = synch.WriteAzionJsonContent(info.Conf, ProjectConf)
	if err != nil {
		logger.Debug("Error while writing azion.json file", zap.Error(err))
		return err
	}

	respResponse, err := client.ListRulesEngine(context.Background(), opts, info.Conf.Application.ID, "response")
	if err != nil {
		return err
	}

	for _, rule := range respResponse.Results {
		if rule.Name == "enable gzip" {
			// we do not add the enable gzip rule created by CLI to azion.json or azion.config
			continue
		}
		mEntry := contracts.RuleEngine{
			Name:        rule.GetName(),
			IsActive:    rule.GetIsActive(),
			Order:       rule.GetOrder(),
			Phase:       rule.GetPhase(),
			Description: rule.Description,
			Behaviors:   rule.GetBehaviors(),
			Criteria:    rule.GetCriteria(),
		}

		for i, beh := range mEntry.Behaviors {
			if beh.RulesEngineBehaviorString != nil && beh.RulesEngineBehaviorString.Name == "set_origin" {
				mEntry.Behaviors[i].RulesEngineBehaviorString.Target = remoteOriginIds[beh.RulesEngineBehaviorString.Target].Name
			} else if beh.RulesEngineBehaviorString != nil && beh.RulesEngineBehaviorString.Name == "set_cache_policy" {
				mEntry.Behaviors[i].RulesEngineBehaviorString.Target = remoteCacheIds[beh.RulesEngineBehaviorString.Target].Name
			} else if beh.RulesEngineBehaviorString != nil && beh.RulesEngineBehaviorString.Name == "run_function" {
				mEntry.Behaviors[i].RulesEngineBehaviorString.Target = info.Conf.Function.Name
			}
		}
		manifest.Rules = append(manifest.Rules, mEntry)

		newRule := contracts.AzionJsonDataRules{
			Id:    rule.GetId(),
			Name:  rule.GetName(),
			Phase: rule.GetPhase(),
		}
		rulesAzion = append(rulesAzion, newRule)
		info.Conf.RulesEngine.Rules = rulesAzion
	}

	err = synch.WriteAzionJsonContent(info.Conf, ProjectConf)
	if err != nil {
		logger.Debug("Error while writing azion.json file", zap.Error(err))
		return err
	}

	return nil
}

func (synch *SyncCmd) syncEnv(f *cmdutil.Factory) error {

	client := varApi.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))
	resp, err := client.List(context.Background())
	if err != nil {
		return err
	}

	// Load the .env file
	envs, err := synch.ReadEnv(synch.EnvPath)
	if err != nil {
		logger.Debug("Error while loading .env file", zap.Error(err))
		return nil // not every project has a .env file... this should not stop the execution
	}

	for _, variable := range resp {
		if v := envs[variable.GetKey()]; v != "" {
			updateRequest := &varApi.Request{}
			updateRequest.SetKey(variable.GetKey())
			updateRequest.SetValue(v)
			updateRequest.Uuid = variable.GetUuid()
			if utils.ContainSubstring(variable.GetKey(), words) {
				logger.FInfo(f.IOStreams.Out, msg.VARIABLESETSECRET)
				updateRequest.SetSecret(true)
			}
			_, err := client.Update(ctx, updateRequest)
			if err != nil {
				return err
			}
			logger.FInfoFlags(synch.Io.Out, fmt.Sprintf(msg.SYNCUPDATEENV, variable.GetKey()), synch.F.Format, synch.F.Out)
			delete(envs, variable.GetKey())
		}
	}

	for key, value := range envs {
		createReq := &varApi.Request{}
		createReq.Key = key
		createReq.Value = value
		if utils.ContainSubstring(key, words) {
			logger.FInfo(f.IOStreams.Out, msg.VARIABLESETSECRET)
			createReq.SetSecret(true)
		}
		_, err := client.Create(ctx, *createReq)
		if err != nil {
			logger.Debug("Error while creating variables during sync process", zap.Error(err))
			return err
		}
		logger.FInfoFlags(synch.Io.Out, fmt.Sprintf(msg.SYNCMESSAGEENV, key), synch.F.Format, synch.F.Out)
	}
	return nil
}

func WriteManifest(manifest *contracts.Manifest, pathMan string) error {

	data, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		logger.Debug("Error marshalling response", zap.Error(err))
		return msg.ERRORMARSHALMANIFEST
	}

	err = os.WriteFile(path.Join(pathMan, "manifesttoconvert.json"), data, 0644)
	if err != nil {
		logger.Debug("Error writing file", zap.Error(err))
		return msg.ERRORWRITEMANIFEST
	}

	return nil
}
