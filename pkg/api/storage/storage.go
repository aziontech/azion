package storage

import (
	"context"
	"fmt"

	sdk "github.com/aziontech/azionapi-go-sdk/storage"
	"go.uber.org/zap"

	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/utils"
)

type RequestBucket struct {
	sdk.BucketCreate
}

func (c *Client) CreateBucket(ctx context.Context, request RequestBucket) error {
	logger.Debug("Creating bucket")
	req := c.apiClient.StorageAPI.StorageApiBucketsCreate(ctx).BucketCreate(request.BucketCreate)
	_, httpResp, err := req.Execute()
	if err != nil {
		if httpResp != nil {
			logger.Debug("Error while creating the project Bucket", zap.Error(err))
			err := utils.LogAndRewindBody(httpResp)
			if err != nil {
				return err
			}
		}
		return utils.ErrorPerStatusCode(httpResp, err)
	}
	return nil
}

func (c *Client) Upload(ctx context.Context, fileOps *contracts.FileOps, conf *contracts.AzionApplicationOptions) error {
	var file string
	if conf.Prefix != "" {
		file = fmt.Sprintf("%s%s", conf.Prefix, fileOps.Path)
		logger.Debug("Object_key: " + file)
	} else {
		file = fileOps.Path
	}
	req := c.apiClient.StorageAPI.StorageApiBucketsObjectsCreate(ctx, conf.Bucket, file).Body(fileOps.FileContent).ContentType(fileOps.MimeType)
	_, httpResp, err := req.Execute()
	if err != nil {
		if httpResp != nil {
			logger.Debug("Error while uploading file <"+fileOps.Path+"> to storage api", zap.Error(err))
			err := utils.LogAndRewindBody(httpResp)
			if err != nil {
				return err
			}
			return utils.ErrorPerStatusCode(httpResp, err)
		}
	}
	return nil 
}

