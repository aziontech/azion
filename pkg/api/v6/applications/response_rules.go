package applications

import (
	"context"

	"github.com/aziontech/azion-cli/pkg/contracts"

	"github.com/aziontech/azion-cli/pkg/api/v6/version"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/utils"
	sdk "github.com/aziontech/azionapi-v6-go-sdk/azion-api"
	"go.uber.org/zap"
)

// ListResponseRules lists the response rules of one version, one page at a time.
func (c *Client) ListResponseRules(ctx context.Context, id int64, versionID string, opts *contracts.ListOptions) (*sdk.PaginatedResponsePhaseRuleList, error) {
	logger.Debug("List Application response rules")

	req := c.apiClient.ApplicationsResponseRulesAPI.
		ListApplicationVersionResponseRules(ctx, id, versionID).Page(opts.Page).PageSize(int32(opts.PageSize))
	if opts.OrderBy != "" {
		req = req.Ordering(opts.OrderBy)
	}
	res, httpResp, err := req.Execute()
	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while listing response rules", zap.Error(err))
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

// CreateResponseRule creates a response rule in a Draft.
func (c *Client) CreateResponseRule(ctx context.Context, id int64, versionID string, req sdk.ResponsePhaseRuleRequest) (version.Item, error) {
	logger.Debug("Create Application response rule")

	res, httpResp, err := c.apiClient.ApplicationsResponseRulesAPI.
		CreateApplicationVersionResponseRules(ctx, id, versionID).ResponsePhaseRuleRequest(req).Execute()
	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while creating a response rule", zap.Error(err))
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

// UpdateResponseRule replaces a response rule in a Draft.
func (c *Client) UpdateResponseRule(ctx context.Context, id int64, versionID string, subID int64, req sdk.ResponsePhaseRuleRequest) (version.Item, error) {
	logger.Debug("Update Application response rule")

	res, httpResp, err := c.apiClient.ApplicationsResponseRulesAPI.
		UpdateApplicationVersionResponseRule(ctx, id, subID, versionID).ResponsePhaseRuleRequest(req).Execute()
	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while updating a response rule", zap.Error(err))
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

// DeleteResponseRule removes a response rule from a Draft.
func (c *Client) DeleteResponseRule(ctx context.Context, id int64, versionID string, subID int64) error {
	logger.Debug("Delete Application response rule")

	httpResp, err := c.apiClient.ApplicationsResponseRulesAPI.
		DeleteApplicationVersionResponseRule(ctx, id, subID, versionID).Execute()
	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while deleting a response rule", zap.Error(err))
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

// OrderResponseRules sets the order of the response rules in a Draft.
func (c *Client) OrderResponseRules(ctx context.Context, id int64, versionID string, ids []int64) error {
	logger.Debug("Order Application response rules")

	_, httpResp, err := c.apiClient.ApplicationsResponseRulesAPI.
		UpdateApplicationVersionResponseRulesOrder(ctx, id, versionID).ApplicationResponsePhaseRuleEngineOrderRequest(sdk.ApplicationResponsePhaseRuleEngineOrderRequest{Order: ids}).Execute()
	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while ordering response rules", zap.Error(err))
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
