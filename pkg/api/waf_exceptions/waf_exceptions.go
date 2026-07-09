package wafexceptions

import (
	"context"

	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/utils"
	sdk "github.com/aziontech/azionapi-v4-go-sdk-dev/azion-api"
	"go.uber.org/zap"
)

func (c *Client) List(ctx context.Context, opts *contracts.ListOptions, wafID int64) (*sdk.PaginatedWAFRuleList, error) {

	logger.Debug("List WAF Exceptions")

	req := c.apiClient.WAFsExceptionsAPI.
		ListWafExceptions(ctx, wafID).
		Page(opts.Page).
		PageSize(opts.PageSize)
	resp, httpResp, err := req.Execute()

	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while listing WAF Exceptions", zap.Error(err))
			errBody, err = utils.LogAndRewindBodyV4(httpResp)
			if err != nil {
				return nil, err
			}
		}
		return nil, utils.ErrorPerStatusCodeV4(errBody, httpResp, err)
	}

	return resp, nil
}

func (c *Client) Get(ctx context.Context, wafID, exceptionID int64) (sdk.WAFRule, error) {

	logger.Debug("Get WAF Exception")
	req := c.apiClient.WAFsExceptionsAPI.RetrieveWafException(ctx, exceptionID, wafID)
	resp, httpResp, err := req.Execute()

	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while getting WAF Exception", zap.Error(err))
			errBody, err = utils.LogAndRewindBodyV4(httpResp)
			if err != nil {
				return sdk.WAFRule{}, err
			}
		}
		return sdk.WAFRule{}, utils.ErrorPerStatusCodeV4(errBody, httpResp, err)
	}

	return resp.Data, nil
}

func (c *Client) Delete(ctx context.Context, wafID, exceptionID int64) error {
	logger.Debug("Delete WAF Exception")
	req := c.apiClient.WAFsExceptionsAPI.DeleteWafException(ctx, wafID, exceptionID)
	_, httpResp, err := req.Execute()

	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while deleting WAF Exception", zap.Error(err))
			errBody, err = utils.LogAndRewindBodyV4(httpResp)
			if err != nil {
				return err
			}
		}
		return utils.ErrorPerStatusCodeV4(errBody, httpResp, err)
	}

	return nil
}

func (c *Client) Create(ctx context.Context, wafID int64, req sdk.WAFRuleRequest) (sdk.WAFRule, error) {
	logger.Debug("Create WAF Exception")
	request := c.apiClient.WAFsExceptionsAPI.CreateWafException(ctx, wafID).WAFRuleRequest(req)
	resp, httpResp, err := request.Execute()

	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while creating WAF Exception", zap.Error(err))
			errBody, err = utils.LogAndRewindBodyV4(httpResp)
			if err != nil {
				return sdk.WAFRule{}, err
			}
		}
		return sdk.WAFRule{}, utils.ErrorPerStatusCodeV4(errBody, httpResp, err)
	}

	return resp.Data, nil
}

func (c *Client) Update(ctx context.Context, wafID, exceptionID int64, req sdk.PatchedWAFRuleRequest) (sdk.WAFRule, error) {
	logger.Debug("Update WAF Exception")
	request := c.apiClient.WAFsExceptionsAPI.PartialUpdateWafException(ctx, wafID, exceptionID).PatchedWAFRuleRequest(req)
	resp, httpResp, err := request.Execute()

	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while updating WAF Exception", zap.Error(err))
			errBody, err = utils.LogAndRewindBodyV4(httpResp)
			if err != nil {
				return sdk.WAFRule{}, err
			}
		}
		return sdk.WAFRule{}, utils.ErrorPerStatusCodeV4(errBody, httpResp, err)
	}

	return resp.Data, nil
}
