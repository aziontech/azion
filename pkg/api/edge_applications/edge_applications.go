package edgeapplications

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/aziontech/azion-cli/pkg/cmd/version"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/utils"
	sdk "github.com/aziontech/azionapi-go-sdk/edgeapplications"
)

type EdgeApplicationResponse interface {
	GetId() int64
	GetName() string
	GetActive() bool
	GetApplicationAcceleration() bool
	GetCaching() bool
	GetDeliveryProtocol() string
	GetDeviceDetection() bool
	GetEdgeFirewall() bool
	GetEdgeFunctions() bool
	GetHttpPort() int64
	GetHttpsPort() int64
	GetImageOptimization() bool
	GetL2Caching() bool
	GetLoadBalancer() bool
	GetMinimumTlsVersion() string
	GetNext() string
	GetRawLogs() bool
	GetWebApplicationFirewall() bool
}

type Client struct {
	apiClient *sdk.APIClient
}

type CreateRequest struct {
	sdk.CreateApplicationRequest
}

type UpdateRequest struct {
	sdk.ApplicationUpdateRequest
	Id int64
}

type UpdateInstanceRequest struct {
	sdk.ApplicationUpdateInstanceRequest
	Id         string
	IdInstace  string
	FunctionId int64
}

type CreateInstanceRequest struct {
	sdk.ApplicationCreateInstanceRequest
	ApplicationId int64
}

type EdgeApplicationsResponse interface {
	GetId() int64
	GetName() string
}

