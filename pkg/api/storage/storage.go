package storage

import (
	"context"
	"net/http"

	"github.com/aziontech/azion-cli/pkg/cmd/version"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/logger"
	"go.uber.org/zap"

	"github.com/aziontech/azion-cli/utils"
	storage "github.com/aziontech/azionapi-go-sdk/storage"
	sdk "github.com/aziontech/azionapi-go-sdk/storageapi"
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
	create := storage.BucketCreate{
		Name:       name,
		EdgeAccess: storage.READ_WRITE,
	}
	_, httpResp, err := c.apiClient.BucketsAPI.ApiV1StorageBucketsCreate(ctx).BucketCreate(create).Execute()
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

func (c *Client) Upload(ctx context.Context, fileOps *contracts.FileOps) error {
	req := c.apiClient.DefaultApi.StorageVersionIdPost(ctx, fileOps.VersionID).XAzionStaticPath(fileOps.Path).Body(fileOps.FileContent).ContentType(fileOps.MimeType)
	_, httpResp, err := req.Execute()
	if err != nil {
		logger.Debug("Error while uploading file <"+fileOps.Path+"> to storage api", zap.Error(err))
		return utils.ErrorPerStatusCode(httpResp, err)
	}
	return nil
}
