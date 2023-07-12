package realtime_purge

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/aziontech/azion-cli/pkg/cmd/version"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/utils"
	sdk "github.com/aziontech/azionapi-go-sdk/realtimepurge"
	"go.uber.org/zap"
)

type Client struct {
	apiClient *sdk.APIClient
}

func NewClient(c *http.Client, url string, token string) *Client {
	conf := sdk.NewConfiguration()
	conf.HTTPClient = c
	conf.AddDefaultHeader("Authorization", "token "+token)
	conf.AddDefaultHeader("Accept", "application/json;version=3")
	conf.UserAgent = "Azion_CLI/" + version.BinVersion
	conf.Servers = sdk.ServerConfigurations{
		{URL: url},
	}
	conf.HTTPClient.Timeout = 30 * time.Second

	return &Client{
		apiClient: sdk.NewAPIClient(conf),
	}
}

func (c *Client) Purge(ctx context.Context, urlToPurge []string) error {
	var purg sdk.PurgeUrlRequest
	purg.SetUrls(urlToPurge)
	purg.SetMethod("delete")
	request := c.apiClient.RealTimePurgeApi.PurgeUrl(ctx).PurgeUrlRequest(purg)

	httpResp, err := c.apiClient.RealTimePurgeApi.PurgeUrlExecute(request)
	if err != nil {
		logger.Debug("Error while purgin a cache", zap.Error(err))
		logger.Debug("Status Code", zap.Any("http", httpResp.StatusCode))
		logger.Debug("Headers", zap.Any("http", httpResp.Header))
		logger.Debug("Response body", zap.Any("http", httpResp.Body))
		return utils.ErrorPerStatusCode(httpResp, err)
	}

	if httpResp.StatusCode != 201 {
		return fmt.Errorf("%w: %s", err, httpResp.Status)
	}

	return nil
}
