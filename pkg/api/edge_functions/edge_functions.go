package edge_functions

import (
	"context"
	"net/http"
	"time"

	"github.com/aziontech/azion-cli/pkg/cmd/version"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/utils"
	sdk "github.com/aziontech/azionapi-go-sdk/edgefunctions"
	"go.uber.org/zap"
)

const javascript = "javascript"

type Client struct {
	apiClient *sdk.APIClient
}

type CreateRequest struct {
	sdk.CreateEdgeFunctionRequest
}

func NewCreateRequest() *CreateRequest {
	return &CreateRequest{}
}

type UpdateRequest struct {
	sdk.PatchEdgeFunctionRequest
	Id int64
}

func NewUpdateRequest(id int64) *UpdateRequest {
	return &UpdateRequest{Id: id}
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
	conf.HTTPClient.Timeout = 30 * time.Second

	return &Client{
		apiClient: sdk.NewAPIClient(conf),
	}
}

func (c *Client) Get(ctx context.Context, id int64) (EdgeFunctionResponse, error) {
	req := c.apiClient.EdgeFunctionsApi.EdgeFunctionsIdGet(ctx, id)
	logger.Debug("request", zap.Any("request", req))

	res, httpResp, err := req.Execute()
	logger.Debug("response", zap.Any("response struct", res), zap.Any("response http", httpResp), zap.Error(err))

	if err != nil {
		logger.Error("Get request.Execute return error", zap.Error(err))
		return nil, utils.ErrorPerStatusCode(httpResp, err)
	}

	return res.Results, nil
}

func (c *Client) Delete(ctx context.Context, id int64) error {
	req := c.apiClient.EdgeFunctionsApi.EdgeFunctionsIdDelete(ctx, id)
	logger.Debug("request", zap.Any("request", req))

	httpResp, err := req.Execute()
	logger.Debug("response", zap.Any("response http", httpResp), zap.Error(err))

	if err != nil {
		logger.Error("Delete request.Execute return error", zap.Error(err))
		return utils.ErrorPerStatusCode(httpResp, err)
	}

	return nil
}

func (c *Client) Create(ctx context.Context, req *CreateRequest) (EdgeFunctionResponse, error) {
	// Although there's only one option, the API requires the `language` field.
	// Hard-coding javascript for now
	req.CreateEdgeFunctionRequest.SetLanguage(javascript)

	request := c.apiClient.EdgeFunctionsApi.EdgeFunctionsPost(ctx).CreateEdgeFunctionRequest(req.CreateEdgeFunctionRequest)
	logger.Debug("request", zap.Any("request", request))

	edgeFuncResponse, httpResp, err := request.Execute()
	logger.Debug("response", zap.Any("response struct", edgeFuncResponse), zap.Any("response http", httpResp), zap.Error(err))
	if err != nil {
		logger.Error("Create request.Execute return error", zap.Error(err))
		return nil, utils.ErrorPerStatusCode(httpResp, err)
	}

	return edgeFuncResponse.Results, nil
}

func (c *Client) Update(ctx context.Context, req *UpdateRequest) (EdgeFunctionResponse, error) {
	request := c.apiClient.EdgeFunctionsApi.EdgeFunctionsIdPatch(ctx, req.Id).PatchEdgeFunctionRequest(req.PatchEdgeFunctionRequest)
	logger.Debug("request", zap.Any("request", request))

	edgeFuncResponse, httpResp, err := request.Execute()
	logger.Debug("response", zap.Any("response struct", edgeFuncResponse), zap.Any("response http", httpResp), zap.Error(err))
	if err != nil {
		logger.Error("Update request.Execute return error", zap.Error(err))
		return nil, utils.ErrorPerStatusCode(httpResp, err)
	}

	return edgeFuncResponse.Results, nil
}

func (c *Client) List(ctx context.Context, opts *contracts.ListOptions) ([]EdgeFunctionResponse, int64, error) {
	resp, httpResp, err := c.apiClient.EdgeFunctionsApi.EdgeFunctionsGet(ctx).
		OrderBy(opts.OrderBy).
		Page(opts.Page).
		PageSize(opts.PageSize).
		Sort(opts.Sort).
		Execute()
	logger.Debug("response", zap.Any("response struct", resp), zap.Any("response http", httpResp), zap.Error(err))

	if err != nil {
		logger.Error("list request.Execute return error", zap.Error(err))
		return nil, 0, utils.ErrorPerStatusCode(httpResp, err)
	}

	var result []EdgeFunctionResponse

	for i := range resp.GetResults() {
		result = append(result, &resp.GetResults()[i])
	}

	return result, *resp.TotalPages, nil
}
