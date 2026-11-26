package environments

import (
	"context"

	sdk "github.com/aziontech/azionapi-v6-go-sdk/azion-api"

	"github.com/aziontech/azion-cli/pkg/api/v6/version"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/utils"
	"go.uber.org/zap"
)

// FindByName returns the Environment with exactly this name, if any.
func (c *Client) FindByName(ctx context.Context, name string) (Environment, bool, error) {
	logger.Debug("List Environments")

	res, httpResp, err := c.apiClient.EnvironmentsAPI.
		ListEnvironments(ctx).Name(name).PageSize(100).Execute()
	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while listing Environments", zap.Error(err))
			var bodyErr error
			errBody, bodyErr = utils.LogAndRewindBodyV4(httpResp)
			if bodyErr != nil {
				return Environment{}, false, bodyErr
			}
		}
		return Environment{}, false, version.ErrorFor(errBody, httpResp, err)
	}

	results := res.GetResults()
	for i := range results {
		if results[i].GetName() == name {
			return environmentFrom(&results[i]), true, nil
		}
	}
	return Environment{}, false, nil
}

// Get reads one Environment.
func (c *Client) Get(ctx context.Context, id string) (Environment, error) {
	logger.Debug("Get Environment")

	res, httpResp, err := c.apiClient.EnvironmentsAPI.
		RetrieveEnvironment(ctx, id).Execute()
	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while getting an Environment", zap.Error(err))
			var bodyErr error
			errBody, bodyErr = utils.LogAndRewindBodyV4(httpResp)
			if bodyErr != nil {
				return Environment{}, bodyErr
			}
		}
		return Environment{}, version.ErrorFor(errBody, httpResp, err)
	}

	return environmentFrom(&res.Data), nil
}

// Create creates an Environment. Its deployment policy cannot change later.
func (c *Client) Create(ctx context.Context, req sdk.EnvironmentCreateRequest) (Environment, error) {
	logger.Debug("Create Environment")

	res, httpResp, err := c.apiClient.EnvironmentsAPI.
		CreateEnvironment(ctx).EnvironmentCreateRequest(req).Execute()
	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while creating an Environment", zap.Error(err))
			var bodyErr error
			errBody, bodyErr = utils.LogAndRewindBodyV4(httpResp)
			if bodyErr != nil {
				return Environment{}, bodyErr
			}
		}
		return Environment{}, version.ErrorFor(errBody, httpResp, err)
	}

	return environmentFrom(&res.Data), nil
}
