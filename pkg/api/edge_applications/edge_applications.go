package edge_applications

import (
	"context"

	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/utils"
	sdk "github.com/aziontech/azionapi-v4-go-sdk-dev/edge-api"
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
	GetModules() sdk.ApplicationModules
}

type RulesEngineResponse interface {
	GetId() int64
	GetDescription() string
	// GetBehaviors() []sdk.RulesEngineBehaviorEntry
	// GetCriteria() [][]sdk.RulesEngineCriteria
	GetActive() bool
	GetOrder() int64
	GetName() string
}

type UpdateRequest struct {
	sdk.PatchedApplicationRequest
	Id int64
}

type UpdateInstanceRequest struct {
	// sdk.ApplicationUpdateInstanceRequest
	sdk.PatchedApplicationFunctionInstanceRequest
}

type CreateInstanceRequest struct {
	sdk.ApplicationFunctionInstanceRequest
	ApplicationId int64
}

type UpdateRulesEngineRequest struct {
	sdk.PatchedApplicationRequestPhaseRuleEngineRequest
	IdApplication string
	Phase         string
	Id            string
}

type UpdateRulesEngineResponse struct {
	sdk.PatchedApplicationResponsePhaseRuleEngineRequest
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
	sdk.ApplicationRequestPhaseRuleEngineRequest
}

type CreateRulesEngineResponse struct {
	sdk.ApplicationResponsePhaseRuleEngineRequest
}

type FunctionsInstancesResponse interface {
	GetId() int64
	GetFunction() int64
	GetName() string
	GetArgs() interface{}
}

type DeviceGroupsResponse interface {
	GetId() int64
	GetName() string
	GetUserAgent() string
}

func (c *Client) Update(ctx context.Context, req *UpdateRequest) (EdgeApplicationsResponse, error) {
	logger.Debug("Update Edge Application")
	request := c.apiClient.ApplicationsAPI.PartialUpdateApplication(ctx, req.Id).PatchedApplicationRequest(req.PatchedApplicationRequest)

	edgeApplicationsResponse, httpResp, err := request.Execute()
	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while updating an Edge Application", zap.Error(err), zap.Any("ID", req.Id), zap.Any("Name", req.Name))
			errBody, err = utils.LogAndRewindBodyV4(httpResp)
			if err != nil {
				return nil, err
			}
		}
		return nil, utils.ErrorPerStatusCodeV4(errBody, httpResp, err)
	}

	return &edgeApplicationsResponse.Data, nil
}

func (c *Client) UpdateInstance(ctx context.Context, req *UpdateInstanceRequest, appID string, instanceID string) (sdk.ApplicationFunctionInstance, error) {
	logger.Debug("Update Instance")
	request := c.apiClient.ApplicationsFunctionAPI.PartialUpdateApplicationFunctionInstance(ctx, appID, instanceID).PatchedApplicationFunctionInstanceRequest(req.PatchedApplicationFunctionInstanceRequest)

	edgeApplicationsResponse, httpResp, err := request.Execute()
	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while updating an Edge Function instance", zap.Error(err), zap.Any("ID", instanceID), zap.Any("Name", req.Name))
			errBody, err = utils.LogAndRewindBodyV4(httpResp)
			if err != nil {
				return sdk.ApplicationFunctionInstance{}, err
			}
		}
		return sdk.ApplicationFunctionInstance{}, utils.ErrorPerStatusCodeV4(errBody, httpResp, err)
	}

	return edgeApplicationsResponse.Data, nil
}

func (c *Client) Delete(ctx context.Context, id int64) error {
	logger.Debug("Delete Edge Application")
	req := c.apiClient.ApplicationsAPI.DestroyApplication(ctx, id)

	_, httpResp, err := req.Execute()
	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while deleting an Edge Application", zap.Error(err), zap.Any("ID", id))
			errBody, err = utils.LogAndRewindBodyV4(httpResp)
			if err != nil {
				return err
			}
		}

		return utils.ErrorPerStatusCodeV4(errBody, httpResp, err)
	}

	return nil
}

func (c *Client) ListRulesEngineResponse(ctx context.Context, opts *contracts.ListOptions, edgeApplicationID string) (*sdk.PaginatedApplicationResponsePhaseRuleEngineList, error) {
	logger.Debug("List Rules Engine")
	if opts.OrderBy == "" {
		opts.OrderBy = "id"
	}

	resp, httpResp, err := c.apiClient.ApplicationsResponseRulesAPI.EdgeApplicationApiApplicationsResponseRulesList(ctx, edgeApplicationID).
		Ordering(opts.OrderBy).
		Page(opts.Page).
		PageSize(opts.PageSize).
		Search(opts.Sort).Execute()

	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while listing Rules Engine", zap.Error(err))
			errBody, err = utils.LogAndRewindBodyV4(httpResp)
			if err != nil {
				return nil, err
			}
		}
		return nil, utils.ErrorPerStatusCodeV4(errBody, httpResp, err)
	}

	return resp, nil
}

