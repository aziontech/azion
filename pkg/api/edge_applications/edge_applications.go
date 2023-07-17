package edgeapplications

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/aziontech/azion-cli/pkg/cmd/version"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/utils"
	sdk "github.com/aziontech/azionapi-go-sdk/edgeapplications"
	"go.uber.org/zap"
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

type CreateCacheSettingsRequest struct {
	sdk.ApplicationCacheCreateRequest
}

type UpdateCacheSettingsRequest struct {
	sdk.ApplicationCachePatchRequest
	Id int64
}

type CreateRulesEngineRequest struct {
	sdk.CreateRulesEngineRequest
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

type CreateDeviceGroupsRequest struct {
	sdk.CreateDeviceGroupsRequest
}

type DeviceGroupsResponse interface {
	GetId() int64
	GetName() string
	GetUserAgent() string
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
	logger.Debug("Get Edge Application")
	req := c.apiClient.EdgeApplicationsMainSettingsApi.EdgeApplicationsIdGet(ctx, id)

	res, httpResp, err := req.Execute()
	if err != nil {
		logger.Debug("Error while getting an edge application", zap.Error(err))
		logger.Debug("Status Code", zap.Any("http", httpResp.StatusCode))
		logger.Debug("Headers", zap.Any("http", httpResp.Header))
		logger.Debug("Response body", zap.Any("http", httpResp.Body))
		return nil, utils.ErrorPerStatusCode(httpResp, err)
	}

	return &res.Results, nil
}

func (c *Client) Create(ctx context.Context, req *CreateRequest) (EdgeApplicationsResponse, error) {
	logger.Debug("Create Edge Application")
	request := c.apiClient.EdgeApplicationsMainSettingsApi.EdgeApplicationsPost(ctx).CreateApplicationRequest(req.CreateApplicationRequest)

	edgeApplicationsResponse, httpResp, err := request.Execute()
	if err != nil {
		logger.Debug("Error while creating an edge application", zap.Error(err))
		logger.Debug("Status Code", zap.Any("http", httpResp.StatusCode))
		logger.Debug("Headers", zap.Any("http", httpResp.Header))
		logger.Debug("Response body", zap.Any("http", httpResp.Body))
		return nil, utils.ErrorPerStatusCode(httpResp, err)
	}

	return &edgeApplicationsResponse.Results, nil
}

func (c *Client) Update(ctx context.Context, req *UpdateRequest) (EdgeApplicationsResponse, error) {
	logger.Debug("Update Edge Application")
	str := strconv.FormatInt(req.Id, 10)
	request := c.apiClient.EdgeApplicationsMainSettingsApi.EdgeApplicationsIdPatch(ctx, str).ApplicationUpdateRequest(req.ApplicationUpdateRequest)

	edgeApplicationsResponse, httpResp, err := request.Execute()
	if err != nil {
		logger.Debug("Error while updating an edge application", zap.Error(err))
		logger.Debug("Status Code", zap.Any("http", httpResp.StatusCode))
		logger.Debug("Headers", zap.Any("http", httpResp.Header))
		logger.Debug("Response body", zap.Any("http", httpResp.Body))
		return nil, utils.ErrorPerStatusCode(httpResp, err)
	}

	return &edgeApplicationsResponse.Results, nil
}

func (c *Client) UpdateInstance(ctx context.Context, req *UpdateInstanceRequest, appID string, instanceID string) (FunctionsInstancesResponse, error) {
	logger.Debug("Update Instance")
	request := c.apiClient.EdgeApplicationsEdgeFunctionsInstancesApi.EdgeApplicationsEdgeApplicationIdFunctionsInstancesFunctionsInstancesIdPatch(ctx, appID, instanceID).ApplicationUpdateInstanceRequest(req.ApplicationUpdateInstanceRequest)

	edgeApplicationsResponse, httpResp, err := request.Execute()
	if err != nil {
		logger.Debug("Error while updating an edge function instance", zap.Error(err))
		logger.Debug("Status Code", zap.Any("http", httpResp.StatusCode))
		logger.Debug("Headers", zap.Any("http", httpResp.Header))
		logger.Debug("Response body", zap.Any("http", httpResp.Body))
		return nil, utils.ErrorPerStatusCode(httpResp, err)
	}

	return edgeApplicationsResponse.Results, nil
}

func (c *Client) CreateInstancePublish(ctx context.Context, req *CreateInstanceRequest) (EdgeApplicationsResponse, error) {
	logger.Debug("Create Instance Publish")
	args := make(map[string]interface{})
	req.SetArgs(args)

	request := c.apiClient.EdgeApplicationsEdgeFunctionsInstancesApi.EdgeApplicationsEdgeApplicationIdFunctionsInstancesPost(ctx, req.ApplicationId).ApplicationCreateInstanceRequest(req.ApplicationCreateInstanceRequest)

	edgeApplicationsResponse, httpResp, err := request.Execute()
	if err != nil {
		logger.Debug("Error while creating an edge function instance", zap.Error(err))
		logger.Debug("Status Code", zap.Any("http", httpResp.StatusCode))
		logger.Debug("Headers", zap.Any("http", httpResp.Header))
		logger.Debug("Response body", zap.Any("http", httpResp.Body))
		return nil, utils.ErrorPerStatusCode(httpResp, err)
	}

	return edgeApplicationsResponse.Results, nil
}

func (c *Client) Delete(ctx context.Context, id int64) error {
	logger.Debug("Delete Edge Application")
	str := strconv.FormatInt(id, 10)
	req := c.apiClient.EdgeApplicationsMainSettingsApi.EdgeApplicationsIdDelete(ctx, str)

	httpResp, err := req.Execute()
	if err != nil {
		logger.Debug("Error while deleting an edge application", zap.Error(err))
		logger.Debug("Status Code", zap.Any("http", httpResp.StatusCode))
		logger.Debug("Headers", zap.Any("http", httpResp.Header))
		logger.Debug("Response body", zap.Any("http", httpResp.Body))
		return utils.ErrorPerStatusCode(httpResp, err)
	}

	return nil
}

func (c *Client) List(ctx context.Context, opts *contracts.ListOptions) (*sdk.GetApplicationsResponse, error) {
	logger.Debug("List Edge Application")
	if opts.OrderBy == "" {
		opts.OrderBy = "id"
	}

	resp, httpResp, err := c.apiClient.EdgeApplicationsMainSettingsApi.EdgeApplicationsGet(ctx).
		OrderBy(opts.OrderBy).
		Page(opts.Page).
		PageSize(opts.PageSize).
		Sort(opts.Sort).Execute()

	if err != nil {
		logger.Debug("Error while listing edge applications", zap.Error(err))
		logger.Debug("Status Code", zap.Any("http", httpResp.StatusCode))
		logger.Debug("Headers", zap.Any("http", httpResp.Header))
		logger.Debug("Response body", zap.Any("http", httpResp.Body))
		return &sdk.GetApplicationsResponse{}, utils.ErrorPerStatusCode(httpResp, err)
	}

	return resp, nil
}

func (c *Client) GetOrigin(ctx context.Context, edgeApplicationID, originID int64) (sdk.OriginsResultResponse, error) {
	logger.Debug("Get Origin")
	resp, httpResp, err := c.apiClient.EdgeApplicationsOriginsApi.EdgeApplicationsEdgeApplicationIdOriginsGet(ctx, edgeApplicationID).Execute()
	if err != nil {
		logger.Debug("Error while getting an origin", zap.Error(err))
		logger.Debug("Status Code", zap.Any("http", httpResp.StatusCode))
		logger.Debug("Headers", zap.Any("http", httpResp.Header))
		logger.Debug("Response body", zap.Any("http", httpResp.Body))
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
	logger.Debug("List Origins")
	resp, httpResp, err := c.apiClient.EdgeApplicationsOriginsApi.EdgeApplicationsEdgeApplicationIdOriginsGet(ctx, edgeApplicationID).Execute()
	if err != nil {
		logger.Debug("Error while listing origins", zap.Error(err))
		logger.Debug("Status Code", zap.Any("http", httpResp.StatusCode))
		logger.Debug("Headers", zap.Any("http", httpResp.Header))
		logger.Debug("Response body", zap.Any("http", httpResp.Body))
		return &sdk.OriginsResponse{}, utils.ErrorPerStatusCode(httpResp, err)
	}
	return resp, nil
}

func (c *Client) CreateOrigins(ctx context.Context, edgeApplicationID int64, req *CreateOriginsRequest) (OriginsResponse, error) {
	logger.Debug("Create Origins")
	resp, httpResp, err := c.apiClient.EdgeApplicationsOriginsApi.EdgeApplicationsEdgeApplicationIdOriginsPost(ctx, edgeApplicationID).CreateOriginsRequest(req.CreateOriginsRequest).Execute()
	if err != nil {
		logger.Debug("Error while creating an origin", zap.Error(err))
		logger.Debug("Status Code", zap.Any("http", httpResp.StatusCode))
		logger.Debug("Headers", zap.Any("http", httpResp.Header))
		logger.Debug("Response body", zap.Any("http", httpResp.Body))
		return nil, utils.ErrorPerStatusCode(httpResp, err)
	}
	return &resp.Results, nil
}

func (c *Client) UpdateOrigins(ctx context.Context, edgeApplicationID int64, originKey string, req *UpdateOriginsRequest) (OriginsResponse, error) {
	logger.Debug("Update Origins")
	resp, httpResp, err := c.apiClient.EdgeApplicationsOriginsApi.
		EdgeApplicationsEdgeApplicationIdOriginsOriginKeyPatch(ctx, edgeApplicationID, originKey).PatchOriginsRequest(req.PatchOriginsRequest).Execute()
	if err != nil {
		logger.Debug("Error while updating an origin", zap.Error(err))
		logger.Debug("Status Code", zap.Any("http", httpResp.StatusCode))
		logger.Debug("Headers", zap.Any("http", httpResp.Header))
		logger.Debug("Response body", zap.Any("http", httpResp.Body))
		return nil, utils.ErrorPerStatusCode(httpResp, err)
	}
	return &resp.Results, nil
}

func (c *Client) DeleteOrigins(ctx context.Context, edgeApplicationID int64, originKey string) error {
	logger.Debug("Delete Origins")
	httpResp, err := c.apiClient.EdgeApplicationsOriginsApi.EdgeApplicationsEdgeApplicationIdOriginsOriginKeyDelete(ctx, edgeApplicationID, originKey).Execute()
	if err != nil {
		logger.Debug("Error while deleting an origin", zap.Error(err))
		logger.Debug("Status Code", zap.Any("http", httpResp.StatusCode))
		logger.Debug("Headers", zap.Any("http", httpResp.Header))
		logger.Debug("Response body", zap.Any("http", httpResp.Body))
		return utils.ErrorPerStatusCode(httpResp, err)
	}
	return nil
}

func (c *Client) CreateCacheSettings(ctx context.Context, req *CreateCacheSettingsRequest, applicationId int64) (CacheSettingsResponse, error) {
	logger.Debug("Create Cache Settings")
	request := c.apiClient.EdgeApplicationsCacheSettingsApi.EdgeApplicationsEdgeApplicationIdCacheSettingsPost(ctx, applicationId).ApplicationCacheCreateRequest(req.ApplicationCacheCreateRequest)

	cacheResponse, httpResp, err := request.Execute()
	if err != nil {
		logger.Debug("Error while creating a cache setting", zap.Error(err))
		logger.Debug("Status Code", zap.Any("http", httpResp.StatusCode))
		logger.Debug("Headers", zap.Any("http", httpResp.Header))
		logger.Debug("Response body", zap.Any("http", httpResp.Body))
		return nil, utils.ErrorPerStatusCode(httpResp, err)
	}

	return cacheResponse.Results, nil
}

func (c *Client) UpdateCacheSettings(ctx context.Context, req *UpdateCacheSettingsRequest, applicationId int64) (CacheSettingsResponse, error) {
	logger.Debug("Update Cache Settings")
	request := c.apiClient.EdgeApplicationsCacheSettingsApi.EdgeApplicationsEdgeApplicationIdCacheSettingsCacheSettingsIdPatch(ctx, applicationId, req.Id).ApplicationCachePatchRequest(req.ApplicationCachePatchRequest)

	cacheResponse, httpResp, err := request.Execute()
	if err != nil {
		logger.Debug("Error while updating a cache setting", zap.Error(err))
		logger.Debug("Status Code", zap.Any("http", httpResp.StatusCode))
		logger.Debug("Headers", zap.Any("http", httpResp.Header))
		logger.Debug("Response body", zap.Any("http", httpResp.Body))
		return nil, utils.ErrorPerStatusCode(httpResp, err)
	}

	return cacheResponse.Results, nil
}

func (c *Client) ListCacheSettings(ctx context.Context, opts *contracts.ListOptions, edgeApplicationID int64) (*sdk.ApplicationCacheGetResponse, error) {
	logger.Debug("List Cache Settings")
	if opts.OrderBy == "" {
		opts.OrderBy = "id"
	}

	resp, httpResp, err := c.apiClient.EdgeApplicationsCacheSettingsApi.EdgeApplicationsEdgeApplicationIdCacheSettingsGet(ctx, edgeApplicationID).
		OrderBy(opts.OrderBy).
		Page(opts.Page).
		PageSize(opts.PageSize).
		Sort(opts.Sort).Execute()

	if err != nil {
		logger.Debug("Error while listing cache settings", zap.Error(err))
		logger.Debug("Status Code", zap.Any("http", httpResp.StatusCode))
		logger.Debug("Headers", zap.Any("http", httpResp.Header))
		logger.Debug("Response body", zap.Any("http", httpResp.Body))
		return &sdk.ApplicationCacheGetResponse{}, utils.ErrorPerStatusCode(httpResp, err)
	}

	return resp, nil
}

func (c *Client) GetCacheSettings(ctx context.Context, edgeApplicationID, cacheSettingsID int64) (CacheSettingsResponse, error) {
	logger.Debug("Get Cache Settings")
	resp, httpResp, err := c.apiClient.EdgeApplicationsCacheSettingsApi.EdgeApplicationsEdgeApplicationIdCacheSettingsCacheSettingsIdGet(ctx, edgeApplicationID, cacheSettingsID).Execute()
	if err != nil {
		logger.Debug("Error while getting a cache setting", zap.Error(err))
		logger.Debug("Status Code", zap.Any("http", httpResp.StatusCode))
		logger.Debug("Headers", zap.Any("http", httpResp.Header))
		logger.Debug("Response body", zap.Any("http", httpResp.Body))
		return nil, utils.ErrorPerStatusCode(httpResp, err)
	}

	return &resp.Results, nil
}

func (c *Client) DeleteCacheSettings(ctx context.Context, edgeApplicationID, cacheSettingsID int64) error {
	logger.Debug("Delete Cache Settings")
	httpResp, err := c.apiClient.EdgeApplicationsCacheSettingsApi.EdgeApplicationsEdgeApplicationIdCacheSettingsCacheSettingsIdDelete(ctx, edgeApplicationID, cacheSettingsID).Execute()
	if err != nil {
		logger.Debug("Error while deleting a cache setting", zap.Error(err))
		logger.Debug("Status Code", zap.Any("http", httpResp.StatusCode))
		logger.Debug("Headers", zap.Any("http", httpResp.Header))
		logger.Debug("Response body", zap.Any("http", httpResp.Body))
		return utils.ErrorPerStatusCode(httpResp, err)
	}
	return nil
}

func (c *Client) ListRulesEngine(ctx context.Context, opts *contracts.ListOptions, edgeApplicationID int64, phase string) (*sdk.RulesEngineResponse, error) {
	logger.Debug("List Rules Engine")
	if opts.OrderBy == "" {
		opts.OrderBy = "id"
	}

	resp, httpResp, err := c.apiClient.EdgeApplicationsRulesEngineApi.EdgeApplicationsEdgeApplicationIdRulesEnginePhaseRulesGet(ctx, edgeApplicationID, phase).
		OrderBy(opts.OrderBy).
		Page(opts.Page).
		PageSize(opts.PageSize).
		Sort(opts.Sort).Execute()

	if err != nil {
		logger.Debug("Error while listing rules in rules engines", zap.Error(err))
		logger.Debug("Status Code", zap.Any("http", httpResp.StatusCode))
		logger.Debug("Headers", zap.Any("http", httpResp.Header))
		logger.Debug("Response body", zap.Any("http", httpResp.Body))
		return &sdk.RulesEngineResponse{}, utils.ErrorPerStatusCode(httpResp, err)
	}

	return resp, nil
}

func (c *Client) GetRulesEngine(ctx context.Context, edgeApplicationID, rulesID int64, phase string) (RulesEngineResponse, error) {
	logger.Debug("Get Rules Engine")
	resp, httpResp, err := c.apiClient.EdgeApplicationsRulesEngineApi.EdgeApplicationsEdgeApplicationIdRulesEnginePhaseRulesRuleIdGet(ctx, edgeApplicationID, phase, rulesID).Execute()
	if err != nil {
		logger.Debug("Error while getting a rule in rules engine", zap.Error(err))
		logger.Debug("Status Code", zap.Any("http", httpResp.StatusCode))
		logger.Debug("Headers", zap.Any("http", httpResp.Header))
		logger.Debug("Response body", zap.Any("http", httpResp.Body))
		return nil, utils.ErrorPerStatusCode(httpResp, err)
	}
	return &resp.Results, nil
}

func (c *Client) DeleteRulesEngine(ctx context.Context, edgeApplicationID int64, phase string, ruleID int64) error {
	logger.Debug("Delete Rules Engine")
	httpResp, err := c.apiClient.EdgeApplicationsRulesEngineApi.EdgeApplicationsEdgeApplicationIdRulesEnginePhaseRulesRuleIdDelete(ctx, edgeApplicationID, phase, ruleID).Execute()
	if err != nil {
		logger.Debug("Error while deleting a rule in rules engine", zap.Error(err))
		logger.Debug("Status Code", zap.Any("http", httpResp.StatusCode))
		logger.Debug("Headers", zap.Any("http", httpResp.Header))
		logger.Debug("Response body", zap.Any("http", httpResp.Body))
		return utils.ErrorPerStatusCode(httpResp, err)
	}
	return nil
}

func (c *Client) UpdateRulesEnginePublish(ctx context.Context, req *UpdateRulesEngineRequest, idFunc int64) (EdgeApplicationsResponse, error) {
	logger.Debug("Update Rules Engine Publish")
	request := c.apiClient.EdgeApplicationsRulesEngineApi.EdgeApplicationsEdgeApplicationIdRulesEnginePhaseRulesGet(ctx, req.IdApplication, "request")

	edgeApplicationRules, httpResp, err := request.Execute()
	if err != nil {
		logger.Debug("Error while updating a rule in rules engine", zap.Error(err))
		logger.Debug("Status Code", zap.Any("http", httpResp.StatusCode))
		logger.Debug("Headers", zap.Any("http", httpResp.Header))
		logger.Debug("Response body", zap.Any("http", httpResp.Body))
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
		logger.Debug("Error while updating a rule in rules engine [run-function step]", zap.Error(err))
		logger.Debug("Status Code", zap.Any("http", httpResp.StatusCode))
		logger.Debug("Headers", zap.Any("http", httpResp.Header))
		logger.Debug("Response body", zap.Any("http", httpResp.Body))
		return nil, utils.ErrorPerStatusCode(httpResp, err)
	}

	return &edgeApplicationsResponse.Results, nil
}

func (c *Client) UpdateRulesEngine(ctx context.Context, req *UpdateRulesEngineRequest) (RulesEngineResponse, error) {
	logger.Debug("Update Rules Engine")
	requestUpdate := c.apiClient.EdgeApplicationsRulesEngineApi.EdgeApplicationsEdgeApplicationIdRulesEnginePhaseRulesRuleIdPatch(ctx, req.IdApplication, req.Phase, req.Id).PatchRulesEngineRequest(req.PatchRulesEngineRequest)

	edgeApplicationsResponse, httpResp, err := requestUpdate.Execute()
	if err != nil {
		logger.Debug("Error while updating a rule in rules engine", zap.Error(err))
		logger.Debug("Status Code", zap.Any("http", httpResp.StatusCode))
		logger.Debug("Headers", zap.Any("http", httpResp.Header))
		logger.Debug("Response body", zap.Any("http", httpResp.Body))
		return nil, utils.ErrorPerStatusCode(httpResp, err)
	}

	return &edgeApplicationsResponse.Results, nil

}

func (c *Client) CreateRulesEngine(ctx context.Context, edgeApplicationID int64, phase string, req *CreateRulesEngineRequest) (RulesEngineResponse, error) {
	logger.Debug("Create Rules Engine")
	resp, httpResp, err := c.apiClient.EdgeApplicationsRulesEngineApi.
		EdgeApplicationsEdgeApplicationIdRulesEnginePhaseRulesPost(ctx, edgeApplicationID, phase).
		CreateRulesEngineRequest(req.CreateRulesEngineRequest).Execute()
	if err != nil {
		logger.Debug("Error while creating a rule in rules engine", zap.Error(err))
		logger.Debug("Status Code", zap.Any("http", httpResp.StatusCode))
		logger.Debug("Headers", zap.Any("http", httpResp.Header))
		logger.Debug("Response body", zap.Any("http", httpResp.Body))
		return nil, utils.ErrorPerStatusCode(httpResp, err)
	}
	return &resp.Results, nil
}

func (c *Client) EdgeFuncInstancesList(ctx context.Context, opts *contracts.ListOptions, edgeApplicationID int64) (*sdk.ApplicationInstancesGetResponse, error) {
	logger.Debug("List Edge Function Instances")
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
		logger.Debug("Error while listing edge function instances", zap.Error(err))
		logger.Debug("Status Code", zap.Any("http", httpResp.StatusCode))
		logger.Debug("Headers", zap.Any("http", httpResp.Header))
		logger.Debug("Response body", zap.Any("http", httpResp.Body))
		return nil, utils.ErrorPerStatusCode(httpResp, err)
	}
	return resp, nil
}

func (c *Client) DeleteFunctionInstance(ctx context.Context, appID string, funcID string) error {
	logger.Debug("Delete Edge Function Instance")
	req := c.apiClient.EdgeApplicationsEdgeFunctionsInstancesApi.EdgeApplicationsEdgeApplicationIdFunctionsInstancesFunctionsInstancesIdDelete(ctx, appID, funcID)

	httpResp, err := req.Execute()
	if err != nil {
		logger.Debug("Error while deleting an edge function instance", zap.Error(err))
		logger.Debug("Status Code", zap.Any("http", httpResp.StatusCode))
		logger.Debug("Headers", zap.Any("http", httpResp.Header))
		logger.Debug("Response body", zap.Any("http", httpResp.Body))
		return utils.ErrorPerStatusCode(httpResp, err)
	}

	return nil
}

func (c *Client) CreateFuncInstances(ctx context.Context, req *CreateFuncInstancesRequest, applicationID int64) (FunctionsInstancesResponse, error) {
	logger.Debug("Create Edge Function Instance")
	resp, httpResp, err := c.apiClient.EdgeApplicationsEdgeFunctionsInstancesApi.EdgeApplicationsEdgeApplicationIdFunctionsInstancesPost(ctx, applicationID).
		ApplicationCreateInstanceRequest(req.ApplicationCreateInstanceRequest).Execute()
	if err != nil {
		logger.Debug("Error while creating an edge function instance", zap.Error(err))
		logger.Debug("Status Code", zap.Any("http", httpResp.StatusCode))
		logger.Debug("Headers", zap.Any("http", httpResp.Header))
		logger.Debug("Response body", zap.Any("http", httpResp.Body))
		return nil, utils.ErrorPerStatusCode(httpResp, err)
	}
	return resp.Results, nil
}

func (c *Client) GetFuncInstance(ctx context.Context, edgeApplicationID, instanceID int64) (FunctionsInstancesResponse, error) {
	logger.Debug("Get Edge Function Instance")
	resp, httpResp, err := c.apiClient.EdgeApplicationsEdgeFunctionsInstancesApi.EdgeApplicationsEdgeApplicationIdFunctionsInstancesFunctionsInstancesIdGet(ctx, edgeApplicationID, instanceID).Execute()
	if err != nil {
		logger.Debug("Error while getting an edge function instance", zap.Error(err))
		logger.Debug("Status Code", zap.Any("http", httpResp.StatusCode))
		logger.Debug("Headers", zap.Any("http", httpResp.Header))
		logger.Debug("Response body", zap.Any("http", httpResp.Body))
		return nil, utils.ErrorPerStatusCode(httpResp, err)
	}
	return &resp.Results, nil
}

func (c *Client) DeviceGroupsList(ctx context.Context, opts *contracts.ListOptions, edgeApplicationID int64) (*sdk.DeviceGroupsResponse, error) {
	logger.Debug("List Device Groups")
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
		logger.Debug("Error while listing device groups", zap.Error(err))
		logger.Debug("Status Code", zap.Any("http", httpResp.StatusCode))
		logger.Debug("Headers", zap.Any("http", httpResp.Header))
		logger.Debug("Response body", zap.Any("http", httpResp.Body))
		return nil, utils.ErrorPerStatusCode(httpResp, err)
	}
	return resp, nil
}

