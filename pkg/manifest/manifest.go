package manifest

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	msgcache "github.com/aziontech/azion-cli/messages/cache_setting"
	msgrule "github.com/aziontech/azion-cli/messages/delete/rules_engine"
	msg "github.com/aziontech/azion-cli/messages/manifest"
	msgorigin "github.com/aziontech/azion-cli/messages/origin"
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

var (
	CacheIds         map[string]int64
	CacheIdsBackup   map[string]int64
	RuleIds          map[string]contracts.RuleIdsStruct
	OriginKeys       map[string]string
	OriginIds        map[string]int64
	manifestFilePath = "/.edge/manifest.json"
)

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

	CacheIds = make(map[string]int64)
	CacheIdsBackup = make(map[string]int64)
	RuleIds = make(map[string]contracts.RuleIdsStruct)
	OriginKeys = make(map[string]string)
	OriginIds = make(map[string]int64)

	for _, cacheConf := range conf.CacheSettings {
		CacheIds[cacheConf.Name] = cacheConf.Id
	}

	for _, ruleConf := range conf.RulesEngine.Rules {
		RuleIds[ruleConf.Name] = contracts.RuleIdsStruct{
			Id:    ruleConf.Id,
			Phase: ruleConf.Phase,
		}
	}

	for _, originConf := range conf.Origin {
		OriginKeys[originConf.Name] = originConf.OriginKey
		OriginIds[originConf.Name] = originConf.OriginId
	}

	originConf := []contracts.AzionJsonDataOrigin{}
	for _, origin := range manifest.Origins {
		if id := OriginIds[origin.Name]; id > 0 {
			requestUpdate := makeOriginUpdateRequest(origin, conf)
			if origin.Name != "" {
				requestUpdate.Name = &origin.Name
			} else {
				requestUpdate.Name = &conf.Name
			}
			updated, err := clientOrigin.Update(ctx, conf.Application.ID, OriginKeys[origin.Name], requestUpdate)
			if err != nil {
				return fmt.Errorf("%w: %s", msg.ErrorUpdateOrigin, err.Error())
			}

			newEntry := contracts.AzionJsonDataOrigin{
				OriginId:  updated.GetOriginId(),
				OriginKey: updated.GetOriginKey(),
				Name:      updated.GetName(),
			}
			originConf = append(originConf, newEntry)
			logger.FInfo(f.IOStreams.Out, fmt.Sprintf(msg.ManifestUpdateOrigin, origin.Name, updated.GetOriginKey()))
		} else {
			requestCreate := makeOriginCreateRequest(origin, conf)
			if origin.Name != "" {
				requestCreate.Name = origin.Name
			} else {
				requestCreate.Name = conf.Name
			}
			created, err := clientOrigin.Create(ctx, conf.Application.ID, requestCreate)
			if err != nil {
				return fmt.Errorf("%w: %s", msg.ErrorCreateOrigin, err.Error())
			}
			newOrigin := contracts.AzionJsonDataOrigin{
				OriginId:  created.GetOriginId(),
				OriginKey: created.GetOriginKey(),
				Name:      created.GetName(),
			}

			originConf = append(originConf, newOrigin)
			OriginIds[created.GetName()] = created.GetOriginId()
			OriginKeys[created.GetName()] = created.GetOriginKey()
			logger.FInfo(f.IOStreams.Out, fmt.Sprintf(msg.ManifestCreateOrigin, origin.Name, created.GetOriginId()))
		}
	}

	conf.Origin = originConf
	err := man.WriteAzionJsonContent(conf)
	if err != nil {
		logger.Debug("Error while writing azion.json file", zap.Error(err))
		return err
	}

	cacheConf := []contracts.AzionJsonDataCacheSettings{}
	for _, cache := range manifest.CacheSettings {
		if id := CacheIds[*cache.Name]; id > 0 {
			requestUpdate := makeCacheRequestUpdate(cache)
			if cache.Name != nil {
				requestUpdate.Name = cache.Name
			} else {
				requestUpdate.Name = &conf.Name
			}
			updated, err := clientCache.Update(ctx, requestUpdate, conf.Application.ID, id)
			if err != nil {
				return fmt.Errorf("%w: %s", msg.ErrorUpdateCache, err.Error())
			}
			newCache := contracts.AzionJsonDataCacheSettings{
				Id:   updated.GetId(),
				Name: updated.GetName(),
			}
			cacheConf = append(cacheConf, newCache)
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
				return fmt.Errorf("%w: %s", msg.ErrorCreateCache, err.Error())
			}
			newCache := contracts.AzionJsonDataCacheSettings{
				Id:   created.GetId(),
				Name: created.GetName(),
			}
			cacheConf = append(cacheConf, newCache)
			CacheIds[newCache.Name] = newCache.Id
			logger.FInfo(f.IOStreams.Out, fmt.Sprintf(msg.ManifestCreateCache, *cache.Name, newCache.Id))
		}
	}

	//backup cache ids
	for k, v := range CacheIds {
		CacheIdsBackup[k] = v
	}

	conf.CacheSettings = cacheConf
	err = man.WriteAzionJsonContent(conf)
	if err != nil {
		logger.Debug("Error while writing azion.json file", zap.Error(err))
		return err
	}

	ruleConf := []contracts.AzionJsonDataRules{}
	for _, rule := range manifest.Rules {
		if r := RuleIds[rule.Name]; r.Id > 0 {
			requestUpdate, err := makeRuleRequestUpdate(rule, conf)
			if err != nil {
				return err
			}
			requestUpdate.Id = r.Id
			requestUpdate.Phase = rule.Phase
			requestUpdate.IsActive = &rule.IsActive
			requestUpdate.Order = &rule.Order
			requestUpdate.IdApplication = conf.Application.ID
			updated, err := client.UpdateRulesEngine(ctx, requestUpdate)
			if err != nil {
				return fmt.Errorf("%w: %s", msg.ErrorUpdateRule, err.Error())
			}
			newRule := contracts.AzionJsonDataRules{
				Id:    updated.GetId(),
				Name:  updated.GetName(),
				Phase: updated.GetPhase(),
			}
			ruleConf = append(ruleConf, newRule)
			delete(RuleIds, rule.Name)
		} else {
			requestCreate, err := makeRuleRequestCreate(rule, conf, client, ctx)
			if err != nil {
				return err
			}
			if rule.Name != "" {
				requestCreate.Name = rule.Name
			} else {
				requestCreate.Name = conf.Name + thoth.GenerateName()
			}
			requestCreate.IsActive = &rule.IsActive
			requestCreate.Order = &rule.Order
			created, err := client.CreateRulesEngine(ctx, conf.Application.ID, rule.Phase, requestCreate)
			if err != nil {
				return fmt.Errorf("%w: %s", msg.ErrorCreateRule, err.Error())
			}
			newRule := contracts.AzionJsonDataRules{
				Id:    created.GetId(),
				Name:  created.GetName(),
				Phase: created.GetPhase(),
			}
			ruleConf = append(ruleConf, newRule)
		}
	}

	conf.RulesEngine.Rules = ruleConf
	err = man.WriteAzionJsonContent(conf)
	if err != nil {
		logger.Debug("Error while writing azion.json file", zap.Error(err))
		return err
	}

	err = deleteResources(ctx, f, conf)
	if err != nil {
		return err
	}

	return nil
}

