package connector

import (
	"context"

	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/utils"
	sdk "github.com/aziontech/azionapi-v4-go-sdk-dev/azion-api"
	"go.uber.org/zap"
)

func (c *Client) Get(ctx context.Context, id int64) (sdk.Connector2, error) {
	logger.Debug("Get Connector")
	request := c.apiClient.ConnectorsAPI.RetrieveConnector(ctx, id)

	res, httpResp, err := request.Execute()
	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while getting a Connector", zap.Error(err))
			errBody, err = utils.LogAndRewindBodyV4(httpResp)
			if err != nil {
				return sdk.Connector2{}, err
			}
		}
		return sdk.Connector2{}, utils.ErrorPerStatusCodeV4(errBody, httpResp, err)
	}

	return res.Data, nil
}

func (c *Client) Delete(ctx context.Context, id int64) error {
	logger.Debug("Delete Connector")
	request := c.apiClient.ConnectorsAPI.DeleteConnector(ctx, id)

	_, httpResp, err := request.Execute()

	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while deleting a Connector", zap.Error(err))
			errBody, err = utils.LogAndRewindBodyV4(httpResp)
			if err != nil {
				return err
			}
		}
		return utils.ErrorPerStatusCodeV4(errBody, httpResp, err)
	}

	return nil
}

func (c *Client) Create(ctx context.Context, req *CreateRequest) (sdk.Connector2, error) {
	logger.Debug("Create Connector")

	request := c.apiClient.ConnectorsAPI.CreateConnector(ctx).ConnectorRequest2(req.ConnectorRequest2)

	response, httpResp, err := request.Execute()
	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while creating a Connector", zap.Error(err))
			errBody, err = utils.LogAndRewindBodyV4(httpResp)
			if err != nil {
				return sdk.Connector2{}, err
			}
		}
		return sdk.Connector2{}, utils.ErrorPerStatusCodeV4(errBody, httpResp, err)
	}

	return response.Data, nil
}

func (c *Client) Update(ctx context.Context, req *UpdateRequest, id int64) (sdk.Connector2, error) {
	logger.Debug("Update Connector")
	request := c.apiClient.ConnectorsAPI.PartialUpdateConnector(ctx, id).PatchedConnectorRequest2(req.PatchedConnectorRequest2)

	response, httpResp, err := request.Execute()
	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while updating a Connector", zap.Error(err), zap.Any("ID", id))
			errBody, err = utils.LogAndRewindBodyV4(httpResp)
			if err != nil {
				return sdk.Connector2{}, err
			}
		}
		return sdk.Connector2{}, utils.ErrorPerStatusCodeV4(errBody, httpResp, err)
	}

	return response.Data, nil
}

func (c *Client) List(ctx context.Context, opts *contracts.ListOptions) (*sdk.PaginatedConnectorList, error) {
	logger.Debug("List Connectors")
	if opts.OrderBy == "" {
		opts.OrderBy = "id"
	}
	resp, httpResp, err := c.apiClient.ConnectorsAPI.ListConnectors(ctx).
		Ordering(opts.OrderBy).
		// Page(opts.Page).
		// PageSize(opts.PageSize).
		Search(opts.Sort).
		Execute()

	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while listing the Connectors", zap.Error(err))
			errBody, err = utils.LogAndRewindBodyV4(httpResp)
			if err != nil {
				return nil, err
			}
		}
		return nil, utils.ErrorPerStatusCodeV4(errBody, httpResp, err)
	}

	return resp, nil
}
