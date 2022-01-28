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
	sdk.CreateEdgeFunctionRequest
}

func NewCreateRequest() *CreateRequest {
	return &CreateRequest{}
}

func (c *Client) Create(ctx context.Context, req *CreateRequest) (EdgeFunctionResponse, error) {
	request := c.apiClient.EdgeFunctionsApi.EdgeFunctionsPost(ctx).CreateEdgeFunctionRequest(req.CreateEdgeFunctionRequest)

	edgeFuncResponse, httpRes, err := request.Execute()
	if err != nil {
		responseBody, _ := ioutil.ReadAll(httpRes.Body)
		return nil, fmt.Errorf("%w: %s", err, responseBody)
	}

	return edgeFuncResponse.Results, nil
}

type UpdateRequest struct {
	sdk.PatchEdgeFunctionRequest
	id int64
}

func NewUpdateRequest(id int64) *UpdateRequest {
	return &UpdateRequest{id: id}
}

func (c *Client) Update(ctx context.Context, req *UpdateRequest) (EdgeFunctionResponse, error) {
	request := c.apiClient.EdgeFunctionsApi.EdgeFunctionsIdPatch(ctx, req.id).PatchEdgeFunctionRequest(req.PatchEdgeFunctionRequest)

	edgeFuncResponse, httpRes, err := request.Execute()
	if err != nil {
		responseBody, _ := ioutil.ReadAll(httpRes.Body)
		return nil, fmt.Errorf("%w: %s", err, responseBody)
	}

	return edgeFuncResponse.Results, nil
}
