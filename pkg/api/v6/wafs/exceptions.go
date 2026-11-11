package wafs

import (
	"context"

	"github.com/aziontech/azion-cli/pkg/contracts"

	"github.com/aziontech/azion-cli/pkg/api/v6/version"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/utils"
	sdk "github.com/aziontech/azionapi-v6-go-sdk/azion-api"
	"go.uber.org/zap"
)

// ListExceptions lists the WAF exceptions of one version, one page at a time.
func (c *Client) ListExceptions(ctx context.Context, id int64, versionID string, opts *contracts.ListOptions) (*sdk.PaginatedWAFRuleList, error) {
	logger.Debug("List WAF WAF exceptions")

	req := c.apiClient.WAFsExceptionsAPI.
		ListWafVersionExceptions(ctx, versionID, id).Page(opts.Page).PageSize(int32(opts.PageSize))
	if opts.OrderBy != "" {
		req = req.Ordering(opts.OrderBy)
	}
	res, httpResp, err := req.Execute()
	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while listing WAF exceptions", zap.Error(err))
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

// CreateException creates a WAF exception in a Draft.
func (c *Client) CreateException(ctx context.Context, id int64, versionID string, req sdk.WAFRuleRequest) (version.Item, error) {
	logger.Debug("Create WAF WAF exception")

	res, httpResp, err := c.apiClient.WAFsExceptionsAPI.
		CreateWafVersionExceptions(ctx, versionID, id).WAFRuleRequest(req).Execute()
	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while creating a WAF exception", zap.Error(err))
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

// UpdateException replaces a WAF exception in a Draft.
func (c *Client) UpdateException(ctx context.Context, id int64, versionID string, subID int64, req sdk.WAFRuleRequest) (version.Item, error) {
	logger.Debug("Update WAF WAF exception")

	res, httpResp, err := c.apiClient.WAFsExceptionsAPI.
		UpdateWafVersionException(ctx, subID, versionID, id).WAFRuleRequest(req).Execute()
	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while updating a WAF exception", zap.Error(err))
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

// DeleteException removes a WAF exception from a Draft.
func (c *Client) DeleteException(ctx context.Context, id int64, versionID string, subID int64) error {
	logger.Debug("Delete WAF WAF exception")

	_, httpResp, err := c.apiClient.WAFsExceptionsAPI.
		DeleteWafVersionException(ctx, subID, versionID, id).Execute()
	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while deleting a WAF exception", zap.Error(err))
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
