package realtime_purge

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/aziontech/azion-cli/pkg/cmd/version"
	sdk "github.com/aziontech/azionapi-go-sdk/realtimepurge"
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
	conf.HTTPClient.Timeout = 10 * time.Second

	return &Client{
		apiClient: sdk.NewAPIClient(conf),
	}
}

func (c *Client) Purge(ctx context.Context, urlToPurge []string) error {
	var purg sdk.PurgeUrlRequest
	purg.SetUrls(urlToPurge)
	purg.SetMethod("delete")
	request := c.apiClient.RealTimePurgeApi.PurgeUrl(ctx).PurgeUrlRequest(purg)

	httpRes, err := c.apiClient.RealTimePurgeApi.PurgeUrlExecute(request)
	if err != nil {
		responseBody, _ := ioutil.ReadAll(httpRes.Body)
		return fmt.Errorf("%w: %s", err, responseBody)
	}

	if httpRes.StatusCode != 201 {
		return fmt.Errorf("%w: %s", err, httpRes.Status)
	}

	return nil
}
