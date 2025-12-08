package sync

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strconv"

	msg "github.com/aziontech/azion-cli/messages/sync"
	edgeApp "github.com/aziontech/azion-cli/pkg/api/applications"
	varApi "github.com/aziontech/azion-cli/pkg/api/variables"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/manifest"
	vulcanPkg "github.com/aziontech/azion-cli/pkg/vulcan"
	"github.com/aziontech/azion-cli/utils"
	edgesdk "github.com/aziontech/azionapi-v4-go-sdk-dev/edge-api"
	"go.uber.org/zap"
)

var (
	opts  *contracts.ListOptions
	ctx   context.Context = context.Background()
	words                 = []string{"PASSWORD", "PWD", "SECRET", "HASH", "ENCRYPTED", "PASSCODE", "AUTH", "TOKEN", "SECRET"}
)

func SyncLocalResources(f *cmdutil.Factory, info contracts.SyncOpts, synch *SyncCmd) error {
	opts = &contracts.ListOptions{
		PageSize: 100,
		Page:     1,
	}

	if info.Conf.Application.ID <= 0 {
		return msg.ERRORNOTDEPLOYED
	}

	var err error
	var manifestStruct *contracts.ManifestV4
	var msgs []string

	interpreter := manifest.NewManifestInterpreter()
	pathManifest, err := interpreter.ManifestPath()
	if err != nil {
		manifestStruct = &contracts.ManifestV4{
			Applications:        []contracts.Applications{},
			Workloads:           []contracts.WorkloadManifest{},
			WorkloadDeployments: []contracts.WorkloadDeployment{},
			Purge:               []contracts.PurgeManifest{},
			Storage:             []contracts.StorageManifest{},
			Functions:           []contracts.Function{},
			Connectors:          []edgesdk.ConnectorPolymorphicRequest{},
		}
	} else {
		manifestStruct, err = interpreter.ReadManifest(pathManifest, f, &msgs)
		if err != nil {
			manifestStruct = &contracts.ManifestV4{
				Applications:        []contracts.Applications{},
				Workloads:           []contracts.WorkloadManifest{},
				WorkloadDeployments: []contracts.WorkloadDeployment{},
				Purge:               []contracts.PurgeManifest{},
				Storage:             []contracts.StorageManifest{},
				Functions:           []contracts.Function{},
				Connectors:          []edgesdk.ConnectorPolymorphicRequest{},
			}
		}
	}

	_, err = synch.syncCache(info, f, manifestStruct)
	if err != nil {
		return fmt.Errorf(msg.ERRORSYNC, err.Error())
	}

	err = synch.syncRules(info, f, manifestStruct)
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

		err = synch.WriteManifest(manifestStruct, "")
		if err != nil {
			return err
		}
		defer os.Remove("manifesttoconvert.json")
		fileName := fmt.Sprintf("azion.config.%s", IaCFormat)

		vul := vulcanPkg.NewVulcan()
		command := vul.Command("", "manifest transform --output %s --entry %s", f)
		err = synch.CommandRunInteractive(f, fmt.Sprintf(command, fileName, "manifesttoconvert.json"))
		if err != nil {
			return err
		}
	}

	return nil
}

