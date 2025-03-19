package workloads

import (
	"context"

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
			logger.Debug("Error while creating a domain", zap.Error(err))
			err := utils.LogAndRewindBody(httpResp)
			if err != nil {
				return nil, err
			}
		}
		return nil, utils.ErrorPerStatusCode(httpResp, err)
	}
	return &workloadsResponse.Data, nil
}
