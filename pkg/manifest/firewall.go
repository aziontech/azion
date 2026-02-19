package manifest

import (
	"context"
	"fmt"

	msg "github.com/aziontech/azion-cli/messages/manifest"
	apiFirewall "github.com/aziontech/azion-cli/pkg/api/firewall"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/logger"
	edgesdk "github.com/aziontech/azionapi-v4-go-sdk-dev/edge-api"
	"go.uber.org/zap"
)

var (
	FirewallIds     map[string]int64
	FirewallRuleIds map[string]firewallRuleIdStruct
)

type firewallRuleIdStruct struct {
	FirewallId int64
	RuleId     int64
}

func transformFirewallRequestCreate(fw contracts.FirewallManifest) *apiFirewall.CreateRequest {
	request := apiFirewall.NewCreateRequest()
	request.SetName(fw.Name)
	if fw.Active != nil {
		request.SetActive(*fw.Active)
	}
	if fw.Debug != nil {
		request.SetDebug(*fw.Debug)
	}
	if fw.Modules != nil {
		request.SetModules(*fw.Modules)
	}
	return request
}

func transformFirewallRequestUpdate(fw contracts.FirewallManifest) *apiFirewall.UpdateRequest {
	request := apiFirewall.NewUpdateRequest()
	request.SetName(fw.Name)
	if fw.Active != nil {
		request.SetActive(*fw.Active)
	}
	if fw.Debug != nil {
		request.SetDebug(*fw.Debug)
	}
	if fw.Modules != nil {
		request.SetModules(*fw.Modules)
	}
	return request
}

func transformFirewallRuleRequestUpdate(rule edgesdk.FirewallRuleRequest) edgesdk.PatchedFirewallRuleRequest {
	request := edgesdk.PatchedFirewallRuleRequest{}
	request.SetName(rule.Name)
	if rule.Active != nil {
		request.SetActive(*rule.Active)
	}
	if rule.Description != nil {
		request.SetDescription(*rule.Description)
	}
	request.SetCriteria(rule.Criteria)
	request.SetBehaviors(rule.Behaviors)
	return request
}

func (man *ManifestInterpreter) CreateFirewalls(
	ctx context.Context,
	f *cmdutil.Factory,
	conf *contracts.AzionApplicationOptions,
	manifest *contracts.ManifestV4,
	projectConf string,
	msgs *[]string,
) error {
	if len(manifest.Firewalls) == 0 {
		return nil
	}

	firewallClient := apiFirewall.NewClient(f.HttpClient, f.Config.GetString("api_v4_url"), f.Config.GetString("token"))

	FirewallIds = make(map[string]int64)
	FirewallRuleIds = make(map[string]firewallRuleIdStruct)

	// Load existing firewall IDs from config
	for _, fwConf := range conf.Firewalls {
		FirewallIds[fwConf.Name] = fwConf.Id
		for _, ruleConf := range fwConf.Rules {
			FirewallRuleIds[ruleConf.Name] = firewallRuleIdStruct{
				FirewallId: fwConf.Id,
				RuleId:     ruleConf.Id,
			}
		}
	}

	firewallConf := []contracts.AzionJsonDataFirewall{}

	for _, fwMan := range manifest.Firewalls {
		var firewallId int64

		if id := FirewallIds[fwMan.Name]; id > 0 {
			// Update existing firewall
			updateReq := transformFirewallRequestUpdate(fwMan)
			updated, err := firewallClient.Update(ctx, updateReq, id)
			if err != nil {
				logger.Debug("Error while updating firewall", zap.Error(err))
				return err
			}
			firewallId = updated.GetId()
			msgf := fmt.Sprintf(msg.ManifestUpdateFirewall, updated.GetName(), updated.GetId())
			logger.FInfoFlags(f.IOStreams.Out, msgf, f.Format, f.Out)
			*msgs = append(*msgs, msgf)
		} else {
			// Create new firewall
			createReq := transformFirewallRequestCreate(fwMan)
			created, err := firewallClient.Create(ctx, createReq)
			if err != nil {
				logger.Debug("Error while creating firewall", zap.Error(err))
				return err
			}
			firewallId = created.GetId()
			msgf := fmt.Sprintf(msg.ManifestCreateFirewall, created.GetName(), created.GetId())
			logger.FInfoFlags(f.IOStreams.Out, msgf, f.Format, f.Out)
			*msgs = append(*msgs, msgf)
		}

		fwRuleConf := []contracts.AzionJsonDataFirewallRule{}
		for _, rule := range fwMan.RulesEngine {
			if ruleRef := FirewallRuleIds[rule.Name]; ruleRef.RuleId > 0 {
				updateReq := transformFirewallRuleRequestUpdate(rule)
				updated, err := firewallClient.UpdateRule(ctx, firewallId, ruleRef.RuleId, updateReq)
				if err != nil {
					logger.Debug("Error while updating firewall rule", zap.Error(err))
					return err
				}
				fwRuleConf = append(fwRuleConf, contracts.AzionJsonDataFirewallRule{
					Id:   updated.GetId(),
					Name: updated.GetName(),
				})
				msgf := fmt.Sprintf(msg.ManifestUpdateFirewallRule, updated.GetName(), updated.GetId())
				logger.FInfoFlags(f.IOStreams.Out, msgf, f.Format, f.Out)
				*msgs = append(*msgs, msgf)
			} else {
				created, err := firewallClient.CreateRule(ctx, firewallId, rule)
				if err != nil {
					logger.Debug("Error while creating firewall rule", zap.Error(err))
					return err
				}
				fwRuleConf = append(fwRuleConf, contracts.AzionJsonDataFirewallRule{
					Id:   created.GetId(),
					Name: created.GetName(),
				})
				msgf := fmt.Sprintf(msg.ManifestCreateFirewallRule, created.GetName(), created.GetId())
				logger.FInfoFlags(f.IOStreams.Out, msgf, f.Format, f.Out)
				*msgs = append(*msgs, msgf)
			}
		}

		firewallConf = append(firewallConf, contracts.AzionJsonDataFirewall{
			Id:    firewallId,
			Name:  fwMan.Name,
			Rules: fwRuleConf,
		})
	}

	conf.Firewalls = firewallConf
	err := man.WriteAzionJsonContent(conf, projectConf)
	if err != nil {
		logger.Debug("Error while writing azion.json file", zap.Error(err))
		return err
	}

	return nil
}
