package firewallrules

import (
	"context"
	"fmt"

	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/utils"
	sdk "github.com/aziontech/azionapi-v4-go-sdk-dev/azion-api"
	"go.uber.org/zap"
)

func (c *Client) List(ctx context.Context, opts *contracts.ListOptions, firewallID int64) (*sdk.PaginatedFirewallRuleList, error) {
	logger.Debug("List Firewall Rules")

	req := c.apiClient.FirewallsRulesEngineAPI.
		ListFirewallRules(ctx, firewallID).
		Page(opts.Page).
		PageSize(opts.PageSize)
	resp, httpResp, err := req.Execute()

	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while listing Firewall Rules", zap.Error(err))
			errBody, err = utils.LogAndRewindBodyV4(httpResp)
			if err != nil {
				return nil, err
			}
		}
		return nil, utils.ErrorPerStatusCodeV4(errBody, httpResp, err)
	}

	return resp, nil
}

func (c *Client) Get(ctx context.Context, firewallID, ruleID int64) (sdk.FirewallRule, error) {
	logger.Debug("Get Firewall Rule")
	req := c.apiClient.FirewallsRulesEngineAPI.RetrieveFirewallRule(ctx, firewallID, ruleID)
	resp, httpResp, err := req.Execute()

	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while getting Firewall Rule", zap.Error(err))
			errBody, err = utils.LogAndRewindBodyV4(httpResp)
			if err != nil {
				return sdk.FirewallRule{}, err
			}
		}
		return sdk.FirewallRule{}, utils.ErrorPerStatusCodeV4(errBody, httpResp, err)
	}

	return resp.Data, nil
}

func (c *Client) Delete(ctx context.Context, firewallID, ruleID int64) error {
	logger.Debug("Delete Firewall Rule")
	req := c.apiClient.FirewallsRulesEngineAPI.DeleteFirewallRule(ctx, firewallID, ruleID)
	_, httpResp, err := req.Execute()

	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while deleting Firewall Rule", zap.Error(err))
			errBody, err = utils.LogAndRewindBodyV4(httpResp)
			if err != nil {
				return err
			}
		}
		return utils.ErrorPerStatusCodeV4(errBody, httpResp, err)
	}

	return nil
}

func (c *Client) Create(ctx context.Context, firewallID int64, req *CreateRequest) (sdk.FirewallRule, error) {
	logger.Debug("Create Firewall Rule")
	request := c.apiClient.FirewallsRulesEngineAPI.CreateFirewallRule(ctx, firewallID).FirewallRuleRequest(req.FirewallRuleRequest)
	resp, httpResp, err := request.Execute()

	if err != nil {
		fmt.Println(err)
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while creating Firewall Rule", zap.Error(err))
			errBody, err = utils.LogAndRewindBodyV4(httpResp)
			fmt.Println(errBody)
			if err != nil {
				return sdk.FirewallRule{}, err
			}
		}
		return sdk.FirewallRule{}, utils.ErrorPerStatusCodeV4(errBody, httpResp, err)
	}

	return resp.Data, nil
}

func (c *Client) Update(ctx context.Context, firewallID, ruleID int64, req *UpdateRequest) (sdk.FirewallRule, error) {
	logger.Debug("Update Firewall Rule")
	request := c.apiClient.FirewallsRulesEngineAPI.PartialUpdateFirewallRule(ctx, firewallID, ruleID).PatchedFirewallRuleRequest(req.PatchedFirewallRuleRequest)
	resp, httpResp, err := request.Execute()

	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while updating Firewall Rule", zap.Error(err))
			errBody, err = utils.LogAndRewindBodyV4(httpResp)
			if err != nil {
				return sdk.FirewallRule{}, err
			}
		}
		return sdk.FirewallRule{}, utils.ErrorPerStatusCodeV4(errBody, httpResp, err)
	}

	return resp.Data, nil
}
