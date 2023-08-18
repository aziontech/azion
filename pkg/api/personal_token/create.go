package personal_token

import (
	"context"
	"time"

	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/utils"
	sdk "github.com/aziontech/azionapi-go-sdk/personal_tokens"
	"go.uber.org/zap"
)

type Response interface {
	GetUuid() string
	GetName() string
	GetKey() string
	GetUserId() float32
	GetCreated() time.Time
	GetExpiresAt() time.Time
	GetDescription() string
}

type Request struct {
	sdk.CreatePersonalTokenRequest
}

func (c *Client) Create(ctx context.Context, req *Request) (Response, error) {
	logger.Debug("Create Personal Token")

	response, httpResp, err := c.apiClient.PersonalTokenApi.CreatePersonalToken(ctx).
		CreatePersonalTokenRequest(req.CreatePersonalTokenRequest).Execute()

	if err != nil {
		logger.Error("Error while creating a personal token", zap.Error(err))
		return nil, utils.ErrorPerStatusCode(httpResp, err)
	}

	return response, nil
}
