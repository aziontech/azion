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
	vulcanPkg "github.com/aziontech/azion-cli/pkg/vulcan"
	"github.com/aziontech/azion-cli/utils"
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
	manifest := &contracts.ManifestV4{}

	_, err = synch.syncCache(info, f)
	if err != nil {
		return fmt.Errorf(msg.ERRORSYNC, err.Error())
	}

	err = synch.syncRules(info, f, manifest)
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
		command := vul.Command("", "manifest transform --output %s --entry %s", f)
		err = synch.CommandRunInteractive(f, fmt.Sprintf(command, fileName, "manifesttoconvert.json"))
		if err != nil {
			return err
		}
	}

	return nil
}

func (synch *SyncCmd) syncCache(info contracts.SyncOpts, f *cmdutil.Factory) (map[string]contracts.AzionJsonDataCacheSettings, error) {
	remoteCacheIds := make(map[string]contracts.AzionJsonDataCacheSettings)
	client := edgeApp.NewClient(f.HttpClient, f.Config.GetString("api_v4_url"), f.Config.GetString("token"))
	str := strconv.FormatInt(info.Conf.Application.ID, 10)
	resp, err := client.ListCacheEdgeApp(context.Background(), str, opts)
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

		newCache := contracts.AzionJsonDataCacheSettings{
			Id:   cache.GetId(),
			Name: cache.GetName(),
		}
		cacheAzion = append(cacheAzion, newCache)
		info.Conf.CacheSettings = cacheAzion
	}
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
	rulesAzion := []contracts.AzionJsonDataRules{}
	info.Conf.RulesEngine.Rules = rulesAzion

	// Initialize Applications slice if it's empty
	if len(manifest.Applications) == 0 {
		manifest.Applications = append(manifest.Applications, contracts.Applications{
			Name:  info.Conf.Application.Name,
			Rules: []contracts.ManifestRulesEngine{},
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

		manifestRule := contracts.ManifestRulesEngine{
			Phase: "request",
			Rule: contracts.ManifestRule{
				Name:        rule.GetName(),
				Description: rule.GetDescription(),
				Active:      rule.GetActive(),
				// Convert criteria and behaviors to appropriate types
			},
		}

		manifest.Applications[0].Rules = append(manifest.Applications[0].Rules, manifestRule)

		newRule := contracts.AzionJsonDataRules{
			Id:    rule.GetId(),
			Name:  rule.GetName(),
			Phase: "request",
		}
		rulesAzion = append(rulesAzion, newRule)
	}

	respResp, err := client.ListRulesEngineResponse(context.Background(), opts, str)
	if err != nil {
		logger.Debug("Error while listing response phase rules", zap.Error(err))
		return fmt.Errorf("failed to list response phase rules: %w", err)
	}

	for _, rule := range respResp.Results {
		if rule.Name == "Default Rule" || rule.Name == "enable gzip" {
			//default rule or enable gzip rule are not added to azion.json or azion.config
			continue
		}

		manifestRule := contracts.ManifestRulesEngine{
			Phase: "response",
			Rule: contracts.ManifestRule{
				Name:        rule.GetName(),
				Description: rule.GetDescription(),
				Active:      rule.GetActive(),
				// Convert criteria and behaviors to appropriate types
			},
		}

		manifest.Applications[0].Rules = append(manifest.Applications[0].Rules, manifestRule)

		newRule := contracts.AzionJsonDataRules{
			Id:    rule.GetId(),
			Name:  rule.GetName(),
			Phase: "response",
		}
		rulesAzion = append(rulesAzion, newRule)
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