// this is called to delete resources no longer present in manifest.json
func deleteResources(ctx context.Context, f *cmdutil.Factory, conf *contracts.AzionApplicationOptions) error {
	client := apiEdgeApplications.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))
	clientCache := apiCache.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))
	clientOrigin := apiOrigin.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))

	for _, value := range RuleIds {
		//since until [UXE-3599] was carried out we'd only cared about "request" phase, this check guarantees that if Phase is empty
		// we are probably dealing with a rule engine from a previous version
		phase := "request"
		if value.Phase != "" {
			phase = value.Phase
		}
		err := client.DeleteRulesEngine(ctx, conf.Application.ID, phase, value.Id)
		if err != nil {
			return err
		}
		logger.FInfo(f.IOStreams.Out, fmt.Sprintf(msgrule.DeleteOutputSuccess+"\n", value.Id))
	}

	for i, value := range OriginKeys {
		if strings.Contains(i, "_single") {
			continue
		}
		err := clientOrigin.DeleteOrigins(ctx, conf.Application.ID, value)
		if err != nil {
			return err
		}
		logger.FInfo(f.IOStreams.Out, fmt.Sprintf(msgorigin.DeleteOutputSuccess+"\n", value))
	}

	for _, value := range CacheIds {
		err := clientCache.Delete(ctx, conf.Application.ID, value)
		if err != nil {
			return err
		}
		logger.FInfo(f.IOStreams.Out, fmt.Sprintf(msgcache.DeleteOutputSuccess+"\n", value))
	}

	return nil
}
