package edge_applications

import (
	"context"
	"fmt"
	"strconv"
	"strings"

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
	GetDeviceGroup() []int32
}

type GetCacheSettingsResponse interface {
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
	GetCookieNames() []*string
	GetEnableCachingForPost() bool
	GetL2CachingEnabled() bool
	GetAdaptiveDeliveryAction() string
	GetDeviceGroup() []int32
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
	GetDescription() string
	GetBehaviors() []sdk.RulesEngineBehaviorEntry
	GetCriteria() [][]sdk.RulesEngineCriteria
	GetIsActive() bool
	GetOrder() int64
	GetName() string
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

type UpdateRulesEngineRequest struct {
	sdk.PatchRulesEngineRequest
	IdApplication int64
	Phase         string
	Id            int64
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

func (c *Client) Update(ctx context.Context, req *UpdateRequest) (EdgeApplicationsResponse, error) {
	logger.Debug("Update Edge Application")
	str := strconv.FormatInt(req.Id, 10)
	request := c.apiClient.EdgeApplicationsMainSettingsAPI.EdgeApplicationsIdPatch(ctx, str).ApplicationUpdateRequest(req.ApplicationUpdateRequest)

	edgeApplicationsResponse, httpResp, err := request.Execute()
	if err != nil {
		if httpResp != nil {
			logger.Debug("Error while updating an Edge Application", zap.Error(err))
			err := utils.LogAndRewindBody(httpResp)
			if err != nil {
				return nil, err
			}
		}
		return nil, utils.ErrorPerStatusCode(httpResp, err)
	}

	return &edgeApplicationsResponse.Results, nil
}

func (c *Client) UpdateInstance(ctx context.Context, req *UpdateInstanceRequest, appID string, instanceID string) (FunctionsInstancesResponse, error) {
	logger.Debug("Update Instance")
	request := c.apiClient.EdgeApplicationsEdgeFunctionsInstancesAPI.EdgeApplicationsEdgeApplicationIdFunctionsInstancesFunctionsInstancesIdPatch(ctx, appID, instanceID).ApplicationUpdateInstanceRequest(req.ApplicationUpdateInstanceRequest)

	edgeApplicationsResponse, httpResp, err := request.Execute()
	if err != nil {
		logger.Debug("Error while updating an Edge Function instance", zap.Error(err))
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

	request := c.apiClient.EdgeApplicationsEdgeFunctionsInstancesAPI.EdgeApplicationsEdgeApplicationIdFunctionsInstancesPost(ctx, req.ApplicationId).ApplicationCreateInstanceRequest(req.ApplicationCreateInstanceRequest)

	edgeApplicationsResponse, httpResp, err := request.Execute()
	if err != nil {
		if httpResp != nil {
			logger.Debug("Error while creating an Edge Function instance", zap.Error(err))
			err := utils.LogAndRewindBody(httpResp)
			if err != nil {
				return nil, err
			}
		}
		return nil, utils.ErrorPerStatusCode(httpResp, err)
	}

	return edgeApplicationsResponse.Results, nil
}

func (c *Client) Delete(ctx context.Context, id int64) error {
	logger.Debug("Delete Edge Application")
	str := strconv.FormatInt(id, 10)
	req := c.apiClient.EdgeApplicationsMainSettingsAPI.EdgeApplicationsIdDelete(ctx, str)

	httpResp, err := req.Execute()
	if err != nil {
		if httpResp != nil {
			logger.Debug("Error while deleting an Edge Application", zap.Error(err))
			err := utils.LogAndRewindBody(httpResp)
			if err != nil {
				return err
			}
		}

		return utils.ErrorPerStatusCode(httpResp, err)
	}

	return nil
}

func (c *Client) ListRulesEngine(ctx context.Context, opts *contracts.ListOptions, edgeApplicationID int64, phase string) (*sdk.RulesEngineResponse, error) {
	logger.Debug("List Rules Engine")
	if opts.OrderBy == "" {
		opts.OrderBy = "id"
	}

	resp, httpResp, err := c.apiClient.EdgeApplicationsRulesEngineAPI.EdgeApplicationsEdgeApplicationIdRulesEnginePhaseRulesGet(ctx, edgeApplicationID, phase).
		OrderBy(opts.OrderBy).
		Page(opts.Page).
		PageSize(opts.PageSize).
		Sort(opts.Sort).Execute()

	if err != nil {
		if httpResp != nil {
			logger.Debug("Error while listing Rules Engine", zap.Error(err))
			err := utils.LogAndRewindBody(httpResp)
			if err != nil {
				return nil, err
			}
		}
		return nil, utils.ErrorPerStatusCode(httpResp, err)
	}

	return resp, nil
}

func (c *Client) GetRulesEngine(ctx context.Context, edgeApplicationID, rulesID int64, phase string) (RulesEngineResponse, error) {
	logger.Debug("Get Rules Engine")
	resp, httpResp, err := c.apiClient.EdgeApplicationsRulesEngineAPI.EdgeApplicationsEdgeApplicationIdRulesEnginePhaseRulesRuleIdGet(ctx, edgeApplicationID, phase, rulesID).Execute()
	if err != nil {
		if httpResp != nil {
			logger.Debug("Error while describing a Rules Engine", zap.Error(err))
			err := utils.LogAndRewindBody(httpResp)
			if err != nil {
				return nil, err
			}
		}
		return nil, utils.ErrorPerStatusCode(httpResp, err)
	}
	return &resp.Results, nil
}

func (c *Client) DeleteRulesEngine(ctx context.Context, edgeApplicationID int64, phase string, ruleID int64) error {
	logger.Debug("Delete Rules Engine")
	httpResp, err := c.apiClient.EdgeApplicationsRulesEngineAPI.EdgeApplicationsEdgeApplicationIdRulesEnginePhaseRulesRuleIdDelete(ctx, edgeApplicationID, phase, ruleID).Execute()
	if err != nil {
		if httpResp != nil {
			logger.Debug("Error while deleting a Rules Engine", zap.Error(err))
			err := utils.LogAndRewindBody(httpResp)
			if err != nil {
				return err
			}
		}
		return utils.ErrorPerStatusCode(httpResp, err)
	}
	return nil
}

func (c *Client) GetRulesDefault(ctx context.Context, applicationID int64, phase string) (int64, error) {
	logger.Debug("Get Rules Engine Default")
	request := c.apiClient.EdgeApplicationsRulesEngineAPI.EdgeApplicationsEdgeApplicationIdRulesEnginePhaseRulesGet(ctx, applicationID, "request")
	rules, httpResp, err := request.Execute()
	if err != nil {
		if httpResp != nil {
			logger.Debug("Error while deleting a Rules Engine", zap.Error(err))
			err := utils.LogAndRewindBody(httpResp)
			if err != nil {
				return 0, err
			}
		}
		return 0, utils.ErrorPerStatusCode(httpResp, err)
	}
	return rules.Results[0].Id, nil
}

func (c *Client) UpdateRulesEnginePublish(ctx context.Context, req *UpdateRulesEngineRequest, idFunc int64) (EdgeApplicationsResponse, error) {
	logger.Debug("Update Rules Engine Publish")
	request := c.apiClient.EdgeApplicationsRulesEngineAPI.EdgeApplicationsEdgeApplicationIdRulesEnginePhaseRulesGet(ctx, req.IdApplication, "request")

	edgeApplicationRules, httpResp, err := request.Execute()
	if err != nil {
		if httpResp != nil {
			logger.Debug("Error while updating a rules engine", zap.Error(err))
			err := utils.LogAndRewindBody(httpResp)
			if err != nil {
				return nil, err
			}
		}
		return nil, utils.ErrorPerStatusCode(httpResp, err)
	}

	idRule := edgeApplicationRules.Results[0].Id

	behaviors := make([]sdk.RulesEngineBehaviorEntry, 0)

	var behString sdk.RulesEngineBehaviorString
	behString.SetName("run_function")
	behString.SetTarget(fmt.Sprintf("%d", idFunc))

	behaviors = append(behaviors, sdk.RulesEngineBehaviorEntry{
		RulesEngineBehaviorString: &behString,
	})

	req.SetBehaviors(behaviors)

	requestUpdate := c.apiClient.EdgeApplicationsRulesEngineAPI.EdgeApplicationsEdgeApplicationIdRulesEnginePhaseRulesRuleIdPatch(ctx, req.IdApplication, "request", idRule).PatchRulesEngineRequest(req.PatchRulesEngineRequest)

	edgeApplicationsResponse, httpResp, err := requestUpdate.Execute()
	if err != nil {
		if httpResp != nil {
			logger.Debug("Error while updating a Rules Engine", zap.Error(err))
			err := utils.LogAndRewindBody(httpResp)
			if err != nil {
				return nil, err
			}
		}
		return nil, utils.ErrorPerStatusCode(httpResp, err)
	}

	return &edgeApplicationsResponse.Results, nil
}

func (c *Client) UpdateRulesEngine(ctx context.Context, req *UpdateRulesEngineRequest) (RulesEngineResponse, error) {
	logger.Debug("Update Rules Engine")
	requestUpdate := c.apiClient.EdgeApplicationsRulesEngineAPI.EdgeApplicationsEdgeApplicationIdRulesEnginePhaseRulesRuleIdPatch(ctx, req.IdApplication, req.Phase, req.Id).PatchRulesEngineRequest(req.PatchRulesEngineRequest)

	edgeApplicationsResponse, httpResp, err := requestUpdate.Execute()
	if err != nil {
		if httpResp != nil {
			logger.Debug("Error while updating a rules engine", zap.Error(err))
			err := utils.LogAndRewindBody(httpResp)
			if err != nil {
				return nil, err
			}
		}
		return nil, utils.ErrorPerStatusCode(httpResp, err)
	}

	return &edgeApplicationsResponse.Results, nil

}

func (c *Client) CreateRulesEngine(ctx context.Context, edgeApplicationID int64, phase string, req *CreateRulesEngineRequest) (RulesEngineResponse, error) {
	logger.Debug("Create Rules Engine")
	resp, httpResp, err := c.apiClient.EdgeApplicationsRulesEngineAPI.
		EdgeApplicationsEdgeApplicationIdRulesEnginePhaseRulesPost(ctx, edgeApplicationID, phase).
		CreateRulesEngineRequest(req.CreateRulesEngineRequest).Execute()
	if err != nil {
		if httpResp != nil {
			logger.Debug("Error while updating a Rules Engine", zap.Error(err))
			errLog := utils.LogAndRewindBody(httpResp)
			if errLog != nil {
				return &sdk.RulesEngineResultResponse{}, errLog
			}
			return nil, utils.ErrorPerStatusCode(httpResp, err)
		}
		return &sdk.RulesEngineResultResponse{}, err
	}
	return &resp.Results, nil
}

func (c *Client) EdgeFuncInstancesList(ctx context.Context, opts *contracts.ListOptions, edgeApplicationID int64) (*sdk.ApplicationInstancesGetResponse, error) {
	logger.Debug("List Edge Function Instances")
	if opts.OrderBy == "" {
		opts.OrderBy = "id"
	}

	resp, httpResp, err := c.apiClient.EdgeApplicationsEdgeFunctionsInstancesAPI.
		EdgeApplicationsEdgeApplicationIdFunctionsInstancesGet(ctx, edgeApplicationID).
		OrderBy(opts.OrderBy).
		Page(opts.Page).
		PageSize(opts.PageSize).
		Sort(opts.Sort).Execute()

	if err != nil {
		logger.Debug("Error while listing Edge Function instances", zap.Error(err))
		logger.Debug("Status Code", zap.Any("http", httpResp.StatusCode))
		logger.Debug("Headers", zap.Any("http", httpResp.Header))
		logger.Debug("Response body", zap.Any("http", httpResp.Body))
		return nil, utils.ErrorPerStatusCode(httpResp, err)
	}
	return resp, nil
}

func (c *Client) DeleteFunctionInstance(ctx context.Context, appID string, funcID string) error {
	logger.Debug("Delete Edge Function Instance")
	req := c.apiClient.EdgeApplicationsEdgeFunctionsInstancesAPI.EdgeApplicationsEdgeApplicationIdFunctionsInstancesFunctionsInstancesIdDelete(ctx, appID, funcID)

	httpResp, err := req.Execute()
	if err != nil {
		logger.Debug("Error while deleting an Edge Function instance", zap.Error(err))
		logger.Debug("Status Code", zap.Any("http", httpResp.StatusCode))
		logger.Debug("Headers", zap.Any("http", httpResp.Header))
		logger.Debug("Response body", zap.Any("http", httpResp.Body))
		return utils.ErrorPerStatusCode(httpResp, err)
	}

	return nil
}

func (c *Client) CreateFuncInstances(ctx context.Context, req *CreateFuncInstancesRequest, applicationID int64) (FunctionsInstancesResponse, error) {
	logger.Debug("Create Edge Function Instance")
	resp, httpResp, err := c.apiClient.EdgeApplicationsEdgeFunctionsInstancesAPI.EdgeApplicationsEdgeApplicationIdFunctionsInstancesPost(ctx, applicationID).
		ApplicationCreateInstanceRequest(req.ApplicationCreateInstanceRequest).Execute()
	if err != nil {
		logger.Debug("Error while creating an Edge Function instance", zap.Error(err))
		logger.Debug("Status Code", zap.Any("http", httpResp.StatusCode))
		logger.Debug("Headers", zap.Any("http", httpResp.Header))
		logger.Debug("Response body", zap.Any("http", httpResp.Body))
		return nil, utils.ErrorPerStatusCode(httpResp, err)
	}
	return resp.Results, nil
}

func (c *Client) GetFuncInstance(ctx context.Context, edgeApplicationID, instanceID int64) (FunctionsInstancesResponse, error) {
	logger.Debug("Get Edge Function Instance")
	resp, httpResp, err := c.apiClient.EdgeApplicationsEdgeFunctionsInstancesAPI.EdgeApplicationsEdgeApplicationIdFunctionsInstancesFunctionsInstancesIdGet(ctx, edgeApplicationID, instanceID).Execute()
	if err != nil {
		logger.Debug("Error while getting an Edge Function instance", zap.Error(err))
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
	resp, httpResp, err := c.apiClient.EdgeApplicationsDeviceGroupsAPI.
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
	req := c.apiClient.EdgeApplicationsDeviceGroupsAPI.EdgeApplicationsEdgeApplicationIdDeviceGroupsDeviceGroupIdDelete(ctx, appID, groupID)

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
	resp, httpResp, err := c.apiClient.EdgeApplicationsDeviceGroupsAPI.EdgeApplicationsEdgeApplicationIdDeviceGroupsDeviceGroupIdGet(ctx, edgeApplicationID, groupID).Execute()
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
	request := c.apiClient.EdgeApplicationsDeviceGroupsAPI.EdgeApplicationsEdgeApplicationIdDeviceGroupsDeviceGroupIdPatch(ctx, appID, groupID).PatchDeviceGroupsRequest(req)

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
	resp, httpResp, err := c.apiClient.EdgeApplicationsDeviceGroupsAPI.EdgeApplicationsEdgeApplicationIdDeviceGroupsPost(ctx, applicationID).
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

func (c *Client) CreateRulesEngineNextApplication(ctx context.Context, applicationId int64, cacheId int64, typeLang string, mode string) error {
	logger.Debug("Create Rules Engine Next Application")

	req := CreateRulesEngineRequest{}
	req.SetName("cache policy")

	behaviors := make([]sdk.RulesEngineBehaviorEntry, 0)

	var behStringCache sdk.RulesEngineBehaviorString
	behStringCache.SetName("set_cache_policy")
	behStringCache.SetTarget(fmt.Sprintf("%d", cacheId))

	behaviors = append(behaviors, sdk.RulesEngineBehaviorEntry{
		RulesEngineBehaviorString: &behStringCache,
	})

	req.SetBehaviors(behaviors)

	criteria := make([][]sdk.RulesEngineCriteria, 1)
	for i := 0; i < 1; i++ {
		criteria[i] = make([]sdk.RulesEngineCriteria, 1)
	}

	criteria[0][0].SetConditional("if")
	criteria[0][0].SetVariable("${uri}")
	criteria[0][0].SetOperator("starts_with")

	if typeLang == "Next" && strings.ToLower(mode) == "compute" {
		criteria[0][0].SetInputValue("/_next/static")
	} else {
		criteria[0][0].SetInputValue("/")
	}

	req.SetCriteria(criteria)

	_, httpResp, err := c.apiClient.EdgeApplicationsRulesEngineAPI.
		EdgeApplicationsEdgeApplicationIdRulesEnginePhaseRulesPost(ctx, applicationId, "request").
		CreateRulesEngineRequest(req.CreateRulesEngineRequest).Execute()
	if err != nil {
		if httpResp != nil {
			logger.Debug("Error while creating a Rules Engine", zap.Error(err))
			err := utils.LogAndRewindBody(httpResp)
			if err != nil {
				return err
			}
			return utils.ErrorPerStatusCode(httpResp, err)
		}
		logger.Debug("", zap.Any("Error", err.Error()))
		return utils.ErrorPerStatusCode(httpResp, err)
	}

	req.SetName("enable gzip")

	behaviorsGZIP := make([]sdk.RulesEngineBehaviorEntry, 0)

	var behString sdk.RulesEngineBehaviorString
	behString.SetName("enable_gzip")
	behString.SetTarget("")

	behaviorsGZIP = append(behaviorsGZIP, sdk.RulesEngineBehaviorEntry{
		RulesEngineBehaviorString: &behString,
	})

	req.SetBehaviors(behaviorsGZIP)

	criteria[0][0].SetConditional("if")
	criteria[0][0].SetVariable("${request_uri}")
	criteria[0][0].SetOperator("exists")
	criteria[0][0].SetInputValue("")
	req.SetCriteria(criteria)

	_, httpResp, err = c.apiClient.EdgeApplicationsRulesEngineAPI.
		EdgeApplicationsEdgeApplicationIdRulesEnginePhaseRulesPost(ctx, applicationId, "response").
		CreateRulesEngineRequest(req.CreateRulesEngineRequest).Execute()
	if err != nil {
		if httpResp != nil {
			logger.Debug("Error while creating a Rules Engine", zap.Error(err))
			err := utils.LogAndRewindBody(httpResp)
			if err != nil {
				return err
			}
			return utils.ErrorPerStatusCode(httpResp, err)
		}
		return utils.ErrorPerStatusCode(httpResp, err)
	}

	return nil
}
