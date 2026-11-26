package deployments

import (
	"context"

	"github.com/aziontech/azion-cli/pkg/api/v6/version"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/utils"
	sdk "github.com/aziontech/azionapi-v6-go-sdk/azion-api"
	"go.uber.org/zap"
)

// CreateRelease opens a release Draft on a Deployment.
func (c *Client) CreateRelease(ctx context.Context, deploymentID string, resources []ReleaseResource, originType string) (Release, error) {
	logger.Debug("Create release")

	req := sdk.ReleaseCreateRequest{Origin: &sdk.ReleaseOrigin{Type: originType}}
	for _, r := range resources {
		req.Resources = append(req.Resources, sdk.ReleaseCreateRequestResourcesInner{
			ResourceId: r.ResourceID, ResourceType: r.ResourceType, VersionId: r.VersionID,
		})
	}
	logger.Debug("Release create request", zap.String("deployment_id", deploymentID), zap.Any("body", req))
	res, httpResp, err := c.releases.
		CreateRelease(ctx, deploymentID).ReleaseCreateRequest(req).Execute()
	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while creating a release", zap.Error(err))
			var bodyErr error
			errBody, bodyErr = utils.LogAndRewindBodyV4(httpResp)
			if bodyErr != nil {
				return Release{}, bodyErr
			}
		}
		return Release{}, version.ErrorFor(errBody, httpResp, err)
	}

	// What the API stored, to compare with the request above.
	logger.Debug("Release create response", zap.Any("body", res))

	return asyncRelease(deploymentID, res), nil
}

// BuildRelease queues the build of a release Draft.
func (c *Client) BuildRelease(ctx context.Context, deploymentID, releaseID string) (Release, error) {
	logger.Debug("Build release")

	// The SDK requires a body; the API takes none.
	res, httpResp, err := c.releases.
		BuildRelease(ctx, deploymentID, releaseID).Body(map[string]interface{}{}).Execute()
	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while building a release", zap.Error(err))
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

// GetRelease reads one release. Waits poll it, so a failure also carries the
// HTTP status (version.NewStatusError) for the retry rule.
func (c *Client) GetRelease(ctx context.Context, deploymentID, releaseID string) (Release, error) {
	logger.Debug("Get release")

	res, httpResp, err := c.releases.
		RetrieveRelease(ctx, deploymentID, releaseID).Execute()
	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while getting a release", zap.Error(err))
			var bodyErr error
			errBody, bodyErr = utils.LogAndRewindBodyV4(httpResp)
			if bodyErr != nil {
				return Release{}, bodyErr
			}
		}
		return Release{}, version.NewStatusError(httpResp, version.ErrorFor(errBody, httpResp, err))
	}

	return releaseFrom(&res.Data), nil
}

// CancelRelease asks the API to stop a release build that has not settled.
func (c *Client) CancelRelease(ctx context.Context, deploymentID, releaseID string) error {
	logger.Debug("Cancel release")

	_, httpResp, err := c.releases.
		CancelRelease(ctx, deploymentID, releaseID).CancelReleaseRequest(sdk.CancelReleaseRequest{}).Execute()
	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while canceling a release", zap.Error(err))
			var bodyErr error
			errBody, bodyErr = utils.LogAndRewindBodyV4(httpResp)
			if bodyErr != nil {
				return bodyErr
			}
		}
		return version.ErrorFor(errBody, httpResp, err)
	}

	return nil
}

// DeleteRelease removes a release that never went live, so failed runs do not
// use up the Deployment's release Drafts.
func (c *Client) DeleteRelease(ctx context.Context, deploymentID, releaseID string) error {
	logger.Debug("Delete release")

	httpResp, err := c.releases.
		DeleteRelease(ctx, deploymentID, releaseID).Execute()
	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while deleting a release", zap.Error(err))
			var bodyErr error
			errBody, bodyErr = utils.LogAndRewindBodyV4(httpResp)
			if bodyErr != nil {
				return bodyErr
			}
		}
		return version.ErrorFor(errBody, httpResp, err)
	}

	return nil
}
