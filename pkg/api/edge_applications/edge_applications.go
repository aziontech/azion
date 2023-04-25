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

type CacheSettingsResponse interface {
	GetId() int64
	GetName() string
	GetBrowserCacheSettings() string
	GetBrowserCacheSettingsMaximumTtl() int64
	GetCdnCacheSettingsMaximumTtl() int64
	GetCdnCacheSettings() string
	GetCacheByQueryString() string
	GetQueryStringFields() []string
	GetEnableQueryStringSort() bool
	GetCacheByCookies() string
	GetCookieNames() []string
	GetEnableCachingForPost() bool
	GetL2CachingEnabled() bool
	GetAdaptiveDeliveryAction() string
	GetDeviceGroup() []string
}

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
	GetHttpPort() interface{}
	GetHttpsPort() interface{}
	GetImageOptimization() bool
	GetL2Caching() bool
	GetLoadBalancer() bool
	GetMinimumTlsVersion() string
	GetNext() string
	GetRawLogs() bool
	GetWebApplicationFirewall() bool
}

type RulesEngineResponse interface {
	GetId() int64
	GetPhase() string
	GetBehaviors() []sdk.RulesEngineResultResponseBehaviors
	GetCriteria() [][]sdk.RulesEngineCriteria
	GetIsActive() bool
	GetOrder() int64
	GetName() string
}