func (synch *SyncCmd) syncCache(info contracts.SyncOpts, f *cmdutil.Factory, manifest *contracts.ManifestV4) (map[string]contracts.AzionJsonDataCacheSettings, error) {
	remoteCacheIds := make(map[string]contracts.AzionJsonDataCacheSettings)
	client := edgeApp.NewClient(f.HttpClient, f.Config.GetString("api_v4_url"), f.Config.GetString("token"))
	str := strconv.FormatInt(info.Conf.Application.ID, 10)
	resp, err := client.ListCacheEdgeApp(context.Background(), str, opts)
	if err != nil {
		return remoteCacheIds, err
	}

	if len(manifest.Applications) == 0 {
		appName := info.Conf.Application.Name
		if appName == "" || appName == "__DEFAULT__" {
			appName = info.Conf.Name
		}
		manifest.Applications = append(manifest.Applications, contracts.Applications{
			Name:          appName,
			Rules:         []contracts.ManifestRulesEngine{},
			CacheSettings: []contracts.ManifestCacheSetting{},
		})
	} else if manifest.Applications[0].CacheSettings == nil {
		manifest.Applications[0].CacheSettings = []contracts.ManifestCacheSetting{}
	}

	cacheAzion := info.Conf.CacheSettings
	existingCacheNames := make(map[string]bool)
	for _, existingCache := range cacheAzion {
		existingCacheNames[existingCache.Name] = true
	}

	for _, cache := range resp {
		remoteCacheIds[strconv.FormatInt(cache.Id, 10)] = contracts.AzionJsonDataCacheSettings{
			Id:   cache.Id,
			Name: cache.Name,
		}

		if !existingCacheNames[cache.GetName()] {
			newCache := contracts.AzionJsonDataCacheSettings{
				Id:   cache.GetId(),
				Name: cache.GetName(),
			}
			cacheAzion = append(cacheAzion, newCache)
			existingCacheNames[cache.GetName()] = true

			cacheManifest := contracts.ManifestCacheSetting{
				Name: cache.GetName(),
			}

			browserCache := cache.GetBrowserCache()
			browserCacheBytes, err := json.Marshal(browserCache)
			if err != nil {
				return remoteCacheIds, err
			}
			var browserCacheRequest edgesdk.BrowserCacheModuleRequest
			err = json.Unmarshal(browserCacheBytes, &browserCacheRequest)
			if err != nil {
				return remoteCacheIds, err
			}
			cacheManifest.BrowserCache = &browserCacheRequest

			// Convert Modules from response type to request type via JSON
			modules := cache.GetModules()
			modulesBytes, err := json.Marshal(modules)
			if err != nil {
				return remoteCacheIds, err
			}
			var modulesRequest edgesdk.CacheSettingsModulesRequest
			err = json.Unmarshal(modulesBytes, &modulesRequest)
			if err != nil {
				return remoteCacheIds, err
			}
			cacheManifest.Modules = &modulesRequest

			manifest.Applications[0].CacheSettings = append(manifest.Applications[0].CacheSettings, cacheManifest)
		}
	}

	info.Conf.CacheSettings = cacheAzion
	err = utils.WriteAzionJsonContentPreserveOrder(info.Conf, ProjectConf)
	if err != nil {
		logger.Debug("Error while writing azion.json file", zap.Error(err))
		return remoteCacheIds, err
	}
	return remoteCacheIds, nil
}

