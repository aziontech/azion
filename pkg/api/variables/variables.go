package variables

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/utils"
)

func (c *Client) List(ctx context.Context) ([]Response, error) {
	logger.Debug("List Environment Variables")

	resp, httpResp, err := c.apiClient.VariablesAPI.ApiVariablesList(ctx).Execute()
	if err != nil {
		fmt.Println(err.Error())
		if httpResp != nil {
			logger.Debug("Error while listing variables", zap.Error(err))
			err := utils.LogAndRewindBody(httpResp)
			if err != nil {
				return nil, err
			}
		}
		return nil, utils.ErrorPerStatusCode(httpResp, err)
	}

	var result []Response

	for i := range resp {
		result = append(result, &resp[i])
	}

	return result, nil
}

func (c *Client) Delete(ctx context.Context, id string) error {
	logger.Debug("Delete Environment Variable")

	req := c.apiClient.VariablesAPI.ApiVariablesDestroy(ctx, id)

	httpResp, err := req.Execute()
	if err != nil {
		if httpResp != nil {
			logger.Debug("Error while deleting a variables", zap.Error(err))
			err := utils.LogAndRewindBody(httpResp)
			if err != nil {
				return err
			}
		}
		return utils.ErrorPerStatusCode(httpResp, err)
	}

	return nil
}

func (c *Client) Update(ctx context.Context, req *Request) (Response, error) {
	logger.Debug("Update Environment Variable")

	request := c.apiClient.VariablesAPI.ApiVariablesUpdate(ctx, req.Uuid).VariableCreate(req.VariableCreate)

	resp, httpResp, err := request.Execute()
	if err != nil {
		if httpResp != nil {
			logger.Debug("Error while updating a variables", zap.Error(err))
			err := utils.LogAndRewindBody(httpResp)
			if err != nil {
				return nil, err
			}
		}
		return nil, utils.ErrorPerStatusCode(httpResp, err)
	}

	return resp, nil
}

func (c *Client) Get(ctx context.Context, id string) (Response, error) {
	logger.Debug("Get Environment Variable")

	req := c.apiClient.VariablesAPI.ApiVariablesRetrieve(ctx, id)
	resp, httpResp, err := req.Execute()
	if err != nil {
		if httpResp != nil {
			logger.Debug("Error while getting a variables", zap.Error(err))
			err := utils.LogAndRewindBody(httpResp)
			if err != nil {
				return nil, err
			}
		}
		return nil, utils.ErrorPerStatusCode(httpResp, err)
	}

	return resp, nil
}

func (c *Client) Create(ctx context.Context, strReq Request) (Response, error) {
	logger.Debug("Create Environment Variable")

	resp, httpResp, err := c.apiClient.VariablesAPI.ApiVariablesCreate(ctx).
		VariableCreate(strReq.VariableCreate).Execute()
	if err != nil {
		if httpResp != nil {
			logger.Debug("Error while creating a variables", zap.Error(err))
			err := utils.LogAndRewindBody(httpResp)
			if err != nil {
				return nil, err
			}
		}
		return nil, utils.ErrorPerStatusCode(httpResp, err)
	}
	return resp, nil
}
