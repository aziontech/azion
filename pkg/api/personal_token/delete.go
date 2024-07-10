package personal_token

import (
	"context"

	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/utils"
	"go.uber.org/zap"
)

func (c *Client) Delete(ctx context.Context, id string) error {
	logger.Debug("Delete personal token")
	httpResp, err := c.apiClient.PersonalTokenApi.DeletePersonalToken(ctx, id).Execute()
	if err != nil {
		if httpResp != nil {
			logger.Debug("Error while delete your personal token", zap.Error(err))
			err := utils.LogAndRewindBody(httpResp)
			if err != nil {
				return err
			}
		}
		return utils.ErrorPerStatusCode(httpResp, err)
	}
	return nil
}
