package manifest

import (
	"context"
	"errors"
	"fmt"
	"strings"

	msgcache "github.com/aziontech/azion-cli/messages/cache_setting"
	msgrule "github.com/aziontech/azion-cli/messages/delete/rules_engine"
	msg "github.com/aziontech/azion-cli/messages/manifest"
	msgorigin "github.com/aziontech/azion-cli/messages/origin"
	apiCache "github.com/aziontech/azion-cli/pkg/api/v3/cache_setting"
	apiDomain "github.com/aziontech/azion-cli/pkg/api/v3/domain"
	apiEdgeApplications "github.com/aziontech/azion-cli/pkg/api/v3/edge_applications"
	apiOrigin "github.com/aziontech/azion-cli/pkg/api/v3/origin"
	purge "github.com/aziontech/azion-cli/pkg/cmd/purge/v3"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/pipeline"
	"github.com/aziontech/azion-cli/utils"
	thoth "github.com/aziontech/go-thoth"
	"go.uber.org/zap"
)

// errPurgeAborted stops the pipeline without reporting a failure.
//
// It preserves an oddity of the sequence this replaces: a failing purge was
// logged at debug level and then returned nil from CreateResources, which
// reported success to the caller and skipped orphan removal entirely. That is
// almost certainly a bug, but changing it would change what a user sees on a
// failed purge, so it is reproduced exactly and left for a separate decision.
var errPurgeAborted = errors.New("purge failed; the run stops and reports success")

// ResourceContext is the state the v3 deploy steps operate on: the request
// context, the project config, the manifest being applied and the four API
// clients the steps share.
//
// The v4 side has carried an equivalent type for a while. The v3 sequence built
// its clients as locals inside one long function instead, so this type is new
// here — it is what lets the steps be separate functions at all.
type ResourceContext struct {
	ctx         context.Context
	factory     *cmdutil.Factory
	conf        *contracts.AzionApplicationOptionsV3
	manifest    *contracts.Manifest
	projectConf string
	msgs        *[]string
	writeConfig func(conf *contracts.AzionApplicationOptionsV3, confPath string) error

	client       *apiEdgeApplications.Client
	clientCache  *apiCache.Client
	clientOrigin *apiOrigin.Client
	clientDomain *apiDomain.Client
}

// newResourceContext builds the state and seeds the package-level id maps from
// the project config.
//
// The maps stay package-level because request.go reads them while building rule
// requests, and deleteResources reads what is left in them to decide what is an
// orphan. Moving them onto this type would change that, so they are seeded here
// exactly where the original seeded them: before any API call.
func newResourceContext(
	f *cmdutil.Factory,
	conf *contracts.AzionApplicationOptionsV3,
	manifest *contracts.Manifest,
	projectConf string,
	msgs *[]string,
	writeConfig func(conf *contracts.AzionApplicationOptionsV3, confPath string) error,
) *ResourceContext {
	url := f.Config.GetString("api_url")
	token := f.Config.GetString("token")

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

	return &ResourceContext{
		ctx:         context.Background(),
		factory:     f,
		conf:        conf,
		manifest:    manifest,
		projectConf: projectConf,
		msgs:        msgs,
		writeConfig: writeConfig,

		client:       apiEdgeApplications.NewClient(f.HttpClient, url, token),
		clientCache:  apiCache.NewClient(f.HttpClient, url, token),
		clientOrigin: apiOrigin.NewClient(f.HttpClient, url, token),
		clientDomain: apiDomain.NewClient(f.HttpClient, url, token),
	}
}

// print writes a line to the user and records it for --format json.
func (rc *ResourceContext) print(line string) {
	logger.FInfoFlags(rc.factory.IOStreams.Out, line, rc.factory.Format, rc.factory.Out)
	*rc.msgs = append(*rc.msgs, line)
}

// V3Pipeline is the deploy pipeline for the v3 API.
//
// Unlike the v4 steps, most of these carry no predicate. In the sequence this
// replaces, the origins, cache and rules blocks were plain loops with no
// enclosing length check, and each ended by assigning its result back to the
// project config and writing azion.json — so an empty list in the manifest
// still emptied that section of azion.json. Gating them on "the manifest
// declares some" would quietly stop that from happening.
type V3Pipeline struct{}

