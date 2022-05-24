package realtime_purge

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/aziontech/azion-cli/pkg/cmd/version"
	sdk "github.com/aziontech/azionapi-go-sdk/realtimepurge"
)

type Client struct {
	apiClient *sdk.APIClient
}

type DomainResponse interface {
	GetId() int64
	GetDomainName() string
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

	return &Client{
		apiClient: sdk.NewAPIClient(conf),
	}
}

func (c *Client) Purge(ctx context.Context, urlToPurge *[]string) (*http.Response, error) {
	var purg sdk.PurgeUrlRequest
	purg.SetUrls(*urlToPurge)
	purg.SetMethod("delete")
	request := c.apiClient.RealTimePurgeApi.PurgeUrl(ctx).PurgeUrlRequest(purg)

	httpRes, err := c.apiClient.RealTimePurgeApi.PurgeUrlExecute(request)
	if err != nil {
		responseBody, _ := ioutil.ReadAll(httpRes.Body)
		return httpRes, fmt.Errorf("%w: %s", err, responseBody)
	}

	return httpRes, nil
}