func (synch *SyncCmd) syncRules(info contracts.SyncOpts, f *cmdutil.Factory, manifest *contracts.ManifestV4) error {
	client := edgeApp.NewClient(f.HttpClient, f.Config.GetString("api_v4_url"), f.Config.GetString("token"))
	str := strconv.FormatInt(info.Conf.Application.ID, 10)
	rulesAzion := info.Conf.RulesEngine.Rules
	existingRuleNames := make(map[string]bool)
	for _, existingRule := range rulesAzion {
		existingRuleNames[existingRule.Name] = true
	}

	// Initialize Applications slice if it's empty
	if len(manifest.Applications) == 0 {
		appName := info.Conf.Application.Name
		if appName == "" || appName == "__DEFAULT__" {
			appName = info.Conf.Name
		}
		manifest.Applications = append(manifest.Applications, contracts.Applications{
			Name:          appName,
			Rules:         []contracts.ManifestRulesEngine{},
			CacheSettings: []contracts.ManifestCacheSetting{},
		})
	}

	// Get request phase rules
	reqResp, err := client.ListRulesEngineRequest(context.Background(), opts, str)
	if err != nil {
		logger.Debug("Error while listing request phase rules", zap.Error(err))
		return fmt.Errorf("failed to list request phase rules: %w", err)
	}

	for _, rule := range reqResp.Results {
		if rule.Name == "Default Rule" || rule.Name == "enable gzip" {
			//default rule or enable gzip rule are not added to azion.json or azion.config
			continue
		}

		var criteria [][]edgesdk.EdgeApplicationCriterionFieldRequest
		criteriaBytes, err := json.Marshal(rule.Criteria)
		if err != nil {
			return fmt.Errorf(msg.ERRORMARSHALCRITERIA, err)
		}
		err = json.Unmarshal(criteriaBytes, &criteria)
		if err != nil {
			return fmt.Errorf(msg.ERRORUNMARSHALCRITERIA, err)
		}

		var behaviors []contracts.ManifestRuleBehavior
		behaviorsBytes, err := json.Marshal(rule.Behaviors)
		if err != nil {
			return fmt.Errorf(msg.ERRORMARSHALBEHAVIORS, err)
		}
		err = json.Unmarshal(behaviorsBytes, &behaviors)
		if err != nil {
			return fmt.Errorf(msg.ERRORUNMARSHALBEHAVIORS, err)
		}

		// Only add rule if it doesn't already exist locally
		if !existingRuleNames[rule.GetName()] {
			manifestRule := contracts.ManifestRulesEngine{
				Phase: "request",
				Rule: contracts.ManifestRule{
					Name:        rule.GetName(),
					Description: rule.GetDescription(),
					Active:      rule.GetActive(),
					Criteria:    criteria,
					Behaviors:   behaviors,
				},
			}

			manifest.Applications[0].Rules = append(manifest.Applications[0].Rules, manifestRule)

			newRule := contracts.AzionJsonDataRules{
				Id:    rule.GetId(),
				Name:  rule.GetName(),
				Phase: "request",
			}
			rulesAzion = append(rulesAzion, newRule)
			existingRuleNames[rule.GetName()] = true
		}
	}

	respResp, err := client.ListRulesEngineResponse(context.Background(), opts, str)
	if err != nil {
		logger.Debug("Error while listing response phase rules", zap.Error(err))
		return fmt.Errorf(msg.ERRORLISTRESPONSERULES, err)
	}

	for _, rule := range respResp.Results {
		if rule.Name == "Default Rule" || rule.Name == "enable gzip" {
			//default rule or enable gzip rule are not added to azion.json or azion.config
			continue
		}

		var criteria [][]edgesdk.EdgeApplicationCriterionFieldRequest
		criteriaBytes, err := json.Marshal(rule.Criteria)
		if err != nil {
			return fmt.Errorf(msg.ERRORMARSHALCRITERIA, err)
		}
		err = json.Unmarshal(criteriaBytes, &criteria)
		if err != nil {
			return fmt.Errorf(msg.ERRORUNMARSHALCRITERIA, err)
		}

		var behaviors []contracts.ManifestRuleBehavior
		behaviorsBytes, err := json.Marshal(rule.Behaviors)
		if err != nil {
			return fmt.Errorf(msg.ERRORMARSHALBEHAVIORS, err)
		}
		err = json.Unmarshal(behaviorsBytes, &behaviors)
		if err != nil {
			return fmt.Errorf(msg.ERRORUNMARSHALBEHAVIORS, err)
		}

		if !existingRuleNames[rule.GetName()] {
			manifestRule := contracts.ManifestRulesEngine{
				Phase: "response",
				Rule: contracts.ManifestRule{
					Name:        rule.GetName(),
					Description: rule.GetDescription(),
					Active:      rule.GetActive(),
					Criteria:    criteria,
					Behaviors:   behaviors,
				},
			}

			manifest.Applications[0].Rules = append(manifest.Applications[0].Rules, manifestRule)

			newRule := contracts.AzionJsonDataRules{
				Id:    rule.GetId(),
				Name:  rule.GetName(),
				Phase: "response",
			}
			rulesAzion = append(rulesAzion, newRule)
			existingRuleNames[rule.GetName()] = true
		}
	}

	// Update the configuration with all rules
	info.Conf.RulesEngine.Rules = rulesAzion
	err = utils.WriteAzionJsonContentPreserveOrder(info.Conf, ProjectConf)
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
		return nil // not every project has a .env file; this should not stop the execution
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

func WriteManifest(manifest *contracts.ManifestV4, pathMan string) error {

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
