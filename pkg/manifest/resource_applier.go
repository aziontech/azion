package manifest

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"

	msg "github.com/aziontech/azion-cli/messages/manifest"
	apiApplications "github.com/aziontech/azion-cli/pkg/api/applications"
	apiCache "github.com/aziontech/azion-cli/pkg/api/cache_setting"
	apiConnector "github.com/aziontech/azion-cli/pkg/api/connector"
	apiFirewall "github.com/aziontech/azion-cli/pkg/api/firewall"
	functionsApi "github.com/aziontech/azion-cli/pkg/api/function"
	apiPurge "github.com/aziontech/azion-cli/pkg/api/realtime_purge"
	apiWorkloads "github.com/aziontech/azion-cli/pkg/api/workloads"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/utils"
	edgesdk "github.com/aziontech/azionapi-v4-go-sdk-dev/azion-api"
	"go.uber.org/zap"
)

type ResourceContext struct {
	Ctx             context.Context
	Factory         *cmdutil.Factory
	Conf            *contracts.AzionApplicationOptions
	Manifest        *contracts.ManifestV4
	ProjectConf     string
	Msgs            *[]string
	WriteConfigFunc func(conf *contracts.AzionApplicationOptions, confPath string) error

	// API Clients
	ApplicationClient *apiApplications.Client
	CacheClient       *apiCache.ClientV4
	WorkloadClient    *apiWorkloads.Client
	ConnectorClient   *apiConnector.Client
	FunctionClient    *functionsApi.Client
	PurgeClient       *apiPurge.Client
	FirewallClient    *apiFirewall.Client

	// ID Mappings - these track created/existing resource IDs
	CacheIds        map[string]int64
	CacheIdsBackup  map[string]int64
	RuleIds         map[string]contracts.RuleIdsStruct
	ConnectorIds    map[string]int64
	DeploymentIds   map[string]int64
	FunctionIds     map[string]contracts.AzionJsonDataFunction
	FirewallIds     map[string]int64
	FirewallRuleIds map[string]firewallRuleIdRef
}

type firewallRuleIdRef struct {
	FirewallId int64
	RuleId     int64
}

func NewResourceContext(
	f *cmdutil.Factory,
	conf *contracts.AzionApplicationOptions,
	manifest *contracts.ManifestV4,
	projectConf string,
	msgs *[]string,
	writeConfigFunc func(conf *contracts.AzionApplicationOptions, confPath string) error,
) *ResourceContext {
	ctx := context.Background()
	apiURL := f.Config.GetString("api_v4_url")
	token := f.Config.GetString("token")

	rc := &ResourceContext{
		Ctx:             ctx,
		Factory:         f,
		Conf:            conf,
		Manifest:        manifest,
		ProjectConf:     projectConf,
		Msgs:            msgs,
		WriteConfigFunc: writeConfigFunc,

		// Initialize clients
		ApplicationClient: apiApplications.NewClient(f.HttpClient, apiURL, token),
		CacheClient:       apiCache.NewClientV4(f.HttpClient, apiURL, token),
		WorkloadClient:    apiWorkloads.NewClient(f.HttpClient, apiURL, token),
		ConnectorClient:   apiConnector.NewClient(f.HttpClient, apiURL, token),
		FunctionClient:    functionsApi.NewClient(f.HttpClient, apiURL, token),
		PurgeClient:       apiPurge.NewClient(f.HttpClient, apiURL, token),
		FirewallClient:    apiFirewall.NewClient(f.HttpClient, apiURL, token),

		// Initialize ID maps
		CacheIds:        make(map[string]int64),
		CacheIdsBackup:  make(map[string]int64),
		RuleIds:         make(map[string]contracts.RuleIdsStruct),
		ConnectorIds:    make(map[string]int64),
		DeploymentIds:   make(map[string]int64),
		FunctionIds:     make(map[string]contracts.AzionJsonDataFunction),
		FirewallIds:     make(map[string]int64),
		FirewallRuleIds: make(map[string]firewallRuleIdRef),
	}

	// Populate ID maps from existing config
	rc.populateIdMapsFromConfig()

	return rc
}

func (rc *ResourceContext) populateIdMapsFromConfig() {
	for _, cacheConf := range rc.Conf.CacheSettings {
		rc.CacheIds[cacheConf.Name] = cacheConf.Id
	}

	for _, funcConf := range rc.Conf.Function {
		rc.FunctionIds[funcConf.Name] = funcConf
	}

	for _, deploymentConf := range rc.Conf.Workloads.Deployments {
		rc.DeploymentIds[deploymentConf.Name] = deploymentConf.Id
	}

	for _, ruleConf := range rc.Conf.RulesEngine.Rules {
		rc.RuleIds[ruleConf.Name] = contracts.RuleIdsStruct{
			Id:    ruleConf.Id,
			Phase: ruleConf.Phase,
		}
	}

	for _, connectorConf := range rc.Conf.Connectors {
		rc.ConnectorIds[connectorConf.Name] = connectorConf.Id
	}

	for _, fwConf := range rc.Conf.Firewalls {
		rc.FirewallIds[fwConf.Name] = fwConf.Id
		for _, ruleConf := range fwConf.Rules {
			rc.FirewallRuleIds[ruleConf.Name] = firewallRuleIdRef{
				FirewallId: fwConf.Id,
				RuleId:     ruleConf.Id,
			}
		}
	}
}

