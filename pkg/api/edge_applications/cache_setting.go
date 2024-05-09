package edge_applications

import (
	"context"

	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/utils"

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
