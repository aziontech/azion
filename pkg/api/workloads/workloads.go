package workloads

import (
	"context"
	"strconv"

	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/utils"
	"go.uber.org/zap"
)

func (c *Client) Create(ctx context.Context, req *CreateRequest) (WorkloadResponse, error) {
	logger.Debug("Create Workload")
	request := c.apiClient.WorkloadsAPI.CreateWorkload(ctx).WorkloadRequest(req.WorkloadRequest)
	// request := c.apiClient.DomainsAPI.CreateDomain(ctx).CreateDomainRequest(req.CreateDomainRequest)
	workloadsResponse, httpResp, err := request.Execute()
	if err != nil {
		if httpResp != nil {
			logger.Debug("Error while creating a workload", zap.Error(err))
			err := utils.LogAndRewindBody(httpResp)
			if err != nil {
				return nil, err
			}
		}
		return nil, utils.ErrorPerStatusCode(httpResp, err)
	}
	return &workloadsResponse.Data, nil
}

func (c *Client) Delete(ctx context.Context, id int64) error {
	logger.Debug("Delete Workload")
	str := strconv.FormatInt(id, 10)
	req := c.apiClient.WorkloadsAPI.DestroyWorkload(ctx, str)

	_, httpResp, err := req.Execute()
	if err != nil {
		if httpResp != nil {
			logger.Debug("Error while deleting a workload", zap.Error(err))
			err := utils.LogAndRewindBody(httpResp)
			if err != nil {
				return err
			}
		}
		return utils.ErrorPerStatusCode(httpResp, err)
	}

	return nil
}

func (c *Client) Get(ctx context.Context, id string) (WorkloadResponse, error) {
	logger.Debug("Get Workload")
	request := c.apiClient.WorkloadsAPI.RetrieveWorkload(ctx, id)
	res, httpResp, err := request.Execute()
	if err != nil {
		if httpResp != nil {
			logger.Debug("Error while describing a Workload", zap.Error(err))
			err := utils.LogAndRewindBody(httpResp)
			if err != nil {
				return nil, err
			}
		}
		return nil, utils.ErrorPerStatusCode(httpResp, err)
	}
	return &res.Data, nil
}