func (c *Client) ListRulesEngineRequest(ctx context.Context, opts *contracts.ListOptions, edgeApplicationID string) (*sdk.PaginatedApplicationRequestPhaseRuleEngineList, error) {
	logger.Debug("List Rules Engine")
	if opts.OrderBy == "" {
		opts.OrderBy = "id"
	}

	resp, httpResp, err := c.apiClient.ApplicationsRequestRulesAPI.EdgeApplicationApiApplicationsRequestRulesList(ctx, edgeApplicationID).
		Ordering(opts.OrderBy).
		Page(opts.Page).
		PageSize(opts.PageSize).
		Search(opts.Sort).Execute()

	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while listing Rules Engine", zap.Error(err))
			errBody, err = utils.LogAndRewindBodyV4(httpResp)
			if err != nil {
				return nil, err
			}
		}
		return nil, utils.ErrorPerStatusCodeV4(errBody, httpResp, err)
	}

	return resp, nil
}

func (c *Client) GetRulesEngineRequest(ctx context.Context, edgeApplicationID, rulesID string) (RulesEngineResponse, error) {
	logger.Debug("Get Rules Engine")
	resp, httpResp, err := c.apiClient.ApplicationsRequestRulesAPI.EdgeApplicationApiApplicationsRequestRulesRetrieve(ctx, edgeApplicationID, rulesID).Execute()
	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while describing a Rules Engine", zap.Error(err))
			errBody, err = utils.LogAndRewindBodyV4(httpResp)
			if err != nil {
				return nil, err
			}
		}
		return nil, utils.ErrorPerStatusCodeV4(errBody, httpResp, err)
	}
	return &resp.Data, nil
}

func (c *Client) GetRulesEngineResponse(ctx context.Context, edgeApplicationID, rulesID string) (RulesEngineResponse, error) {
	logger.Debug("Get Rules Engine")
	resp, httpResp, err := c.apiClient.ApplicationsResponseRulesAPI.EdgeApplicationApiApplicationsResponseRulesRetrieve(ctx, edgeApplicationID, rulesID).Execute()
	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while describing a Rules Engine", zap.Error(err))
			errBody, err = utils.LogAndRewindBodyV4(httpResp)
			if err != nil {
				return nil, err
			}
		}
		return nil, utils.ErrorPerStatusCodeV4(errBody, httpResp, err)
	}
	return &resp.Data, nil
}

func (c *Client) DeleteRulesEngineRequest(ctx context.Context, edgeApplicationID string, phase string, ruleID string) (int, error) {
	logger.Debug("Delete Rules Engine")
	_, httpResp, err := c.apiClient.ApplicationsRequestRulesAPI.EdgeApplicationApiApplicationsRequestRulesDestroy(ctx, edgeApplicationID, ruleID).Execute()
	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while deleting a Rules Engine", zap.Error(err), zap.Any("ID", ruleID))
			errBody, err = utils.LogAndRewindBodyV4(httpResp)
			if err != nil {
				return httpResp.StatusCode, err
			}
		}
		return httpResp.StatusCode, utils.ErrorPerStatusCodeV4(errBody, httpResp, err)
	}
	return 0, nil
}

func (c *Client) DeleteRulesEngineResponse(ctx context.Context, edgeApplicationID string, phase string, ruleID string) (int, error) {
	logger.Debug("Delete Rules Engine")
	_, httpResp, err := c.apiClient.ApplicationsResponseRulesAPI.EdgeApplicationApiApplicationsResponseRulesDestroy(ctx, edgeApplicationID, ruleID).Execute()
	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while deleting a Rules Engine", zap.Error(err), zap.Any("ID", ruleID))
			errBody, err = utils.LogAndRewindBodyV4(httpResp)
			if err != nil {
				return httpResp.StatusCode, err
			}
		}
		return httpResp.StatusCode, utils.ErrorPerStatusCodeV4(errBody, httpResp, err)
	}
	return 0, nil
}

func (c *Client) GetRulesDefault(ctx context.Context, applicationID string, phase string) (int64, error) {
	logger.Debug("Get Rules Engine Default")
	request := c.apiClient.ApplicationsRequestRulesAPI.EdgeApplicationApiApplicationsRequestRulesList(ctx, applicationID)
	rules, httpResp, err := request.Execute()
	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while retrieving a Rule Engine", zap.Error(err), zap.Any("Application ID", applicationID))
			errBody, err = utils.LogAndRewindBodyV4(httpResp)
			if err != nil {
				return 0, err
			}
		}
		return 0, utils.ErrorPerStatusCodeV4(errBody, httpResp, err)
	}
	return rules.Results[0].Id, nil
}

