package edge_functions

import (
	"context"
	"io"
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
	logger.Debug("Get Edge Function")
	request := c.apiClient.EdgeFunctionsApi.EdgeFunctionsIdGet(ctx, id)

	res, httpResp, err := request.Execute()
	if err != nil {
		logger.Debug("Error while getting an edge function", zap.Error(err))
		logger.Debug("Status Code", zap.Any("http", httpResp.StatusCode))
		logger.Debug("Headers", zap.Any("http", httpResp.Header))
		logger.Debug("Response body", zap.Any("http", httpResp.Body))
		return nil, utils.ErrorPerStatusCode(httpResp, err)
	}

	return res.Results, nil
}

func (c *Client) Delete(ctx context.Context, id int64) error {
	logger.Debug("Delete Edge Function")
	request := c.apiClient.EdgeFunctionsApi.EdgeFunctionsIdDelete(ctx, id)

	httpResp, err := request.Execute()

	if err != nil {
		logger.Debug("Error while deleting an edge function", zap.Error(err))
		logger.Debug("Status Code", zap.Any("http", httpResp.StatusCode))
		logger.Debug("Headers", zap.Any("http", httpResp.Header))
		logger.Debug("Response body", zap.Any("http", httpResp.Body))
		return utils.ErrorPerStatusCode(httpResp, err)
	}

	return nil
}

func (c *Client) Create(ctx context.Context, req *CreateRequest) (EdgeFunctionResponse, error) {
	// Although there's only one option, the API requires the `language` field.
	// Hard-coding javascript for now
	logger.Debug("Create Edge Function")
	req.CreateEdgeFunctionRequest.SetLanguage(javascript)

	request := c.apiClient.EdgeFunctionsApi.EdgeFunctionsPost(ctx).CreateEdgeFunctionRequest(req.CreateEdgeFunctionRequest)

	edgeFuncResponse, httpResp, err := request.Execute()
	if err != nil {
		if httpResp != nil {
			logger.Debug("Error while creating an edge function", zap.Error(err))
			logger.Debug("", zap.Any("Status Code", httpResp.StatusCode))
			logger.Debug("", zap.Any("Headers", httpResp.Header))
			body, err := io.ReadAll(httpResp.Body)
			if err != nil {
				logger.Debug("Error while reading body of the http response", zap.Error(err))
				return nil, utils.ErrorPerStatusCode(httpResp, err)
			}
			logger.Debug("", zap.Any("Body", string(body)))
		}
		return nil, utils.ErrorPerStatusCode(httpResp, err)
	}

	return edgeFuncResponse.Results, nil
}

func (c *Client) Update(ctx context.Context, req *UpdateRequest) (EdgeFunctionResponse, error) {
	logger.Debug("Update Edge Function")
	request := c.apiClient.EdgeFunctionsApi.EdgeFunctionsIdPatch(ctx, req.Id).PatchEdgeFunctionRequest(req.PatchEdgeFunctionRequest)

	edgeFuncResponse, httpResp, err := request.Execute()
	if err != nil {
		if httpResp != nil {
			logger.Debug("Error while updating an edge function", zap.Error(err))
			logger.Debug("", zap.Any("Status Code", httpResp.StatusCode))
			logger.Debug("", zap.Any("Headers", httpResp.Header))
			body, err := io.ReadAll(httpResp.Body)
			if err != nil {
				logger.Debug("Error while reading body of the http response", zap.Error(err))
				return nil, utils.ErrorPerStatusCode(httpResp, err)
			}
			logger.Debug("", zap.Any("Body", string(body)))
		}
		return nil, utils.ErrorPerStatusCode(httpResp, err)
	}

	return edgeFuncResponse.Results, nil
}

func (c *Client) List(ctx context.Context, opts *contracts.ListOptions) ([]EdgeFunctionResponse, int64, error) {
	logger.Debug("List Edge Functions")
	resp, httpResp, err := c.apiClient.EdgeFunctionsApi.EdgeFunctionsGet(ctx).
		OrderBy(opts.OrderBy).
		Page(opts.Page).
		PageSize(opts.PageSize).
		Sort(opts.Sort).
		Execute()

	if err != nil {
		logger.Debug("Error while listing edge functions", zap.Error(err))
		logger.Debug("Status Code", zap.Any("http", httpResp.StatusCode))
		logger.Debug("Headers", zap.Any("http", httpResp.Header))
		logger.Debug("Response body", zap.Any("http", httpResp.Body))
		return nil, 0, utils.ErrorPerStatusCode(httpResp, err)
	}

	var result []EdgeFunctionResponse

	for i := range resp.GetResults() {
		result = append(result, &resp.GetResults()[i])
	}

	return result, *resp.TotalPages, nil
}