func (V3Pipeline) Name() string { return "v3" }

func (V3Pipeline) Steps() []pipeline.Step[ResourceContext] {
	return []pipeline.Step[ResourceContext]{
		{
			Name: "ManifestDomain",
			Applies: func(rc *ResourceContext) bool {
				return rc.manifest.Domain != nil && rc.manifest.Domain.Name != ""
			},
			Run: func(_ context.Context, rc *ResourceContext) error { return rc.applyDomain() },
		},
		{
			Name: "ManifestOrigins",
			Run:  func(_ context.Context, rc *ResourceContext) error { return rc.applyOrigins() },
		},
		{
			Name: "ManifestCacheSettings",
			Run:  func(_ context.Context, rc *ResourceContext) error { return rc.applyCacheSettings() },
		},
		{
			Name: "ManifestRulesEngine",
			Run:  func(_ context.Context, rc *ResourceContext) error { return rc.applyRulesEngine() },
		},
		{
			Name: "ManifestPurge",
			Run:  func(_ context.Context, rc *ResourceContext) error { return rc.applyPurge() },
		},
		{
			Name: "ManifestDeleteOrphanedResources",
			Run:  func(_ context.Context, rc *ResourceContext) error { return rc.deleteOrphanedResources() },
		},
	}
}

func (rc *ResourceContext) applyDomain() error {
	if rc.conf.Domain.Id > 0 {
		requestUpdate := makeDomainUpdateRequest(rc.manifest.Domain, rc.conf)
		updated, err := rc.clientDomain.Update(rc.ctx, requestUpdate)
		if err != nil {
			return fmt.Errorf("%w - '%s': %s", msg.ErrorUpdateDomain, *requestUpdate.Name, err.Error())
		}
		rc.conf.Domain.Name = updated.GetName()
		rc.conf.Domain.DomainName = updated.GetDomainName()
		rc.conf.Domain.Url = utils.Concat("https://", updated.GetDomainName())
		return nil
	}

	requestCreate := makeDomainCreateRequest(rc.manifest.Domain, rc.conf)
	created, err := rc.clientDomain.Create(rc.ctx, requestCreate)
	if err != nil {
		return fmt.Errorf("%w - '%s': %s", msg.ErrorUpdateDomain, requestCreate.Name, err.Error())
	}
	rc.conf.Domain.Name = created.GetName()
	rc.conf.Domain.DomainName = created.GetDomainName()
	rc.conf.Domain.Url = utils.Concat("https://", created.GetDomainName())
	rc.conf.Domain.Id = created.GetId()
	return nil
}

func (rc *ResourceContext) applyOrigins() error {
	originConf := []contracts.AzionJsonDataOrigin{}
	for _, origin := range rc.manifest.Origins {
		if id := OriginIds[origin.Name]; id > 0 {
			requestUpdate := makeOriginUpdateRequest(origin, rc.conf)
			if origin.Name != "" {
				requestUpdate.Name = &origin.Name
			} else {
				requestUpdate.Name = &rc.conf.Name
			}
			updated, err := rc.clientOrigin.Update(rc.ctx, rc.conf.Application.ID, OriginKeys[origin.Name], requestUpdate)
			if errors.Is(err, utils.ErrorNotFound404) {
				logger.Debug("Origin not found. Skipping update", zap.Any("Error", err))
				continue
			}
			if err != nil {
				return fmt.Errorf("%w - '%s': %s", msg.ErrorUpdateOrigin, origin.Name, err.Error())
			}

			newEntry := contracts.AzionJsonDataOrigin{
				OriginId:  updated.GetOriginId(),
				OriginKey: updated.GetOriginKey(),
				Name:      updated.GetName(),
			}
			originConf = append(originConf, newEntry)

			rc.print(fmt.Sprintf(msg.ManifestUpdateOrigin, origin.Name, updated.GetOriginKey()))
		} else {
			requestCreate := makeOriginCreateRequest(origin, rc.conf)
			if origin.Name != "" {
				requestCreate.Name = origin.Name
			} else {
				requestCreate.Name = rc.conf.Name
			}
			created, err := rc.clientOrigin.Create(rc.ctx, rc.conf.Application.ID, requestCreate)
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
			rc.print(fmt.Sprintf(msg.ManifestCreateOrigin, origin.Name, created.GetOriginId()))
		}
	}

	rc.conf.Origin = originConf
	if err := rc.writeConfig(rc.conf, rc.projectConf); err != nil {
		logger.Debug("Error while writing azion.json file", zap.Error(err))
		return err
	}
	return nil
}