func (c *Client) DeleteDeviceGroup(ctx context.Context, appID int64, groupID int64) error {
	logger.Debug("Delete Device Group")
	req := c.apiClient.EdgeApplicationsDeviceGroupsApi.EdgeApplicationsEdgeApplicationIdDeviceGroupsDeviceGroupIdDelete(ctx, appID, groupID)

	httpResp, err := req.Execute()
	if err != nil {
		logger.Debug("Error while deleting a device group", zap.Error(err))
		logger.Debug("Status Code", zap.Any("http", httpResp.StatusCode))
		logger.Debug("Headers", zap.Any("http", httpResp.Header))
		logger.Debug("Response body", zap.Any("http", httpResp.Body))
		return utils.ErrorPerStatusCode(httpResp, err)
	}

	return nil
}

func (c *Client) GetDeviceGroups(ctx context.Context, edgeApplicationID, groupID int64) (DeviceGroupsResponse, error) {
	logger.Debug("Get Device Groups")
	resp, httpResp, err := c.apiClient.EdgeApplicationsDeviceGroupsApi.EdgeApplicationsEdgeApplicationIdDeviceGroupsDeviceGroupIdGet(ctx, edgeApplicationID, groupID).Execute()
	if err != nil {
		logger.Debug("Error while getting a device group", zap.Error(err))
		logger.Debug("Status Code", zap.Any("http", httpResp.StatusCode))
		logger.Debug("Headers", zap.Any("http", httpResp.Header))
		logger.Debug("Response body", zap.Any("http", httpResp.Body))
		return nil, utils.ErrorPerStatusCode(httpResp, err)
	}
	return &resp.Results, nil
}

