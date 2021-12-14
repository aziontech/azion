package requests

import (
	"net/http"

	sdk "github.com/aziontech/edgeservices-go-sdk"
)

//TODO: receives token as an argument
//TODO: URL is passed during build-time

func CreateClient(client *http.Client, token string) (*sdk.APIClient, error) {

	conf := sdk.NewConfiguration()
	conf.HTTPClient = client
	conf.AddDefaultHeader("Authorization", "token "+token)
	conf.Servers = sdk.ServerConfigurations{
		{
			URL:         "https://stage-api.azion.net",
			Description: "User supplied",
		},
	}
	cli := sdk.NewAPIClient(conf)

	return cli, nil
}
