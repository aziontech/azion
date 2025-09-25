package function

import (
	"context"

	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/utils"
	sdk "github.com/aziontech/azionapi-v4-go-sdk-dev/edge-api"
	"go.uber.org/zap"
)

func (c *Client) Get(ctx context.Context, id string) (sdk.EdgeFunctions, error) {
	logger.Debug("Get Function")
	request := c.apiClient.FunctionsAPI.RetrieveFunction(ctx, id)

	res, httpResp, err := request.Execute()
	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while getting a Function", zap.Error(err))
			errBody, err = utils.LogAndRewindBodyV4(httpResp)
			if err != nil {
				return sdk.EdgeFunctions{}, err
			}
		}
		return sdk.EdgeFunctions{}, utils.ErrorPerStatusCodeV4(errBody, httpResp, err)
	}

	return res.Data, nil
}

func (c *Client) Delete(ctx context.Context, id string) error {
	logger.Debug("Delete Function")
	request := c.apiClient.FunctionsAPI.DestroyFunction(ctx, id)

	_, httpResp, err := request.Execute()

	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while deleting a Function", zap.Error(err))
			errBody, err = utils.LogAndRewindBodyV4(httpResp)
			if err != nil {
				return err
			}
		}
		return utils.ErrorPerStatusCodeV4(errBody, httpResp, err)
	}

	return nil
}

func (c *Client) Create(ctx context.Context, req *CreateRequest) (sdk.EdgeFunctions, error) {
	// Although there's only one option, the API requires the `language` field.
	// Hard-coding javascript for now
	logger.Debug("Create Function")

	request := c.apiClient.FunctionsAPI.CreateFunction(ctx).EdgeFunctionsRequest(req.EdgeFunctionsRequest)

	edgeFuncResponse, httpResp, err := request.Execute()
	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while creating a Function", zap.Error(err), zap.Any("Name", req.Name))
			errBody, err = utils.LogAndRewindBodyV4(httpResp)
			if err != nil {
				return sdk.EdgeFunctions{}, err
			}
		}
		return sdk.EdgeFunctions{}, utils.ErrorPerStatusCodeV4(errBody, httpResp, err)
	}

	return edgeFuncResponse.Data, nil
}

func (c *Client) Update(ctx context.Context, req *UpdateRequest, id string) (sdk.EdgeFunctions, error) {
	logger.Debug("Update Function", zap.Any("Function ID", id), zap.Any("Function name", req.Name))
	request := c.apiClient.FunctionsAPI.PartialUpdateFunction(ctx, id).PatchedEdgeFunctionsRequest(req.PatchedEdgeFunctionsRequest)

	edgeFuncResponse, httpResp, err := request.Execute()
	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while updating a Function", zap.Error(err), zap.Any("ID", id), zap.Any("Name", req.Name))
			errBody, err = utils.LogAndRewindBodyV4(httpResp)
			if err != nil {
				return sdk.EdgeFunctions{}, err
			}
		}
		return sdk.EdgeFunctions{}, utils.ErrorPerStatusCodeV4(errBody, httpResp, err)
	}

	return edgeFuncResponse.Data, nil
}

func (c *Client) List(ctx context.Context, opts *contracts.ListOptions) (*sdk.PaginatedEdgeFunctionsList, error) {
	logger.Debug("List Functions")
	if opts.OrderBy == "" {
		opts.OrderBy = "id"
	}
	resp, httpResp, err := c.apiClient.FunctionsAPI.ListFunctions(ctx).
		Ordering(opts.OrderBy).
		Page(opts.Page).
		PageSize(opts.PageSize).
		Search(opts.Sort).
		Execute()

	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while listing the Functions", zap.Error(err))
			errBody, err = utils.LogAndRewindBodyV4(httpResp)
			if err != nil {
				return nil, err
			}
		}
		return nil, utils.ErrorPerStatusCodeV4(errBody, httpResp, err)
	}

	return resp, nil
}
