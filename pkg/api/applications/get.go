package applications

import (
	"context"

	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/utils"
	"go.uber.org/zap"
)

func (c *Client) Get(ctx context.Context, id int64) (ApplicationResponse, error) {
	logger.Debug("Get Application")

	res, httpResp, err := c.apiClient.ApplicationsAPI.
		RetrieveApplication(ctx, id).Execute()

	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while getting an Application", zap.Error(err))
			errBody, err = utils.LogAndRewindBodyV4(httpResp)
			if err != nil {
				return nil, err
			}
		}
		return nil, utils.ErrorPerStatusCodeV4(errBody, httpResp, err)
	}

	return &res.Data, nil
}
