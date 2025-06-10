package edgeconnector

import (
	"context"

	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/utils"
	sdk "github.com/aziontech/azionapi-v4-go-sdk/edge"
	"go.uber.org/zap"
)

func (c *Client) Get(ctx context.Context, id string) (sdk.BaseEdgeConnector, error) {
	logger.Debug("Get Edge Connector")
	request := c.apiClient.EdgeConnectorsAPI.RetrieveEdgeConnector(ctx, id)

	res, httpResp, err := request.Execute()
	if err != nil {
		if httpResp != nil {
			logger.Debug("Error while getting an Edge Connector", zap.Error(err))
			err := utils.LogAndRewindBody(httpResp)
			if err != nil {
				return sdk.BaseEdgeConnector{}, err
			}
		}
		return sdk.BaseEdgeConnector{}, utils.ErrorPerStatusCode(httpResp, err)
	}

	return res.Data, nil
}

func (c *Client) Delete(ctx context.Context, id string) error {
	logger.Debug("Delete Edge Connector")
	request := c.apiClient.EdgeConnectorsAPI.DestroyEdgeConnector(ctx, id)

	_, httpResp, err := request.Execute()

	if err != nil {
		if httpResp != nil {
			logger.Debug("Error while deleting an Edge Connector", zap.Error(err))
			err := utils.LogAndRewindBody(httpResp)
			if err != nil {
				return err
			}
		}
		return utils.ErrorPerStatusCode(httpResp, err)
	}

	return nil
}

func (c *Client) Create(ctx context.Context, req *CreateRequest) (sdk.BaseEdgeConnector, error) {
	logger.Debug("Create Edge Connector")

	request := c.apiClient.EdgeConnectorsAPI.CreateEdgeConnector(ctx).EdgeConnectorPolymorphicRequest(req.EdgeConnectorPolymorphicRequest)

	response, httpResp, err := request.Execute()
	if err != nil {
		if httpResp != nil {
			logger.Debug("Error while creating an Edge Connector", zap.Error(err))
			err := utils.LogAndRewindBody(httpResp)
			if err != nil {
				return sdk.BaseEdgeConnector{}, err
			}
		}
		return sdk.BaseEdgeConnector{}, utils.ErrorPerStatusCode(httpResp, err)
	}

	return response.Data, nil
}

func (c *Client) Update(ctx context.Context, req *UpdateRequest, id string) (sdk.BaseEdgeConnector, error) {
	logger.Debug("Update Edge Connector")
	request := c.apiClient.EdgeConnectorsAPI.PartialUpdateEdgeConnector(ctx, id).PatchedEdgeConnectorPolymorphicRequest(req.PatchedEdgeConnectorPolymorphicRequest)

	response, httpResp, err := request.Execute()
	if err != nil {
		if httpResp != nil {
			logger.Debug("Error while updating an Edge Connector", zap.Error(err))
			err := utils.LogAndRewindBody(httpResp)
			if err != nil {
				return sdk.BaseEdgeConnector{}, err
			}
		}
		return sdk.BaseEdgeConnector{}, utils.ErrorPerStatusCode(httpResp, err)
	}

	return response.Data, nil
}

func (c *Client) List(ctx context.Context, opts *contracts.ListOptions) (*sdk.PaginatedResponseListBaseEdgeConnectorList, error) {
	logger.Debug("List Edge Connectors")
	if opts.OrderBy == "" {
		opts.OrderBy = "id"
	}
	resp, httpResp, err := c.apiClient.EdgeConnectorsAPI.ListEdgeConnectors(ctx).
		Ordering(opts.OrderBy).
		Page(opts.Page).
		PageSize(opts.PageSize).
		Search(opts.Sort).
		Execute()

	if err != nil {
		if httpResp != nil {
			logger.Debug("Error while listing the Edge Connectors", zap.Error(err))
			err := utils.LogAndRewindBody(httpResp)
			if err != nil {
				return nil, err
			}
		}
		return nil, utils.ErrorPerStatusCode(httpResp, err)
	}

	return resp, nil
}
