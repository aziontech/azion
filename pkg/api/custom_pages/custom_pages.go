package custompages

import (
	"context"

	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/utils"
	sdk "github.com/aziontech/azionapi-v4-go-sdk-dev/azion-api"
	"go.uber.org/zap"
)

func (c *Client) Get(ctx context.Context, id int64) (sdk.CustomPage, error) {
	logger.Debug("Get Custom Page")
	request := c.apiClient.CustomPagesAPI.RetrieveCustomPage(ctx, id)

	res, httpResp, err := request.Execute()
	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while getting a Custom Page", zap.Error(err))
			errBody, err = utils.LogAndRewindBodyV4(httpResp)
			if err != nil {
				return sdk.CustomPage{}, err
			}
		}
		return sdk.CustomPage{}, utils.ErrorPerStatusCodeV4(errBody, httpResp, err)
	}

	return res.Data, nil
}

func (c *Client) Delete(ctx context.Context, id int64) error {
	logger.Debug("Delete Custom Page")
	request := c.apiClient.CustomPagesAPI.DeleteCustomPage(ctx, id)

	_, httpResp, err := request.Execute()

	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while deleting a Custom Page", zap.Error(err))
			errBody, err = utils.LogAndRewindBodyV4(httpResp)
			if err != nil {
				return err
			}
		}
		return utils.ErrorPerStatusCodeV4(errBody, httpResp, err)
	}

	return nil
}

func (c *Client) Create(ctx context.Context, req *CreateRequest) (sdk.CustomPage, error) {
	logger.Debug("Create Custom Page")

	request := c.apiClient.CustomPagesAPI.CreateCustomPage(ctx).CustomPageRequest(req.CustomPageRequest)

	response, httpResp, err := request.Execute()
	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while creating a Custom Page", zap.Error(err))
			errBody, err = utils.LogAndRewindBodyV4(httpResp)
			if err != nil {
				return sdk.CustomPage{}, err
			}
		}
		return sdk.CustomPage{}, utils.ErrorPerStatusCodeV4(errBody, httpResp, err)
	}

	return response.Data, nil
}

func (c *Client) Update(ctx context.Context, req *UpdateRequest, id int64) (sdk.CustomPage, error) {
	logger.Debug("Update Custom Page")
	request := c.apiClient.CustomPagesAPI.PartialUpdateCustomPage(ctx, id).PatchedCustomPageRequest(req.PatchedCustomPageRequest)

	response, httpResp, err := request.Execute()
	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while updating a Custom Page", zap.Error(err), zap.Any("ID", id))
			errBody, err = utils.LogAndRewindBodyV4(httpResp)
			if err != nil {
				return sdk.CustomPage{}, err
			}
		}
		return sdk.CustomPage{}, utils.ErrorPerStatusCodeV4(errBody, httpResp, err)
	}

	return response.Data, nil
}

func (c *Client) List(ctx context.Context, opts *contracts.ListOptions) (*sdk.PaginatedCustomPageList, error) {
	logger.Debug("List Custom Pages")
	if opts.OrderBy == "" {
		opts.OrderBy = "id"
	}
	resp, httpResp, err := c.apiClient.CustomPagesAPI.ListCustomPages(ctx).
		Ordering(opts.OrderBy).
		Page(opts.Page).
		PageSize(opts.PageSize).
		Search(opts.Sort).
		Execute()

	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while listing the Custom Pages", zap.Error(err))
			errBody, err = utils.LogAndRewindBodyV4(httpResp)
			if err != nil {
				return nil, err
			}
		}
		return nil, utils.ErrorPerStatusCodeV4(errBody, httpResp, err)
	}

	return resp, nil
}
