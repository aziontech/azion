package datastream

import (
	"context"

	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/utils"
	sdk "github.com/aziontech/azionapi-v4-go-sdk-dev/data-stream-api"
	"go.uber.org/zap"
)

func (c *Client) GetTemplate(ctx context.Context, id int64) (sdk.Template, error) {
	logger.Debug("Get Template")
	request := c.apiClient.DataStreamTemplatesAPI.RetrieveTemplate(ctx, id)

	res, httpResp, err := request.Execute()
	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while getting a template", zap.Error(err))
			errBody, err = utils.LogAndRewindBodyV4(httpResp)
			if err != nil {
				return sdk.Template{}, err
			}
		}
		return sdk.Template{}, utils.ErrorPerStatusCodeV4(errBody, httpResp, err)
	}

	return res.Data, nil
}

func (c *Client) DeleteTemplate(ctx context.Context, id int64) error {
	logger.Debug("Delete Template")
	request := c.apiClient.DataStreamTemplatesAPI.DeleteTemplate(ctx, id)

	_, httpResp, err := request.Execute()

	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while deleting a template", zap.Error(err))
			errBody, err = utils.LogAndRewindBodyV4(httpResp)
			if err != nil {
				return err
			}
		}
		return utils.ErrorPerStatusCodeV4(errBody, httpResp, err)
	}

	return nil
}

func (c *Client) CreateTemplate(ctx context.Context, req *CreateTemplateRequest) (sdk.Template, error) {
	logger.Debug("Create Template")

	request := c.apiClient.DataStreamTemplatesAPI.CreateTemplate(ctx).TemplateRequest(req.TemplateRequest)

	response, httpResp, err := request.Execute()
	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while creating a template", zap.Error(err))
			errBody, err = utils.LogAndRewindBodyV4(httpResp)
			if err != nil {
				return sdk.Template{}, err
			}
		}
		return sdk.Template{}, utils.ErrorPerStatusCodeV4(errBody, httpResp, err)
	}

	return response.Data, nil
}

func (c *Client) UpdateTemplate(ctx context.Context, req *UpdateTemplateRequest, id int64) (sdk.Template, error) {
	logger.Debug("Update Template")
	request := c.apiClient.DataStreamTemplatesAPI.PartialUpdateTemplate(ctx, id).PatchedTemplateRequest(req.PatchedTemplateRequest)

	response, httpResp, err := request.Execute()
	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while updating a template", zap.Error(err), zap.Any("ID", id))
			errBody, err = utils.LogAndRewindBodyV4(httpResp)
			if err != nil {
				return sdk.Template{}, err
			}
		}
		return sdk.Template{}, utils.ErrorPerStatusCodeV4(errBody, httpResp, err)
	}

	return response.Data, nil
}

func (c *Client) ListTemplates(ctx context.Context, opts *contracts.ListOptions) (*sdk.PaginatedTemplateList, error) {
	logger.Debug("List Templates")
	if opts.OrderBy == "" {
		opts.OrderBy = "id"
	}
	resp, httpResp, err := c.apiClient.DataStreamTemplatesAPI.ListTemplates(ctx).
		Ordering(opts.OrderBy).
		Page(opts.Page).
		PageSize(opts.PageSize).
		Search(opts.Sort).
		Execute()

	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while listing the templates", zap.Error(err))
			errBody, err = utils.LogAndRewindBodyV4(httpResp)
			if err != nil {
				return nil, err
			}
		}
		return nil, utils.ErrorPerStatusCodeV4(errBody, httpResp, err)
	}

	return resp, nil
}
