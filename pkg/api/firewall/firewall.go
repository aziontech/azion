package firewall

import (
	"context"

	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/utils"
	sdk "github.com/aziontech/azionapi-v4-go-sdk-dev/edge-api"
	"go.uber.org/zap"
)

func (c *Client) List(ctx context.Context, opts *contracts.ListOptions) (*sdk.PaginatedFirewallList, error) {
	logger.Debug("List firewalls")
	if opts.OrderBy == "" {
		opts.OrderBy = "id"
	}
	resp, httpResp, err := c.apiClient.FirewallsAPI.ListFirewalls(ctx).
		Ordering(opts.OrderBy).
		Page(opts.Page).
		PageSize(opts.PageSize).
		Search(opts.Sort).
		Execute()

	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while listing the firewalls", zap.Error(err))
			errBody, err = utils.LogAndRewindBodyV4(httpResp)
			if err != nil {
				return nil, err
			}
		}
		return nil, utils.ErrorPerStatusCodeV4(errBody, httpResp, err)
	}

	return resp, nil
}

func (c *Client) Get(ctx context.Context, id int64) (sdk.Firewall, error) {
	logger.Debug("Get Firewall")
	request := c.apiClient.FirewallsAPI.RetrieveFirewall(ctx, id)

	res, httpResp, err := request.Execute()
	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while getting a Firewall", zap.Error(err))
			errBody, err = utils.LogAndRewindBodyV4(httpResp)
			if err != nil {
				return sdk.Firewall{}, err
			}
		}
		return sdk.Firewall{}, utils.ErrorPerStatusCodeV4(errBody, httpResp, err)
	}

	return res.Data, nil
}

func (c *Client) Delete(ctx context.Context, id int64) error {
	logger.Debug("Delete Firewall")
	req := c.apiClient.FirewallsAPI.DeleteFirewall(ctx, id)

	_, httpResp, err := req.Execute()
	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while deleting a Firewall", zap.Error(err), zap.Any("ID", id))
			errBody, err = utils.LogAndRewindBodyV4(httpResp)
			if err != nil {
				return err
			}
		}

		return utils.ErrorPerStatusCodeV4(errBody, httpResp, err)
	}

	return nil
}

func (c *Client) Create(ctx context.Context, req *CreateRequest) (sdk.Firewall, error) {
	logger.Debug("Create Firewall")

	request := c.apiClient.FirewallsAPI.CreateFirewall(ctx).FirewallRequest(req.FirewallRequest)

	firewallResponse, httpResp, err := request.Execute()
	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while creating a Firewall", zap.Error(err), zap.Any("Name", req.Name))
			errBody, err = utils.LogAndRewindBodyV4(httpResp)
			if err != nil {
				return sdk.Firewall{}, err
			}
		}
		return sdk.Firewall{}, utils.ErrorPerStatusCodeV4(errBody, httpResp, err)
	}

	return firewallResponse.Data, nil
}

func (c *Client) Update(ctx context.Context, req *UpdateRequest, id int64) (sdk.Firewall, error) {
	logger.Debug("Update Firewall", zap.Any("Firewall ID", id), zap.Any("Firewall name", req.Name))
	request := c.apiClient.FirewallsAPI.PartialUpdateFirewall(ctx, id).PatchedFirewallRequest(req.PatchedFirewallRequest)

	firewallResponse, httpResp, err := request.Execute()
	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while updating a Firewall", zap.Error(err), zap.Any("ID", id), zap.Any("Name", req.Name))
			errBody, err = utils.LogAndRewindBodyV4(httpResp)
			if err != nil {
				return sdk.Firewall{}, err
			}
		}
		return sdk.Firewall{}, utils.ErrorPerStatusCodeV4(errBody, httpResp, err)
	}

	return firewallResponse.Data, nil
}

func (c *Client) CreateRule(ctx context.Context, firewallId int64, req sdk.FirewallRuleRequest) (sdk.FirewallRule, error) {
	logger.Debug("Create Firewall Rule", zap.Any("Firewall ID", firewallId), zap.Any("Rule name", req.Name))
	request := c.apiClient.FirewallsRulesEngineAPI.CreateFirewallRule(ctx, firewallId).FirewallRuleRequest(req)

	ruleResponse, httpResp, err := request.Execute()
	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while creating a Firewall Rule", zap.Error(err), zap.Any("Name", req.Name))
			errBody, err = utils.LogAndRewindBodyV4(httpResp)
			if err != nil {
				return sdk.FirewallRule{}, err
			}
		}
		return sdk.FirewallRule{}, utils.ErrorPerStatusCodeV4(errBody, httpResp, err)
	}

	return ruleResponse.Data, nil
}

func (c *Client) UpdateRule(ctx context.Context, firewallId int64, ruleId int64, req sdk.PatchedFirewallRuleRequest) (sdk.FirewallRule, error) {
	logger.Debug("Update Firewall Rule", zap.Any("Firewall ID", firewallId), zap.Any("Rule ID", ruleId))
	request := c.apiClient.FirewallsRulesEngineAPI.PartialUpdateFirewallRule(ctx, firewallId, ruleId).PatchedFirewallRuleRequest(req)

	ruleResponse, httpResp, err := request.Execute()
	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while updating a Firewall Rule", zap.Error(err), zap.Any("Rule ID", ruleId))
			errBody, err = utils.LogAndRewindBodyV4(httpResp)
			if err != nil {
				return sdk.FirewallRule{}, err
			}
		}
		return sdk.FirewallRule{}, utils.ErrorPerStatusCodeV4(errBody, httpResp, err)
	}

	return ruleResponse.Data, nil
}

func (c *Client) DeleteRule(ctx context.Context, firewallId int64, ruleId int64) error {
	logger.Debug("Delete Firewall Rule", zap.Any("Firewall ID", firewallId), zap.Any("Rule ID", ruleId))
	req := c.apiClient.FirewallsRulesEngineAPI.DeleteFirewallRule(ctx, firewallId, ruleId)

	_, httpResp, err := req.Execute()
	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while deleting a Firewall Rule", zap.Error(err), zap.Any("Rule ID", ruleId))
			errBody, err = utils.LogAndRewindBodyV4(httpResp)
			if err != nil {
				return err
			}
		}
		return utils.ErrorPerStatusCodeV4(errBody, httpResp, err)
	}

	return nil
}
