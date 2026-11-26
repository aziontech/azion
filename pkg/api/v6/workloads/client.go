// Package workloads is the API v6 client for Workloads. In v6 a Workload does
// not point at an Application and a Firewall: it carries one binding per
// Environment, naming the Deployment whose active release it serves and the
// domains it answers on.
package workloads

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
