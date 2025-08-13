package storage

import (
	"context"

	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/utils"
	sdk "github.com/aziontech/azionapi-v4-go-sdk/storage-api"
	"go.uber.org/zap"
)

type RequestCredentials struct {
	sdk.CredentialCreateRequest
}

func (c *Client) CreateCredentials(ctx context.Context, request RequestCredentials) (*sdk.CredentialCreate, error) {
	logger.Debug("Creating s3 credentials ", zap.Any("name", request.Name))
	req := c.apiClient.EdgeStorageCredentialsAPI.CreateCredential(ctx).CredentialCreateRequest(request.CredentialCreateRequest)
	resp, httpResp, err := req.Execute()
	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while creating the user's s3 credentials", zap.Error(err))
			errBody, err = utils.LogAndRewindBodyV4(httpResp)
			if err != nil {
				return nil, err
			}
		}

		return nil, utils.ErrorPerStatusCodeV4(errBody, httpResp, err)
	}
	return resp, nil
}
