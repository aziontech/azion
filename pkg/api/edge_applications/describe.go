package edge_applications

import (
	"context"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/utils"
	"go.uber.org/zap"
)

func (c *Client) Get(ctx context.Context, id string) (EdgeApplicationResponse, error) {
	logger.Debug("Get Edge Application")

	res, httpResp, err := c.apiClient.EdgeApplicationsMainSettingsAPI.
		EdgeApplicationsIdGet(ctx, id).Execute()

	if err != nil {
		logger.Debug("Error while getting an edge application", zap.Error(err))
		return nil, utils.ErrorPerStatusCode(httpResp, err)
	}

	return &res.Results, nil
}
