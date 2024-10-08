package storage

import (
	"context"

	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/utils"
	sdk "github.com/aziontech/azionapi-go-sdk/storage"
	"go.uber.org/zap"
)

type RequestCredentials struct {
	sdk.S3CredentialCreate
}

func (c *Client) CreateCredentials(ctx context.Context, request RequestCredentials) (*sdk.ResponseS3Credential, error) {
	logger.Debug("Creating s3 credentials ", zap.Any("name", request.Name))
	req := c.apiClient.StorageAPI.StorageApiS3CredentialsCreate(ctx).S3CredentialCreate(request.S3CredentialCreate)
	resp, httpResp, err := req.Execute()
	if err != nil {
		logger.Debug("Error while creating the user's s3 credentials", zap.Error(err))
		return nil, utils.ErrorPerStatusCode(httpResp, err)
	}
	return resp, nil
}