func (rc *ResourceContext) applyCacheSettings() error {
	cacheConf := []contracts.AzionJsonDataCacheSettings{}
	for _, cache := range rc.manifest.CacheSettings {
		if id := CacheIds[*cache.Name]; id > 0 {
			requestUpdate := makeCacheRequestUpdate(cache)
			if cache.Name != nil {
				requestUpdate.Name = cache.Name
			} else {
				requestUpdate.Name = &rc.conf.Name
			}
			updated, err := rc.clientCache.Update(rc.ctx, requestUpdate, rc.conf.Application.ID, id)
			if errors.Is(err, utils.ErrorNotFound404) {
				logger.Debug("Cache Setting not found. Skipping update", zap.Any("Error", err))
				continue
			}
			if err != nil {
				return fmt.Errorf("%w - '%s': %s", msg.ErrorUpdateCache, *cache.Name, err.Error())
			}
			newCache := contracts.AzionJsonDataCacheSettings{
				Id:   updated.GetId(),
				Name: updated.GetName(),
			}
			cacheConf = append(cacheConf, newCache)
			rc.print(fmt.Sprintf(msg.ManifestUpdateCache, *cache.Name, id))
		} else {
			requestCreate := makeCacheRequestCreate(cache)
			if cache.Name != nil {
				requestCreate.Name = *cache.Name
			} else {
				requestCreate.Name = rc.conf.Name + thoth.GenerateName()
			}
			created, err := rc.clientCache.Create(rc.ctx, requestCreate, rc.conf.Application.ID)
			if err != nil {
				return fmt.Errorf("%w - '%s': %s", msg.ErrorCreateCache, requestCreate.Name, err.Error())
			}
			newCache := contracts.AzionJsonDataCacheSettings{
				Id:   created.GetId(),
				Name: created.GetName(),
			}
			cacheConf = append(cacheConf, newCache)
			CacheIds[newCache.Name] = newCache.Id
			rc.print(fmt.Sprintf(msg.ManifestCreateCache, *cache.Name, newCache.Id))
		}
	}

	// Backup the cache ids: request.go reads CacheIdsBackup while building rule
	// requests, after this step has deleted consumed entries from CacheIds.
	for k, v := range CacheIds {
		CacheIdsBackup[k] = v
	}

	rc.conf.CacheSettings = cacheConf
	if err := rc.writeConfig(rc.conf, rc.projectConf); err != nil {
		logger.Debug("Error while writing azion.json file", zap.Error(err))
		return err
	}
	return nil
}

