package rules_engine

import (
	"context"
	"io"

	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/utils"
	sdk "github.com/aziontech/azionapi-go-sdk/edgeapplications"
	"go.uber.org/zap"
)

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
			logger.Debug("Error while deleting a rule engine", zap.Error(err))
			err := utils.LogAndRewindBody(httpResp)
			if err != nil {
				return err
			}
		}
		return utils.ErrorPerStatusCode(httpResp, err)
	}
	return nil
}

func (c *Client) Create(ctx context.Context, edgeApplicationID int64, phase string, req sdk.CreateRulesEngineRequest) (RulesEngineResponse, error) {
	logger.Debug("Create Rules Engine")
	resp, httpResp, err := c.apiClient.EdgeApplicationsRulesEngineAPI.
		EdgeApplicationsEdgeApplicationIdRulesEnginePhaseRulesPost(ctx, edgeApplicationID, phase).
		CreateRulesEngineRequest(req).Execute()

	if err != nil {
		if httpResp != nil {
			logger.Debug("Error while updating a rules engine", zap.Error(err))
			logger.Debug("", zap.Any("Status Code", httpResp.StatusCode))
			logger.Debug("", zap.Any("Headers", httpResp.Header))
			body, errReadAll := io.ReadAll(httpResp.Body)
			if errReadAll != nil {
				logger.Debug("Error while reading body of the http response", zap.Error(errReadAll))
				return nil, utils.ErrorPerStatusCode(httpResp, errReadAll)
			}
			logger.Debug("", zap.Any("Body", string(body)))
		}
		return nil, utils.ErrorPerStatusCode(httpResp, err)
	}
	return &resp.Results, nil
}
