package rules_engine

import (
	"context"

	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/utils"
	sdkv3 "github.com/aziontech/azionapi-go-sdk/edgeapplications"
	sdk "github.com/aziontech/azionapi-v4-go-sdk-dev/azion-api"
	"go.uber.org/zap"
)

type UpdateRulesEngineRequest struct {
	sdk.PatchedRequestPhaseRule
	ApplicationID int64
	RulesID       int64
	Phase         string
}

type UpdateRulesEngineResponse struct {
	sdk.PatchedResponsePhaseRuleRequest
	ApplicationID int64
	RulesID       int64
	Phase         string
}

type CreateRulesEngineRequest struct {
	sdk.RequestPhaseRuleRequest
}

type CreateRulesEngineResponse struct {
	sdk.ResponsePhaseRuleRequest
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

func (c *Client) DeleteRequest(ctx context.Context, applicationID int64, ruleID int64) error {
	logger.Debug("Delete Rules Engine")
	_, httpResp, err := c.apiClient.ApplicationsRequestRulesAPI.DeleteApplicationRequestRule(ctx, applicationID, ruleID).Execute()
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

func (c *Client) DeleteResponse(ctx context.Context, edgeApplicationID int64, ruleID int64) error {
	logger.Debug("Delete Rules Engine")
	_, httpResp, err := c.apiClient.ApplicationsResponseRulesAPI.DeleteApplicationResponseRule(ctx, edgeApplicationID, ruleID).Execute()
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
	requestUpdate := c.apiClient.ApplicationsRequestRulesAPI.PartialUpdateApplicationRequestRule(ctx, req.ApplicationID, req.RulesID).PatchedRequestPhaseRule(req.PatchedRequestPhaseRule)

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
	requestUpdate := c.apiClient.ApplicationsResponseRulesAPI.PartialUpdateApplicationResponseRule(ctx, req.ApplicationID, req.RulesID).PatchedResponsePhaseRuleRequest(req.PatchedResponsePhaseRuleRequest)

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

func (c *Client) CreateRequest(ctx context.Context, edgeApplicationID int64, req sdk.RequestPhaseRuleRequest) (RulesEngineResponse, error) {
	logger.Debug("Create Rules Engine")
	resp, httpResp, err := c.apiClient.ApplicationsRequestRulesAPI.
		CreateApplicationRequestRule(ctx, edgeApplicationID).
		RequestPhaseRuleRequest(req).Execute()

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

func (c *Client) CreateResponse(ctx context.Context, edgeApplicationID int64, req sdk.ResponsePhaseRuleRequest) (RulesEngineResponse, error) {
	logger.Debug("Create Rules Engine")
	resp, httpResp, err := c.apiClient.ApplicationsResponseRulesAPI.
		CreateApplicationResponseRule(ctx, edgeApplicationID).
		ResponsePhaseRuleRequest(req).Execute()

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
