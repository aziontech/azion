package edge_funtions

import (
	"context"
	"net/http"

	sdk "github.com/aziontech/azionapi-go-sdk/edgefunctions"
)

type Client struct {
	apiClient *sdk.APIClient
}

type EdgeFunction interface {
	GetId() float32 // Should be uint64
	GetName() string
	GetLanguage() string
	GetReferenceCount() float32 // Should be uint64
	GetModified() string
	GetInitiatorType() string
	GetLastEditor() string
	GetFunctionToRun() string
	GetJsonArgs() map[string]interface{}
	GetCode() string
}

func NewClient(c *http.Client, url string, token string) *Client {
	conf := sdk.NewConfiguration()
	conf.HTTPClient = c
	conf.AddDefaultHeader("Authorization", "token "+token)
	conf.AddDefaultHeader("Accept", "application/json;version=3")
	conf.Servers = sdk.ServerConfigurations{
		{URL: url},
	}

	return &Client{
		apiClient: sdk.NewAPIClient(conf),
	}
}

func (c *Client) Get(ctx context.Context, id string) (EdgeFunction, error) {
	req := c.apiClient.EdgeFunctionsApi.EdgeFunctionsIdGet(ctx, id)

	res, _, err := req.Execute()

	if err != nil {
		return nil, err
	}

	return res.Results, nil
}

func (c *Client) Delete(ctx context.Context, id string) error {
	req := c.apiClient.EdgeFunctionsApi.EdgeFunctionsIdDelete(ctx, id)

	_, err := req.Execute()

	if err != nil {
		return err
	}

	return nil

}