func (rc *ResourceContext) WriteConfig() error {
	err := rc.WriteConfigFunc(rc.Conf, rc.ProjectConf)
	if err != nil {
		logger.Debug("Error while writing azion.json file", zap.Error(err))
		return err
	}
	return nil
}

func (rc *ResourceContext) ApplyFunctions(functions []contracts.Function) error {
	for _, funcMan := range functions {
		code, err := os.ReadFile(path.Join(".edge", funcMan.Path))
		if err != nil {
			return fmt.Errorf(msg.ErrorReadCodeFile.Error(), err)
		}

		if funcConf := rc.FunctionIds[funcMan.Name]; funcConf.ID > 0 {
			request := functionsApi.UpdateRequest{}
			request.SetActive(true)
			request.SetDefaultArgs(funcMan.DefaultArgs)
			request.SetName(funcMan.Name)
			request.SetCode(string(code))
			updated, err := rc.FunctionClient.Update(rc.Ctx, &request, funcConf.ID)
			if err != nil {
				return err
			}
			msgf := fmt.Sprintf(msg.ManifestUpdateFunction, updated.GetName(), updated.GetId())
			logger.FInfoFlags(rc.Factory.IOStreams.Out, msgf, rc.Factory.Format, rc.Factory.Out)
			*rc.Msgs = append(*rc.Msgs, msgf)
		} else {
			request := functionsApi.CreateRequest{}
			request.SetActive(true)
			request.SetDefaultArgs(funcMan.DefaultArgs)
			request.SetName(funcMan.Name)
			request.SetCode(string(code))
			resp, err := rc.FunctionClient.Create(rc.Ctx, &request)
			if err != nil {
				return err
			}
			newFunc := contracts.AzionJsonDataFunction{
				ID:   resp.GetId(),
				Name: resp.GetName(),
				File: funcMan.Argument,
				Args: "./azion/args.json",
			}
			rc.FunctionIds[funcMan.Name] = newFunc
			rc.Conf.Function = append(rc.Conf.Function, newFunc)
			msgf := fmt.Sprintf(msg.ManifestCreateFunction, resp.GetName(), resp.GetId())
			logger.FInfoFlags(rc.Factory.IOStreams.Out, msgf, rc.Factory.Format, rc.Factory.Out)
			*rc.Msgs = append(*rc.Msgs, msgf)
		}
	}

	return rc.WriteConfig()
}

func (rc *ResourceContext) ApplyFunctionInstances(instances []contracts.FunctionInstance) error {
	if rc.Conf.Application.ID == 0 {
		return msg.ErrorApplicationIDRequired
	}

	for _, funcMan := range instances {
		if funcConf := rc.FunctionIds[funcMan.Function]; funcConf.InstanceID > 0 {
			request := apiApplications.UpdateInstanceRequest{}
			request.SetActive(funcMan.Active)
			request.SetFunction(funcConf.ID)
			if len(funcMan.Args) > 0 {
				request.SetArgs(funcMan.Args)
			} else {
				args, err := unmarshalJsonArgs(funcConf.Args)
				if err != nil {
					return err
				}
				request.SetArgs(args)
			}
			request.SetName(funcMan.Name)
			updated, err := rc.ApplicationClient.UpdateInstance(rc.Ctx, &request, rc.Conf.Application.ID, funcConf.InstanceID)
			if err != nil {
				return err
			}
			msgf := fmt.Sprintf(msg.ManifestUpdateFunctionInstance, updated.GetName(), updated.GetId())
			logger.FInfoFlags(rc.Factory.IOStreams.Out, msgf, rc.Factory.Format, rc.Factory.Out)
			*rc.Msgs = append(*rc.Msgs, msgf)
		} else {
			request := apiApplications.CreateInstanceRequest{}
			request.SetActive(true)
			if len(funcMan.Args) > 0 {
				request.SetArgs(funcMan.Args)
			} else {
				args, err := unmarshalJsonArgs(funcConf.Args)
				if err != nil {
					return err
				}
				request.SetArgs(args)
			}
			request.SetName(funcMan.Name)
			request.SetFunction(funcConf.ID)
			resp, err := rc.ApplicationClient.CreateFuncInstances(rc.Ctx, &request, rc.Conf.Application.ID)
			if err != nil {
				return err
			}
			newFunc := contracts.AzionJsonDataFunction{
				ID:         funcConf.ID,
				CacheId:    funcConf.CacheId,
				Name:       funcConf.Name,
				File:       funcConf.File,
				Args:       funcConf.Args,
				InstanceID: resp.GetId(),
			}
			rc.FunctionIds[funcConf.Name] = newFunc
			msgf := fmt.Sprintf(msg.ManifestCreateFunctionInstance, resp.GetName(), resp.GetId())
			logger.FInfoFlags(rc.Factory.IOStreams.Out, msgf, rc.Factory.Format, rc.Factory.Out)
			*rc.Msgs = append(*rc.Msgs, msgf)
		}
	}

	// Update config with all function data
	funcsToWrite := []contracts.AzionJsonDataFunction{}
	for _, funcs := range rc.FunctionIds {
		funcsToWrite = append(funcsToWrite, funcs)
	}
	rc.Conf.Function = funcsToWrite

	return rc.WriteConfig()
}

