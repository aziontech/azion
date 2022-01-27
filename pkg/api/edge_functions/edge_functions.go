package edge_funtions

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"

	sdk "github.com/aziontech/azionapi-go-sdk/edgefunctions"
)

type Client struct {
	apiClient *sdk.APIClient
}

type EdgeFunctionResponse interface {
	GetId() int64
	GetName() string
	GetLanguage() string
	GetReferenceCount() int64
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

func (c *Client) Get(ctx context.Context, id int64) (EdgeFunctionResponse, error) {
	req := c.apiClient.EdgeFunctionsApi.EdgeFunctionsIdGet(ctx, id)

	res, _, err := req.Execute()

	if err != nil {
		return nil, err
	}

	return res.Results, nil
}

func (c *Client) Delete(ctx context.Context, id int64) error {
	req := c.apiClient.EdgeFunctionsApi.EdgeFunctionsIdDelete(ctx, id)

	_, err := req.Execute()

	if err != nil {
		return err
	}

	return nil
}

type CreateRequest struct {
	Name          string
	Language      string
	InitiatorType string
	Active        bool
	Code          string
	JsonArgs      map[string]interface{}
}

func (c *Client) Create(ctx context.Context, req CreateRequest) (EdgeFunctionResponse, error) {
	var body sdk.CreateEdgeFunctionRequest

	body.SetActive(req.Active)
	body.SetCode(req.Code)
	body.SetName(req.Name)
	body.SetLanguage(req.Language)
	body.SetJsonArgs(req.JsonArgs)

	request := c.apiClient.EdgeFunctionsApi.EdgeFunctionsPost(ctx).CreateEdgeFunctionRequest(body)

	edgeFuncResponse, httpRes, err := request.Execute()
	if err != nil {
		responseBody, _ := ioutil.ReadAll(httpRes.Body)
		return nil, fmt.Errorf("%w: %s", err, responseBody)
	}

	return edgeFuncResponse.Results, nil
}
