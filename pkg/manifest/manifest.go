package manifest

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	msgcache "github.com/aziontech/azion-cli/messages/cache_setting"
	msgrule "github.com/aziontech/azion-cli/messages/delete/rules_engine"
	msg "github.com/aziontech/azion-cli/messages/manifest"
	msgorigin "github.com/aziontech/azion-cli/messages/origin"
	apiCache "github.com/aziontech/azion-cli/pkg/api/cache_setting"
	apipurge "github.com/aziontech/azion-cli/pkg/api/realtime_purge"

	apiEdgeApplications "github.com/aziontech/azion-cli/pkg/api/edge_applications"
	apiOrigin "github.com/aziontech/azion-cli/pkg/api/origin"
	apiDomain "github.com/aziontech/azion-cli/pkg/api/workloads"
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
	WriteAzionJsonContent func(conf *contracts.AzionApplicationOptions, confPath string) error
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

func (man *ManifestInterpreter) ReadManifest(path string, f *cmdutil.Factory, msgs *[]string) (*contracts.ManifestV4, error) {
	logger.FInfoFlags(f.IOStreams.Out, msg.ReadingManifest, f.Format, f.Out)
	*msgs = append(*msgs, msg.ReadingManifest)
	manifest := &contracts.ManifestV4{}

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

func (man *ManifestInterpreter) CreateResources(conf *contracts.AzionApplicationOptions, manifest *contracts.ManifestV4, f *cmdutil.Factory, projectConf string, msgs *[]string) error {

	logger.FInfoFlags(f.IOStreams.Out, msg.CreatingManifest, f.Format, f.Out)
	*msgs = append(*msgs, msg.CreatingManifest)

	client := apiEdgeApplications.NewClient(f.HttpClient, f.Config.GetString("api_v4_url"), f.Config.GetString("token"))
	clientCache := apiCache.NewClientV3(f.HttpClient, f.Config.GetString("api_v4_url"), f.Config.GetString("token"))
	clientOrigin := apiOrigin.NewClient(f.HttpClient, f.Config.GetString("api_v4_url"), f.Config.GetString("token"))
	clientDomain := apiDomain.NewClient(f.HttpClient, f.Config.GetString("api_v4_url"), f.Config.GetString("token"))
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

	if len(manifest.EdgeApplications) > 0 {
		edgeappman := manifest.EdgeApplications[0]
		if conf.Application.ID > 0 {
			req := transformEdgeApplicationRequest(edgeappman.EdgeApplicationRequest)
			_, err := client.Update(ctx, req)
			if err != nil {
				return err
			}
		} else {
			req := &apiEdgeApplications.CreateRequest{
				EdgeApplicationRequest: edgeappman.EdgeApplicationRequest,
			}
			resp, err := client.Create(ctx, req)
			if err != nil {
				return err
			}
			conf.Application.ID = resp.GetId()
		}

		if len(edgeappman.Rules) > 0 {
			for _, rule := range edgeappman.Rules {
				if r := RuleIds[rule.Name]; r.Id > 0 {
					req := transformRuleRequest(rule)
					_, err := client.UpdateRulesEngine(ctx, req)
					if err != nil {
						return err
					}
				} else {
					req := &apiEdgeApplications.CreateRulesEngineRequest{
						EdgeApplicationRuleEngineRequest: rule,
					}
					appstring := strconv.FormatInt(conf.Application.ID, 10)
					_, err := client.CreateRulesEngine(ctx, appstring, rule.Phase, req)
					if err != nil {
						return err
					}
				}
			}
		}
	}

	if len(manifest.EdgeStorage) > 0 {
		storageman := manifest.EdgeStorage[0]
	}

	if manifest.Domain != nil && manifest.Domain.Name != "" {
		if conf.Workloads.Id > 0 {
			requestUpdate := makeDomainUpdateRequest(manifest.Domain, conf)
			req := &apiDomain.UpdateRequest{}
			updated, err := clientDomain.Update(ctx, req)
			if err != nil {
				return fmt.Errorf("%w - '%s': %s", msg.ErrorUpdateDomain, *requestUpdate.Name, err.Error())
			}
			conf.Workloads.Domains = updated.GetDomains()
			conf.Workloads.Url = utils.Concat("https://", conf.Workloads.Domains[0].GetDomain())
		} else {
			requestCreate := makeDomainCreateRequest(manifest.Domain, conf)
			req := &apiDomain.CreateRequest{}
			created, err := clientDomain.Create(ctx, req)
			if err != nil {
				return fmt.Errorf("%w - '%s': %s", msg.ErrorUpdateDomain, requestCreate.Name, err.Error())
			}
			conf.Workloads.Name = created.GetName()
			conf.Workloads.Domains = created.GetDomains()
			conf.Workloads.Url = utils.Concat("https://", created.GetDomains()[0].GetDomain())
			conf.Workloads.Id = created.GetId()
		}
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
				return fmt.Errorf("%w - '%s': %s", msg.ErrorUpdateOrigin, origin.Name, err.Error())
			}

			newEntry := contracts.AzionJsonDataOrigin{
				OriginId:  updated.GetOriginId(),
				OriginKey: updated.GetOriginKey(),
				Name:      updated.GetName(),
			}
			originConf = append(originConf, newEntry)

			msgf := fmt.Sprintf(msg.ManifestUpdateOrigin, origin.Name, updated.GetOriginKey())
			logger.FInfoFlags(f.IOStreams.Out, msgf, f.Format, f.Out)
			*msgs = append(*msgs, msgf)
		} else {
			requestCreate := makeOriginCreateRequest(origin, conf)
			if origin.Name != "" {
				requestCreate.Name = origin.Name
			} else {
				requestCreate.Name = conf.Name
			}
			created, err := clientOrigin.Create(ctx, conf.Application.ID, requestCreate)
			if err != nil {
				return fmt.Errorf("%w - '%s': %s", msg.ErrorCreateOrigin, requestCreate.Name, err.Error())
			}
			newOrigin := contracts.AzionJsonDataOrigin{
				OriginId:  created.GetOriginId(),
				OriginKey: created.GetOriginKey(),
				Name:      created.GetName(),
			}

			originConf = append(originConf, newOrigin)
			OriginIds[created.GetName()] = created.GetOriginId()
			OriginKeys[created.GetName()] = created.GetOriginKey()
			msgf := fmt.Sprintf(msg.ManifestCreateOrigin, origin.Name, created.GetOriginId())
			logger.FInfoFlags(f.IOStreams.Out, msgf, f.Format, f.Out)
			*msgs = append(*msgs, msgf)
		}
	}

	conf.Origin = originConf
	err := man.WriteAzionJsonContent(conf, projectConf)
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
				return fmt.Errorf("%w - '%s': %s", msg.ErrorUpdateCache, *cache.Name, err.Error())
			}
			newCache := contracts.AzionJsonDataCacheSettings{
				Id:   updated.GetId(),
				Name: updated.GetName(),
			}
			cacheConf = append(cacheConf, newCache)
			msgf := fmt.Sprintf(msg.ManifestUpdateCache, *cache.Name, id)
			logger.FInfoFlags(f.IOStreams.Out, msgf, f.Format, f.Out)
			*msgs = append(*msgs, msgf)
		} else {
			requestCreate := makeCacheRequestCreate(cache)
			if cache.Name != nil {
				requestCreate.Name = *cache.Name
			} else {
				requestCreate.Name = conf.Name + thoth.GenerateName()
			}
			created, err := clientCache.Create(ctx, requestCreate, conf.Application.ID)
			if err != nil {
				return fmt.Errorf("%w - '%s': %s", msg.ErrorCreateCache, requestCreate.Name, err.Error())
			}
			newCache := contracts.AzionJsonDataCacheSettings{
				Id:   created.GetId(),
				Name: created.GetName(),
			}
			cacheConf = append(cacheConf, newCache)
			CacheIds[newCache.Name] = newCache.Id
			msgf := fmt.Sprintf(msg.ManifestCreateCache, *cache.Name, newCache.Id)
			logger.FInfoFlags(f.IOStreams.Out, msgf, f.Format, f.Out)
			*msgs = append(*msgs, msgf)
		}
	}

	//backup cache ids
	for k, v := range CacheIds {
		CacheIdsBackup[k] = v
	}

	conf.CacheSettings = cacheConf
	err = man.WriteAzionJsonContent(conf, projectConf)
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
			// requestUpdate.Id = r.Id
			requestUpdate.Phase = rule.Phase
			// requestUpdate.IsActive = &rule.IsActive
			// requestUpdate.Order = &rule.Order
			// requestUpdate.IdApplication = conf.Application.ID //TODO: correct these fields
			updated, err := client.UpdateRulesEngine(ctx, requestUpdate)
			if err != nil {
				return fmt.Errorf("%w - '%s': %s", msg.ErrorUpdateRule, rule.Name, err.Error())
			}
			newRule := contracts.AzionJsonDataRules{
				Id:    updated.GetId(),
				Name:  updated.GetName(),
				Phase: updated.GetPhase(),
			}
			msgf := fmt.Sprintf(msg.ManifestUpdateRule, newRule.Name, newRule.Id)
			logger.FInfoFlags(f.IOStreams.Out, msgf, f.Format, f.Out)
			*msgs = append(*msgs, msgf)
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
			// requestCreate.IsActive = &rule.IsActive
			// requestCreate.Order = &rule.Order //TODO: correct these fields
			created, err := client.CreateRulesEngine(ctx, "conf.Application.ID", rule.Phase, requestCreate)
			if err != nil {
				return fmt.Errorf("%w - '%s': %s", msg.ErrorCreateRule, requestCreate.Name, err.Error())
			}
			newRule := contracts.AzionJsonDataRules{
				Id:    created.GetId(),
				Name:  created.GetName(),
				Phase: created.GetPhase(),
			}
			ruleConf = append(ruleConf, newRule)
			msgf := fmt.Sprintf(msg.ManifestCreateRule, newRule.Name, newRule.Id)
			logger.FInfoFlags(f.IOStreams.Out, msgf, f.Format, f.Out)
			*msgs = append(*msgs, msgf)
		}
	}

	conf.RulesEngine.Rules = ruleConf
	err = man.WriteAzionJsonContent(conf, projectConf)
	if err != nil {
		logger.Debug("Error while writing azion.json file", zap.Error(err))
		return err
	}

	clipurge := apipurge.NewClient(f.HttpClient, f.Config.GetString("api_v4_url"), f.Config.GetString("token"))
	for _, purgeObj := range manifest.Purge {
		switch purgeObj.Type {
		case "url":
			err = clipurge.PurgeCache(ctx, purgeObj.Urls, "url", "edge_cache")
			if err != nil {
				logger.Debug("Error while purging domains", zap.Error(err))
				return err
			}
		case "cachekey":
			err = clipurge.PurgeCache(ctx, purgeObj.Urls, "cachekey", "edge_cache")
			if err != nil {
				logger.Debug("Error while purging domains", zap.Error(err))
				return err
			}
		case "wildcard":
			err = clipurge.PurgeCache(ctx, purgeObj.Urls, "wildcard", "edge_cache")
			if err != nil {
				logger.Debug("Error while purging domains", zap.Error(err))
				return err
			}
		}
	}

	err = deleteResources(ctx, f, conf, msgs)
	if err != nil {
		return err
	}

	return nil
}