func (rc *ResourceContext) ApplyEdgeApplication(app contracts.Applications) error {
	if rc.Conf.Application.ID > 0 {
		req := transformEdgeApplicationRequestUpdate(app)
		req.Id = rc.Conf.Application.ID
		updated, err := rc.ApplicationClient.Update(rc.Ctx, req)
		if err != nil {
			return err
		}
		msgf := fmt.Sprintf(msg.ManifestUpdateEdgeApplication, updated.GetName(), updated.GetId())
		logger.FInfoFlags(rc.Factory.IOStreams.Out, msgf, rc.Factory.Format, rc.Factory.Out)
		*rc.Msgs = append(*rc.Msgs, msgf)
	} else {
		createreq := transformEdgeApplicationRequestCreate(app)
		resp, err := rc.ApplicationClient.Create(rc.Ctx, createreq)
		if err != nil {
			return err
		}
		rc.Conf.Application.ID = resp.GetId()
		msgf := fmt.Sprintf(msg.ManifestCreateEdgeApplication, resp.GetName(), resp.GetId())
		logger.FInfoFlags(rc.Factory.IOStreams.Out, msgf, rc.Factory.Format, rc.Factory.Out)
		*rc.Msgs = append(*rc.Msgs, msgf)
	}

	return rc.WriteConfig()
}

func (rc *ResourceContext) ApplyCacheSettings(cacheSettings []contracts.ManifestCacheSetting) error {
	if rc.Conf.Application.ID == 0 {
		return msg.ErrorApplicationIDRequired
	}

	CacheConf := []contracts.AzionJsonDataCacheSettings{}
	for _, cache := range cacheSettings {
		if r := rc.CacheIds[cache.Name]; r > 0 {
			updated, err := rc.updateCache(cache, r)
			if err != nil {
				return err
			}
			CacheConf = append(CacheConf, updated)
		} else {
			newCache, err := rc.createCache(cache)
			if err != nil {
				return err
			}
			CacheConf = append(CacheConf, newCache)
			rc.CacheIds[newCache.Name] = newCache.Id
		}
	}

	// Backup cache IDs for behavior transformations
	for k, v := range rc.CacheIds {
		rc.CacheIdsBackup[k] = v
	}

	rc.Conf.CacheSettings = CacheConf
	return rc.WriteConfig()
}

func (rc *ResourceContext) updateCache(cache contracts.ManifestCacheSetting, cacheId int64) (contracts.AzionJsonDataCacheSettings, error) {
	request := transformCacheRequest(cache)
	updated, err := rc.CacheClient.Update(rc.Ctx, request, rc.Conf.Application.ID, cacheId)
	if errors.Is(err, utils.ErrorNotFound404) {
		logger.Debug("Cache Setting not found. Trying to create", zap.Any("Error", err))
		logger.FInfoFlags(rc.Factory.IOStreams.Out, msg.MessageDeleteResource+"\n", rc.Factory.Format, rc.Factory.Out)
		return rc.createCache(cache)
	}
	if err != nil {
		return contracts.AzionJsonDataCacheSettings{}, err
	}
	newCache := contracts.AzionJsonDataCacheSettings{
		Id:   updated.GetData().Id,
		Name: updated.GetData().Name,
	}
	msgf := fmt.Sprintf(msg.ManifestUpdateCache, newCache.Name, newCache.Id)
	logger.FInfoFlags(rc.Factory.IOStreams.Out, msgf, rc.Factory.Format, rc.Factory.Out)
	*rc.Msgs = append(*rc.Msgs, msgf)
	return newCache, nil
}

