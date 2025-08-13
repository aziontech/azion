package sync

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strconv"

	msg "github.com/aziontech/azion-cli/messages/sync"
	edgeApp "github.com/aziontech/azion-cli/pkg/api/edge_applications"
	varApi "github.com/aziontech/azion-cli/pkg/api/variables"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/logger"
	vulcanPkg "github.com/aziontech/azion-cli/pkg/vulcan"
	"github.com/aziontech/azion-cli/utils"
	edgesdk "github.com/aziontech/azionapi-v4-go-sdk/edge-api"
	"go.uber.org/zap"
)

var (
	opts  *contracts.ListOptions
	ctx   context.Context = context.Background()
	words                 = []string{"PASSWORD", "PWD", "SECRET", "HASH", "ENCRYPTED", "PASSCODE", "AUTH", "TOKEN", "SECRET"}
)

func SyncLocalResources(f *cmdutil.Factory, info contracts.SyncOpts, synch *SyncCmd) error {
	opts = &contracts.ListOptions{
		PageSize: 1000,
		Page:     1,
	}

	if info.Conf.Application.ID <= 0 {
		return msg.ERRORNOTDEPLOYED
	}

	var err error
	manifest := &contracts.ManifestV4{}

	// remoteCaches, err := synch.syncCache(info, f, manifest)
	// if err != nil {
	// 	return fmt.Errorf(msg.ERRORSYNC, err.Error())
	// }

	// remoteConnectors, err := synch.syncConnector(info, f, manifest)
	// if err != nil {
	// 	return fmt.Errorf(msg.ERRORSYNC, err.Error())
	// }

	// err = synch.syncRules(info, f, manifest, remoteCaches)
	// if err != nil {
	// 	return fmt.Errorf(msg.ERRORSYNC, err.Error())
	// }

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

// func (synch *SyncCmd) syncConnector(info contracts.SyncOpts, f *cmdutil.Factory, manifest *contracts.ManifestV4) (map[string]contracts.AzionJsonDataConnectors, error) {
// 	remoteConnectorIds := make(map[string]contracts.AzionJsonDataConnectors)
// 	client := connectorApi.NewClient(f.HttpClient, f.Config.GetString("api_v4_url"), f.Config.GetString("token"))
// 	resp, err := client.List(ctx, opts)
// 	if err != nil {
// 		return remoteConnectorIds, err
// 	}

// 	ConnectorsAzion := []contracts.AzionJsonDataConnectors{}
// 	info.Conf.Connectors = ConnectorsAzion

// 	for _, connector := range resp.Results {
// 		remoteConnectorIds[strconv.FormatInt(*&connector.Id, 10)] = contracts.AzionJsonDataConnectors{
// 			Id:      connector.Id,
// 			Name:    connector.Name,
// 			Address: connector.Addresses,
// 		}
// 		jsonBytes, err := json.Marshal(connector)
// 		if err != nil {
// 			return remoteConnectorIds, err
// 		}
// 		oEntry := edgesdk.EdgeConnectorPolymorphicRequest{}
// 		err = json.Unmarshal(jsonBytes, &oEntry)
// 		if err != nil {
// 			return remoteConnectorIds, err
// 		}
// 		manifest.EdgeConnectors = append(manifest.EdgeConnectors, oEntry)
// 		newConnector := contracts.AzionJsonDataConnectors{
// 			Id:      connector.Id,
// 			Name:    connector.Name,
// 			Address: connector.Addresses,
// 		}
// 		ConnectorsAzion = append(ConnectorsAzion, newConnector)
// 		info.Conf.Connectors = ConnectorsAzion
// 	}
// 	err = synch.WriteAzionJsonContent(info.Conf, ProjectConf)
// 	if err != nil {
// 		logger.Debug("Error while writing azion.json file", zap.Error(err))
// 		return remoteConnectorIds, err
// 	}
// 	return remoteConnectorIds, nil
// }

func (synch *SyncCmd) syncCache(info contracts.SyncOpts, f *cmdutil.Factory, manifest *contracts.ManifestV4) (map[string]contracts.AzionJsonDataCacheSettings, error) {
	remoteCacheIds := make(map[string]contracts.AzionJsonDataCacheSettings)
	client := edgeApp.NewClient(f.HttpClient, f.Config.GetString("api_v4_url"), f.Config.GetString("token"))
	str := strconv.FormatInt(info.Conf.Application.ID, 10)
	resp, err := client.ListCacheEdgeApp(context.Background(), str)
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
		cEntry := edgesdk.CacheSettingRequest{}
		jsonBytes, err := json.Marshal(cache)
		if err != nil {
			return remoteCacheIds, err
		}
		err = json.Unmarshal(jsonBytes, &cEntry)
		if err != nil {
			return remoteCacheIds, err
		}
		// manifest.EdgeApplications[0].Cache = append(manifest.EdgeApplications[0].Cache, cEntry)

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

// func (synch *SyncCmd) syncRules(info contracts.SyncOpts, f *cmdutil.Factory, manifest *contracts.ManifestV4,
// 	remoteCacheIds map[string]contracts.AzionJsonDataCacheSettings) error {
// 	// Get request rules first
// 	client := edgeApp.NewClient(f.HttpClient, f.Config.GetString("api_v4_url"), f.Config.GetString("token"))
// 	str := strconv.FormatInt(info.Conf.Application.ID, 10)
// 	resp, err := client.ListRulesEngine(context.Background(), opts, str)
// 	if err != nil {
// 		return err
// 	}
// 	rulesAzion := []contracts.AzionJsonDataRules{}
// 	info.Conf.RulesEngine.Rules = rulesAzion
// 	for _, rule := range resp.Results {
// 		if rule.Name == "Default Rule" || rule.Name == "enable gzip" {
// 			//default rule or enable gzip rule are not added to azion.json or azion.config
// 			continue
// 		}
// 		jsonBytes, err := json.Marshal(rule)
// 		if err != nil {
// 			return err
// 		}
// 		rEntry := edgesdk.EdgeApplicationRuleEngineRequest{}
// 		err = json.Unmarshal(jsonBytes, &rEntry)
// 		if err != nil {
// 			return err
// 		}

// 		manifest.EdgeApplications[0].Rules = append(manifest.EdgeApplications[0].Rules, rEntry)

// 		newRule := contracts.AzionJsonDataRules{
// 			Id:    rule.GetId(),
// 			Name:  rule.GetName(),
// 			Phase: rule.GetPhase(),
// 		}
// 		rulesAzion = append(rulesAzion, newRule)
// 		info.Conf.RulesEngine.Rules = rulesAzion
// 	}
// 	err = synch.WriteAzionJsonContent(info.Conf, ProjectConf)
// 	if err != nil {
// 		logger.Debug("Error while writing azion.json file", zap.Error(err))
// 		return err
// 	}

// 	return nil
// }

func (synch *SyncCmd) syncEnv(f *cmdutil.Factory) error {

	client := varApi.NewClient(f.HttpClient, f.Config.GetString("api_v4_url"), f.Config.GetString("token"))
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