func (c *Client) UpdateDeviceGroup(ctx context.Context, req sdk.PatchDeviceGroupsRequest, appID int64, groupID int64) (DeviceGroupsResponse, error) {
	logger.Debug("Update Device Group")
	request := c.apiClient.EdgeApplicationsDeviceGroupsApi.EdgeApplicationsEdgeApplicationIdDeviceGroupsDeviceGroupIdPatch(ctx, appID, groupID).PatchDeviceGroupsRequest(req)

	deviceGroup, httpResp, err := request.Execute()
	if err != nil {
		logger.Debug("Error while updating a device group", zap.Error(err))
		logger.Debug("Status Code", zap.Any("http", httpResp.StatusCode))
		logger.Debug("Headers", zap.Any("http", httpResp.Header))
		logger.Debug("Response body", zap.Any("http", httpResp.Body))
		return nil, utils.ErrorPerStatusCode(httpResp, err)
	}

	return &deviceGroup.Results, nil
}

func (c *Client) CreateDeviceGroups(ctx context.Context, req *CreateDeviceGroupsRequest, applicationID int64) (DeviceGroupsResponse, error) {
	logger.Debug("Create Device Groups")
	resp, httpResp, err := c.apiClient.EdgeApplicationsDeviceGroupsApi.EdgeApplicationsEdgeApplicationIdDeviceGroupsPost(ctx, applicationID).
		CreateDeviceGroupsRequest(req.CreateDeviceGroupsRequest).Execute()
	if err != nil {
		logger.Debug("Error while creating a device group", zap.Error(err))
		logger.Debug("Status Code", zap.Any("http", httpResp.StatusCode))
		logger.Debug("Headers", zap.Any("http", httpResp.Header))
		logger.Debug("Response body", zap.Any("http", httpResp.Body))
		return nil, utils.ErrorPerStatusCode(httpResp, err)
	}

	return &resp.Results, nil
}

