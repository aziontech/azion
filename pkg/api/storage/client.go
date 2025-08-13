package storage

import (
	"net/http"

	sdk "github.com/aziontech/azionapi-v4-go-sdk/storage-api"

	"github.com/aziontech/azion-cli/pkg/cmd/version"
)

type Client struct {
	apiClient *sdk.APIClient
}

func NewClient(c *http.Client, url string, token string) *Client {
	conf := sdk.NewConfiguration()
	conf.HTTPClient = c
	conf.AddDefaultHeader("Authorization", "Token "+token)
	conf.UserAgent = "Azion_CLI/" + version.BinVersion
	conf.Servers = sdk.ServerConfigurations{
		{URL: url},
	}
	return &Client{
		apiClient: sdk.NewAPIClient(conf),
	}
}
