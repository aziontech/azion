package edge_applications

import (
	"context"

	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/utils"
	"go.uber.org/zap"
)

func (c *Client) Get(ctx context.Context, id int64) (EdgeApplicationResponse, error) {
	logger.Debug("Get Edge Application")

	res, httpResp, err := c.apiClient.ApplicationsAPI.
		RetrieveApplication(ctx, id).Execute()

	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while getting an Edge Application", zap.Error(err))
			errBody, err = utils.LogAndRewindBodyV4(httpResp)
			if err != nil {
				return nil, err
			}
		}
		return nil, utils.ErrorPerStatusCodeV4(errBody, httpResp, err)
	}

	return &res.Data, nil
}
