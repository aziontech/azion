package devicegroups

import (
	"context"

	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/utils"
	sdk "github.com/aziontech/azionapi-v4-go-sdk-dev/azion-api"
	"go.uber.org/zap"
)

func (c *Client) Create(ctx context.Context, req sdk.DeviceGroupRequest, applicationID int64) (sdk.DeviceGroup, error) {
	logger.Debug("Create Device Group")

	request := c.apiClient.ApplicationsDeviceGroupsAPI.
		CreateDeviceGroup(ctx, applicationID).
		DeviceGroupRequest(req)
	resp, httpResp, err := request.Execute()
	if err != nil {
		logger.Debug("Error while creating a Device Group", zap.Error(err))
		errBody, err := utils.LogAndRewindBodyV4(httpResp)
		if err != nil {
			return sdk.DeviceGroup{}, err
		}

		return sdk.DeviceGroup{}, utils.ErrorPerStatusCodeV4(errBody, httpResp, err)
	}

	return resp.GetData(), nil
}

func (c *Client) Update(ctx context.Context, req sdk.PatchedDeviceGroupRequest, applicationID, deviceGroupID int64) (sdk.DeviceGroup, error) {
	logger.Debug("Update Device Group")

	request := c.apiClient.ApplicationsDeviceGroupsAPI.
		PartialUpdateDeviceGroup(ctx, applicationID, deviceGroupID).
		PatchedDeviceGroupRequest(req)
	resp, httpResp, err := request.Execute()
	if err != nil {
		logger.Debug("Error while updating a Device Group", zap.Any("ID", deviceGroupID), zap.Error(err))
		errBody, err := utils.LogAndRewindBodyV4(httpResp)
		if err != nil {
			return sdk.DeviceGroup{}, err
		}

		return sdk.DeviceGroup{}, utils.ErrorPerStatusCodeV4(errBody, httpResp, err)
	}

	return resp.GetData(), nil
}

func (c *Client) List(ctx context.Context, opts *contracts.ListOptions, applicationID int64) (*sdk.PaginatedDeviceGroupList, error) {
	logger.Debug("List Device Groups")
	if opts.OrderBy == "" {
		opts.OrderBy = "id"
	}

	resp, httpResp, err := c.apiClient.ApplicationsDeviceGroupsAPI.
		ListDeviceGroups(ctx, applicationID).
		Ordering(opts.OrderBy).
		Page(opts.Page).
		PageSize(opts.PageSize).
		Search(opts.Filter).
		Execute()
	if err != nil {
		logger.Debug("Error while listing Device Groups", zap.Error(err))
		errBody, err := utils.LogAndRewindBodyV4(httpResp)
		if err != nil {
			return nil, err
		}

		return nil, utils.ErrorPerStatusCodeV4(errBody, httpResp, err)
	}

	return resp, nil
}

func (c *Client) Get(ctx context.Context, applicationID, deviceGroupID int64) (sdk.DeviceGroup, error) {
	logger.Debug("Get Device Group")

	resp, httpResp, err := c.apiClient.ApplicationsDeviceGroupsAPI.
		RetrieveDeviceGroup(ctx, applicationID, deviceGroupID).
		Execute()
	if err != nil {
		logger.Debug("Error while getting a Device Group", zap.Error(err))
		errBody, err := utils.LogAndRewindBodyV4(httpResp)
		if err != nil {
			return sdk.DeviceGroup{}, err
		}

		return sdk.DeviceGroup{}, utils.ErrorPerStatusCodeV4(errBody, httpResp, err)
	}

	return resp.GetData(), nil
}

func (c *Client) Delete(ctx context.Context, applicationID, deviceGroupID int64) error {
	logger.Debug("Delete Device Group")

	req := c.apiClient.ApplicationsDeviceGroupsAPI.
		DeleteDeviceGroup(ctx, applicationID, deviceGroupID)
	_, httpResp, err := req.Execute()
	if err != nil {
		logger.Debug("Error while deleting a Device Group", zap.Error(err))
		errBody, err := utils.LogAndRewindBodyV4(httpResp)
		if err != nil {
			return err
		}

		return utils.ErrorPerStatusCodeV4(errBody, httpResp, err)
	}

	return nil
}