func (rc *ResourceContext) createCache(cache contracts.ManifestCacheSetting) (contracts.AzionJsonDataCacheSettings, error) {
	request := transformCacheRequestCreate(cache)
	responseCache, err := rc.CacheClient.Create(rc.Ctx, request.CacheSettingRequest, rc.Conf.Application.ID)
	if err != nil {
		return contracts.AzionJsonDataCacheSettings{}, err
	}
	newCache := contracts.AzionJsonDataCacheSettings{
		Id:   responseCache.GetId(),
		Name: responseCache.GetName(),
	}
	rc.CacheIds[newCache.Name] = newCache.Id
	msgf := fmt.Sprintf(msg.ManifestCreateCache, newCache.Name, newCache.Id)
	logger.FInfoFlags(rc.Factory.IOStreams.Out, msgf, rc.Factory.Format, rc.Factory.Out)
	*rc.Msgs = append(*rc.Msgs, msgf)
	return newCache, nil
}

func (rc *ResourceContext) ApplyConnectors(connectors []edgesdk.ConnectorRequest) error {
	connectorConf := []contracts.AzionJsonDataConnectors{}

	for _, connector := range connectors {
		connName, connType := getConnectorName(connector, rc.Conf.Name)
		if id := rc.ConnectorIds[connName]; id > 0 {
			request := transformEdgeConnectorRequest(connector)
			connectorResp, err := rc.ConnectorClient.Update(rc.Ctx, request, id)
			if err != nil {
				return err
			}
			conn := contracts.AzionJsonDataConnectors{}
			switch connType {
			case "http":
				http := connectorResp.ConnectorHTTP
				conn.Id = http.GetId()
				conn.Name = http.GetName()
				conn.Address = http.Attributes.Addresses
			case "storage":
				storage := connectorResp.ConnectorBase
				conn.Id = storage.GetId()
				conn.Name = storage.GetName()
			default:
				return msg.ErrorConnectorTypeNotFound
			}
			msgf := fmt.Sprintf(msg.ManifestUpdateConnector, conn.Name, conn.Id)
			logger.FInfoFlags(rc.Factory.IOStreams.Out, msgf, rc.Factory.Format, rc.Factory.Out)
			*rc.Msgs = append(*rc.Msgs, msgf)
			connectorConf = append(connectorConf, conn)
		} else {
			request := apiConnector.CreateRequest{
				ConnectorRequest: connector,
			}
			connectorResp, err := rc.ConnectorClient.Create(rc.Ctx, &request)
			if err != nil {
				return err
			}
			conn := contracts.AzionJsonDataConnectors{}
			switch connType {
			case "http":
				http := connectorResp.ConnectorHTTP
				conn.Id = http.GetId()
				conn.Name = http.GetName()
			case "storage":
				storage := connectorResp.ConnectorBase
				conn.Id = storage.GetId()
				conn.Name = storage.GetName()
			default:
				return msg.ErrorConnectorTypeNotFound
			}
			rc.ConnectorIds[conn.Name] = conn.Id
			msgf := fmt.Sprintf(msg.ManifestCreateConnector, conn.Name, conn.Id)
			logger.FInfoFlags(rc.Factory.IOStreams.Out, msgf, rc.Factory.Format, rc.Factory.Out)
			*rc.Msgs = append(*rc.Msgs, msgf)
			connectorConf = append(connectorConf, conn)
		}
	}

	rc.Conf.Connectors = connectorConf
	return rc.WriteConfig()
}

func (rc *ResourceContext) ApplyRulesEngine(rules []contracts.ManifestRulesEngine) error {
	if rc.Conf.Application.ID == 0 {
		return msg.ErrorApplicationIDRequired
	}

	ruleConf := []contracts.AzionJsonDataRules{}

	for _, rule := range rules {
		if r := rc.RuleIds[rule.Rule.Name]; r.Id > 0 {
			newRule, err := rc.updateRule(rule, r)
			if err != nil {
				if errors.Is(err, utils.ErrorNotFound404) {
					logger.Debug("Rule not found. Skipping update", zap.Any("Error", err))
					delete(rc.RuleIds, rule.Rule.Name)
					continue
				}
				return err
			}
			msgf := fmt.Sprintf(msg.ManifestUpdateRule, newRule.Name, newRule.Id)
			logger.FInfoFlags(rc.Factory.IOStreams.Out, msgf, rc.Factory.Format, rc.Factory.Out)
			*rc.Msgs = append(*rc.Msgs, msgf)
			ruleConf = append(ruleConf, newRule)
			delete(rc.RuleIds, newRule.Name)
		} else {
			newRule, err := rc.createRule(rule)
			if err != nil {
				return err
			}
			msgf := fmt.Sprintf(msg.ManifestCreateRule, newRule.Name, newRule.Id)
			logger.FInfoFlags(rc.Factory.IOStreams.Out, msgf, rc.Factory.Format, rc.Factory.Out)
			*rc.Msgs = append(*rc.Msgs, msgf)
			ruleConf = append(ruleConf, newRule)
		}
	}

	rc.Conf.RulesEngine.Rules = ruleConf
	return rc.WriteConfig()
}

