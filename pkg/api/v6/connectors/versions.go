package connectors

import (
	"context"

	"github.com/aziontech/azion-cli/pkg/api/v6/version"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/utils"
	sdk "github.com/aziontech/azionapi-v6-go-sdk/azion-api"
	"go.uber.org/zap"
)

var _ version.Versioned = (*Client)(nil)

// CreateVersion opens a new Draft. An empty from clones the latest ready
// version. It returns the new version id.
func (c *Client) CreateVersion(ctx context.Context, id int64, from, comment string) (string, error) {
	logger.Debug("Create Connector version")

	body := sdk.ConnectorVersionCreateRequest{}
	if from != "" {
		body.SetSourceVersion(from)
	}
	if comment != "" {
		body.SetComment(comment)
	}
	res, httpResp, err := c.apiClient.ConnectorVersionsAPI.
		CreateConnectorVersion(ctx, id).ConnectorVersionCreateRequest(body).Execute()
	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while creating a Connector version", zap.Error(err))
			var bodyErr error
			errBody, bodyErr = utils.LogAndRewindBodyV4(httpResp)
			if bodyErr != nil {
				return "", bodyErr
			}
		}
		return "", version.ErrorFor(errBody, httpResp, err)
	}

	return versionMeta(res).GetVersionId(), nil
}

// GetVersion reads one version. Waits poll it, so a failure also carries the
// HTTP status (version.NewStatusError) for the retry rule.
func (c *Client) GetVersion(ctx context.Context, id int64, versionID string) (version.Meta, error) {
	logger.Debug("Get Connector version")

	res, httpResp, err := c.apiClient.ConnectorVersionsAPI.
		RetrieveConnectorVersion(ctx, id, versionID).Execute()
	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while getting a Connector version", zap.Error(err))
			var bodyErr error
			errBody, bodyErr = utils.LogAndRewindBodyV4(httpResp)
			if bodyErr != nil {
				return nil, bodyErr
			}
		}
		return nil, version.NewStatusError(httpResp, version.ErrorFor(errBody, httpResp, err))
	}

	return versionMeta(res), nil
}

// UpdateVersion replaces the configuration of a Draft.
func (c *Client) UpdateVersion(ctx context.Context, id int64, versionID string, req sdk.ConnectorRequest) (version.Meta, error) {
	logger.Debug("Update Connector version")

	res, httpResp, err := c.apiClient.ConnectorVersionsAPI.
		UpdateConnectorVersion(ctx, id, versionID).ConnectorRequest(req).Execute()
	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while updating a Connector version", zap.Error(err))
			var bodyErr error
			errBody, bodyErr = utils.LogAndRewindBodyV4(httpResp)
			if bodyErr != nil {
				return nil, bodyErr
			}
		}
		return nil, version.ErrorFor(errBody, httpResp, err)
	}

	return versionMeta(res), nil
}

// Build queues the build of a Draft. It returns once the API accepts it.
func (c *Client) Build(ctx context.Context, id int64, versionID, comment string) error {
	logger.Debug("Build Connector version")

	body := sdk.ConnectorVersionBuildRequest{}
	if comment != "" {
		body.SetComment(comment)
	}
	_, httpResp, err := c.apiClient.ConnectorVersionsAPI.
		BuildConnectorVersion(ctx, id, versionID).ConnectorVersionBuildRequest(body).Execute()
	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while building a Connector version", zap.Error(err))
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

// CancelBuild asks the API to stop a build that has not settled.
func (c *Client) CancelBuild(ctx context.Context, id int64, versionID string) error {
	logger.Debug("Cancel Connector version build")

	_, httpResp, err := c.apiClient.ConnectorVersionsAPI.
		CancelConnectorVersionBuild(ctx, id, versionID).Execute()
	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while canceling a Connector version build", zap.Error(err))
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

// VersionState reads the build state of one version.
func (c *Client) VersionState(ctx context.Context, id int64, versionID string) (version.State, error) {
	obj, err := c.GetVersion(ctx, id, versionID)
	if err != nil {
		return "", err
	}
	return version.State(obj.GetVersionState()), nil
}
