package storage

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"

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
	logger.Debug("Creating bucket ", zap.Any("name", request.Name))
	req := c.apiClient.StorageAPI.StorageApiBucketsCreate(ctx).BucketCreate(request.BucketCreate)
	_, httpResp, err := req.Execute()
	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while creating the project Bucket", zap.Error(err), zap.Any("Name", request.Name))
			errBody, err = utils.LogAndRewindBodyV4(httpResp)
			if err != nil {
				return err
			}
		}

		return utils.ErrorPerStatusCodeV4(errBody, httpResp, err)
	}
	return nil
}

func (c *Client) ListBucket(ctx context.Context, opts *contracts.ListOptions) (*sdk.PaginatedBucketList, error) {
	logger.Debug("Listing bucket")
	resp, httpResp, err := c.apiClient.StorageAPI.StorageApiBucketsList(ctx).
		Page(int32(opts.Page)).PageSize(int32(opts.PageSize)).Execute()
	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Error("Error while listing buckets", zap.Error(err))
			errBody, err = utils.LogAndRewindBodyV4(httpResp)
			if err != nil {
				return nil, err
			}
		}

		return nil, utils.ErrorPerStatusCodeV4(errBody, httpResp, err)
	}

	return resp, nil
}

func (c *Client) DeleteBucket(ctx context.Context, name string) error {
	logger.Debug("Delete bucket", zap.Any("bucket-name", name))
	_, httpResp, err := c.apiClient.StorageAPI.
		StorageApiBucketsDestroy(ctx, name).Execute()
	if err != nil {
		errBody := ""
		if httpResp != nil {
			errBody, err = utils.LogAndRewindBodyV4(httpResp)
			if err != nil {
				return err
			}
		}
		return utils.ErrorPerStatusCodeV4(errBody, httpResp, err)
	}
	return nil
}

func (c *Client) UpdateBucket(ctx context.Context, name string, edgeAccess sdk.EdgeAccessEnum) error {
	logger.Debug("Updating bucket")
	bucket := sdk.BucketUpdate{
		EdgeAccess: edgeAccess,
	}
	_, httpResp, err := c.apiClient.StorageAPI.
		StorageApiBucketsPartialUpdate(ctx, name).BucketUpdate(bucket).Execute()
	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while updating the project Bucket", zap.Error(err), zap.Any("Name", name))
			errBody, err = utils.LogAndRewindBodyV4(httpResp)
			if err != nil {
				return err
			}
		}
		return utils.ErrorPerStatusCodeV4(errBody, httpResp, err)
	}
	return nil
}

func (c *Client) CreateObject(ctx context.Context, fileOps *contracts.FileOps, bucketName, objectKey string) error {
	logger.Debug("Creating object")
	req := c.apiClient.StorageAPI.StorageApiBucketsObjectsCreate(ctx, bucketName, objectKey).
		Body(fileOps.FileContent).ContentType(fileOps.MimeType)
	_, httpResp, err := req.Execute()
	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while creating object in the edge storage", zap.Error(err))
			errBody, err = utils.LogAndRewindBodyV4(httpResp)
			if err != nil {
				return err
			}
		}
		return utils.ErrorPerStatusCodeV4(errBody, httpResp, err)
	}

	return nil
}

func (c *Client) ListObject(ctx context.Context, bucketName string, opts *contracts.ListOptions) (*sdk.PaginatedBucketObjectList, error) {
	logger.Debug("Listing bucket")
	req := c.apiClient.StorageAPI.StorageApiBucketsObjectsList(ctx, bucketName).
		MaxObjectCount(int32(opts.PageSize)).ContinuationToken(opts.ContinuationToken)
	resp, httpResp, err := req.Execute()
	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while listing Objects from Bucket", zap.Error(err))
			errBody, err = utils.LogAndRewindBodyV4(httpResp)
			if err != nil {
				return nil, err
			}
		}
		return nil, utils.ErrorPerStatusCodeV4(errBody, httpResp, err)
	}
	return resp, nil
}

func (c *Client) Upload(ctx context.Context, fileOps *contracts.FileOps, conf *contracts.AzionApplicationOptions, bucket string) error {
	file := fileOps.Path
	if conf.Prefix != "" {
		file = fmt.Sprintf("%s%s", conf.Prefix, fileOps.Path)
	}
	logger.Debug("Object_key: " + file)
	req := c.apiClient.StorageAPI.StorageApiBucketsObjectsCreate(ctx, bucket, file).Body(fileOps.FileContent).ContentType(fileOps.MimeType)
	_, httpResp, err := req.Execute()
	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while uploading file <"+fileOps.Path+"> to storage api", zap.Error(err))
			errBody, err = utils.LogAndRewindBodyV4(httpResp)
			if err != nil {
				return err
			}
		}
		return utils.ErrorPerStatusCodeV4(errBody, httpResp, err)
	}

	return nil
}

func (c *Client) GetObject(ctx context.Context, bucketName, objectKey string) ([]byte, error) {
	logger.Debug("Getting bucket")
	httpResp, err := c.apiClient.StorageAPI.StorageApiBucketsObjectsRetrieve(ctx, bucketName, objectKey).Execute()
	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while updating the project Bucket", zap.Error(err))
			errBody, err = utils.LogAndRewindBodyV4(httpResp)
			if err != nil {
				return nil, err
			}
		}
		return nil, utils.ErrorPerStatusCodeV4(errBody, httpResp, err)
	}
	byteObject, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, errors.New("Error reading edge storage objects file")
	}
	return byteObject, nil
}

func (c *Client) DeleteObject(ctx context.Context, bucketName, objectKey string) error {
	logger.Debug("Delete object", zap.Any("object-key", objectKey))
	_, httpResp, err := c.apiClient.StorageAPI.
		StorageApiBucketsObjectsDestroy(ctx, bucketName, objectKey).Execute()
	if err != nil {
		errBody := ""
		if httpResp != nil {
			errBody, err = utils.LogAndRewindBodyV4(httpResp)
			if err != nil {
				return err
			}
		}
		return utils.ErrorPerStatusCodeV4(errBody, httpResp, err)
	}
	return nil
}

func (c *Client) UpdateObject(ctx context.Context, bucketName, objectKey, contentType string, body *os.File) error {
	logger.Debug("Updating objects")
	_, httpResp, err := c.apiClient.StorageAPI.StorageApiBucketsObjectsUpdate(ctx, bucketName, objectKey).
		ContentType(contentType).Body(body).Execute()
	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while updating the object of the bucket", zap.Error(err))
			errBody, err = utils.LogAndRewindBodyV4(httpResp)
			if err != nil {
				return err
			}
		}

		return utils.ErrorPerStatusCodeV4(errBody, httpResp, err)
	}
	return nil
}