func (c *Client) UpdateRulesEngineRequest(ctx context.Context, req *UpdateRulesEngineRequest) (RulesEngineResponse, error) {
	logger.Debug("Update Rules Engine", zap.Any("ID", req.Id), zap.Any("Application ID", req.IdApplication), zap.Any("Name", req.Name))
	requestUpdate := c.apiClient.ApplicationsRequestRulesAPI.EdgeApplicationApiApplicationsRequestRulesPartialUpdate(ctx, req.IdApplication, req.Id).PatchedApplicationRequestPhaseRuleEngineRequest(req.PatchedApplicationRequestPhaseRuleEngineRequest)

	edgeApplicationsResponse, httpResp, err := requestUpdate.Execute()
	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while updating a rules engine", zap.Error(err), zap.Any("ID", req.Id), zap.Any("Application ID", req.IdApplication), zap.Any("Name", req.Name))
			errBody, err = utils.LogAndRewindBodyV4(httpResp)
			if err != nil {
				return nil, err
			}
		}
		return nil, utils.ErrorPerStatusCodeV4(errBody, httpResp, err)
	}

	return &edgeApplicationsResponse.Data, nil

}

func (c *Client) UpdateRulesEngineResponse(ctx context.Context, req *UpdateRulesEngineResponse) (RulesEngineResponse, error) {
	logger.Debug("Update Rules Engine")
	requestUpdate := c.apiClient.ApplicationsResponseRulesAPI.EdgeApplicationApiApplicationsResponseRulesPartialUpdate(ctx, req.IdApplication, req.Id).PatchedApplicationResponsePhaseRuleEngineRequest(req.PatchedApplicationResponsePhaseRuleEngineRequest)

	edgeApplicationsResponse, httpResp, err := requestUpdate.Execute()
	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while updating a rules engine", zap.Error(err), zap.Any("ID", req.Id), zap.Any("Name", req.Name))
			errBody, err = utils.LogAndRewindBodyV4(httpResp)
			if err != nil {
				return nil, err
			}
		}
		return nil, utils.ErrorPerStatusCodeV4(errBody, httpResp, err)
	}

	return &edgeApplicationsResponse.Data, nil

}

func (c *Client) Clone(ctx context.Context, name, id string) error {
	logger.Debug("Cloning Edge Application")
	req := sdk.CloneApplicationRequest{
		Name: name,
	}
	request := c.apiClient.ApplicationsAPI.CloneApplication(ctx, id).CloneApplicationRequest(req)
	_, httpResp, err := request.Execute()
	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while cloning an Edge Application", zap.Error(err), zap.Any("Name", req.Name))
			errBody, err = utils.LogAndRewindBodyV4(httpResp)
			if err != nil {
				return err
			}
		}
		return utils.ErrorPerStatusCodeV4(errBody, httpResp, err)
	}
	return nil
}

func (c *Client) CreateRulesEngineRequest(ctx context.Context, edgeApplicationID string, phase string, req *CreateRulesEngineRequest) (RulesEngineResponse, error) {
	logger.Debug("Create Rules Engine")
	resp, httpResp, err := c.apiClient.ApplicationsRequestRulesAPI.
		EdgeApplicationApiApplicationsRequestRulesCreate(ctx, edgeApplicationID).
		ApplicationRequestPhaseRuleEngineRequest(req.ApplicationRequestPhaseRuleEngineRequest).Execute()
	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while creating a Rules Engine", zap.Error(err), zap.Any("Name", req.Name))
			errBody, err = utils.LogAndRewindBodyV4(httpResp)
			if err != nil {
				return nil, err
			}
			return nil, utils.ErrorPerStatusCodeV4(errBody, httpResp, err)
		}
		return nil, err
	}
	return &resp.Data, nil
}

func (c *Client) CreateRulesEngineResponse(ctx context.Context, edgeApplicationID string, phase string, req *CreateRulesEngineResponse) (RulesEngineResponse, error) {
	logger.Debug("Create Rules Engine")
	resp, httpResp, err := c.apiClient.ApplicationsResponseRulesAPI.
		EdgeApplicationApiApplicationsResponseRulesCreate(ctx, edgeApplicationID).
		ApplicationResponsePhaseRuleEngineRequest(req.ApplicationResponsePhaseRuleEngineRequest).Execute()
	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while creating a Rules Engine", zap.Error(err), zap.Any("Name", req.Name))
			errBody, err = utils.LogAndRewindBodyV4(httpResp)
			if err != nil {
				return nil, err
			}
			return nil, utils.ErrorPerStatusCodeV4(errBody, httpResp, err)
		}
		return nil, err
	}
	return &resp.Data, nil
}

