package edge_applications

import (
	"context"

	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/utils"
	sdk "github.com/aziontech/azionapi-go-sdk/edgeapplications"

	"go.uber.org/zap"
)

// CreateCacheSettingsNextApplication this function creates the necessary Cache Settings for next applications
// to work correctly on the edge
func (c *Client) CreateCacheSettingsNextApplication(
	ctx context.Context, req *CreateCacheSettingsRequest,
	applicationID int64,
) (CacheSettingsResponse, error) {
	logger.Debug("Create Cache Settings Next Application")

	req.SetBrowserCacheSettings("override")
	req.SetBrowserCacheSettingsMaximumTtl(7200)
	req.SetCdnCacheSettings("override")
	req.SetCdnCacheSettingsMaximumTtl(7200)

	resp, err := c.CreateCacheEdgeApplication(ctx, req, applicationID)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *Client) CreateCacheEdgeApplication(
	ctx context.Context, req *CreateCacheSettingsRequest, edgeApplicationID int64,
) (CacheSettingsResponse, error) {
	logger.Debug("Create Cache Edge Application")
	resp, httpResp, err := c.apiClient.EdgeApplicationsCacheSettingsAPI.
		EdgeApplicationsEdgeApplicationIdCacheSettingsPost(ctx, edgeApplicationID).
		ApplicationCacheCreateRequest(req.ApplicationCacheCreateRequest).Execute()
	if err != nil {
		if httpResp != nil {
			logger.Debug("Error while creating a cache setting", zap.Error(err))
			err := utils.LogAndRewindBody(httpResp)
			if err != nil {
				return nil, err
			}
		}
		return nil, utils.ErrorPerStatusCode(httpResp, err)
	}
	return resp.Results, nil
}

func (c *Client) ListCacheEdgeApp(
	ctx context.Context, edgeApplicationID int64,
) ([]sdk.ApplicationCacheResults, error) {
	logger.Debug("List Cache Edge Application")
	resp, httpResp, err := c.apiClient.EdgeApplicationsCacheSettingsAPI.
		EdgeApplicationsEdgeApplicationIdCacheSettingsGet(ctx, edgeApplicationID).OrderBy("id").Execute()
	if err != nil {
		if httpResp != nil {
			logger.Debug("Error while listing a cache setting", zap.Error(err))
			err := utils.LogAndRewindBody(httpResp)
			if err != nil {
				return nil, err
			}
		}
		return nil, utils.ErrorPerStatusCode(httpResp, err)
	}
	return resp.Results, nil
}
