package rules_engine

import (
	"context"
	"io"

	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/utils"
	"go.uber.org/zap"
)

func (c *Client) Delete(ctx context.Context, edgeApplicationID int64, phase string, ruleID int64) error {
	logger.Debug("Delete Rules Engine")
	httpResp, err := c.apiClient.EdgeApplicationsRulesEngineAPI.EdgeApplicationsEdgeApplicationIdRulesEnginePhaseRulesRuleIdDelete(ctx, edgeApplicationID, phase, ruleID).Execute()
	if err != nil {
		if httpResp != nil {
			logger.Debug("Error while deleting a rule engine", zap.Error(err))
			logger.Debug("", zap.Any("Status Code", httpResp.StatusCode))
			logger.Debug("", zap.Any("Headers", httpResp.Header))
			body, err := io.ReadAll(httpResp.Body)
			if err != nil {
				logger.Debug("Error while reading body of the http response", zap.Error(err))
				return utils.ErrorPerStatusCode(httpResp, err)
			}
			logger.Debug("", zap.Any("Body", string(body)))
		}
		return utils.ErrorPerStatusCode(httpResp, err)
	}
	return nil
}