func (c *Client) EdgeFuncInstancesList(ctx context.Context, opts *contracts.ListOptions, edgeApplicationID string) (*sdk.PaginatedApplicationFunctionInstanceList, error) {
	logger.Debug("List Edge Function Instances")
	if opts.OrderBy == "" {
		opts.OrderBy = "id"
	}

	resp, httpResp, err := c.apiClient.ApplicationsFunctionAPI.
		ListApplicationFunctionInstances(ctx, edgeApplicationID).
		Ordering(opts.OrderBy).
		Page(opts.Page).
		PageSize(opts.PageSize).
		Search(opts.Sort).Execute()

	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while listing Edge Function instances", zap.Error(err))
			errBody, err = utils.LogAndRewindBodyV4(httpResp)
			if err != nil {
				return nil, err
			}
		}
		return nil, utils.ErrorPerStatusCodeV4(errBody, httpResp, err)
	}
	return resp, nil
}

func (c *Client) DeleteFunctionInstance(ctx context.Context, appID string, funcID string) error {
	logger.Debug("Delete Edge Function Instance")
	req := c.apiClient.ApplicationsFunctionAPI.DestroyApplicationFunctionInstance(ctx, appID, funcID)

	_, httpResp, err := req.Execute()
	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while deleting an Edge Function instance", zap.Error(err))
			errBody, err = utils.LogAndRewindBodyV4(httpResp)
			if err != nil {
				return err
			}
		}
		return utils.ErrorPerStatusCodeV4(errBody, httpResp, err)
	}

	return nil
}

func (c *Client) CreateFuncInstances(ctx context.Context, req *CreateInstanceRequest, applicationID string) (sdk.ApplicationFunctionInstance, error) {
	logger.Debug("Create Edge Function Instance")
	resp, httpResp, err := c.apiClient.ApplicationsFunctionAPI.CreateApplicationFunctionInstance(ctx, applicationID).
		ApplicationFunctionInstanceRequest(req.ApplicationFunctionInstanceRequest).Execute()
	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while creating an Edge Function instance", zap.Error(err), zap.Any("Name", req.Name))
			errBody, err = utils.LogAndRewindBodyV4(httpResp)
			if err != nil {
				return sdk.ApplicationFunctionInstance{}, err
			}
		}
		return sdk.ApplicationFunctionInstance{}, utils.ErrorPerStatusCodeV4(errBody, httpResp, err)
	}
	return resp.Data, nil
}

func (c *Client) GetFuncInstance(ctx context.Context, edgeApplicationID, instanceID string) (FunctionsInstancesResponse, error) {
	logger.Debug("Get Edge Function Instance")
	resp, httpResp, err := c.apiClient.ApplicationsFunctionAPI.RetrieveApplicationFunctionInstance(ctx, edgeApplicationID, instanceID).Execute()
	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while getting an Edge Function instance", zap.Error(err))
			errBody, err = utils.LogAndRewindBodyV4(httpResp)
			if err != nil {
				return nil, err
			}
		}
		return nil, utils.ErrorPerStatusCodeV4(errBody, httpResp, err)
	}

	return &resp.Data, nil
}

func (c *Client) CreateRulesEngineNextApplication(ctx context.Context, applicationId string, cacheId int64, typeLang string, authorize bool) error {
	logger.Debug("Create Rules Engine Next Application")

	req := CreateRulesEngineResponse{}
	criteria := make([][]sdk.EdgeApplicationCriterionFieldRequest, 1)
	for i := 0; i < 1; i++ {
		criteria[i] = make([]sdk.EdgeApplicationCriterionFieldRequest, 1)
	}

	req.SetName("enable gzip")

	behaviors := make([]sdk.ApplicationRuleEngineResponsePhaseBehaviorsRequest, 0)

	var behString sdk.ApplicationRuleEngineResponsePhaseBehaviorsRequest
	var behSet sdk.ApplicationResponsePhaseBehaviorWithoutArgsRequest
	behSet.SetType("enable_gzip")
	behString.ApplicationResponsePhaseBehaviorWithoutArgsRequest = &behSet

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

	_, httpResp, err := c.apiClient.ApplicationsResponseRulesAPI.
		EdgeApplicationApiApplicationsResponseRulesCreate(ctx, applicationId).
		ApplicationResponsePhaseRuleEngineRequest(req.ApplicationResponsePhaseRuleEngineRequest).Execute()
	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while creating a Rules Engine", zap.Error(err), zap.Any("Name", req.Name))
			errBody, err = utils.LogAndRewindBodyV4(httpResp)
			if err != nil {
				return err
			}
			return utils.ErrorPerStatusCode(httpResp, err)
		}
		return utils.ErrorPerStatusCodeV4(errBody, httpResp, err)
	}

	return nil
}
