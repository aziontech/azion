package workloads

import (
	"context"

	"github.com/aziontech/azion-cli/pkg/api/v6/version"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/utils"
	sdk "github.com/aziontech/azionapi-v6-go-sdk/azion-api"
	"go.uber.org/zap"
)

// Get reads one Workload.
func (c *Client) Get(ctx context.Context, id int64) (sdk.Workload, error) {
	logger.Debug("Get Workload")

	res, httpResp, err := c.apiClient.WorkloadsAPI.
		RetrieveWorkload(ctx, id).Execute()
	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while getting a Workload", zap.Error(err))
			var bodyErr error
			errBody, bodyErr = utils.LogAndRewindBodyV4(httpResp)
			if bodyErr != nil {
				return sdk.Workload{}, bodyErr
			}
		}
		return sdk.Workload{}, version.ErrorFor(errBody, httpResp, err)
	}

	return res.Data, nil
}

// FindByName returns the Workload with exactly this name, if any.
func (c *Client) FindByName(ctx context.Context, name string) (sdk.Workload, bool, error) {
	logger.Debug("List Workloads")

	res, httpResp, err := c.apiClient.WorkloadsAPI.
		ListWorkloads(ctx).Name(name).PageSize(100).Execute()
	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while listing Workloads", zap.Error(err))
			var bodyErr error
			errBody, bodyErr = utils.LogAndRewindBodyV4(httpResp)
			if bodyErr != nil {
				return sdk.Workload{}, false, bodyErr
			}
		}
		return sdk.Workload{}, false, version.ErrorFor(errBody, httpResp, err)
	}

	for _, w := range res.GetResults() {
		if w.GetName() == name {
			return w, true, nil
		}
	}
	return sdk.Workload{}, false, nil
}

// Create creates a Workload.
func (c *Client) Create(ctx context.Context, req sdk.WorkloadRequest) (sdk.Workload, error) {
	logger.Debug("Create Workload")

	res, httpResp, err := c.apiClient.WorkloadsAPI.
		CreateWorkload(ctx).WorkloadRequest(req).Execute()
	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while creating a Workload", zap.Error(err))
			var bodyErr error
			errBody, bodyErr = utils.LogAndRewindBodyV4(httpResp)
			if bodyErr != nil {
				return sdk.Workload{}, bodyErr
			}
		}
		return sdk.Workload{}, version.ErrorFor(errBody, httpResp, err)
	}

	return res.Data, nil
}

// Update replaces a Workload. The request is complete: fields left out are
// reset, which is why callers start from Get.
func (c *Client) Update(ctx context.Context, id int64, req sdk.WorkloadRequest) (sdk.Workload, error) {
	logger.Debug("Update Workload")

	res, httpResp, err := c.apiClient.WorkloadsAPI.
		UpdateWorkload(ctx, id).WorkloadRequest(req).Execute()
	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while updating a Workload", zap.Error(err))
			var bodyErr error
			errBody, bodyErr = utils.LogAndRewindBodyV4(httpResp)
			if bodyErr != nil {
				return sdk.Workload{}, bodyErr
			}
		}
		return sdk.Workload{}, version.ErrorFor(errBody, httpResp, err)
	}

	return res.Data, nil
}