// this is called to delete resources no longer present in manifest.json
func deleteResources(ctx context.Context, f *cmdutil.Factory, conf *contracts.AzionApplicationOptions, msgs *[]string) error {
	client := apiEdgeApplications.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))
	clientCache := apiCache.NewClientV3(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))
	clientOrigin := apiOrigin.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))

	for _, value := range RuleIds {
		//since until [UXE-3599] was carried out we'd only cared about "request" phase, this check guarantees that if Phase is empty
		// we are probably dealing with a rule engine from a previous version
		phase := "request"
		if value.Phase != "" {
			phase = value.Phase
		}
		err := client.DeleteRulesEngine(ctx, "conf.Application.ID", phase, "value.Id")
		if err != nil {
			return err
		}
		msgf := fmt.Sprintf(msgrule.DeleteOutputSuccess+"\n", value.Id)
		logger.FInfoFlags(f.IOStreams.Out, msgf, f.Format, f.Out)
		*msgs = append(*msgs, msgf)
	}

	for i, value := range OriginKeys {
		if strings.Contains(i, "_single") {
			continue
		}
		err := clientOrigin.DeleteOrigins(ctx, conf.Application.ID, value)
		if err != nil {
			return err
		}
		msgf := fmt.Sprintf(msgorigin.DeleteOutputSuccess+"\n", value)
		logger.FInfoFlags(f.IOStreams.Out, msgf, f.Format, f.Out)
		*msgs = append(*msgs, msgf)
	}

	for _, value := range CacheIds {
		err := clientCache.Delete(ctx, conf.Application.ID, value)
		if err != nil {
			return err
		}
		msgf := fmt.Sprintf(msgcache.DeleteOutputSuccess+"\n", value)
		logger.FInfoFlags(f.IOStreams.Out, msgf, f.Format, f.Out)
		*msgs = append(*msgs, msgf)
	}

	return nil
}