type DeviceGroupsResponse interface {
	GetId() int64
	GetName() string
	GetUserAgent() string
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
	Phase         string
	Id            int64
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

func (c *Client) UpdateInstance(ctx context.Context, req *UpdateInstanceRequest, appID string, instanceID string) (FunctionsInstancesResponse, error) {
	request := c.apiClient.EdgeApplicationsEdgeFunctionsInstancesApi.EdgeApplicationsEdgeApplicationIdFunctionsInstancesFunctionsInstancesIdPatch(ctx, appID, instanceID).ApplicationUpdateInstanceRequest(req.ApplicationUpdateInstanceRequest)
	edgeApplicationsResponse, httpResp, err := request.Execute()
	if err != nil {
		return nil, utils.ErrorPerStatusCode(httpResp, err)
	}

	return edgeApplicationsResponse.Results, nil
}

func (c *Client) CreateInstancePublish(ctx context.Context, req *CreateInstanceRequest) (EdgeApplicationsResponse, error) {

	args := make(map[string]interface{})
	req.SetArgs(args)

	request := c.apiClient.EdgeApplicationsEdgeFunctionsInstancesApi.EdgeApplicationsEdgeApplicationIdFunctionsInstancesPost(ctx, req.ApplicationId).ApplicationCreateInstanceRequest(req.ApplicationCreateInstanceRequest)

	edgeApplicationsResponse, httpResp, err := request.Execute()
	if err != nil {
		return nil, utils.ErrorPerStatusCode(httpResp, err)
	}

	return edgeApplicationsResponse.Results, nil
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

type CreateCacheSettingsRequest struct {
	sdk.ApplicationCacheCreateRequest
}

type UpdateCacheSettingsRequest struct {
	sdk.ApplicationCachePatchRequest
	Id int64
}

func (c *Client) CreateCacheSettings(ctx context.Context, req *CreateCacheSettingsRequest, applicationId int64) (CacheSettingsResponse, error) {

	request := c.apiClient.EdgeApplicationsCacheSettingsApi.EdgeApplicationsEdgeApplicationIdCacheSettingsPost(ctx, applicationId).ApplicationCacheCreateRequest(req.ApplicationCacheCreateRequest)

	cacheResponse, httpResp, err := request.Execute()
	if err != nil {
		return nil, utils.ErrorPerStatusCode(httpResp, err)
	}

	return cacheResponse.Results, nil
}

func (c *Client) UpdateCacheSettings(ctx context.Context, req *UpdateCacheSettingsRequest, applicationId int64) (CacheSettingsResponse, error) {

	request := c.apiClient.EdgeApplicationsCacheSettingsApi.EdgeApplicationsEdgeApplicationIdCacheSettingsCacheSettingsPatch(ctx, applicationId, req.Id).ApplicationCachePatchRequest(req.ApplicationCachePatchRequest)

	cacheResponse, httpResp, err := request.Execute()
	if err != nil {
		return nil, utils.ErrorPerStatusCode(httpResp, err)
	}

	return cacheResponse.Results, nil
}

func (c *Client) ListCacheSettings(ctx context.Context, opts *contracts.ListOptions, edgeApplicationID int64) (*sdk.ApplicationCacheGetResponse, error) {
	if opts.OrderBy == "" {
		opts.OrderBy = "id"
	}

	resp, httpResp, err := c.apiClient.EdgeApplicationsCacheSettingsApi.EdgeApplicationsEdgeApplicationIdCacheSettingsGet(ctx, edgeApplicationID).
		OrderBy(opts.OrderBy).
		Page(opts.Page).
		PageSize(opts.PageSize).
		Sort(opts.Sort).Execute()

	if err != nil {
		return &sdk.ApplicationCacheGetResponse{}, utils.ErrorPerStatusCode(httpResp, err)
	}

	return resp, nil
}

func (c *Client) GetCacheSettings(ctx context.Context, edgeApplicationID, cacheSettingsID int64) (CacheSettingsResponse, error) {
	resp, httpResp, err := c.apiClient.EdgeApplicationsCacheSettingsApi.EdgeApplicationsEdgeApplicationIdCacheSettingsCacheSettingsIdGet(ctx, edgeApplicationID, cacheSettingsID).Execute()
	if err != nil {
		return nil, utils.ErrorPerStatusCode(httpResp, err)
	}

	return &resp.Results, nil
}

func (c *Client) DeleteCacheSettings(ctx context.Context, edgeApplicationID, cacheSettingsID int64) error {
	httpResp, err := c.apiClient.EdgeApplicationsCacheSettingsApi.EdgeApplicationsEdgeApplicationIdCacheSettingsCacheSettingsDelete(ctx, edgeApplicationID, cacheSettingsID).Execute()
	if err != nil {
		return utils.ErrorPerStatusCode(httpResp, err)
	}
	return nil
}

type CreateRulesEngineRequest struct {
	sdk.CreateRulesEngineRequest
}

func (c *Client) ListRulesEngine(ctx context.Context, opts *contracts.ListOptions, edgeApplicationID int64, phase string) (*sdk.RulesEngineResponse, error) {
	if opts.OrderBy == "" {
		opts.OrderBy = "id"
	}

	resp, httpResp, err := c.apiClient.EdgeApplicationsRulesEngineApi.EdgeApplicationsEdgeApplicationIdRulesEnginePhaseRulesGet(ctx, edgeApplicationID, phase).
		OrderBy(opts.OrderBy).
		Page(opts.Page).
		PageSize(opts.PageSize).
		Sort(opts.Sort).Execute()

	if err != nil {
		return &sdk.RulesEngineResponse{}, utils.ErrorPerStatusCode(httpResp, err)
	}

	return resp, nil
}

func (c *Client) GetRulesEngine(ctx context.Context, edgeApplicationID, rulesID int64, phase string) (RulesEngineResponse, error) {
	resp, httpResp, err := c.apiClient.EdgeApplicationsRulesEngineApi.EdgeApplicationsEdgeApplicationIdRulesEnginePhaseRulesRuleIdGet(ctx, edgeApplicationID, phase, rulesID).Execute()
	if err != nil {
		return nil, utils.ErrorPerStatusCode(httpResp, err)
	}
	return &resp.Results, nil
}

func (c *Client) DeleteRulesEngine(ctx context.Context, edgeApplicationID int64, phase string, ruleID int64) error {
	httpResp, err := c.apiClient.EdgeApplicationsRulesEngineApi.EdgeApplicationsEdgeApplicationIdRulesEnginePhaseRulesRuleIdDelete(ctx, edgeApplicationID, phase, ruleID).Execute()
	if err != nil {
		return utils.ErrorPerStatusCode(httpResp, err)
	}
	return nil
}

func (c *Client) UpdateRulesEnginePublish(ctx context.Context, req *UpdateRulesEngineRequest, idFunc int64) (EdgeApplicationsResponse, error) {

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

func (c *Client) UpdateRulesEngine(ctx context.Context, req *UpdateRulesEngineRequest) (RulesEngineResponse, error) {

	requestUpdate := c.apiClient.EdgeApplicationsRulesEngineApi.EdgeApplicationsEdgeApplicationIdRulesEnginePhaseRulesRuleIdPatch(ctx, req.IdApplication, req.Phase, req.Id).PatchRulesEngineRequest(req.PatchRulesEngineRequest)

	edgeApplicationsResponse, httpResp, err := requestUpdate.Execute()
	if err != nil {
		return nil, utils.ErrorPerStatusCode(httpResp, err)
	}

	return &edgeApplicationsResponse.Results, nil

}

func (c *Client) CreateRulesEngine(ctx context.Context, edgeApplicationID int64, phase string, req *CreateRulesEngineRequest) (RulesEngineResponse, error) {
	resp, httpResp, err := c.apiClient.EdgeApplicationsRulesEngineApi.
		EdgeApplicationsEdgeApplicationIdRulesEnginePhaseRulesPost(ctx, edgeApplicationID, phase).
		CreateRulesEngineRequest(req.CreateRulesEngineRequest).Execute()
	if err != nil {
		return nil, utils.ErrorPerStatusCode(httpResp, err)
	}
	return &resp.Results, nil
}

type FunctionsInstancesResponse interface {
	GetId() int64
	GetEdgeFunctionId() int64
	GetName() string
	GetArgs() interface{}
}

type CreateFuncInstancesRequest struct {
	sdk.ApplicationCreateInstanceRequest
}

func (c *Client) EdgeFuncInstancesList(ctx context.Context, opts *contracts.ListOptions, edgeApplicationID int64) (*sdk.ApplicationInstancesGetResponse, error) {
	if opts.OrderBy == "" {
		opts.OrderBy = "id"
	}

	resp, httpResp, err := c.apiClient.EdgeApplicationsEdgeFunctionsInstancesApi.
		EdgeApplicationsEdgeApplicationIdFunctionsInstancesGet(ctx, edgeApplicationID).
		OrderBy(opts.OrderBy).
		Page(opts.Page).
		PageSize(opts.PageSize).
		Sort(opts.Sort).Execute()

	if err != nil {
		return nil, utils.ErrorPerStatusCode(httpResp, err)
	}
	return resp, nil
}

func (c *Client) DeleteFunctionInstance(ctx context.Context, appID string, funcID string) error {
	req := c.apiClient.EdgeApplicationsEdgeFunctionsInstancesApi.EdgeApplicationsEdgeApplicationIdFunctionsInstancesFunctionsInstancesIdDelete(ctx, appID, funcID)

	httpResp, err := req.Execute()

	if err != nil {
		return utils.ErrorPerStatusCode(httpResp, err)
	}

	return nil
}

func (c *Client) CreateFuncInstances(ctx context.Context, req *CreateFuncInstancesRequest, applicationID int64) (FunctionsInstancesResponse, error) {
	resp, httpResp, err := c.apiClient.EdgeApplicationsEdgeFunctionsInstancesApi.EdgeApplicationsEdgeApplicationIdFunctionsInstancesPost(ctx, applicationID).
		ApplicationCreateInstanceRequest(req.ApplicationCreateInstanceRequest).Execute()
	if err != nil {
		return nil, utils.ErrorPerStatusCode(httpResp, err)
	}
	return resp.Results, nil
}

func (c *Client) GetFuncInstance(ctx context.Context, edgeApplicationID, instanceID int64) (FunctionsInstancesResponse, error) {
	resp, httpResp, err := c.apiClient.EdgeApplicationsEdgeFunctionsInstancesApi.EdgeApplicationsEdgeApplicationIdFunctionsInstancesFunctionsInstancesIdGet(ctx, edgeApplicationID, instanceID).Execute()
	if err != nil {
		return nil, utils.ErrorPerStatusCode(httpResp, err)
	}
	return &resp.Results, nil
}

type UpdateDeviceGroupRequest struct {
	sdk.PatchDeviceGroupsRequest
}

func (c *Client) DeviceGroupsList(ctx context.Context, opts *contracts.ListOptions, edgeApplicationID int64) (*sdk.DeviceGroupsResponse, error) {
	if opts.OrderBy == "" {
		opts.OrderBy = "id"
	}
	resp, httpResp, err := c.apiClient.EdgeApplicationsDeviceGroupsApi.
		EdgeApplicationsEdgeApplicationIdDeviceGroupsGet(ctx, edgeApplicationID).
		OrderBy(opts.OrderBy).
		Page(opts.Page).
		PageSize(opts.PageSize).
		Sort(opts.Sort).Execute()

	if err != nil {
		return nil, utils.ErrorPerStatusCode(httpResp, err)
	}
	return resp, nil
}

func (c *Client) DeleteDeviceGroup(ctx context.Context, appID int64, groupID int64) error {
	req := c.apiClient.EdgeApplicationsDeviceGroupsApi.EdgeApplicationsEdgeApplicationIdDeviceGroupsDeviceGroupIdDelete(ctx, appID, groupID)

	httpResp, err := req.Execute()

	if err != nil {
		return utils.ErrorPerStatusCode(httpResp, err)
	}

	return nil
}

func (c *Client) GetDeviceGroups(ctx context.Context, edgeApplicationID, groupID int64) (DeviceGroupsResponse, error) {
	resp, httpResp, err := c.apiClient.EdgeApplicationsDeviceGroupsApi.EdgeApplicationsEdgeApplicationIdDeviceGroupsDeviceGroupIdGet(ctx, edgeApplicationID, groupID).Execute()
	if err != nil {
		return nil, utils.ErrorPerStatusCode(httpResp, err)
	}
	return &resp.Results, nil
}

func (c *Client) UpdateDeviceGroup(ctx context.Context, req sdk.PatchDeviceGroupsRequest, appID int64, groupID int64) (DeviceGroupsResponse, error) {
	request := c.apiClient.EdgeApplicationsDeviceGroupsApi.EdgeApplicationsEdgeApplicationIdDeviceGroupsDeviceGroupIdPatch(ctx, appID, groupID).PatchDeviceGroupsRequest(req)
	deviceGroup, httpResp, err := request.Execute()
	if err != nil {
		return nil, utils.ErrorPerStatusCode(httpResp, err)
	}

	return &deviceGroup.Results, nil
}
