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

// ListCacheSettings lists the cache settings of one version, one page at a time.
func (c *Client) ListCacheSettings(ctx context.Context, id int64, versionID string, opts *contracts.ListOptions) (*sdk.PaginatedCacheSettingList, error) {
	logger.Debug("List Application cache settings")

	req := c.apiClient.ApplicationsCacheSettingsAPI.
		ListApplicationVersionCacheSettings(ctx, id, versionID).Page(opts.Page).PageSize(int32(opts.PageSize))
	if opts.OrderBy != "" {
		req = req.Ordering(opts.OrderBy)
	}
	res, httpResp, err := req.Execute()
	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while listing cache settings", zap.Error(err))
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

// CreateCacheSetting creates a cache setting in a Draft.
func (c *Client) CreateCacheSetting(ctx context.Context, id int64, versionID string, req sdk.CacheSettingRequest) (version.Item, error) {
	logger.Debug("Create Application cache setting")

	res, httpResp, err := c.apiClient.ApplicationsCacheSettingsAPI.
		CreateApplicationVersionCacheSettings(ctx, id, versionID).CacheSettingRequest(req).Execute()
	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while creating a cache setting", zap.Error(err))
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

// UpdateCacheSetting replaces a cache setting in a Draft.
func (c *Client) UpdateCacheSetting(ctx context.Context, id int64, versionID string, subID int64, req sdk.CacheSettingRequest) (version.Item, error) {
	logger.Debug("Update Application cache setting")

	res, httpResp, err := c.apiClient.ApplicationsCacheSettingsAPI.
		UpdateApplicationVersionCacheSetting(ctx, id, subID, versionID).CacheSettingRequest(req).Execute()
	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while updating a cache setting", zap.Error(err))
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

// DeleteCacheSetting removes a cache setting from a Draft.
func (c *Client) DeleteCacheSetting(ctx context.Context, id int64, versionID string, subID int64) error {
	logger.Debug("Delete Application cache setting")

	httpResp, err := c.apiClient.ApplicationsCacheSettingsAPI.
		DeleteApplicationVersionCacheSetting(ctx, id, subID, versionID).Execute()
	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while deleting a cache setting", zap.Error(err))
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
