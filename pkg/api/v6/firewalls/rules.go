package firewalls

import (
	"context"

	"github.com/aziontech/azion-cli/pkg/contracts"

	"github.com/aziontech/azion-cli/pkg/api/v6/version"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/utils"
	sdk "github.com/aziontech/azionapi-v6-go-sdk/azion-api"
	"go.uber.org/zap"
)

// ListRules lists the firewall rules of one version, one page at a time.
func (c *Client) ListRules(ctx context.Context, id int64, versionID string, opts *contracts.ListOptions) (*sdk.PaginatedFirewallRuleList, error) {
	logger.Debug("List Firewall firewall rules")

	req := c.apiClient.FirewallsRulesEngineAPI.
		ListFirewallVersionRequestRules(ctx, id, versionID).Page(opts.Page).PageSize(int32(opts.PageSize))
	if opts.OrderBy != "" {
		req = req.Ordering(opts.OrderBy)
	}
	res, httpResp, err := req.Execute()
	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while listing firewall rules", zap.Error(err))
			var bodyErr error
			errBody, bodyErr = utils.LogAndRewindBodyV4(httpResp)
			if bodyErr != nil {
				return nil, bodyErr
			}
		}
		return nil, version.ErrorFor(errBody, httpResp, err)
	}

	return res, nil
}

// CreateRule creates a firewall rule in a Draft.
func (c *Client) CreateRule(ctx context.Context, id int64, versionID string, req sdk.FirewallRuleRequest) (version.Item, error) {
	logger.Debug("Create Firewall firewall rule")

	res, httpResp, err := c.apiClient.FirewallsRulesEngineAPI.
		CreateFirewallVersionRequestRules(ctx, id, versionID).FirewallRuleRequest(req).Execute()
	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while creating a firewall rule", zap.Error(err))
			var bodyErr error
			errBody, bodyErr = utils.LogAndRewindBodyV4(httpResp)
			if bodyErr != nil {
				return nil, bodyErr
			}
		}
		return nil, version.ErrorFor(errBody, httpResp, err)
	}

	return &res.Data, nil
}

// UpdateRule replaces a firewall rule in a Draft.
func (c *Client) UpdateRule(ctx context.Context, id int64, versionID string, subID int64, req sdk.FirewallRuleRequest) (version.Item, error) {
	logger.Debug("Update Firewall firewall rule")

	res, httpResp, err := c.apiClient.FirewallsRulesEngineAPI.
		UpdateFirewallVersionRequestRule(ctx, id, subID, versionID).FirewallRuleRequest(req).Execute()
	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while updating a firewall rule", zap.Error(err))
			var bodyErr error
			errBody, bodyErr = utils.LogAndRewindBodyV4(httpResp)
			if bodyErr != nil {
				return nil, bodyErr
			}
		}
		return nil, version.ErrorFor(errBody, httpResp, err)
	}

	return &res.Data, nil
}

// DeleteRule removes a firewall rule from a Draft.
func (c *Client) DeleteRule(ctx context.Context, id int64, versionID string, subID int64) error {
	logger.Debug("Delete Firewall firewall rule")

	_, httpResp, err := c.apiClient.FirewallsRulesEngineAPI.
		DeleteFirewallVersionRequestRule(ctx, id, subID, versionID).Execute()
	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while deleting a firewall rule", zap.Error(err))
			var bodyErr error
			errBody, bodyErr = utils.LogAndRewindBodyV4(httpResp)
			if bodyErr != nil {
				return bodyErr
			}
		}
		return version.ErrorFor(errBody, httpResp, err)
	}

	return nil
}

// OrderRules sets the order of the firewall rules in a Draft.
func (c *Client) OrderRules(ctx context.Context, id int64, versionID string, ids []int64) error {
	logger.Debug("Order Firewall firewall rules")

	_, httpResp, err := c.apiClient.FirewallsRulesEngineAPI.
		UpdateFirewallVersionRequestRulesOrder(ctx, id, versionID).FirewallRuleEngineOrderRequest(sdk.FirewallRuleEngineOrderRequest{Order: ids}).Execute()
	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while ordering firewall rules", zap.Error(err))
			var bodyErr error
			errBody, bodyErr = utils.LogAndRewindBodyV4(httpResp)
			if bodyErr != nil {
				return bodyErr
			}
		}
		return version.ErrorFor(errBody, httpResp, err)
	}

	return nil
}