func (rc *ResourceContext) updateRule(rule contracts.ManifestRulesEngine, existing contracts.RuleIdsStruct) (contracts.AzionJsonDataRules, error) {
	switch rule.Phase {
	case "request":
		req := transformRuleRequest(rule.Rule)
		req.IdApplication = rc.Conf.Application.ID
		req.Id = existing.Id
		behs, err := rc.transformBehaviorsRequest(rule.Rule.Behaviors)
		if err != nil {
			return contracts.AzionJsonDataRules{}, err
		}
		req.Behaviors = behs
		updated, err := rc.ApplicationClient.UpdateRulesEngineRequest(rc.Ctx, req)
		if err != nil {
			return contracts.AzionJsonDataRules{}, err
		}
		return contracts.AzionJsonDataRules{
			Id:    updated.GetId(),
			Name:  updated.GetName(),
			Phase: rule.Phase,
		}, nil
	case "response":
		req := transformRuleResponse(rule.Rule)
		req.IdApplication = rc.Conf.Application.ID
		req.Id = existing.Id
		behs, err := rc.transformBehaviorsResponse(rule.Rule.Behaviors)
		if err != nil {
			return contracts.AzionJsonDataRules{}, err
		}
		req.Behaviors = behs
		updated, err := rc.ApplicationClient.UpdateRulesEngineResponse(rc.Ctx, req)
		if err != nil {
			return contracts.AzionJsonDataRules{}, err
		}
		return contracts.AzionJsonDataRules{
			Id:    updated.GetId(),
			Name:  updated.GetName(),
			Phase: rule.Phase,
		}, nil
	default:
		return contracts.AzionJsonDataRules{}, msg.ErrorInvalidPhase
	}
}

func (rc *ResourceContext) createRule(rule contracts.ManifestRulesEngine) (contracts.AzionJsonDataRules, error) {
	switch rule.Phase {
	case "request":
		req := &apiApplications.CreateRulesEngineRequest{}
		createRequest := transformRuleRequestCreate(rule.Rule)
		bh, err := rc.transformBehaviorsRequest(rule.Rule.Behaviors)
		if err != nil {
			return contracts.AzionJsonDataRules{}, err
		}
		req.RequestPhaseRuleRequest = createRequest
		req.Behaviors = bh
		created, err := rc.ApplicationClient.CreateRulesEngineRequest(rc.Ctx, rc.Conf.Application.ID, rule.Phase, req)
		if err != nil {
			return contracts.AzionJsonDataRules{}, err
		}
		return contracts.AzionJsonDataRules{
			Id:    created.GetId(),
			Name:  created.GetName(),
			Phase: rule.Phase,
		}, nil
	case "response":
		req := &apiApplications.CreateRulesEngineResponse{}
		createRequest := transformRuleResponseCreate(rule.Rule)
		bh, err := rc.transformBehaviorsResponse(rule.Rule.Behaviors)
		if err != nil {
			return contracts.AzionJsonDataRules{}, err
		}
		req.ResponsePhaseRuleRequest = createRequest
		req.Behaviors = bh
		created, err := rc.ApplicationClient.CreateRulesEngineResponse(rc.Ctx, rc.Conf.Application.ID, rule.Phase, req)
		if err != nil {
			return contracts.AzionJsonDataRules{}, err
		}
		return contracts.AzionJsonDataRules{
			Id:    created.GetId(),
			Name:  created.GetName(),
			Phase: rule.Phase,
		}, nil
	default:
		return contracts.AzionJsonDataRules{}, msg.ErrorInvalidPhase
	}
}

func (rc *ResourceContext) ApplyWorkloads(workloads []contracts.WorkloadManifest) error {
	if len(workloads) == 0 {
		return nil
	}

	workloadMan := workloads[0]
	if rc.Conf.Workloads.Id > 0 {
		request := transformWorkloadRequestUpdate(workloadMan)
		request.Id = rc.Conf.Workloads.Id
		updated, err := rc.WorkloadClient.Update(rc.Ctx, request)
		if err != nil {
			return err
		}
		rc.Conf.Workloads.Domains = updated.GetDomains()
		rc.Conf.Workloads.Url = utils.Concat("https://", updated.GetWorkloadDomain())
		msgf := fmt.Sprintf(msg.ManifestUpdateWorkload, updated.GetName(), updated.GetId())
		logger.FInfoFlags(rc.Factory.IOStreams.Out, msgf, rc.Factory.Format, rc.Factory.Out)
		*rc.Msgs = append(*rc.Msgs, msgf)
	} else {
		request := transformWorkloadRequestCreate(workloadMan, rc.Conf.Application.ID)
		resp, err := rc.WorkloadClient.Create(rc.Ctx, request)
		if err != nil {
			return err
		}
		rc.Conf.Workloads.Id = resp.GetId()
		rc.Conf.Workloads.Name = resp.GetName()
		rc.Conf.Workloads.Domains = resp.GetDomains()
		rc.Conf.Workloads.Url = utils.Concat("https://", resp.GetWorkloadDomain())
		msgf := fmt.Sprintf(msg.ManifestCreateWorkload, resp.GetName(), resp.GetId())
		logger.FInfoFlags(rc.Factory.IOStreams.Out, msgf, rc.Factory.Format, rc.Factory.Out)
		*rc.Msgs = append(*rc.Msgs, msgf)
	}

	return rc.WriteConfig()
}

