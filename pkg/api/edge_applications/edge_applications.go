package edge_applications

import (
	"context"
	"fmt"
	"strconv"

	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/utils"
	sdk "github.com/aziontech/azionapi-v4-go-sdk/edge"
	"go.uber.org/zap"
)

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
	GetDebug() bool
	GetLastEditor() string
	GetModules() sdk.EdgeApplicationModules
}

type RulesEngineResponse interface {
	GetId() int64
	GetPhase() string
	GetDescription() string
	// GetBehaviors() []sdk.RulesEngineBehaviorEntry
	// GetCriteria() [][]sdk.RulesEngineCriteria
	GetActive() bool
	GetOrder() int64
	GetName() string
}

type UpdateRequest struct {
	sdk.PatchedEdgeApplicationRequest
	Id int64
}

type UpdateInstanceRequest struct {
	// sdk.ApplicationUpdateInstanceRequest
	sdk.PatchedEdgeApplicationFunctionInstanceRequest
}

type CreateInstanceRequest struct {
	sdk.EdgeApplicationFunctionInstanceRequest
	ApplicationId int64
}

type UpdateRulesEngineRequest struct {
	sdk.PatchedEdgeApplicationRuleEngineRequest
	IdApplication string
	Phase         string
	Id            string
}

type CreateCacheSettingsRequest struct {
	sdk.CacheSettingRequest
}

type UpdateCacheSettingsRequest struct {
	sdk.ApiPartialUpdateCacheSettingRequest
	Id int64
}

type CreateRulesEngineRequest struct {
	sdk.EdgeApplicationRuleEngineRequest
}

type FunctionsInstancesResponse interface {
	GetId() int64
	GetEdgeFunction() int64
	GetName() string
	GetJsonArgs() interface{}
}

type CreateDeviceGroupsRequest struct {
	sdk.EdgeApplicationDeviceGroupsRequest
}

type DeviceGroupsResponse interface {
	GetId() int64
	GetName() string
	GetUserAgent() string
}

func (c *Client) Update(ctx context.Context, req *UpdateRequest) (EdgeApplicationsResponse, error) {
	logger.Debug("Update Edge Application")
	str := strconv.FormatInt(req.Id, 10)
	request := c.apiClient.EdgeApplicationsAPI.PartialUpdateEdgeApplication(ctx, str).PatchedEdgeApplicationRequest(req.PatchedEdgeApplicationRequest)

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

	return &edgeApplicationsResponse.Data, nil
}

func (c *Client) UpdateInstance(ctx context.Context, req *UpdateInstanceRequest, appID string, instanceID string) (sdk.EdgeApplicationFunctionInstance, error) {
	logger.Debug("Update Instance")
	request := c.apiClient.EdgeApplicationsFunctionAPI.PartialUpdateEdgeApplicationFunctionInstance(ctx, appID, instanceID).PatchedEdgeApplicationFunctionInstanceRequest(req.PatchedEdgeApplicationFunctionInstanceRequest)

	edgeApplicationsResponse, httpResp, err := request.Execute()
	if err != nil {
		logger.Debug("Error while updating an Edge Function instance", zap.Error(err))
		logger.Debug("Status Code", zap.Any("http", httpResp.StatusCode))
		logger.Debug("Headers", zap.Any("http", httpResp.Header))
		logger.Debug("Response body", zap.Any("http", httpResp.Body))
		return sdk.EdgeApplicationFunctionInstance{}, utils.ErrorPerStatusCode(httpResp, err)
	}

	return edgeApplicationsResponse.Data, nil
}

