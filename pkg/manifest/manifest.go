package manifest

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	msg "github.com/aziontech/azion-cli/messages/manifest"
	apiCache "github.com/aziontech/azion-cli/pkg/api/cache_setting"
	apiEdgeApplications "github.com/aziontech/azion-cli/pkg/api/edge_applications"
	apiOrigin "github.com/aziontech/azion-cli/pkg/api/origin"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/utils"
	thoth "github.com/aziontech/go-thoth"
	"go.uber.org/zap"
)

var manifestFilePath = "/.edge/manifest.json"

type ManifestInterpreter struct {
	FileReader            func(path string) ([]byte, error)
	GetWorkDir            func() (string, error)
	WriteAzionJsonContent func(conf *contracts.AzionApplicationOptions) error
}

func NewManifestInterpreter() *ManifestInterpreter {
	return &ManifestInterpreter{
		FileReader:            os.ReadFile,
		GetWorkDir:            utils.GetWorkingDir,
		WriteAzionJsonContent: utils.WriteAzionJsonContent,
	}
}

func (man *ManifestInterpreter) ManifestPath() (string, error) {
	pathWorkingDir, err := man.GetWorkDir()
	if err != nil {
		return "", err
	}

	return utils.Concat(pathWorkingDir, manifestFilePath), nil
}

func (man *ManifestInterpreter) ReadManifest(path string, f *cmdutil.Factory) (*contracts.Manifest, error) {
	logger.FInfo(f.IOStreams.Out, msg.ReadingManifest)
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

func (man *ManifestInterpreter) CreateResources(conf *contracts.AzionApplicationOptions, manifest *contracts.Manifest, f *cmdutil.Factory) error {
	logger.FInfo(f.IOStreams.Out, msg.CreatingManifest)

	client := apiEdgeApplications.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))
	clientCache := apiCache.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))
	clientOrigin := apiOrigin.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))
	ctx := context.Background()

	cacheIds := make(map[string]int64)
	ruleIds := make(map[string]int64)
	originKeys := make(map[string]int64)

	for _, cacheConf := range conf.CacheSettings {
		cacheIds[cacheConf.Name] = cacheConf.Id
	}

	for _, ruleConf := range conf.RulesEngine.Rules {
		ruleIds[ruleConf.Name] = ruleConf.Id
	}

	for _, originConf := range conf.Origin {
		originKeys[originConf.Name] = originConf.OriginId
	}

	for _, origin := range manifest.Origins {
		if id := originKeys[origin.Name]; id == 0 {
			requestCreate := makeOriginCreateRequest(origin, conf)
			if origin.Name != "" {
				requestCreate.Name = origin.Name
			} else {
				requestCreate.Name = conf.Name
			}
			created, err := clientOrigin.Create(ctx, conf.Application.ID, requestCreate)
			if err != nil {
				return err
			}
			newOrigin := contracts.AzionJsonDataOrigin{
				OriginId:  created.GetOriginId(),
				OriginKey: created.GetOriginKey(),
				Name:      created.GetName(),
			}
			conf.Origin = append(conf.Origin, newOrigin)
			originKeys[created.GetName()] = created.GetOriginId()
			logger.FInfo(f.IOStreams.Out, fmt.Sprintf(msg.ManifestCreateOrigin, origin.Name, created.GetOriginId()))
		}
	}

	for _, cache := range manifest.CacheSettings {
		if id := cacheIds[*cache.Name]; id > 0 {
			requestUpdate := makeCacheRequestUpdate(cache)
			if cache.Name != nil {
				requestUpdate.Name = cache.Name
			} else {
				requestUpdate.Name = &conf.Name
			}
			_, err := clientCache.Update(ctx, requestUpdate, conf.Application.ID, id)
			if err != nil {
				return err
			}
			logger.FInfo(f.IOStreams.Out, fmt.Sprintf(msg.ManifestUpdateCache, *cache.Name, id))
		} else {
			requestUpdate := makeCacheRequestCreate(cache)
			if cache.Name != nil {
				requestUpdate.Name = *cache.Name
			} else {
				requestUpdate.Name = conf.Name + thoth.GenerateName()
			}
			created, err := clientCache.Create(ctx, requestUpdate, conf.Application.ID)
			if err != nil {
				return err
			}
			newCache := contracts.AzionJsonDataCacheSettings{
				Id:   created.GetId(),
				Name: created.GetName(),
			}
			conf.CacheSettings = append(conf.CacheSettings, newCache)
			logger.FInfo(f.IOStreams.Out, fmt.Sprintf(msg.ManifestCreateCache, *cache.Name, id))
		}
	}

	err := man.WriteAzionJsonContent(conf)
	if err != nil {
		logger.Debug("Error while writing azion.json file", zap.Error(err))
		return err
	}

	for _, rule := range manifest.Rules {
		if id := ruleIds[rule.Name]; id > 0 {
			requestUpdate, err := makeRuleRequestUpdate(rule, cacheIds, conf, originKeys)
			if err != nil {
				return err
			}
			requestUpdate.Id = id
			requestUpdate.Phase = "request"
			requestUpdate.IdApplication = conf.Application.ID
			_, err = client.UpdateRulesEngine(ctx, requestUpdate)
			if err != nil {
				return err
			}

		} else {
			requestCreate, err := makeRuleRequestCreate(rule, cacheIds, conf, originKeys, client, ctx)
			if err != nil {
				return err
			}
			if rule.Name != "" {
				requestCreate.Name = rule.Name
			} else {
				requestCreate.Name = conf.Name + thoth.GenerateName()
			}
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
	}

	err = man.WriteAzionJsonContent(conf)
	if err != nil {
		logger.Debug("Error while writing azion.json file", zap.Error(err))
		return err
	}

	return nil
}
