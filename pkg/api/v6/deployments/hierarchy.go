package deployments

import (
	"context"

	"github.com/aziontech/azion-cli/pkg/api/v6/version"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/utils"
	"go.uber.org/zap"
)

// Hierarchy reads which resource versions the Deployment serves.
func (c *Client) Hierarchy(ctx context.Context, deploymentID string) (Hierarchy, error) {
	logger.Debug("Get Deployment resource hierarchy")

	res, httpResp, err := c.usage.
		RetrieveDeploymentResourceHierarchy(ctx, deploymentID).Execute()
	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while getting a Deployment's resource hierarchy", zap.Error(err))
			var bodyErr error
			errBody, bodyErr = utils.LogAndRewindBodyV4(httpResp)
			if bodyErr != nil {
				return Hierarchy{}, bodyErr
			}
		}
		return Hierarchy{}, version.ErrorFor(errBody, httpResp, err)
	}

	return hierarchyFrom(&res.Data), nil
}
