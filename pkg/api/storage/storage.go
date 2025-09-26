package storage

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	sdk "github.com/aziontech/azionapi-v4-go-sdk-dev/storage-api"
	"go.uber.org/zap"

	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/utils"
)

type RequestBucket struct {
	sdk.BucketCreateRequest
}

func (c *Client) CreateBucket(ctx context.Context, request RequestBucket) error {
	logger.Debug("Creating bucket ", zap.Any("name", request.Name))
	req := c.apiClient.StorageBucketsAPI.CreateBucket(ctx).BucketCreateRequest(request.BucketCreateRequest)
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
	resp, httpResp, err := c.apiClient.StorageBucketsAPI.ListBuckets(ctx).Page(opts.Page).PageSize(opts.PageSize).Execute()
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
	_, httpResp, err := c.apiClient.StorageBucketsAPI.DeleteBucket(ctx, name).Execute()
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

func (c *Client) UpdateBucket(ctx context.Context, name string, edgeAccess string) error {
	logger.Debug("Updating bucket")
	bucket := sdk.PatchedBucketRequest{
		EdgeAccess: &edgeAccess,
	}
	_, httpResp, err := c.apiClient.StorageBucketsAPI.UpdateBucket(ctx, name).PatchedBucketRequest(bucket).Execute()
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
	c.apiClient.GetConfig().DefaultHeader["Content-Type"] = fileOps.MimeType
	req := c.apiClient.StorageObjectsAPI.CreateObjectKey(ctx, bucketName, objectKey).
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

func (c *Client) ListObject(ctx context.Context, bucketName string, opts *contracts.ListOptions) (*sdk.ResponseBucketObject, error) {
	logger.Debug("Listing bucket")
	// The storage objects list API now returns an array where each element contains
	// continuation_token and results. We perform a raw request to handle this shape
	// and then map the first element back into the SDK struct for downstream usage.

	// Build URL: <base>/v4/edge_storage/buckets/{bucketName}/objects
	base := c.apiClient.GetConfig().Servers[0].URL
	u, err := url.Parse(base)
	if err != nil {
		return nil, err
	}
	u.Path = fmt.Sprintf("%s/v4/edge_storage/buckets/%s/objects", u.Path, bucketName)

	// Query params
	q := u.Query()
	if opts.PageSize > 0 {
		q.Set("max_object_count", fmt.Sprintf("%d", opts.PageSize))
	}
	if opts.ContinuationToken != "" {
		q.Set("continuation_token", opts.ContinuationToken)
	}
	u.RawQuery = q.Encode()

	httpClient := c.apiClient.GetConfig().HTTPClient
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}
	// Headers
	for k, v := range c.apiClient.GetConfig().DefaultHeader {
		req.Header.Add(k, v)
	}
	req.Header.Set("Accept", "application/json")

	httpResp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	// Handle error status codes similarly to previous implementation
	if httpResp.StatusCode >= 400 {
		errBody, readErr := utils.LogAndRewindBodyV4(httpResp)
		if readErr != nil {
			return nil, readErr
		}
		return nil, utils.ErrorPerStatusCodeV4(errBody, httpResp, fmt.Errorf("status %d", httpResp.StatusCode))
	}

	// Read and decode array response
	body, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, err
	}

	// Try to decode as an array of generic pages to accommodate the new format
	var pages []map[string]any
	if err := json.Unmarshal(body, &pages); err != nil {
		// Fallback: try to unmarshal as single object (backward compatibility)
		var single sdk.ResponseBucketObject
		if err2 := json.Unmarshal(body, &single); err2 == nil {
			return &single, nil
		}
		logger.Debug("Error while listing Objects from Bucket", zap.Error(err))
		return nil, err
	}
	if len(pages) == 0 {
		// Return empty object
		empty := &sdk.ResponseBucketObject{}
		return empty, nil
	}

	// Map the first element into the SDK struct
	mapped := sdk.ResponseBucketObject{}
	// Marshal the first page back to JSON and unmarshal into the SDK struct to ensure compatibility
	firstPageBytes, err := json.Marshal(pages[0])
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(firstPageBytes, &mapped); err != nil {
		return nil, err
	}
	return &mapped, nil
}

func (c *Client) Upload(ctx context.Context, fileOps *contracts.FileOps, conf *contracts.AzionApplicationOptions, bucket string) error {
	file := fileOps.Path
	if conf.Prefix != "" {
		file = fmt.Sprintf("%s%s", conf.Prefix, fileOps.Path)
	}
	logger.Debug("Object_key: " + file)
	req := c.apiClient.StorageObjectsAPI.CreateObjectKey(ctx, bucket, file).Body(fileOps.FileContent).ContentType(fileOps.MimeType)

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
	logger.Debug("Getting object")
	_, httpResp, err := c.apiClient.StorageObjectsAPI.DownloadObject(ctx, bucketName, objectKey).Execute()
	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while getting the object", zap.Error(err))
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
	_, httpResp, err := c.apiClient.StorageObjectsAPI.DeleteObjectKey(ctx, bucketName, objectKey).Execute()
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
	_, httpResp, err := c.apiClient.StorageObjectsAPI.UpdateObjectKey(ctx, bucketName, objectKey).Body(body).Execute()
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
