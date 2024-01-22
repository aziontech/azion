package storage

import (
	"context"
	"fmt"
	"net/http"

	"github.com/aziontech/azion-cli/pkg/cmd/version"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/logger"
	"go.uber.org/zap"

	"github.com/aziontech/azion-cli/utils"
	"github.com/aziontech/azionapi-go-sdk/storage"
	sdk "github.com/aziontech/azionapi-go-sdk/storage"
)

type Client struct {
	apiClient *sdk.APIClient
}

func NewClient(c *http.Client, url string, token string) *Client {
	conf := sdk.NewConfiguration()
	conf.AddDefaultHeader("Authorization", "Token "+token)
	conf.UserAgent = "Azion_CLI/" + version.BinVersion
	conf.Servers = sdk.ServerConfigurations{
		{URL: url},
	}
	return &Client{
		apiClient: sdk.NewAPIClient(conf),
	}
}

type ClientStorage struct {
	apiClient *storage.APIClient
}

func NewClientStorage(c *http.Client, url string, token string) *ClientStorage {
	conf := storage.NewConfiguration()
	conf.AddDefaultHeader("Authorization", "Token "+token)
	conf.UserAgent = "Azion_CLI/" + version.BinVersion
	conf.Servers = storage.ServerConfigurations{
		{URL: url},
	}
	return &ClientStorage{
		apiClient: storage.NewAPIClient(conf),
	}
}

func (c *ClientStorage) CreateBucket(ctx context.Context, name string) error {
	logger.Debug("Creating bucket")
	create := storage.BucketCreate{
		Name:       name,
		EdgeAccess: storage.READ_WRITE,
	}

	req := c.apiClient.StorageAPI.StorageApiBucketsCreate(ctx).BucketCreate(create)
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
