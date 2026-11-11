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
func (c *Client) ActivateRelease(ctx context.Context, deploymentID, releaseID string) (Release, error) {
	logger.Debug("Activate release")

	res, httpResp, err := c.routing.
		ActivateRelease(ctx, deploymentID, releaseID).ActivateReleaseRequest(sdk.ActivateReleaseRequest{}).Execute()
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
