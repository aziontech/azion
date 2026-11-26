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

// ListRequestRules lists the request rules of one version, one page at a time.
func (c *Client) ListRequestRules(ctx context.Context, id int64, versionID string, opts *contracts.ListOptions) (*sdk.PaginatedRequestPhaseRuleList, error) {
	logger.Debug("List Application request rules")

	req := c.apiClient.ApplicationsRequestRulesAPI.
		ListApplicationVersionRequestRules(ctx, id, versionID).Page(opts.Page).PageSize(int32(opts.PageSize))
	if opts.OrderBy != "" {
		req = req.Ordering(opts.OrderBy)
	}
	res, httpResp, err := req.Execute()
	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while listing request rules", zap.Error(err))
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

// CreateRequestRule creates a request rule in a Draft.
func (c *Client) CreateRequestRule(ctx context.Context, id int64, versionID string, req sdk.RequestPhaseRuleRequest) (version.Item, error) {
	logger.Debug("Create Application request rule")

	res, httpResp, err := c.apiClient.ApplicationsRequestRulesAPI.
		CreateApplicationVersionRequestRules(ctx, id, versionID).RequestPhaseRuleRequest(req).Execute()
	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while creating a request rule", zap.Error(err))
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

// UpdateRequestRule replaces a request rule in a Draft.
func (c *Client) UpdateRequestRule(ctx context.Context, id int64, versionID string, subID int64, req sdk.RequestPhaseRuleRequest) (version.Item, error) {
	logger.Debug("Update Application request rule")

	res, httpResp, err := c.apiClient.ApplicationsRequestRulesAPI.
		UpdateApplicationVersionRequestRule(ctx, id, subID, versionID).RequestPhaseRuleRequest(req).Execute()
	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while updating a request rule", zap.Error(err))
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

// DeleteRequestRule removes a request rule from a Draft.
func (c *Client) DeleteRequestRule(ctx context.Context, id int64, versionID string, subID int64) error {
	logger.Debug("Delete Application request rule")

	httpResp, err := c.apiClient.ApplicationsRequestRulesAPI.
		DeleteApplicationVersionRequestRule(ctx, id, subID, versionID).Execute()
	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while deleting a request rule", zap.Error(err))
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

// OrderRequestRules sets the order of the request rules in a Draft.
func (c *Client) OrderRequestRules(ctx context.Context, id int64, versionID string, ids []int64) error {
	logger.Debug("Order Application request rules")

	_, httpResp, err := c.apiClient.ApplicationsRequestRulesAPI.
		UpdateApplicationVersionRequestRulesOrder(ctx, id, versionID).ApplicationRequestPhaseRuleEngineOrderRequest(sdk.ApplicationRequestPhaseRuleEngineOrderRequest{Order: ids}).Execute()
	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while ordering request rules", zap.Error(err))
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
