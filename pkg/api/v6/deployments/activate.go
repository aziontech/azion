package deployments

import (
	"context"

	"github.com/aziontech/azion-cli/pkg/api/v6/version"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/utils"
	sdk "github.com/aziontech/azionapi-v6-go-sdk/azion-api"
	"go.uber.org/zap"
)

// ActivateRelease routes the Deployment's traffic to a Ready release.
// deploymentVersionID is the Deployment's current version, as read right
// before; empty leaves it out of the request.
func (c *Client) ActivateRelease(ctx context.Context, deploymentID, releaseID, deploymentVersionID string) (Release, error) {
	logger.Debug("Activate release", zap.String("deployment_version_id", deploymentVersionID))

	req := sdk.ActivateReleaseRequest{}
	if deploymentVersionID != "" {
		req.SetDeploymentVersionId(deploymentVersionID)
	}
	res, httpResp, err := c.routing.
		ActivateRelease(ctx, deploymentID, releaseID).ActivateReleaseRequest(req).Execute()
	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while activating a release", zap.Error(err))
			var bodyErr error
			errBody, bodyErr = utils.LogAndRewindBodyV4(httpResp)
			if bodyErr != nil {
				return Release{}, bodyErr
			}
		}
		return Release{}, version.ErrorFor(errBody, httpResp, err)
	}

	return asyncRelease(deploymentID, res), nil
}
