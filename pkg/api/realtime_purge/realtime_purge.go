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

func (c *Client) PurgeWildcard(ctx context.Context, urlToPurge []string) error {
	logger.Debug("Purge wildcard", zap.Any("url", urlToPurge))
	var purge sdk.PurgeWildcardRequest
	purge.SetUrls(urlToPurge)
	purge.SetMethod("delete")
	request := c.apiClient.RealTimePurgeApi.PurgeWildcard(ctx).PurgeWildcardRequest(purge)

	httpResp, err := c.apiClient.RealTimePurgeApi.PurgeWildcardExecute(request)
	if err != nil {
		logger.Debug("Error while purging wildcard", zap.Error(err))
		err = utils.LogAndRewindBody(httpResp)
		if err != nil {
			return err
		}

		return utils.ErrorPerStatusCode(httpResp, err)
	}

	if httpResp.StatusCode != 201 {
		return fmt.Errorf("%w: %s", err, httpResp.Status)
	}

	return nil
}

func (c *Client) PurgeUrls(ctx context.Context, urlToPurge []string) error {
	logger.Debug("Purge urls", zap.Any("url", urlToPurge))
	var purge sdk.PurgeUrlRequest
	purge.SetUrls(urlToPurge)
	purge.SetMethod("delete")
	request := c.apiClient.RealTimePurgeApi.PurgeUrl(ctx).PurgeUrlRequest(purge)

	httpResp, err := c.apiClient.RealTimePurgeApi.PurgeUrlExecute(request)
	if err != nil {
		logger.Debug("Error while purging urls", zap.Error(err))
		err = utils.LogAndRewindBody(httpResp)
		if err != nil {
			return err
		}

		return utils.ErrorPerStatusCode(httpResp, err)
	}

	if httpResp.StatusCode != 201 {
		return fmt.Errorf("%w: %s", err, httpResp.Status)
	}

	return nil
}

func (c *Client) PurgeCacheKey(ctx context.Context, urlToPurge []string, layer string) error {
	logger.Debug("Purge cache-key")
	var purge sdk.PurgeCacheKeyRequest
	purge.SetUrls(urlToPurge)
	purge.SetMethod("delete")
	purge.SetLayer(layer)
	request := c.apiClient.RealTimePurgeApi.PurgeCacheKey(ctx).PurgeCacheKeyRequest(purge)

	httpResp, err := c.apiClient.RealTimePurgeApi.PurgeCacheKeyExecute(request)
	if err != nil {
		logger.Debug("Error while purging cache keys", zap.Error(err))
		err = utils.LogAndRewindBody(httpResp)
		if err != nil {
			return err
		}

		return utils.ErrorPerStatusCode(httpResp, err)
	}

	return nil
}
