package personal_token

import (
	"context"

	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/utils"
	"go.uber.org/zap"

	sdk "github.com/aziontech/azionapi-go-sdk/personal_tokens"
)

func (c *Client) List(ctx context.Context) ([]sdk.PersonalTokenResponseGet, error) {
	logger.Debug("List Personal Tokens")
	resp, httpResp, err := c.apiClient.PersonalTokenApi.ListPersonalToken(ctx).Execute()
	if err != nil {
		if httpResp != nil {
			logger.Debug("Error while listing your personal token", zap.Error(err))
			err := utils.LogAndRewindBody(httpResp)
			if err != nil {
				return nil, err
			}
		}
		return nil, utils.ErrorPerStatusCode(httpResp, err)
	}
	return resp.Results, nil
}
