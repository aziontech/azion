package applications

import (
	"context"

	"github.com/aziontech/azion-cli/pkg/contracts"

	"github.com/aziontech/azion-cli/pkg/api/v6/version"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/utils"
	sdk "github.com/aziontech/azionapi-v6-go-sdk/azion-api"
	"go.uber.org/zap"
)

// ListDeviceGroups lists the device groups of one version, one page at a time.
func (c *Client) ListDeviceGroups(ctx context.Context, id int64, versionID string, opts *contracts.ListOptions) (*sdk.PaginatedDeviceGroupList, error) {
	logger.Debug("List Application device groups")

	req := c.apiClient.ApplicationsDeviceGroupsAPI.
		ListApplicationVersionDeviceGroups(ctx, id, versionID).Page(opts.Page).PageSize(int32(opts.PageSize))
	if opts.OrderBy != "" {
		req = req.Ordering(opts.OrderBy)
	}
	res, httpResp, err := req.Execute()
	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while listing device groups", zap.Error(err))
			var bodyErr error
			errBody, bodyErr = utils.LogAndRewindBodyV4(httpResp)
			if bodyErr != nil {
				return nil, bodyErr
			}
		}
		return nil, version.ErrorFor(errBody, httpResp, err)
	}

	return res, nil
}

// CreateDeviceGroup creates a device group in a Draft.
func (c *Client) CreateDeviceGroup(ctx context.Context, id int64, versionID string, req sdk.DeviceGroupRequest) (version.Item, error) {
	logger.Debug("Create Application device group")

	res, httpResp, err := c.apiClient.ApplicationsDeviceGroupsAPI.
		CreateApplicationVersionDeviceGroups(ctx, id, versionID).DeviceGroupRequest(req).Execute()
	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while creating a device group", zap.Error(err))
			var bodyErr error
			errBody, bodyErr = utils.LogAndRewindBodyV4(httpResp)
			if bodyErr != nil {
				return nil, bodyErr
			}
		}
		return nil, version.ErrorFor(errBody, httpResp, err)
	}

	return &res.Data, nil
}

// UpdateDeviceGroup replaces a device group in a Draft.
func (c *Client) UpdateDeviceGroup(ctx context.Context, id int64, versionID string, subID int64, req sdk.DeviceGroupRequest) (version.Item, error) {
	logger.Debug("Update Application device group")

	res, httpResp, err := c.apiClient.ApplicationsDeviceGroupsAPI.
		UpdateApplicationVersionDeviceGroup(ctx, id, subID, versionID).DeviceGroupRequest(req).Execute()
	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while updating a device group", zap.Error(err))
			var bodyErr error
			errBody, bodyErr = utils.LogAndRewindBodyV4(httpResp)
			if bodyErr != nil {
				return nil, bodyErr
			}
		}
		return nil, version.ErrorFor(errBody, httpResp, err)
	}

	return &res.Data, nil
}

// DeleteDeviceGroup removes a device group from a Draft.
func (c *Client) DeleteDeviceGroup(ctx context.Context, id int64, versionID string, subID int64) error {
	logger.Debug("Delete Application device group")

	_, httpResp, err := c.apiClient.ApplicationsDeviceGroupsAPI.
		DeleteApplicationVersionDeviceGroup(ctx, id, subID, versionID).Execute()
	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while deleting a device group", zap.Error(err))
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
