package realtime_purge

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/aziontech/azion-cli/pkg/cmd/version"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/utils"
	sdk "github.com/aziontech/azionapi-v4-go-sdk-dev/edge-api"
	"go.uber.org/zap"
)

type Client struct {
	apiClient *sdk.APIClient
}

func NewClient(c *http.Client, url string, token string) *Client {
	conf := sdk.NewConfiguration()
	conf.HTTPClient = c
	conf.AddDefaultHeader("Authorization", "token "+token)
	conf.AddDefaultHeader("Accept", "application/json")
	conf.UserAgent = "Azion_CLI/" + version.BinVersion
	conf.Servers = sdk.ServerConfigurations{
		{URL: url},
	}
	conf.HTTPClient.Timeout = 50 * time.Second

	return &Client{
		apiClient: sdk.NewAPIClient(conf),
	}
}

func (c *Client) PurgeCache(ctx context.Context, urlToPurge []string, purgeType, layer string) error {
	logger.Debug("Purge cache", zap.Any("purge type", purgeType), zap.Any("urls", urlToPurge))
	request := c.apiClient.PurgeAPI.CreatePurgeRequest(ctx, purgeType)
	purgeRequest := *sdk.NewPurgeInputRequest(urlToPurge)
	purgeRequest.SetLayer(layer)

	_, httpResp, err := request.PurgeInputRequest(purgeRequest).Execute()
	if err != nil {
		errBody := ""
		logger.Debug("Error while purging cache", zap.Error(err))
		errBody, err = utils.LogAndRewindBodyV4(httpResp)
		if err != nil {
			return err
		}

		return utils.ErrorPerStatusCodeV4(errBody, httpResp, err)
	}

	if httpResp.StatusCode != 201 {
		return fmt.Errorf("%w: %s", err, httpResp.Status)
	}

	return nil
}