func (rc *ResourceContext) ApplyWorkloadDeployments(deployments []contracts.WorkloadDeployment) error {
	if rc.Conf.Workloads.Id == 0 {
		return msg.ErrorWorkloadIDRequired
	}

	for _, deployment := range deployments {
		if id := rc.DeploymentIds[deployment.Name]; id > 0 {
			request := transformWorkloadDeploymentRequestUpdate(deployment, rc.Conf)
			updated, err := rc.WorkloadClient.UpdateDeployment(rc.Ctx, request, rc.Conf.Workloads.Id, id)
			if err != nil {
				return err
			}
			msgf := fmt.Sprintf(msg.ManifestUpdateWorkloadDeployment, updated.GetName(), updated.GetId())
			logger.FInfoFlags(rc.Factory.IOStreams.Out, msgf, rc.Factory.Format, rc.Factory.Out)
			*rc.Msgs = append(*rc.Msgs, msgf)
		} else {
			request := transformWorkloadDeploymentRequestCreate(deployment, rc.Conf)
			resp, err := rc.WorkloadClient.CreateDeployment(rc.Ctx, request, rc.Conf.Workloads.Id)
			if err != nil {
				return err
			}
			rc.Conf.Workloads.Deployments = append(rc.Conf.Workloads.Deployments, contracts.Deployments{
				Id:   resp.GetId(),
				Name: resp.GetName(),
			})
			msgf := fmt.Sprintf(msg.ManifestCreateWorkloadDeployment, resp.GetName(), resp.GetId())
			logger.FInfoFlags(rc.Factory.IOStreams.Out, msgf, rc.Factory.Format, rc.Factory.Out)
			*rc.Msgs = append(*rc.Msgs, msgf)
		}
	}

	return rc.WriteConfig()
}

func (rc *ResourceContext) ApplyFirewalls(firewalls []contracts.FirewallManifest) error {
	if len(firewalls) == 0 {
		return nil
	}

	firewallConf := []contracts.AzionJsonDataFirewall{}

	for _, fwMan := range firewalls {
		var firewallId int64

		if id := rc.FirewallIds[fwMan.Name]; id > 0 {
			updateReq := apiFirewall.NewUpdateRequest()
			updateReq.SetName(fwMan.Name)
			if fwMan.Active != nil {
				updateReq.SetActive(*fwMan.Active)
			}
			if fwMan.Debug != nil {
				updateReq.SetDebug(*fwMan.Debug)
			}
			if fwMan.Modules != nil {
				updateReq.SetModules(*fwMan.Modules)
			}
			updated, err := rc.FirewallClient.Update(rc.Ctx, updateReq, id)
			if err != nil {
				logger.Debug("Error while updating firewall", zap.Error(err))
				return err
			}
			firewallId = updated.GetId()
			msgf := fmt.Sprintf(msg.ManifestUpdateFirewall, updated.GetName(), updated.GetId())
			logger.FInfoFlags(rc.Factory.IOStreams.Out, msgf, rc.Factory.Format, rc.Factory.Out)
			*rc.Msgs = append(*rc.Msgs, msgf)
		} else {
			createReq := apiFirewall.NewCreateRequest()
			createReq.SetName(fwMan.Name)
			if fwMan.Active != nil {
				createReq.SetActive(*fwMan.Active)
			}
			if fwMan.Debug != nil {
				createReq.SetDebug(*fwMan.Debug)
			}
			if fwMan.Modules != nil {
				createReq.SetModules(*fwMan.Modules)
			}
			created, err := rc.FirewallClient.Create(rc.Ctx, createReq)
			if err != nil {
				logger.Debug("Error while creating firewall", zap.Error(err))
				return err
			}
			firewallId = created.GetId()
			msgf := fmt.Sprintf(msg.ManifestCreateFirewall, created.GetName(), created.GetId())
			logger.FInfoFlags(rc.Factory.IOStreams.Out, msgf, rc.Factory.Format, rc.Factory.Out)
			*rc.Msgs = append(*rc.Msgs, msgf)
		}

		fwRuleConf := []contracts.AzionJsonDataFirewallRule{}
		for _, rule := range fwMan.RulesEngine {
			if ruleRef := rc.FirewallRuleIds[rule.Name]; ruleRef.RuleId > 0 {
				patchReq := edgesdk.PatchedFirewallRuleRequest{}
				patchReq.SetName(rule.Name)
				if rule.Active != nil {
					patchReq.SetActive(*rule.Active)
				}
				if rule.Description != nil {
					patchReq.SetDescription(*rule.Description)
				}
				patchReq.SetCriteria(rule.Criteria)
				patchReq.SetBehaviors(rule.Behaviors)

				updated, err := rc.FirewallClient.UpdateRule(rc.Ctx, firewallId, ruleRef.RuleId, patchReq)
				if err != nil {
					logger.Debug("Error while updating firewall rule", zap.Error(err))
					return err
				}
				fwRuleConf = append(fwRuleConf, contracts.AzionJsonDataFirewallRule{
					Id:   updated.GetId(),
					Name: updated.GetName(),
				})
				msgf := fmt.Sprintf(msg.ManifestUpdateFirewallRule, updated.GetName(), updated.GetId())
				logger.FInfoFlags(rc.Factory.IOStreams.Out, msgf, rc.Factory.Format, rc.Factory.Out)
				*rc.Msgs = append(*rc.Msgs, msgf)
			} else {
				created, err := rc.FirewallClient.CreateRule(rc.Ctx, firewallId, rule)
				if err != nil {
					logger.Debug("Error while creating firewall rule", zap.Error(err))
					return err
				}
				fwRuleConf = append(fwRuleConf, contracts.AzionJsonDataFirewallRule{
					Id:   created.GetId(),
					Name: created.GetName(),
				})
				msgf := fmt.Sprintf(msg.ManifestCreateFirewallRule, created.GetName(), created.GetId())
				logger.FInfoFlags(rc.Factory.IOStreams.Out, msgf, rc.Factory.Format, rc.Factory.Out)
				*rc.Msgs = append(*rc.Msgs, msgf)
			}
		}

		firewallConf = append(firewallConf, contracts.AzionJsonDataFirewall{
			Id:    firewallId,
			Name:  fwMan.Name,
			Rules: fwRuleConf,
		})
	}

	rc.Conf.Firewalls = firewallConf
	return rc.WriteConfig()
}

