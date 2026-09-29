package rules_engine

import (
	"context"

	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/utils"
	sdk "github.com/aziontech/azionapi-go-sdk/edgeapplications"
	"go.uber.org/zap"
)

type UpdateRulesEngineRequest struct {
	sdk.PatchRulesEngineRequest
	ApplicationID int64
	RulesID       int64
	Phase         string
}

type CreateRulesEngineRequest struct {
	sdk.CreateRulesEngineRequest
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

func (c *Client) Delete(ctx context.Context, edgeApplicationID int64, phase string, ruleID int64) error {
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

func (c *Client) Update(ctx context.Context, req *UpdateRulesEngineRequest) (RulesEngineResponse, error) {
	logger.Debug("Update Rules Engine")
	requestUpdate := c.apiClient.EdgeApplicationsRulesEngineAPI.EdgeApplicationsEdgeApplicationIdRulesEnginePhaseRulesRuleIdPatch(ctx, req.ApplicationID, req.Phase, req.RulesID).PatchRulesEngineRequest(req.PatchRulesEngineRequest)

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

func (c *Client) Create(ctx context.Context, edgeApplicationID int64, phase string, req sdk.CreateRulesEngineRequest) (RulesEngineResponse, error) {
	logger.Debug("Create Rules Engine")
	resp, httpResp, err := c.apiClient.EdgeApplicationsRulesEngineAPI.
		EdgeApplicationsEdgeApplicationIdRulesEnginePhaseRulesPost(ctx, edgeApplicationID, phase).
		CreateRulesEngineRequest(req).Execute()

	if err != nil {
		if httpResp != nil {
			logger.Debug("Error while creating a Rules Engine", zap.Error(err))
			err := utils.LogAndRewindBody(httpResp)
			if err != nil {
				return nil, err
			}
		}
		return nil, utils.ErrorPerStatusCode(httpResp, err)
	}
	return &resp.Results, nil
}
