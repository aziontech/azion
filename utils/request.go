package utils

import (
	sdk "github.com/aziontech/edgeservices-go-sdk"
)

//TODO: receives token as an argument

func CreateClient() (*sdk.APIClient, error) {

	conf := sdk.NewConfiguration()
	conf.AddDefaultHeader("Authorization", "token 137dd1d1564efba730356a1d2cf35a5f866b6d9c")
	conf.Servers = sdk.ServerConfigurations{
		{
			URL:         "https://stage-api.azion.net",
			Description: "User supplied",
		},
	}
	cli := sdk.NewAPIClient(conf)

	return cli, nil
}
