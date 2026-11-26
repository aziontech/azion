// Package deployments is the API v6 client for Deployments and their
// releases. A release is the unit that goes live: it names one version of each
// resource, is built like a resource version, and is then activated.
package deployments

import (
	"net/http"

	"github.com/aziontech/azion-cli/pkg/cmd/version"
	sdk "github.com/aziontech/azionapi-v6-go-sdk/azion-api"
)

type Client struct {
	apiClient *sdk.APIClient

	// The Deployment services. SDK 0.4.0 has them but no longer attaches
	// them to APIClient, so the client builds them itself (see services).
	deployments *sdk.DeploymentsAPIService
	releases    *sdk.DeploymentsReleasesAPIService
	routing     *sdk.DeploymentsReleaseRoutingAPIService
	usage       *sdk.DeploymentsResourceUsageAPIService
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

	apiClient := sdk.NewAPIClient(conf)
	cl := &Client{apiClient: apiClient}
	cl.services()
	return cl
}

// services builds the Deployment services that SDK 0.4.0 leaves off APIClient.
// Every SDK service is the same struct, a pointer to the APIClient, under a
// different name, and NewAPIClient points them all at one shared value. Any
// service that is still attached is that value, so converting it gives a
// Deployment service bound to this client, as NewAPIClient would. Remove this
// once the SDK attaches them again.
func (c *Client) services() {
	common := c.apiClient.WorkloadsAPI
	c.deployments = (*sdk.DeploymentsAPIService)(common)
	c.releases = (*sdk.DeploymentsReleasesAPIService)(common)
	c.routing = (*sdk.DeploymentsReleaseRoutingAPIService)(common)
	c.usage = (*sdk.DeploymentsResourceUsageAPIService)(common)
}
