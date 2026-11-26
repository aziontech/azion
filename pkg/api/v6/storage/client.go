// Package storage is the API v6 client for the buckets and S3 credentials the v6
// deploy needs. The other storage commands keep the v4 client in pkg/api/storage.
package storage

import (
	"net/http"

	"github.com/aziontech/azion-cli/pkg/cmd/version"
	sdk "github.com/aziontech/azionapi-v6-go-sdk/azion-api"
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

	return &Client{
		apiClient: sdk.NewAPIClient(conf),
	}
}
