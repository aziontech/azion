package deployments

import (
	"context"

	"github.com/aziontech/azion-cli/pkg/api/v6/version"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/utils"
	sdk "github.com/aziontech/azionapi-v6-go-sdk/azion-api"
	"go.uber.org/zap"
)

// Get reads one Deployment. Waits poll it, so a failure also carries the HTTP
// status (version.NewStatusError) for the retry rule.
func (c *Client) Get(ctx context.Context, id string) (Deployment, error) {
	logger.Debug("Get Deployment")

	res, httpResp, err := c.deployments.
		RetrieveDeployment(ctx, id).Execute()
	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while getting a Deployment", zap.Error(err))
			var bodyErr error
			errBody, bodyErr = utils.LogAndRewindBodyV4(httpResp)
			if bodyErr != nil {
				return Deployment{}, bodyErr
			}
		}
		return Deployment{}, version.NewStatusError(httpResp, version.ErrorFor(errBody, httpResp, err))
	}

	return deploymentFrom(&res.Data), nil
}

// FindByName returns the Deployment with exactly this name, if any.
func (c *Client) FindByName(ctx context.Context, name string) (Deployment, bool, error) {
	logger.Debug("List Deployments")

	res, httpResp, err := c.deployments.
		ListDeployments(ctx).Name(name).PageSize(100).Execute()
	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while listing Deployments", zap.Error(err))
			var bodyErr error
			errBody, bodyErr = utils.LogAndRewindBodyV4(httpResp)
			if bodyErr != nil {
				return Deployment{}, false, bodyErr
			}
		}
		return Deployment{}, false, version.ErrorFor(errBody, httpResp, err)
	}

	results := res.GetResults()
	for i := range results {
		if results[i].GetName() == name {
			return deploymentFrom(&results[i]), true, nil
		}
	}
	return Deployment{}, false, nil
}

// Create creates a Deployment. Its binding policy cannot change later.
func (c *Client) Create(ctx context.Context, req sdk.DeploymentCreateRequest) (Deployment, error) {
	logger.Debug("Create Deployment")

	res, httpResp, err := c.deployments.
		CreateDeployment(ctx).DeploymentCreateRequest(req).Execute()
	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while creating a Deployment", zap.Error(err))
			var bodyErr error
			errBody, bodyErr = utils.LogAndRewindBodyV4(httpResp)
			if bodyErr != nil {
				return Deployment{}, bodyErr
			}
		}
		return Deployment{}, version.ErrorFor(errBody, httpResp, err)
	}

	return deploymentFrom(&res.Data), nil
}