func (rc *ResourceContext) ApplyPurge(purges []contracts.PurgeManifest) error {
	for _, purgeObj := range purges {
		err := rc.PurgeClient.PurgeCache(rc.Ctx, purgeObj.Items, purgeObj.Type, *purgeObj.Layer)
		if err != nil {
			logger.Debug("Error while purging domains", zap.Error(err))
			return err
		}
		msgf := fmt.Sprintf(msg.ManifestPurgeSuccess, purgeObj.Type)
		logger.FInfoFlags(rc.Factory.IOStreams.Out, msgf, rc.Factory.Format, rc.Factory.Out)
		*rc.Msgs = append(*rc.Msgs, msgf)
	}
	return nil
}

func (rc *ResourceContext) DeleteOrphanedResources() error {
	return deleteResources(rc.Ctx, rc.Factory, rc.Conf, rc.Msgs)
}

func (rc *ResourceContext) transformBehaviorsRequest(behaviors []contracts.ManifestRuleBehavior) ([]edgesdk.RequestPhaseBehaviorRequest, error) {
	behaviorsRequest := make([]edgesdk.RequestPhaseBehaviorRequest, 0, len(behaviors))
	for _, behavior := range behaviors {
		var withArgs edgesdk.BehaviorArgs
		var withoutArgs edgesdk.BehaviorNoArgs
		var captureMatchGroups edgesdk.BehaviorCapture
		var beh edgesdk.RequestPhaseBehaviorRequest
		switch behavior.Type {
		case "run_function":
			attributesJSON, err := json.Marshal(behavior.Attributes)
			if err != nil {
				return nil, err
			}
			var attributes edgesdk.BehaviorArgsAttributes
			if err := json.Unmarshal(attributesJSON, &attributes); err != nil {
				return nil, err
			}
			withArgs.SetType("run_function")
			if attributes.Value.Int64 != nil {
				withArgs.SetAttributes(attributes)
			} else if attributes.Value.String != nil {
				funcName := *attributes.Value.String
				if _, ok := rc.FunctionIds[funcName]; !ok {
					return nil, msg.ErrorFuncNotFound
				}
				funcToWorkWith := rc.FunctionIds[funcName]
				v := edgesdk.BehaviorArgsAttributesValue{
					Int64: &funcToWorkWith.InstanceID,
				}
				attributes.SetValue(v)
				withArgs.SetAttributes(attributes)
			}
			beh.BehaviorArgs = &withArgs
			behaviorsRequest = append(behaviorsRequest, beh)
		case "set_cache_policy":
			attributesJSON, err := json.Marshal(behavior.Attributes)
			if err != nil {
				return nil, err
			}
			var attributes edgesdk.BehaviorArgsAttributes
			if err := json.Unmarshal(attributesJSON, &attributes); err != nil {
				return nil, err
			}
			withArgs.SetType("set_cache_policy")
			if attributes.Value.Int64 != nil {
				withArgs.SetAttributes(attributes)
			} else if attributes.Value.String != nil {
				cacheName := *attributes.Value.String
				if id := rc.CacheIdsBackup[cacheName]; id > 0 {
					v := edgesdk.BehaviorArgsAttributesValue{
						Int64: &id,
					}
					attributes.SetValue(v)
					withArgs.SetAttributes(attributes)
					delete(rc.CacheIds, cacheName)
				} else {
					logger.Debug("Cache Setting not found", zap.Any("Target", *attributes.Value.String))
					return nil, msg.ErrorCacheNotFound
				}
			}
			beh.BehaviorArgs = &withArgs
			behaviorsRequest = append(behaviorsRequest, beh)
		case "set_connector":
			attributesJSON, err := json.Marshal(behavior.Attributes)
			if err != nil {
				return nil, err
			}
			var attributes edgesdk.BehaviorArgsAttributes
			if err := json.Unmarshal(attributesJSON, &attributes); err != nil {
				return nil, err
			}
			withArgs.SetType("set_connector")
			if attributes.Value.Int64 != nil {
				withArgs.SetAttributes(attributes)
			} else if attributes.Value.String != nil {
				connectorName := *attributes.Value.String
				if id := rc.ConnectorIds[connectorName]; id > 0 {
					v := edgesdk.BehaviorArgsAttributesValue{
						Int64: &id,
					}
					attributes.SetValue(v)
					withArgs.SetAttributes(attributes)
				} else {
					logger.Debug("Connector not found", zap.Any("Target", connectorName))
					return nil, msg.ErrorConnectorNotFound
				}
			}
			beh.BehaviorArgs = &withArgs
			behaviorsRequest = append(behaviorsRequest, beh)
		case "capture_match_groups":
			attributesJSON, err := json.Marshal(behavior.Attributes)
			if err != nil {
				return nil, err
			}
			var attributes edgesdk.BehaviorCaptureMatchGroupsAttributes
			if err := json.Unmarshal(attributesJSON, &attributes); err != nil {
				return nil, err
			}
			captureMatchGroups.SetType("capture_match_groups")
			captureMatchGroups.SetAttributes(attributes)
			beh.BehaviorCapture = &captureMatchGroups
			behaviorsRequest = append(behaviorsRequest, beh)
		case "redirect_to_301", "redirect_to_302", "filter_request_cookie", "rewrite_request", "add_request_header", "filter_request_header", "add_request_cookie":
			attributesJSON, err := json.Marshal(behavior.Attributes)
			if err != nil {
				return nil, err
			}
			var attributes edgesdk.BehaviorArgsAttributes
			if err := json.Unmarshal(attributesJSON, &attributes); err != nil {
				return nil, err
			}
			withArgs.SetType(behavior.Type)
			withArgs.SetAttributes(attributes)
			beh.BehaviorArgs = &withArgs
			behaviorsRequest = append(behaviorsRequest, beh)
		default:
			withoutArgs.SetType(behavior.Type)
			beh.BehaviorNoArgs = &withoutArgs
			behaviorsRequest = append(behaviorsRequest, beh)
		}
	}

	return behaviorsRequest, nil
}

