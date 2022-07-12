package requests

import (
	"time"

	"github.com/aziontech/azion-cli/pkg/cmd/version"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	sdk "github.com/aziontech/azionapi-go-sdk/edgeservices"
)

func CreateClient(f *cmdutil.Factory) (*sdk.APIClient, error) {
	conf := sdk.NewConfiguration()
	conf.HTTPClient = f.HttpClient
	conf.AddDefaultHeader("Authorization", "token "+f.Config.GetString("token"))
	conf.UserAgent = "Azion_CLI/" + version.BinVersion
	conf.Servers = sdk.ServerConfigurations{
		{
			URL: f.Config.GetString("api_url"),
		},
	}
	conf.HTTPClient.Timeout = 10 * time.Second

	return sdk.NewAPIClient(conf), nil
}
