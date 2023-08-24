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
		logger.Debug("Error while deleting a personal token", zap.Error(err))
		return utils.ErrorPerStatusCode(httpResp, err)
	}

	return nil
}
