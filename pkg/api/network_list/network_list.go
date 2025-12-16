package networklist

import (
	"context"

	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/utils"
	sdk "github.com/aziontech/azionapi-v4-go-sdk-dev/edge-api"
	"go.uber.org/zap"
)

func (c *Client) List(ctx context.Context, opts *contracts.ListOptions) (*sdk.PaginatedNetworkListList, error) {
	logger.Debug("List network lists")
	if opts.OrderBy == "" {
		opts.OrderBy = "id"
	}
	resp, httpResp, err := c.apiClient.NetworkListsAPI.ListNetworkLists(ctx).
		Ordering(opts.OrderBy).
		Page(opts.Page).
		PageSize(opts.PageSize).
		Search(opts.Sort).
		Execute()

	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while listing the network lists", zap.Error(err))
			errBody, err = utils.LogAndRewindBodyV4(httpResp)
			if err != nil {
				return nil, err
			}
		}
		return nil, utils.ErrorPerStatusCodeV4(errBody, httpResp, err)
	}

	return resp, nil
}

func (c *Client) Delete(ctx context.Context, id int64) error {
	logger.Debug("Delete network list")
	request := c.apiClient.NetworkListsAPI.DeleteNetworkList(ctx, id)

	_, httpResp, err := request.Execute()

	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while deleting a network list", zap.Error(err))
			errBody, err = utils.LogAndRewindBodyV4(httpResp)
			if err != nil {
				return err
			}
		}
		return utils.ErrorPerStatusCodeV4(errBody, httpResp, err)
	}

	return nil
}

func (c *Client) Get(ctx context.Context, id int64) (sdk.NetworkListDetail, error) {
	logger.Debug("Get network list")
	request := c.apiClient.NetworkListsAPI.RetrieveNetworkList(ctx, id)

	res, httpResp, err := request.Execute()
	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while getting a network list", zap.Error(err))
			errBody, err = utils.LogAndRewindBodyV4(httpResp)
			if err != nil {
				return sdk.NetworkListDetail{}, err
			}
		}
		return sdk.NetworkListDetail{}, utils.ErrorPerStatusCodeV4(errBody, httpResp, err)
	}

	return res.Data, nil
}

func (c *Client) Create(ctx context.Context, req *CreateRequest) (sdk.NetworkListDetail, error) {
	logger.Debug("Create network list")

	request := c.apiClient.NetworkListsAPI.CreateNetworkList(ctx).NetworkListDetailRequest(req.NetworkListDetailRequest)

	edgeFuncResponse, httpResp, err := request.Execute()
	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while creating a network list", zap.Error(err), zap.Any("Name", req.Name))
			errBody, err = utils.LogAndRewindBodyV4(httpResp)
			if err != nil {
				return sdk.NetworkListDetail{}, err
			}
		}
		return sdk.NetworkListDetail{}, utils.ErrorPerStatusCodeV4(errBody, httpResp, err)
	}

	return edgeFuncResponse.Data, nil
}

func (c *Client) Update(ctx context.Context, req *UpdateRequest, id int64) (sdk.NetworkListDetail, error) {
	logger.Debug("Update network list", zap.Any("ID", id), zap.Any("Name", req.Name))
	request := c.apiClient.NetworkListsAPI.PartialUpdateNetworkList(ctx, id).PatchedNetworkListDetailRequest(req.PatchedNetworkListDetailRequest)

	edgeFuncResponse, httpResp, err := request.Execute()
	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while updating a network list", zap.Error(err), zap.Any("ID", id), zap.Any("Name", req.Name))
			errBody, err = utils.LogAndRewindBodyV4(httpResp)
			if err != nil {
				return sdk.NetworkListDetail{}, err
			}
		}
		return sdk.NetworkListDetail{}, utils.ErrorPerStatusCodeV4(errBody, httpResp, err)
	}

	return edgeFuncResponse.Data, nil
}