func (rc *ResourceContext) transformBehaviorsResponse(behaviors []contracts.ManifestRuleBehavior) ([]edgesdk.ResponsePhaseBehaviorRequest, error) {
	behaviorsResponse := make([]edgesdk.ResponsePhaseBehaviorRequest, 0, len(behaviors))

	for _, behavior := range behaviors {
		var withArgs edgesdk.BehaviorArgs
		var withoutArgs edgesdk.BehaviorNoArgs
		var captureMatchGroups edgesdk.BehaviorCapture
		var beh edgesdk.ResponsePhaseBehaviorRequest

		switch behavior.Type {
		case "capture_match_groups":
			attributesJSON, err := json.Marshal(behavior.Attributes)
			if err != nil {
				return nil, err
			}
			var attributes edgesdk.BehaviorCaptureMatchGroupsAttributes
			if err := json.Unmarshal(attributesJSON, &attributes); err != nil {
				return nil, err
			}
			captureMatchGroups.SetType("capture_match_groups")
			captureMatchGroups.SetAttributes(attributes)
			beh.BehaviorCapture = &captureMatchGroups
			behaviorsResponse = append(behaviorsResponse, beh)
		case "enable_gzip", "deliver":
			withoutArgs.SetType(behavior.Type)
			beh.BehaviorNoArgs = &withoutArgs
			behaviorsResponse = append(behaviorsResponse, beh)
		default:
			attributesJSON, err := json.Marshal(behavior.Attributes)
			if err != nil {
				return nil, err
			}
			var attributes edgesdk.BehaviorArgsAttributes
			if err := json.Unmarshal(attributesJSON, &attributes); err != nil {
				return nil, err
			}
			withArgs.SetType(behavior.Type)
			withArgs.SetAttributes(attributes)
			beh.BehaviorArgs = &withArgs
			behaviorsResponse = append(behaviorsResponse, beh)
		}
	}

	return behaviorsResponse, nil
}