func (rc *ResourceContext) applyRulesEngine() error {
	ruleConf := []contracts.AzionJsonDataRules{}
	for _, rule := range rc.manifest.Rules {
		if r := RuleIds[rule.Name]; r.Id > 0 {
			requestUpdate, err := makeRuleRequestUpdate(rule, rc.conf)
			if err != nil {
				return err
			}
			requestUpdate.Id = r.Id
			requestUpdate.Phase = rule.Phase
			requestUpdate.IsActive = &rule.IsActive
			requestUpdate.Order = &rule.Order
			requestUpdate.IdApplication = rc.conf.Application.ID
			updated, err := rc.client.UpdateRulesEngine(rc.ctx, requestUpdate)
			if errors.Is(err, utils.ErrorNotFound404) {
				logger.Debug("Rule not found. Skipping update", zap.Any("Error", err))
				continue
			}
			if err != nil {
				return fmt.Errorf("%w - '%s': %s", msg.ErrorUpdateRule, rule.Name, err.Error())
			}
			newRule := contracts.AzionJsonDataRules{
				Id:    updated.GetId(),
				Name:  updated.GetName(),
				Phase: updated.GetPhase(),
			}
			rc.print(fmt.Sprintf(msg.ManifestUpdateRule, newRule.Name, newRule.Id))
			ruleConf = append(ruleConf, newRule)
			delete(RuleIds, rule.Name)
		} else {
			requestCreate, err := makeRuleRequestCreate(rule, rc.conf, rc.client, rc.ctx)
			if err != nil {
				return err
			}
			if rule.Name != "" {
				requestCreate.Name = rule.Name
			} else {
				requestCreate.Name = rc.conf.Name + thoth.GenerateName()
			}
			requestCreate.IsActive = &rule.IsActive
			requestCreate.Order = &rule.Order
			created, err := rc.client.CreateRulesEngine(rc.ctx, rc.conf.Application.ID, rule.Phase, requestCreate)
			if err != nil {
				return fmt.Errorf("%w - '%s': %s", msg.ErrorCreateRule, requestCreate.Name, err.Error())
			}
			newRule := contracts.AzionJsonDataRules{
				Id:    created.GetId(),
				Name:  created.GetName(),
				Phase: created.GetPhase(),
			}
			ruleConf = append(ruleConf, newRule)
			rc.print(fmt.Sprintf(msg.ManifestCreateRule, newRule.Name, newRule.Id))
		}
	}

	rc.conf.RulesEngine.Rules = ruleConf
	if err := rc.writeConfig(rc.conf, rc.projectConf); err != nil {
		logger.Debug("Error while writing azion.json file", zap.Error(err))
		return err
	}
	return nil
}

func (rc *ResourceContext) applyPurge() error {
	purgeCmd := purge.NewPurgeCmd(rc.factory)
	for _, purgeObj := range rc.manifest.Purge {
		switch purgeObj.Type {
		case "url":
			if err := purgeCmd.PurgeUrls(purgeObj.Urls, rc.factory); err != nil {
				logger.Debug("Error while purging urls", zap.Error(err))
				return errPurgeAborted
			}
		case "cachekey":
			if err := purgeCmd.PurgeCacheKeys(purgeObj.Urls, rc.factory, purgeCmd.Layer); err != nil {
				logger.Debug("Error while purging cache keys", zap.Error(err))
				return errPurgeAborted
			}
		case "wildcard":
			if err := purgeCmd.PurgeWildcard(purgeObj.Urls, rc.factory); err != nil {
				logger.Debug("Error while purging wildcards", zap.Error(err))
				return errPurgeAborted
			}
		}
	}
	return nil
}

// deleteOrphanedResources removes what the manifest no longer declares: the
// rules, origins and cache settings still left in the package-level maps after
// the steps above consumed the ones they used.
func (rc *ResourceContext) deleteOrphanedResources() error {
	if rc.conf.SkipDeletion != nil && *rc.conf.SkipDeletion {
		rc.print(msg.SkipDeletion)
		return nil
	}

	for _, value := range RuleIds {
		// since until [UXE-3599] was carried out we'd only cared about "request" phase, this check guarantees that if Phase is empty
		// we are probably dealing with a rule engine from a previous version
		phase := "request"
		if value.Phase != "" {
			phase = value.Phase
		}
		if err := rc.client.DeleteRulesEngine(rc.ctx, rc.conf.Application.ID, phase, value.Id); err != nil {
			return err
		}
		rc.print(fmt.Sprintf(msgrule.DeleteOutputSuccess+"\n", value.Id))
	}

	for i, value := range OriginKeys {
		if strings.Contains(i, "_single") {
			continue
		}
		if err := rc.clientOrigin.DeleteOrigins(rc.ctx, rc.conf.Application.ID, value); err != nil {
			return err
		}
		rc.print(fmt.Sprintf(msgorigin.DeleteOutputSuccess+"\n", value))
	}

	for _, value := range CacheIds {
		if err := rc.clientCache.Delete(rc.ctx, rc.conf.Application.ID, value); err != nil {
			return err
		}
		rc.print(fmt.Sprintf(msgcache.DeleteOutputSuccess+"\n", value))
	}

	return nil
}
