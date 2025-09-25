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
	apiCache "github.com/aziontech/azion-cli/pkg/v3api/cache_setting"
	"github.com/aziontech/azion-cli/pkg/v3commands/purge"

	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/logger"
	apiDomain "github.com/aziontech/azion-cli/pkg/v3api/domain"
	apiEdgeApplications "github.com/aziontech/azion-cli/pkg/v3api/edge_applications"
	apiOrigin "github.com/aziontech/azion-cli/pkg/v3api/origin"
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
	WriteAzionJsonContent func(conf *contracts.AzionApplicationOptionsV3, confPath string) error
}

func NewManifestInterpreter() *ManifestInterpreter {
	return &ManifestInterpreter{
		FileReader:            os.ReadFile,
		GetWorkDir:            utils.GetWorkingDir,
		WriteAzionJsonContent: utils.WriteAzionJsonContentV3,
	}
}

func (man *ManifestInterpreter) ManifestPath() (string, error) {
	pathWorkingDir, err := man.GetWorkDir()
	if err != nil {
		return "", err
	}

	return utils.Concat(pathWorkingDir, manifestFilePath), nil
}

func (man *ManifestInterpreter) ReadManifest(
	path string, f *cmdutil.Factory, msgs *[]string) (*contracts.Manifest, error) {
	logger.FInfoFlags(f.IOStreams.Out, msg.ReadingManifest, f.Format, f.Out)
	*msgs = append(*msgs, msg.ReadingManifest)
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

func (man *ManifestInterpreter) CreateResources(
	conf *contracts.AzionApplicationOptionsV3,
	manifest *contracts.Manifest,
	f *cmdutil.Factory,
	projectConf string,
	msgs *[]string) error {

	logger.FInfoFlags(f.IOStreams.Out, msg.CreatingManifest, f.Format, f.Out)
	*msgs = append(*msgs, msg.CreatingManifest)

	client := apiEdgeApplications.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))
	clientCache := apiCache.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))
	clientOrigin := apiOrigin.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))
	clientDomain := apiDomain.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))
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

	if manifest.Domain != nil && manifest.Domain.Name != "" {
		if conf.Domain.Id > 0 {
			requestUpdate := makeDomainUpdateRequest(manifest.Domain, conf)
			updated, err := clientDomain.Update(ctx, requestUpdate)
			if err != nil {
				return fmt.Errorf("%w - '%s': %s", msg.ErrorUpdateDomain, *requestUpdate.Name, err.Error())
			}
			conf.Domain.Name = updated.GetName()
			conf.Domain.DomainName = updated.GetDomainName()
			conf.Domain.Url = utils.Concat("https://", updated.GetDomainName())
		} else {
			requestCreate := makeDomainCreateRequest(manifest.Domain, conf)
			created, err := clientDomain.Create(ctx, requestCreate)
			if err != nil {
				return fmt.Errorf("%w - '%s': %s", msg.ErrorUpdateDomain, requestCreate.Name, err.Error())
			}
			conf.Domain.Name = created.GetName()
			conf.Domain.DomainName = created.GetDomainName()
			conf.Domain.Url = utils.Concat("https://", created.GetDomainName())
			conf.Domain.Id = created.GetId()
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
			requestUpdate.Id = r.Id
			requestUpdate.Phase = rule.Phase
			requestUpdate.IsActive = &rule.IsActive
			requestUpdate.Order = &rule.Order
			requestUpdate.IdApplication = conf.Application.ID
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
			requestCreate.IsActive = &rule.IsActive
			requestCreate.Order = &rule.Order
			created, err := client.CreateRulesEngine(ctx, conf.Application.ID, rule.Phase, requestCreate)
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

	purgeCmd := purge.NewPurgeCmd(f)
	for _, purgeObj := range manifest.Purge {
		switch purgeObj.Type {
		case "url":
			err := purgeCmd.PurgeUrls(purgeObj.Urls, f)
			if err != nil {
				logger.Debug("Error while purging urls", zap.Error(err))
				return nil
			}
		case "cachekey":
			err := purgeCmd.PurgeCacheKeys(purgeObj.Urls, f)
			if err != nil {
				logger.Debug("Error while purging cache keys", zap.Error(err))
				return nil
			}
		case "wildcard":
			err := purgeCmd.PurgeWildcard(purgeObj.Urls, f)
			if err != nil {
				logger.Debug("Error while purging wildcards", zap.Error(err))
				return nil
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
func deleteResources(ctx context.Context, f *cmdutil.Factory, conf *contracts.AzionApplicationOptionsV3, msgs *[]string) error {
	client := apiEdgeApplications.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))
	clientCache := apiCache.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))
	clientOrigin := apiOrigin.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))

	if conf.SkipDeletion {
		logger.FInfoFlags(f.IOStreams.Out, msg.SkipDeletion, f.Format, f.Out)
		*msgs = append(*msgs, msg.SkipDeletion)
		return nil
	}

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
