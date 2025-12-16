package functioninstance

import (
	"context"

	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/utils"
	sdk "github.com/aziontech/azionapi-v4-go-sdk-dev/edge-api"
	"go.uber.org/zap"
)

func (c *Client) List(ctx context.Context, opts *contracts.ListOptions, edgeApplicationID int64) (*sdk.PaginatedApplicationFunctionInstanceList, error) {

	logger.Debug("List Function Instances")

	req := c.apiClient.ApplicationsFunctionAPI.
		ListApplicationFunctionInstances(ctx, edgeApplicationID).
		Page(opts.Page).
		PageSize(opts.PageSize)
	resp, httpResp, err := req.Execute()

	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while listing Function Instances", zap.Error(err))
			errBody, err = utils.LogAndRewindBodyV4(httpResp)
			if err != nil {
				return nil, err
			}
		}
		return nil, utils.ErrorPerStatusCodeV4(errBody, httpResp, err)
	}

	return resp, nil
}

func (c *Client) Get(ctx context.Context, edgeApplicationID, functionInstanceID int64) (sdk.ApplicationFunctionInstance, error) {

	logger.Debug("Get Function Instance")
	req := c.apiClient.ApplicationsFunctionAPI.RetrieveApplicationFunctionInstance(ctx, edgeApplicationID, functionInstanceID)
	resp, httpResp, err := req.Execute()

	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while getting Function Instance", zap.Error(err))
			errBody, err = utils.LogAndRewindBodyV4(httpResp)
			if err != nil {
				return sdk.ApplicationFunctionInstance{}, err
			}
		}
		return sdk.ApplicationFunctionInstance{}, utils.ErrorPerStatusCodeV4(errBody, httpResp, err)
	}

	return resp.Data, nil
}

func (c *Client) Delete(ctx context.Context, edgeApplicationID, functionInstanceID int64) error {
	logger.Debug("Delete Function Instance")
	req := c.apiClient.ApplicationsFunctionAPI.DeleteApplicationFunctionInstance(ctx, edgeApplicationID, functionInstanceID)
	_, httpResp, err := req.Execute()

	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while deleting Function Instance", zap.Error(err))
			errBody, err = utils.LogAndRewindBodyV4(httpResp)
			if err != nil {
				return err
			}
		}
		return utils.ErrorPerStatusCodeV4(errBody, httpResp, err)
	}

	return nil
}

func (c *Client) Create(ctx context.Context, edgeApplicationID int64, req sdk.ApplicationFunctionInstanceRequest) (sdk.ApplicationFunctionInstance, error) {
	logger.Debug("Create Function Instance")
	request := c.apiClient.ApplicationsFunctionAPI.CreateApplicationFunctionInstance(ctx, edgeApplicationID).ApplicationFunctionInstanceRequest(req)
	resp, httpResp, err := request.Execute()

	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while creating Function Instance", zap.Error(err))
			errBody, err = utils.LogAndRewindBodyV4(httpResp)
			if err != nil {
				return sdk.ApplicationFunctionInstance{}, err
			}
		}
		return sdk.ApplicationFunctionInstance{}, utils.ErrorPerStatusCodeV4(errBody, httpResp, err)
	}

	return resp.Data, nil
}

func (c *Client) Update(ctx context.Context, edgeApplicationID, instanceID int64, req sdk.PatchedApplicationFunctionInstanceRequest) (sdk.ApplicationFunctionInstance, error) {
	logger.Debug("Update Function Instance")
	request := c.apiClient.ApplicationsFunctionAPI.PartialUpdateApplicationFunctionInstance(ctx, edgeApplicationID, instanceID).PatchedApplicationFunctionInstanceRequest(req)
	resp, httpResp, err := request.Execute()

	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while updating Function Instance", zap.Error(err))
			errBody, err = utils.LogAndRewindBodyV4(httpResp)
			if err != nil {
				return sdk.ApplicationFunctionInstance{}, err
			}
		}
		return sdk.ApplicationFunctionInstance{}, utils.ErrorPerStatusCodeV4(errBody, httpResp, err)
	}

	return resp.Data, nil
}
