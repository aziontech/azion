package storage

import (
	"context"
	"time"

	"github.com/aziontech/azion-cli/pkg/api/v6/version"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/utils"
	sdk "github.com/aziontech/azionapi-v6-go-sdk/azion-api"
	"go.uber.org/zap"
)

// CreateBucket creates a bucket. A name already in use comes back as
// utils.ErrorNameInUse.
func (c *Client) CreateBucket(ctx context.Context, name, workloadsAccess string) error {
	logger.Debug("Creating bucket", zap.String("name", name))

	_, httpResp, err := c.apiClient.StorageBucketsAPI.
		CreateBucket(ctx).
		BucketCreateRequest(sdk.BucketCreateRequest{Name: name, WorkloadsAccess: workloadsAccess}).
		Execute()
	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while creating the project Bucket", zap.Error(err), zap.String("name", name))
			var bodyErr error
			errBody, bodyErr = utils.LogAndRewindBodyV4(httpResp)
			if bodyErr != nil {
				return bodyErr
			}
		}
		return version.ErrorFor(errBody, httpResp, err)
	}

	return nil
}

// CredentialRequest is what CreateCredentials needs.
type CredentialRequest struct {
	Name         string
	Capabilities []string
	Buckets      []string
	ExpiresAt    *time.Time
}

// CreateCredentials creates S3 credentials and returns them.
func (c *Client) CreateCredentials(ctx context.Context, req CredentialRequest) (*sdk.Credential, error) {
	logger.Debug("Creating s3 credentials", zap.String("name", req.Name))

	res, httpResp, err := c.apiClient.StorageCredentialsAPI.
		CreateCredential(ctx).
		CredentialCreateRequest(sdk.CredentialCreateRequest{
			Name:           req.Name,
			Capabilities:   req.Capabilities,
			Buckets:        req.Buckets,
			ExpirationDate: req.ExpiresAt,
		}).
		Execute()
	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while creating the user's s3 credentials", zap.Error(err))
			var bodyErr error
			errBody, bodyErr = utils.LogAndRewindBodyV4(httpResp)
			if bodyErr != nil {
				return nil, bodyErr
			}
		}
		return nil, version.ErrorFor(errBody, httpResp, err)
	}

	return &res.Data, nil
}