// this function creates the necessary cache settings for next applications to work correctly on the edge
func (c *Client) CreateCacheSettingsNextApplication(ctx context.Context, req *CreateCacheSettingsRequest, applicationId int64) (CacheSettingsResponse, error) {
	logger.Debug("Create Cache Settings Next Application")
	req.SetBrowserCacheSettings("override")
	req.SetBrowserCacheSettingsMaximumTtl(31536000)
	req.SetCdnCacheSettings("override")
	req.SetCdnCacheSettingsMaximumTtl(31536000)

	request := c.apiClient.EdgeApplicationsCacheSettingsApi.EdgeApplicationsEdgeApplicationIdCacheSettingsPost(ctx, applicationId).ApplicationCacheCreateRequest(req.ApplicationCacheCreateRequest)

	resp, httpResp, err := request.Execute()
	if err != nil {
		logger.Debug("Error while creating a cache setting", zap.Error(err))
		logger.Debug("Status Code", zap.Any("http", httpResp.StatusCode))
		logger.Debug("Headers", zap.Any("http", httpResp.Header))
		logger.Debug("Response body", zap.Any("http", httpResp.Body))
		return nil, utils.ErrorPerStatusCode(httpResp, err)
	}

	return resp.Results, nil
}

func (c *Client) CreateRulesEngineNextApplication(ctx context.Context, applicationId int64, cacheId int64, typeLang string) error {
	logger.Debug("Create Rules Engine Next Application")
	req := CreateRulesEngineRequest{}
	req.SetName("cache policy")
	behavior := make([]sdk.RulesEngineBehavior, 1)
	behavior[0].SetName("set_cache_policy")
	behavior[0].SetTarget(cacheId)
	req.SetBehaviors(behavior)

	criteria := make([][]sdk.RulesEngineCriteria, 1)
	for i := 0; i < 1; i++ {
		criteria[i] = make([]sdk.RulesEngineCriteria, 1)
	}

	criteria[0][0].SetConditional("if")
	criteria[0][0].SetVariable("${uri}")
	criteria[0][0].SetOperator("starts_with")
	if typeLang == "nextjs" {
		criteria[0][0].SetInputValue("/_next/static")
	} else {
		criteria[0][0].SetInputValue("/")
	}

	req.SetCriteria(criteria)

	_, httpResp, err := c.apiClient.EdgeApplicationsRulesEngineApi.
		EdgeApplicationsEdgeApplicationIdRulesEnginePhaseRulesPost(ctx, applicationId, "request").
		CreateRulesEngineRequest(req.CreateRulesEngineRequest).Execute()
	if err != nil {
		logger.Debug("Error while creating a cache setting", zap.Error(err))
		logger.Debug("Status Code", zap.Any("http", httpResp.StatusCode))
		logger.Debug("Headers", zap.Any("http", httpResp.Header))
		logger.Debug("Response body", zap.Any("http", httpResp.Body))
		return utils.ErrorPerStatusCode(httpResp, err)
	}

	req.SetName("enable gzip")

	behavior[0].SetName("enable_gzip")
	behavior[0].SetTarget("")
	req.SetBehaviors(behavior)

	criteria[0][0].SetConditional("if")
	criteria[0][0].SetVariable("${request_uri}")
	criteria[0][0].SetOperator("exists")
	criteria[0][0].SetInputValue("")
	req.SetCriteria(criteria)

	_, httpResp, err = c.apiClient.EdgeApplicationsRulesEngineApi.
		EdgeApplicationsEdgeApplicationIdRulesEnginePhaseRulesPost(ctx, applicationId, "response").
		CreateRulesEngineRequest(req.CreateRulesEngineRequest).Execute()
	if err != nil {
		logger.Debug("Error while creating a cache setting", zap.Error(err))
		logger.Debug("Status Code", zap.Any("http", httpResp.StatusCode))
		logger.Debug("Headers", zap.Any("http", httpResp.Header))
		logger.Debug("Response body", zap.Any("http", httpResp.Body))
		return utils.ErrorPerStatusCode(httpResp, err)
	}

	return nil
}
