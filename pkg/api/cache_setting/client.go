package cachesetting

import (
	"net/http"
	"time"

	"github.com/aziontech/azion-cli/pkg/cmd/version"
	sdkV3 "github.com/aziontech/azionapi-go-sdk/edgeapplications"
	sdk "github.com/aziontech/azionapi-v4-go-sdk-dev/edge-api"
)

type ClientV4 struct {
	apiClient *sdk.APIClient
}

func NewClientV4(c *http.Client, url string, token string) *ClientV4 {
	conf := sdk.NewConfiguration()
	conf.HTTPClient = c
	conf.AddDefaultHeader("Authorization", "token "+token)
	conf.AddDefaultHeader("Accept", "application/json;version=3")
	conf.UserAgent = "Azion_CLI/" + version.BinVersion
	conf.Servers = sdk.ServerConfigurations{
		{URL: url},
	}
	conf.HTTPClient.Timeout = 50 * time.Second

	return &ClientV4{
		apiClient: sdk.NewAPIClient(conf),
	}
}

type ClientV3 struct {
	apiClient *sdkV3.APIClient
}

func NewClientV3(c *http.Client, url string, token string) *ClientV3 {
	conf := sdkV3.NewConfiguration()
	conf.HTTPClient = c
	conf.AddDefaultHeader("Authorization", "token "+token)
	conf.AddDefaultHeader("Accept", "application/json;version=3")
	conf.UserAgent = "Azion_CLI/" + version.BinVersion
	conf.Servers = sdkV3.ServerConfigurations{
		{URL: url},
	}
	conf.HTTPClient.Timeout = 50 * time.Second

	return &ClientV3{
		apiClient: sdkV3.NewAPIClient(conf),
	}
}
