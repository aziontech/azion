package requests

import (
	"fmt"

	"github.com/aziontech/azion-cli/pkg/cmdutil"
	sdk "github.com/aziontech/azionapi-go-sdk/edgeservices"
)

func CreateClient(f *cmdutil.Factory) (*sdk.APIClient, error) {
	httpClient, err := f.HttpClient()
	if err != nil {
		return nil, fmt.Errorf("failed to get http client: %w", err)
	}

	conf := sdk.NewConfiguration()
	conf.HTTPClient = httpClient
	conf.AddDefaultHeader("Authorization", "token "+f.Config.GetString("token"))
	conf.Servers = sdk.ServerConfigurations{
		{
			URL: f.Config.GetString("api_url"),
		},
	}

	return sdk.NewAPIClient(conf), nil
}
