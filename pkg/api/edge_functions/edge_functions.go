package edge_functions

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/aziontech/azion-cli/pkg/cmd/version"
	"github.com/aziontech/azion-cli/pkg/contracts"
	sdk "github.com/aziontech/azionapi-go-sdk/edgefunctions"
)

const javascript = "javascript"

type Client struct {
	apiClient *sdk.APIClient
}

type EdgeFunctionResponse interface {
	GetId() int64
	GetName() string
	GetActive() bool
	GetLanguage() string
	GetReferenceCount() int64
	GetModified() string
	GetInitiatorType() string
	GetLastEditor() string
	GetFunctionToRun() string
	GetJsonArgs() interface{}
	GetCode() string
}

func NewClient(c *http.Client, url string, token string) *Client {
	conf := sdk.NewConfiguration()
	conf.HTTPClient = c
	conf.AddDefaultHeader("Authorization", "token "+token)
	conf.AddDefaultHeader("Accept", "application/json;version=3")
	conf.UserAgent = "Azion_CLI/" + version.BinVersion
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
	// Although there's only one option, the API requires the `language` field.
	// Hard-coding javascript for now
	req.CreateEdgeFunctionRequest.SetLanguage(javascript)

	request := c.apiClient.EdgeFunctionsApi.EdgeFunctionsPost(ctx).CreateEdgeFunctionRequest(req.CreateEdgeFunctionRequest)

	edgeFuncResponse, httpRes, err := request.Execute()
	if err != nil {
		if httpRes == nil {
			return nil, err
		}
		responseBody, _ := ioutil.ReadAll(httpRes.Body)
		return nil, fmt.Errorf("%w: %s", err, responseBody)
	}

	return edgeFuncResponse.Results, nil
}

type UpdateRequest struct {
	sdk.PatchEdgeFunctionRequest
	Id int64
}

func NewUpdateRequest(id int64) *UpdateRequest {
	return &UpdateRequest{Id: id}
}

func (c *Client) Update(ctx context.Context, req *UpdateRequest) (EdgeFunctionResponse, error) {
	request := c.apiClient.EdgeFunctionsApi.EdgeFunctionsIdPatch(ctx, req.Id).PatchEdgeFunctionRequest(req.PatchEdgeFunctionRequest)

	edgeFuncResponse, httpRes, err := request.Execute()
	if err != nil {
		if httpRes == nil {
			return nil, err
		}
		responseBody, _ := ioutil.ReadAll(httpRes.Body)
		return nil, fmt.Errorf("%w: %s", err, responseBody)
	}

	return edgeFuncResponse.Results, nil
}

func (c *Client) List(ctx context.Context, opts *contracts.ListOptions) ([]EdgeFunctionResponse, error) {
	resp, httpResp, err := c.apiClient.EdgeFunctionsApi.EdgeFunctionsGet(ctx).
		OrderBy(opts.OrderBy).
		Page(opts.Page).
		PageSize(opts.PageSize).
		Sort(opts.Sort).
		Execute()

	if err != nil {
		if httpResp == nil {
			return nil, err
		}
		responseBody, _ := ioutil.ReadAll(httpResp.Body)
		return nil, fmt.Errorf("%w: %s", err, responseBody)
	}

	var result []EdgeFunctionResponse

	for i := range resp.GetResults() {
		result = append(result, &resp.GetResults()[i])
	}

	return result, nil
}
