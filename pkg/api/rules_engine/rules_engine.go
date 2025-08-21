package rules_engine

import (
	"context"

	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/utils"
	sdkv3 "github.com/aziontech/azionapi-go-sdk/edgeapplications"
	sdk "github.com/aziontech/azionapi-v4-go-sdk-dev/edge-api"
	"go.uber.org/zap"
)

type UpdateRulesEngineRequest struct {
	sdk.PatchedEdgeApplicationRequestPhaseRuleEngineRequest
	ApplicationID string
	RulesID       string
	Phase         string
}

type UpdateRulesEngineResponse struct {
	sdk.PatchedEdgeApplicationResponsePhaseRuleEngineRequest
	ApplicationID string
	RulesID       string
	Phase         string
}

type CreateRulesEngineRequest struct {
	sdk.EdgeApplicationRequestPhaseRuleEngineRequest
}

type CreateRulesEngineResponse struct {
	sdk.EdgeApplicationResponsePhaseRuleEngineRequest
}

type CreateRulesEngineRequestV3 struct {
	sdkv3.CreateRulesEngineRequest
}

type RulesEngineResponse interface {
	GetId() int64
	GetDescription() string
	// GetBehaviors() []sdk.EdgeApplicationBehaviorFieldRequest
	// GetCriteria() [][]sdk.EdgeApplicationCriterionFieldRequest
	GetActive() bool
	GetOrder() int64
	GetName() string
}

func (c *Client) DeleteRequest(ctx context.Context, edgeApplicationID string, ruleID string) error {
	logger.Debug("Delete Rules Engine")
	_, httpResp, err := c.apiClient.EdgeApplicationsRequestRulesAPI.EdgeApplicationApiApplicationsRequestRulesDestroy(ctx, edgeApplicationID, ruleID).Execute()
	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while deleting a Rules Engine", zap.Error(err))
			errBody, err = utils.LogAndRewindBodyV4(httpResp)
			if err != nil {
				return err
			}
		}
		return utils.ErrorPerStatusCodeV4(errBody, httpResp, err)
	}
	return nil
}

func (c *Client) DeleteResponse(ctx context.Context, edgeApplicationID string, ruleID string) error {
	logger.Debug("Delete Rules Engine")
	_, httpResp, err := c.apiClient.EdgeApplicationsResponseRulesAPI.EdgeApplicationApiApplicationsResponseRulesDestroy(ctx, edgeApplicationID, ruleID).Execute()
	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while deleting a Rules Engine", zap.Error(err))
			errBody, err = utils.LogAndRewindBodyV4(httpResp)
			if err != nil {
				return err
			}
		}
		return utils.ErrorPerStatusCodeV4(errBody, httpResp, err)
	}
	return nil
}

func (c *Client) UpdateRequest(ctx context.Context, req *UpdateRulesEngineRequest) (RulesEngineResponse, error) {
	logger.Debug("Update Rules Engine")
	requestUpdate := c.apiClient.EdgeApplicationsRequestRulesAPI.EdgeApplicationApiApplicationsRequestRulesPartialUpdate(ctx, req.ApplicationID, req.RulesID).PatchedEdgeApplicationRequestPhaseRuleEngineRequest(req.PatchedEdgeApplicationRequestPhaseRuleEngineRequest)

	edgeApplicationsResponse, httpResp, err := requestUpdate.Execute()
	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while updating a Rules Engine", zap.Error(err))
			errBody, err = utils.LogAndRewindBodyV4(httpResp)
			if err != nil {
				return nil, err
			}
		}
		return nil, utils.ErrorPerStatusCodeV4(errBody, httpResp, err)
	}

	return &edgeApplicationsResponse.Data, nil
}

func (c *Client) UpdateResponse(ctx context.Context, req *UpdateRulesEngineResponse) (RulesEngineResponse, error) {
	logger.Debug("Update Rules Engine")
	requestUpdate := c.apiClient.EdgeApplicationsResponseRulesAPI.EdgeApplicationApiApplicationsResponseRulesPartialUpdate(ctx, req.ApplicationID, req.RulesID).PatchedEdgeApplicationResponsePhaseRuleEngineRequest(req.PatchedEdgeApplicationResponsePhaseRuleEngineRequest)

	edgeApplicationsResponse, httpResp, err := requestUpdate.Execute()
	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while updating a Rules Engine", zap.Error(err))
			errBody, err = utils.LogAndRewindBodyV4(httpResp)
			if err != nil {
				return nil, err
			}
		}
		return nil, utils.ErrorPerStatusCodeV4(errBody, httpResp, err)
	}

	return &edgeApplicationsResponse.Data, nil
}

func (c *Client) CreateRequest(ctx context.Context, edgeApplicationID string, req sdk.EdgeApplicationRequestPhaseRuleEngineRequest) (RulesEngineResponse, error) {
	logger.Debug("Create Rules Engine")
	resp, httpResp, err := c.apiClient.EdgeApplicationsRequestRulesAPI.
		EdgeApplicationApiApplicationsRequestRulesCreate(ctx, edgeApplicationID).
		EdgeApplicationRequestPhaseRuleEngineRequest(req).Execute()

	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while creating a Rules Engine", zap.Error(err))
			errBody, err = utils.LogAndRewindBodyV4(httpResp)
			if err != nil {
				return nil, err
			}
		}
		return nil, utils.ErrorPerStatusCodeV4(errBody, httpResp, err)
	}
	return &resp.Data, nil
}

func (c *Client) CreateResponse(ctx context.Context, edgeApplicationID string, req sdk.EdgeApplicationResponsePhaseRuleEngineRequest) (RulesEngineResponse, error) {
	logger.Debug("Create Rules Engine")
	resp, httpResp, err := c.apiClient.EdgeApplicationsResponseRulesAPI.
		EdgeApplicationApiApplicationsResponseRulesCreate(ctx, edgeApplicationID).
		EdgeApplicationResponsePhaseRuleEngineRequest(req).Execute()

	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while creating a Rules Engine", zap.Error(err))
			errBody, err = utils.LogAndRewindBodyV4(httpResp)
			if err != nil {
				return nil, err
			}
		}
		return nil, utils.ErrorPerStatusCodeV4(errBody, httpResp, err)
	}
	return &resp.Data, nil
}
