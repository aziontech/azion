package storage

import (
	"context"

	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/utils"
	sdk "github.com/aziontech/azionapi-v4-go-sdk-dev/storage-api"
	"go.uber.org/zap"
)

type RequestCredentials struct {
	sdk.CredentialCreateRequest
}

func (c *Client) CreateCredentials(ctx context.Context, request RequestCredentials) (*sdk.ResponseCredential, error) {
	logger.Debug("Creating s3 credentials ", zap.Any("name", request.Name))
	req := c.apiClient.StorageCredentialsAPI.CreateCredential(ctx).CredentialCreateRequest(request.CredentialCreateRequest)
	resp, httpResp, err := req.Execute()
	if err != nil {
		logger.Debug("Error while creating the user's s3 credentials", zap.Error(err))
		return nil, utils.ErrorPerStatusCode(httpResp, err)
	}
	return resp, nil
}
