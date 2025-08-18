package workloads

import (
	"context"

	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/utils"
	sdk "github.com/aziontech/azionapi-v4-go-sdk/edge-api"
	"go.uber.org/zap"
)

func (c *Client) Create(ctx context.Context, req *CreateRequest) (WorkloadResponse, error) {
	logger.Debug("Create Workload")
	request := c.apiClient.WorkloadsAPI.CreateWorkload(ctx).WorkloadRequest(req.WorkloadRequest)
	// request := c.apiClient.DomainsAPI.CreateDomain(ctx).CreateDomainRequest(req.CreateDomainRequest)
	workloadsResponse, httpResp, err := request.Execute()
	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while creating a workload", zap.Error(err), zap.Any("Name", req.Name))
			errBody, err = utils.LogAndRewindBodyV4(httpResp)
			if err != nil {
				return nil, err
			}
		}
		return nil, utils.ErrorPerStatusCodeV4(errBody, httpResp, err)
	}
	return &workloadsResponse.Data, nil
}

func (c *Client) Delete(ctx context.Context, id int64) error {
	logger.Debug("Delete Workload")
	req := c.apiClient.WorkloadsAPI.DestroyWorkload(ctx, id)

	_, httpResp, err := req.Execute()
	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while deleting a workload", zap.Error(err))
			errBody, err = utils.LogAndRewindBodyV4(httpResp)
			if err != nil {
				return err
			}
		}
		return utils.ErrorPerStatusCodeV4(errBody, httpResp, err)
	}

	return nil
}

func (c *Client) Get(ctx context.Context, id int64) (WorkloadResponse, error) {
	logger.Debug("Get Workload")
	request := c.apiClient.WorkloadsAPI.RetrieveWorkload(ctx, id)
	res, httpResp, err := request.Execute()
	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while describing a Workload", zap.Error(err))
			errBody, err = utils.LogAndRewindBodyV4(httpResp)
			if err != nil {
				return nil, err
			}
		}
		return nil, utils.ErrorPerStatusCodeV4(errBody, httpResp, err)
	}
	return &res.Data, nil
}

func (c *Client) Update(ctx context.Context, req *UpdateRequest) (WorkloadResponse, error) {
	logger.Debug("Update Workload (PATCH)")
	request := c.apiClient.WorkloadsAPI.PartialUpdateWorkload(ctx, req.Id).PatchedWorkloadRequest(req.PatchedWorkloadRequest)

	workloadsResponse, httpResp, err := request.Execute()

	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while updating a workload (PATCH)", zap.Error(err), zap.Any("ID", req.Id), zap.Any("Name", req.Name))
			errBody, err = utils.LogAndRewindBodyV4(httpResp)
			if err != nil {
				return nil, err
			}
		}
		return nil, utils.ErrorPerStatusCodeV4(errBody, httpResp, err)
	}

	return &workloadsResponse.Data, nil
}

func (c *Client) List(ctx context.Context, opts *contracts.ListOptions) (*sdk.PaginatedWorkloadList, error) {
	logger.Debug("List Workloads")
	if opts.OrderBy == "" {
		opts.OrderBy = "id"
	}
	resp, httpResp, err := c.apiClient.WorkloadsAPI.ListWorkloads(ctx).
		Ordering(opts.OrderBy).
		Page(opts.Page).
		PageSize(opts.PageSize).
		Search(opts.Sort).
		Execute()

	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while listing workloads", zap.Error(err))
			errBody, err = utils.LogAndRewindBodyV4(httpResp)
			if err != nil {
				return nil, err
			}
		}
		return nil, utils.ErrorPerStatusCodeV4(errBody, httpResp, err)
	}

	return resp, nil
}

func (c *Client) ListDeployments(ctx context.Context, opts *contracts.ListOptions, id int64) (*sdk.PaginatedWorkloadDeploymentList, error) {
	logger.Debug("List Workload Deployments")
	if opts.OrderBy == "" {
		opts.OrderBy = "id"
	}
	resp, httpResp, err := c.apiClient.WorkloadDeploymentsAPI.ListWorkloadDeployments(ctx, id).
		Ordering(opts.OrderBy).
		Page(opts.Page).
		PageSize(opts.PageSize).
		Search(opts.Sort).
		Execute()

	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while listing workload deployments", zap.Error(err))
			errBody, err = utils.LogAndRewindBodyV4(httpResp)
			if err != nil {
				return nil, err
			}
		}
		return nil, utils.ErrorPerStatusCodeV4(errBody, httpResp, err)
	}

	return resp, nil
}

func (c *Client) GetDeployment(ctx context.Context, id, deploymentid int64) (DeploymentResponse, error) {
	logger.Debug("Get Workload Deployment")
	request := c.apiClient.WorkloadDeploymentsAPI.RetrieveWorkloadDeployment(ctx, id, deploymentid)
	res, httpResp, err := request.Execute()
	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while describing a Workload Deployment", zap.Error(err))
			errBody, err = utils.LogAndRewindBodyV4(httpResp)
			if err != nil {
				return nil, err
			}
		}
		return nil, utils.ErrorPerStatusCodeV4(errBody, httpResp, err)
	}
	return &res.Data, nil
}

func (c *Client) CreateDeployment(ctx context.Context, req sdk.WorkloadDeploymentRequest, id int64) (DeploymentResponse, error) {
	logger.Debug("Create Workload Deployment")
	request := c.apiClient.WorkloadDeploymentsAPI.CreateWorkloadDeployment(ctx, id).WorkloadDeploymentRequest(req)
	// request := c.apiClient.DomainsAPI.CreateDomain(ctx).CreateDomainRequest(req.CreateDomainRequest)
	workloadDeploymentsResponse, httpResp, err := request.Execute()
	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while creating a workload deployment", zap.Error(err), zap.Any("Name", req.Name))
			errBody, err = utils.LogAndRewindBodyV4(httpResp)
			if err != nil {
				return nil, err
			}
		}
		return nil, utils.ErrorPerStatusCodeV4(errBody, httpResp, err)
	}
	return &workloadDeploymentsResponse.Data, nil
}

func (c *Client) UpdateDeployment(ctx context.Context, req sdk.PatchedWorkloadDeploymentRequest, id int64, deploymentid int64) (DeploymentResponse, error) {
	logger.Debug("Update Workload Deployment")
	request := c.apiClient.WorkloadDeploymentsAPI.PartialUpdateWorkloadDeployment(ctx, deploymentid, id).PatchedWorkloadDeploymentRequest(req)
	// request := c.apiClient.DomainsAPI.CreateDomain(ctx).CreateDomainRequest(req.CreateDomainRequest)
	workloadDeploymentsResponse, httpResp, err := request.Execute()
	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while updating a workload deployment", zap.Error(err), zap.Any("Name", req.Name))
			errBody, err = utils.LogAndRewindBodyV4(httpResp)
			if err != nil {
				return nil, err
			}
		}
		return nil, utils.ErrorPerStatusCodeV4(errBody, httpResp, err)
	}
	return &workloadDeploymentsResponse.Data, nil
}
