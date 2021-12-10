package utils

import (
	sdk "github.com/aziontech/edgeservices-go-sdk"
)

//TODO: receives token as an argument
//TODO: URL is passed during build-time

func CreateClient() (*sdk.APIClient, error) {

	conf := sdk.NewConfiguration()
	conf.AddDefaultHeader("Authorization", "token d64b71607c8bf3e897b7c45b0420b88dfde8420b")
	conf.Servers = sdk.ServerConfigurations{
		{
			URL:         "https://stage-api.azion.net",
			Description: "User supplied",
		},
	}
	cli := sdk.NewAPIClient(conf)

	return cli, nil
}
