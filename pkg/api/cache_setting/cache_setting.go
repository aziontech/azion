package cachesetting

import (
	"context"

	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/utils"
	sdk "github.com/aziontech/azionapi-go-sdk/edgeapplications"
	"go.uber.org/zap"
)

func (c *Client) Create(ctx context.Context, req *CreateRequest, applicationId int64) (Response, error) {
	logger.Debug("Create Cache Settings")

	request := c.apiClient.EdgeApplicationsCacheSettingsAPI.
		EdgeApplicationsEdgeApplicationIdCacheSettingsPost(ctx, applicationId).
		ApplicationCacheCreateRequest(req.ApplicationCacheCreateRequest)
	cacheResponse, httpResp, err := request.Execute()
	if err != nil {
		logger.Debug("Error while creating a cache setting", zap.Error(err))
		err = utils.LogAndRewindBody(httpResp)
		if err != nil {
			return nil, err
		}

		return nil, utils.ErrorPerStatusCode(httpResp, err)
	}

	return cacheResponse.Results, nil
}

func (c *Client) Update(ctx context.Context, req *UpdateRequest, applicationID, cacheSettingID int64) (Response, error) {
	logger.Debug("Update Cache Settings")

	request := c.apiClient.EdgeApplicationsCacheSettingsAPI.
		EdgeApplicationsEdgeApplicationIdCacheSettingsCacheSettingsIdPatch(ctx, applicationID, cacheSettingID).
		ApplicationCachePatchRequest(req.ApplicationCachePatchRequest)
	cacheResponse, httpResp, err := request.Execute()
	if err != nil {
		logger.Debug("Error while updating a cache setting", zap.Error(err))
		err = utils.LogAndRewindBody(httpResp)
		if err != nil {
			return nil, err
		}

		return nil, utils.ErrorPerStatusCode(httpResp, err)
	}

	return cacheResponse.Results, nil
}

func (c *Client) List(ctx context.Context, opts *contracts.ListOptions, edgeApplicationID int64,
) (*sdk.ApplicationCacheGetResponse, error) {
	logger.Debug("List Cache Settings")
	if opts.OrderBy == "" {
		opts.OrderBy = "id"
	}

	resp, httpResp, err := c.apiClient.EdgeApplicationsCacheSettingsAPI.
		EdgeApplicationsEdgeApplicationIdCacheSettingsGet(ctx, edgeApplicationID).
		OrderBy(opts.OrderBy).
		Page(opts.Page).
		PageSize(opts.PageSize).
		Sort(opts.Sort).Execute()

	if err != nil {
		logger.Debug("Error while listing cache settings", zap.Error(err))
		err = utils.LogAndRewindBody(httpResp)
		if err != nil {
			return nil, err
		}
		return &sdk.ApplicationCacheGetResponse{}, utils.ErrorPerStatusCode(httpResp, err)
	}

	return resp, nil
}
func (c *Client) Get(ctx context.Context, edgeApplicationID, cacheSettingsID int64) (GetResponse, error) {
	logger.Debug("Get Cache Settings")
	resp, httpResp, err := c.apiClient.EdgeApplicationsCacheSettingsAPI.
		EdgeApplicationsEdgeApplicationIdCacheSettingsCacheSettingsIdGet(
			ctx, edgeApplicationID, cacheSettingsID).Execute()
	if err != nil {
		logger.Debug("Error while getting a cache setting", zap.Error(err))
		err = utils.LogAndRewindBody(httpResp)
		if err != nil {
			return nil, err
		}

		return nil, utils.ErrorPerStatusCode(httpResp, err)
	}
	return &resp.Results, nil
}

func (c *Client) Delete(ctx context.Context, edgeApplicationID, cacheSettingsID int64) error {
	logger.Debug("Delete Cache Settings")
	httpResp, err := c.apiClient.EdgeApplicationsCacheSettingsAPI.
		EdgeApplicationsEdgeApplicationIdCacheSettingsCacheSettingsIdDelete(
			ctx, edgeApplicationID, cacheSettingsID).Execute()
	if err != nil {
		logger.Debug("Error while deleting a cache setting", zap.Error(err))
		err = utils.LogAndRewindBody(httpResp)
		if err != nil {
			return err
		}

		return utils.ErrorPerStatusCode(httpResp, err)
	}
	return nil
}
