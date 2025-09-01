package edge_applications

import (
	"context"

	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/utils"
	sdk "github.com/aziontech/azionapi-v4-go-sdk-dev/edge-api"

	"go.uber.org/zap"
)

// CreateCacheSettingsNextApplication this function creates the necessary Cache Settings for next applications
// to work correctly on the edge
func (c *Client) CreateCacheSettingsNextApplication(ctx context.Context, req *CreateCacheSettingsRequest, applicationID string) (sdk.CacheSetting, error) {
	logger.Debug("Create Cache Settings Next Application")

	BCache := sdk.BrowserCacheModuleRequest{}
	BCache.SetBehavior("override")
	BCache.SetMaxAge(7200)
	ECache := sdk.CacheSettingsEdgeCacheModuleRequest{}
	ECache.SetBehavior("override")
	ECache.SetMaxAge(7200)

	cacheModules := sdk.CacheSettingsModulesRequest{
		EdgeCache: &ECache,
	}

	req.SetModules(cacheModules)
	req.SetBrowserCache(BCache)

	resp, err := c.CreateCacheEdgeApplication(ctx, req, applicationID)
	if err != nil {
		return sdk.CacheSetting{}, err
	}

	return resp, nil
}

func (c *Client) CreateCacheEdgeApplication(
	ctx context.Context, req *CreateCacheSettingsRequest, edgeApplicationID string,
) (sdk.CacheSetting, error) {
	logger.Debug("Create Cache - Application")
	resp, httpResp, err := c.apiClient.ApplicationsCacheSettingsAPI.
		CreateCacheSetting(ctx, edgeApplicationID).
		CacheSettingRequest(req.CacheSettingRequest).Execute()
	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while creating a cache setting", zap.Error(err))
			errBody, err = utils.LogAndRewindBodyV4(httpResp)
			if err != nil {
				return sdk.CacheSetting{}, err
			}
		}
		return sdk.CacheSetting{}, utils.ErrorPerStatusCodeV4(errBody, httpResp, err)
	}
	return resp.Data, nil
}

func (c *Client) ListCacheEdgeApp(
	ctx context.Context, edgeApplicationID string,
) ([]sdk.ResponseListCacheSetting, error) {
	logger.Debug("List Cache - Application")
	resp, httpResp, err := c.apiClient.ApplicationsCacheSettingsAPI.
		ListCacheSettings(ctx, edgeApplicationID).Execute()
	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while listing a cache setting", zap.Error(err))
			errBody, err = utils.LogAndRewindBodyV4(httpResp)
			if err != nil {
				return nil, err
			}
		}
		return nil, utils.ErrorPerStatusCodeV4(errBody, httpResp, err)
	}
	return resp.Results, nil
}
