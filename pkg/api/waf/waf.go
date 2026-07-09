package waf

import (
	"context"

	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/utils"
	sdk "github.com/aziontech/azionapi-v4-go-sdk-dev/azion-api"
	"go.uber.org/zap"
)

func (c *Client) List(ctx context.Context, opts *contracts.ListOptions) (*sdk.PaginatedWAFList, error) {
	logger.Debug("List wafs")
	if opts.OrderBy == "" {
		opts.OrderBy = "id"
	}
	resp, httpResp, err := c.apiClient.WAFsAPI.ListWafs(ctx).
		Ordering(opts.OrderBy).
		Page(opts.Page).
		PageSize(opts.PageSize).
		Search(opts.Sort).
		Execute()

	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while listing wafs", zap.Error(err))
			errBody, err = utils.LogAndRewindBodyV4(httpResp)
			if err != nil {
				return nil, err
			}
		}
		return nil, utils.ErrorPerStatusCodeV4(errBody, httpResp, err)
	}

	return resp, nil
}

func (c *Client) Get(ctx context.Context, id int64) (sdk.WAF, error) {
	logger.Debug("Get WAF")
	request := c.apiClient.WAFsAPI.RetrieveWaf(ctx, id)

	res, httpResp, err := request.Execute()
	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while getting a Firewall", zap.Error(err))
			errBody, err = utils.LogAndRewindBodyV4(httpResp)
			if err != nil {
				return sdk.WAF{}, err
			}
		}
		return sdk.WAF{}, utils.ErrorPerStatusCodeV4(errBody, httpResp, err)
	}

	return res.Data, nil
}

func (c *Client) Delete(ctx context.Context, id int64) error {
	logger.Debug("Delete WAF", zap.Any("ID", id))
	req := c.apiClient.WAFsAPI.DeleteWaf(ctx, id)

	_, httpResp, err := req.Execute()
	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while deleting a WAF", zap.Error(err), zap.Any("ID", id))
			errBody, err = utils.LogAndRewindBodyV4(httpResp)
			if err != nil {
				return err
			}
		}

		return utils.ErrorPerStatusCodeV4(errBody, httpResp, err)
	}

	return nil
}

func (c *Client) Create(ctx context.Context, req *CreateRequest) (sdk.WAF, error) {
	logger.Debug("Create WAF")

	request := c.apiClient.WAFsAPI.CreateWaf(ctx).WAFRequest(req.WAFRequest)

	wafResponse, httpResp, err := request.Execute()
	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while creating a Firewall", zap.Error(err), zap.Any("Name", req.Name))
			errBody, err = utils.LogAndRewindBodyV4(httpResp)
			if err != nil {
				return sdk.WAF{}, err
			}
		}
		return sdk.WAF{}, utils.ErrorPerStatusCodeV4(errBody, httpResp, err)
	}

	return wafResponse.Data, nil
}

func (c *Client) Update(ctx context.Context, req *UpdateRequest, id int64) (sdk.WAF, error) {
	logger.Debug("Update WAF", zap.Any("WAF ID", id), zap.Any("WAF name", req.Name))
	request := c.apiClient.WAFsAPI.PartialUpdateWaf(ctx, id).PatchedWAFRequest(req.PatchedWAFRequest)

	wafResponse, httpResp, err := request.Execute()
	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while updating a WAF", zap.Error(err), zap.Any("ID", id), zap.Any("Name", req.Name))
			errBody, err = utils.LogAndRewindBodyV4(httpResp)
			if err != nil {
				return sdk.WAF{}, err
			}
		}
		return sdk.WAF{}, utils.ErrorPerStatusCodeV4(errBody, httpResp, err)
	}

	return wafResponse.Data, nil
}
