package requests

import (
	sdk "github.com/aziontech/edgeservices-go-sdk"
)

//TODO: receives token as an argument
//TODO: URL is passed during build-time

func CreateClient() (*sdk.APIClient, error) {

	conf := sdk.NewConfiguration()
	conf.AddDefaultHeader("Authorization", "token 364d8f40562c20608c671760c447ab08aa91c62b")
	conf.Servers = sdk.ServerConfigurations{
		{
			URL:         "https://stage-api.azion.net",
			Description: "User supplied",
		},
	}
	cli := sdk.NewAPIClient(conf)

	return cli, nil
}
