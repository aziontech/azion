package rules_engine

import (
	"context"

	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/utils"
	sdk "github.com/aziontech/azionapi-v4-go-sdk/edge"
	"go.uber.org/zap"
)

type UpdateRulesEngineRequest struct {
	sdk.PatchedEdgeApplicationRuleEngineRequest
	ApplicationID string
	RulesID       string
	Phase         string
}

type CreateRulesEngineRequest struct {
	sdk.EdgeApplicationRuleEngineRequest
}

type RulesEngineResponse interface {
	GetId() int64
	GetPhase() string
	GetDescription() string
	// GetBehaviors() []sdk.EdgeApplicationBehaviorFieldRequest
	// GetCriteria() [][]sdk.EdgeApplicationCriterionFieldRequest
	GetActive() bool
	GetOrder() int64
	GetName() string
}

// behaviors := make([]sdk.EdgeApplicationBehaviorFieldRequest, 0)

// var behString sdk.EdgeApplicationBehaviorFieldRequest
// var behSet sdk.EdgeApplicationBehaviorPolymorphicArgumentRequest
// funcId := fmt.Sprintf("%d", idFunc)
// behSet.String = &funcId
// behString.SetName("run_function")
// behString.SetArgument(behSet)

// req.SetBehaviors(behaviors)

func (c *Client) Delete(ctx context.Context, edgeApplicationID string, ruleID string) error {
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

func (c *Client) Update(ctx context.Context, req *UpdateRulesEngineRequest) (RulesEngineResponse, error) {
	logger.Debug("Update Rules Engine")
	requestUpdate := c.apiClient.EdgeApplicationsRulesAPI.PartialUpdateEdgeApplicationRule(ctx, req.ApplicationID, req.RulesID).PatchedEdgeApplicationRuleEngineRequest(req.PatchedEdgeApplicationRuleEngineRequest)

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

func (c *Client) Create(ctx context.Context, edgeApplicationID string, req sdk.EdgeApplicationRuleEngineRequest) (RulesEngineResponse, error) {
	logger.Debug("Create Rules Engine")
	resp, httpResp, err := c.apiClient.EdgeApplicationsRulesAPI.
		CreateEdgeApplicationRule(ctx, edgeApplicationID).
		EdgeApplicationRuleEngineRequest(req).Execute()

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
	return &resp.Data, nil
}
