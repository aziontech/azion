package applications

import (
	"context"

	"go.uber.org/zap"

	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/utils"
)

func (c *Client) Delete(ctx context.Context, id int64) error {
	logger.Debug("Delete Application")
	req := c.apiClient.ApplicationsAPI.DeleteApplication(ctx, id)

	_, httpResp, err := req.Execute()
	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while deleting an Application", zap.Error(err), zap.Any("ID", id))
			var rewindErr error
			errBody, rewindErr = utils.LogAndRewindBodyV4(httpResp)
			if rewindErr != nil {
				return rewindErr
			}
		}

		return utils.ErrorPerStatusCodeV4(errBody, httpResp, err)
	}

	return nil
}