func (c *Client) Delete(ctx context.Context, id int64) error {
	logger.Debug("Delete Edge Application")
	str := strconv.FormatInt(id, 10)
	req := c.apiClient.EdgeApplicationsAPI.DestroyEdgeApplication(ctx, str)

	_, httpResp, err := req.Execute()
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

func (c *Client) ListRulesEngine(ctx context.Context, opts *contracts.ListOptions, edgeApplicationID string) (*sdk.PaginatedResponseListEdgeApplicationRuleEngineList, error) {
	logger.Debug("List Rules Engine")
	if opts.OrderBy == "" {
		opts.OrderBy = "id"
	}

	resp, httpResp, err := c.apiClient.EdgeApplicationsRulesAPI.ListEdgeApplicationRule(ctx, edgeApplicationID).
		Ordering(opts.OrderBy).
		Page(opts.Page).
		PageSize(opts.PageSize).
		Search(opts.Sort).Execute()

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

func (c *Client) GetRulesEngine(ctx context.Context, edgeApplicationID, rulesID string) (RulesEngineResponse, error) {
	logger.Debug("Get Rules Engine")
	resp, httpResp, err := c.apiClient.EdgeApplicationsRulesAPI.RetrieveEdgeApplicationRule(ctx, edgeApplicationID, rulesID).Execute()
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
	return &resp.Data, nil
}

func (c *Client) DeleteRulesEngine(ctx context.Context, edgeApplicationID string, phase string, ruleID string) error {
	logger.Debug("Delete Rules Engine")
	_, httpResp, err := c.apiClient.EdgeApplicationsRulesAPI.DestroyEdgeApplicationRule(ctx, edgeApplicationID, ruleID).Execute()
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

func (c *Client) GetRulesDefault(ctx context.Context, applicationID string, phase string) (int64, error) {
	logger.Debug("Get Rules Engine Default")
	request := c.apiClient.EdgeApplicationsRulesAPI.ListEdgeApplicationRule(ctx, applicationID)
	rules, httpResp, err := request.Execute()
	if err != nil {
		if httpResp != nil {
			logger.Debug("Error while retrieving a Rule Engine", zap.Error(err))
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
	request := c.apiClient.EdgeApplicationsRulesAPI.ListEdgeApplicationRule(ctx, req.IdApplication)

	edgeApplicationRules, httpResp, err := request.Execute()
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

	idRule := edgeApplicationRules.Results[0].Id
	ruleId := fmt.Sprintf("%d", idRule)

	behaviors := make([]sdk.EdgeApplicationBehaviorFieldRequest, 0)

	var behString sdk.EdgeApplicationBehaviorFieldRequest
	var behSet sdk.EdgeApplicationBehaviorPolymorphicArgumentRequest
	funcId := fmt.Sprintf("%d", idFunc)
	behSet.String = &funcId
	behString.SetName("run_function")
	behString.SetArgument(behSet)

	req.SetBehaviors(behaviors)

	requestUpdate := c.apiClient.EdgeApplicationsRulesAPI.PartialUpdateEdgeApplicationRule(ctx, req.IdApplication, ruleId).PatchedEdgeApplicationRuleEngineRequest(req.PatchedEdgeApplicationRuleEngineRequest)

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

	return &edgeApplicationsResponse.Data, nil
}

func (c *Client) UpdateRulesEngine(ctx context.Context, req *UpdateRulesEngineRequest) (RulesEngineResponse, error) {
	logger.Debug("Update Rules Engine")
	requestUpdate := c.apiClient.EdgeApplicationsRulesAPI.PartialUpdateEdgeApplicationRule(ctx, req.IdApplication, req.Id).PatchedEdgeApplicationRuleEngineRequest(req.PatchedEdgeApplicationRuleEngineRequest)

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

	return &edgeApplicationsResponse.Data, nil

}

func (c *Client) CreateRulesEngine(ctx context.Context, edgeApplicationID string, phase string, req *CreateRulesEngineRequest) (RulesEngineResponse, error) {
	logger.Debug("Create Rules Engine")
	resp, httpResp, err := c.apiClient.EdgeApplicationsRulesAPI.
		CreateEdgeApplicationRule(ctx, edgeApplicationID).
		EdgeApplicationRuleEngineRequest(req.EdgeApplicationRuleEngineRequest).Execute()
	if err != nil {
		if httpResp != nil {
			logger.Debug("Error while creating a Rules Engine", zap.Error(err))
			errLog := utils.LogAndRewindBody(httpResp)
			if errLog != nil {
				return nil, errLog
			}
			return nil, utils.ErrorPerStatusCode(httpResp, err)
		}
		return nil, err
	}
	return &resp.Data, nil
}

func (c *Client) EdgeFuncInstancesList(ctx context.Context, opts *contracts.ListOptions, edgeApplicationID string) (*sdk.PaginatedResponseListEdgeApplicationFunctionInstanceList, error) {
	logger.Debug("List Edge Function Instances")
	if opts.OrderBy == "" {
		opts.OrderBy = "id"
	}

	resp, httpResp, err := c.apiClient.EdgeApplicationsFunctionAPI.
		ListEdgeApplicationFunctionInstances(ctx, edgeApplicationID).
		Ordering(opts.OrderBy).
		Page(opts.Page).
		PageSize(opts.PageSize).
		Search(opts.Sort).Execute()

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
	req := c.apiClient.EdgeApplicationsFunctionAPI.DestroyEdgeApplicationFunctionInstance(ctx, appID, funcID)

	_, httpResp, err := req.Execute()
	if err != nil {
		logger.Debug("Error while deleting an Edge Function instance", zap.Error(err))
		logger.Debug("Status Code", zap.Any("http", httpResp.StatusCode))
		logger.Debug("Headers", zap.Any("http", httpResp.Header))
		logger.Debug("Response body", zap.Any("http", httpResp.Body))
		return utils.ErrorPerStatusCode(httpResp, err)
	}

	return nil
}

func (c *Client) CreateFuncInstances(ctx context.Context, req *CreateInstanceRequest, applicationID string) (sdk.EdgeApplicationFunctionInstance, error) {
	logger.Debug("Create Edge Function Instance")
	resp, httpResp, err := c.apiClient.EdgeApplicationsFunctionAPI.CreateEdgeFirewallFunctionInstance(ctx, applicationID).
		EdgeApplicationFunctionInstanceRequest(req.EdgeApplicationFunctionInstanceRequest).Execute()
	if err != nil {
		logger.Debug("Error while creating an Edge Function instance", zap.Error(err))
		logger.Debug("Status Code", zap.Any("http", httpResp.StatusCode))
		logger.Debug("Headers", zap.Any("http", httpResp.Header))
		logger.Debug("Response body", zap.Any("http", httpResp.Body))
		return sdk.EdgeApplicationFunctionInstance{}, utils.ErrorPerStatusCode(httpResp, err)
	}
	return resp.Data, nil
}

func (c *Client) GetFuncInstance(ctx context.Context, edgeApplicationID, instanceID string) (FunctionsInstancesResponse, error) {
	logger.Debug("Get Edge Function Instance")
	resp, httpResp, err := c.apiClient.EdgeApplicationsFunctionAPI.RetrieveFunctionInstance(ctx, edgeApplicationID, instanceID).Execute()
	if err != nil {
		logger.Debug("Error while getting an Edge Function instance", zap.Error(err))
		logger.Debug("Status Code", zap.Any("http", httpResp.StatusCode))
		logger.Debug("Headers", zap.Any("http", httpResp.Header))
		logger.Debug("Response body", zap.Any("http", httpResp.Body))
		return nil, utils.ErrorPerStatusCode(httpResp, err)
	}
	return &resp.Data, nil
}

func (c *Client) DeviceGroupsList(ctx context.Context, opts *contracts.ListOptions, edgeApplicationID string) (*sdk.PaginatedResponseListEdgeApplicationDeviceGroupsList, error) {
	logger.Debug("List Device Groups")
	if opts.OrderBy == "" {
		opts.OrderBy = "id"
	}
	resp, httpResp, err := c.apiClient.EdgeApplicationsDeviceGroupsAPI.
		ListDeviceGroups(ctx, edgeApplicationID).
		Ordering(opts.OrderBy).
		Page(opts.Page).
		PageSize(opts.PageSize).
		Search(opts.Sort).Execute()
	if err != nil {
		logger.Debug("Error while listing device groups", zap.Error(err))
		logger.Debug("Status Code", zap.Any("http", httpResp.StatusCode))
		logger.Debug("Headers", zap.Any("http", httpResp.Header))
		logger.Debug("Response body", zap.Any("http", httpResp.Body))
		return nil, utils.ErrorPerStatusCode(httpResp, err)
	}
	return resp, nil
}

func (c *Client) DeleteDeviceGroup(ctx context.Context, appID string, groupID string) error {
	logger.Debug("Delete Device Group")
	req := c.apiClient.EdgeApplicationsDeviceGroupsAPI.DestroyDeviceGroup(ctx, appID, groupID)

	_, httpResp, err := req.Execute()
	if err != nil {
		logger.Debug("Error while deleting a device group", zap.Error(err))
		logger.Debug("Status Code", zap.Any("http", httpResp.StatusCode))
		logger.Debug("Headers", zap.Any("http", httpResp.Header))
		logger.Debug("Response body", zap.Any("http", httpResp.Body))
		return utils.ErrorPerStatusCode(httpResp, err)
	}

	return nil
}

func (c *Client) GetDeviceGroups(ctx context.Context, edgeApplicationID, groupID string) (DeviceGroupsResponse, error) {
	logger.Debug("Get Device Groups")
	resp, httpResp, err := c.apiClient.EdgeApplicationsDeviceGroupsAPI.RetrieveDeviceGroup(ctx, edgeApplicationID, groupID).Execute()
	if err != nil {
		logger.Debug("Error while getting a device group", zap.Error(err))
		logger.Debug("Status Code", zap.Any("http", httpResp.StatusCode))
		logger.Debug("Headers", zap.Any("http", httpResp.Header))
		logger.Debug("Response body", zap.Any("http", httpResp.Body))
		return nil, utils.ErrorPerStatusCode(httpResp, err)
	}
	return &resp.Data, nil
}

func (c *Client) UpdateDeviceGroup(ctx context.Context, req sdk.PatchedEdgeApplicationDeviceGroupsRequest, appID string, groupID string) (DeviceGroupsResponse, error) {
	logger.Debug("Update Device Group")
	request := c.apiClient.EdgeApplicationsDeviceGroupsAPI.PartialUpdateDeviceGroup(ctx, appID, groupID).PatchedEdgeApplicationDeviceGroupsRequest(req)

	deviceGroup, httpResp, err := request.Execute()
	if err != nil {
		logger.Debug("Error while updating a device group", zap.Error(err))
		logger.Debug("Status Code", zap.Any("http", httpResp.StatusCode))
		logger.Debug("Headers", zap.Any("http", httpResp.Header))
		logger.Debug("Response body", zap.Any("http", httpResp.Body))
		return nil, utils.ErrorPerStatusCode(httpResp, err)
	}

	return &deviceGroup.Data, nil
}

func (c *Client) CreateDeviceGroups(ctx context.Context, req *CreateDeviceGroupsRequest, applicationID string) (DeviceGroupsResponse, error) {
	logger.Debug("Create Device Groups")
	resp, httpResp, err := c.apiClient.EdgeApplicationsDeviceGroupsAPI.CreateDeviceGroup(ctx, applicationID).
		EdgeApplicationDeviceGroupsRequest(req.EdgeApplicationDeviceGroupsRequest).Execute()
	if err != nil {
		logger.Debug("Error while creating a device group", zap.Error(err))
		logger.Debug("Status Code", zap.Any("http", httpResp.StatusCode))
		logger.Debug("Headers", zap.Any("http", httpResp.Header))
		logger.Debug("Response body", zap.Any("http", httpResp.Body))
		return nil, utils.ErrorPerStatusCode(httpResp, err)
	}

	return &resp.Data, nil
}

func (c *Client) CreateRulesEngineNextApplication(ctx context.Context, applicationId string, cacheId int64, typeLang string, authorize bool) error {
	logger.Debug("Create Rules Engine Next Application")

	req := CreateRulesEngineRequest{}
	criteria := make([][]sdk.EdgeApplicationCriterionFieldRequest, 1)
	for i := 0; i < 1; i++ {
		criteria[i] = make([]sdk.EdgeApplicationCriterionFieldRequest, 1)
	}

	req.SetName("enable gzip")

	behaviors := make([]sdk.EdgeApplicationBehaviorFieldRequest, 0)

	var behString sdk.EdgeApplicationBehaviorFieldRequest
	behString.SetName("enable_gzip")

	// var behString sdk.EdgeApplicationBehaviorFieldRequest
	var behSet sdk.EdgeApplicationBehaviorPolymorphicArgumentRequest
	behString.SetArgument(behSet)
	behaviors = append(behaviors, behString)

	req.SetBehaviors(behaviors)

	emptyString := ""
	arg := sdk.EdgeApplicationCriterionPolymorphicArgumentRequest{
		String: &emptyString,
	}

	criteria[0][0].SetConditional("if")
	criteria[0][0].SetVariable("${request_uri}")
	criteria[0][0].SetOperator("exists")
	criteria[0][0].SetArgument(arg)
	req.SetCriteria(criteria)

	_, httpResp, err := c.apiClient.EdgeApplicationsRulesAPI.
		CreateEdgeApplicationRule(ctx, applicationId).
		EdgeApplicationRuleEngineRequest(req.EdgeApplicationRuleEngineRequest).Execute()
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