type UpdateRulesEngineRequest struct {
	sdk.PatchRulesEngineRequest
	IdApplication int64
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

func (c *Client) Get(ctx context.Context, id string) (EdgeApplicationResponse, error) {
	req := c.apiClient.EdgeApplicationsMainSettingsApi.EdgeApplicationsIdGet(ctx, id)

	res, httpResp, err := req.Execute()
	if err != nil {
		return nil, utils.ErrorPerStatusCode(httpResp, err)
	}

	return &res.Results, nil
}

func (c *Client) Create(ctx context.Context, req *CreateRequest) (EdgeApplicationsResponse, error) {

	request := c.apiClient.EdgeApplicationsMainSettingsApi.EdgeApplicationsPost(ctx).CreateApplicationRequest(req.CreateApplicationRequest)

	edgeApplicationsResponse, httpResp, err := request.Execute()
	if err != nil {
		return nil, utils.ErrorPerStatusCode(httpResp, err)
	}

	return &edgeApplicationsResponse.Results, nil
}

func (c *Client) Update(ctx context.Context, req *UpdateRequest) (EdgeApplicationsResponse, error) {
	str := strconv.FormatInt(req.Id, 10)
	request := c.apiClient.EdgeApplicationsMainSettingsApi.EdgeApplicationsIdPatch(ctx, str).ApplicationUpdateRequest(req.ApplicationUpdateRequest)

	edgeApplicationsResponse, httpResp, err := request.Execute()
	if err != nil {
		return nil, utils.ErrorPerStatusCode(httpResp, err)
	}

	return &edgeApplicationsResponse.Results, nil
}

func (c *Client) UpdateInstance(ctx context.Context, req *UpdateInstanceRequest) (EdgeApplicationsResponse, error) {
	request := c.apiClient.EdgeApplicationsEdgeFunctionsInstancesApi.EdgeApplicationsEdgeApplicationIdFunctionsInstancesFunctionsInstancesIdPatch(ctx, req.Id, req.IdInstace).ApplicationUpdateInstanceRequest(req.ApplicationUpdateInstanceRequest)

	req.ApplicationUpdateInstanceRequest.SetName("justfortests2")
	req.SetEdgeFunctionId(req.FunctionId)

	edgeApplicationsResponse, httpResp, err := request.Execute()
	if err != nil {
		return nil, utils.ErrorPerStatusCode(httpResp, err)
	}

	return edgeApplicationsResponse.Results, nil
}

func (c *Client) CreateInstance(ctx context.Context, req *CreateInstanceRequest) (EdgeApplicationsResponse, error) {

	args := make(map[string]interface{})
	req.SetArgs(args)

	request := c.apiClient.EdgeApplicationsEdgeFunctionsInstancesApi.EdgeApplicationsEdgeApplicationIdFunctionsInstancesPost(ctx, req.ApplicationId).ApplicationCreateInstanceRequest(req.ApplicationCreateInstanceRequest)

	edgeApplicationsResponse, httpResp, err := request.Execute()
	if err != nil {
		return nil, utils.ErrorPerStatusCode(httpResp, err)
	}

	return edgeApplicationsResponse.Results, nil
}

func (c *Client) UpdateRulesEngine(ctx context.Context, req *UpdateRulesEngineRequest, idFunc int64) (EdgeApplicationsResponse, error) {

	request := c.apiClient.EdgeApplicationsRulesEngineApi.EdgeApplicationsEdgeApplicationIdRulesEnginePhaseRulesGet(ctx, req.IdApplication, "request")

	edgeApplicationRules, httpResp, err := request.Execute()
	if err != nil {
		return nil, utils.ErrorPerStatusCode(httpResp, err)
	}

	idRule := edgeApplicationRules.Results[0].Id

	b := make([]sdk.RulesEngineBehavior, 1)
	b[0].SetName("run_function")
	b[0].SetTarget(idFunc)
	req.SetBehaviors(b)

	requestUpdate := c.apiClient.EdgeApplicationsRulesEngineApi.EdgeApplicationsEdgeApplicationIdRulesEnginePhaseRulesRuleIdPatch(ctx, req.IdApplication, "request", idRule).PatchRulesEngineRequest(req.PatchRulesEngineRequest)

	edgeApplicationsResponse, httpResp, err := requestUpdate.Execute()
	if err != nil {
		return nil, utils.ErrorPerStatusCode(httpResp, err)
	}

	return &edgeApplicationsResponse.Results, nil
}

func (c *Client) Delete(ctx context.Context, id int64) error {
	str := strconv.FormatInt(id, 10)
	req := c.apiClient.EdgeApplicationsMainSettingsApi.EdgeApplicationsIdDelete(ctx, str)

	httpResp, err := req.Execute()

	if err != nil {
		return utils.ErrorPerStatusCode(httpResp, err)
	}

	return nil
}

func (c *Client) List(ctx context.Context, opts *contracts.ListOptions) (*sdk.GetApplicationsResponse, error) {
	if opts.OrderBy == "" {
		opts.OrderBy = "id"
	}

	resp, httpResp, err := c.apiClient.EdgeApplicationsMainSettingsApi.EdgeApplicationsGet(ctx).
		OrderBy(opts.OrderBy).
		Page(opts.Page).
		PageSize(opts.PageSize).
		Sort(opts.Sort).Execute()

	if err != nil {
		return &sdk.GetApplicationsResponse{}, utils.ErrorPerStatusCode(httpResp, err)
	}

	return resp, nil
}

type CreateOriginsRequest struct {
	sdk.CreateOriginsRequest
}

type UpdateOriginsRequest struct {
	sdk.PatchOriginsRequest
}

type OriginsResponse interface {
	GetOriginKey() string
	GetOriginId() int64
	GetName() string
}

func (c *Client) GetOrigin(ctx context.Context, edgeApplicationID, originID int64) (sdk.OriginsResultResponse, error) {
	resp, httpResp, err := c.apiClient.EdgeApplicationsOriginsApi.EdgeApplicationsEdgeApplicationIdOriginsGet(ctx, edgeApplicationID).Execute()
	if err != nil {
		return sdk.OriginsResultResponse{}, utils.ErrorPerStatusCode(httpResp, err)
	}
	if len(resp.Results) > 0 {
		for _, result := range resp.Results {
			if result.OriginId == originID {
				return result, nil
			}
		}
	}
	return sdk.OriginsResultResponse{}, utils.ErrorPerStatusCode(&http.Response{Status: "404 Not Found", StatusCode: http.StatusNotFound}, errors.New("404 Not Found"))
}

func (c *Client) ListOrigins(ctx context.Context, opts *contracts.ListOptions, edgeApplicationID int64) (*sdk.OriginsResponse, error) {
	resp, httpResp, err := c.apiClient.EdgeApplicationsOriginsApi.EdgeApplicationsEdgeApplicationIdOriginsGet(ctx, edgeApplicationID).Execute()
	if err != nil {
		return &sdk.OriginsResponse{}, utils.ErrorPerStatusCode(httpResp, err)
	}
	return resp, nil
}

func (c *Client) CreateOrigins(ctx context.Context, edgeApplicationID int64, req *CreateOriginsRequest) (OriginsResponse, error) {
	resp, httpResp, err := c.apiClient.EdgeApplicationsOriginsApi.EdgeApplicationsEdgeApplicationIdOriginsPost(ctx, edgeApplicationID).CreateOriginsRequest(req.CreateOriginsRequest).Execute()
	if err != nil {
		return nil, utils.ErrorPerStatusCode(httpResp, err)
	}
	return &resp.Results, nil
}

func (c *Client) UpdateOrigins(ctx context.Context, edgeApplicationID int64, originKey string, req *UpdateOriginsRequest) (OriginsResponse, error) {
	resp, httpResp, err := c.apiClient.EdgeApplicationsOriginsApi.
		EdgeApplicationsEdgeApplicationIdOriginsOriginKeyPatch(ctx, edgeApplicationID, originKey).PatchOriginsRequest(req.PatchOriginsRequest).Execute()
	if err != nil {
		return nil, utils.ErrorPerStatusCode(httpResp, err)
	}
	return &resp.Results, nil
}

func (c *Client) DeleteOrigins(ctx context.Context, edgeApplicationID int64, originKey string) error {
	httpResp, err := c.apiClient.EdgeApplicationsOriginsApi.EdgeApplicationsEdgeApplicationIdOriginsOriginKeyDelete(ctx, edgeApplicationID, originKey).Execute()
	if err != nil {
		return utils.ErrorPerStatusCode(httpResp, err)
	}
	return nil
}
