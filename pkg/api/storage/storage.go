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

func (c *Client) Upload(ctx context.Context, fileOps *contracts.FileOps) error {
	c.apiClient.GetConfig().DefaultHeader["Content-Disposition"] = fmt.Sprintf("attachment; filename=\"%s\"", fileOps.Path)
	req := c.apiClient.DefaultApi.StorageVersionIdPost(ctx, fileOps.VersionID).XAzionStaticPath(fileOps.Path).Body(fileOps.FileContent).ContentType(fileOps.MimeType)
	_, httpResp, err := req.Execute()
	if err != nil {
		logger.Debug("Error while uploading file <"+fileOps.Path+"> to storage api", zap.Error(err))
		return utils.ErrorPerStatusCode(httpResp, err)
	}
	return nil
}
