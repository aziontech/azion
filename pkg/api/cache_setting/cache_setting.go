package cachesetting

import (
	"context"

	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/utils"
	sdk "github.com/aziontech/azionapi-v4-go-sdk-dev/edge-api"
	"go.uber.org/zap"
)

func (c *ClientV4) Create(ctx context.Context, req sdk.CacheSettingRequest, applicationID int64) (sdk.CacheSetting, error) {
	logger.Debug("Create Cache Settings")

	request := c.apiClient.ApplicationsCacheSettingsAPI.
		CreateCacheSetting(ctx, applicationID).
		CacheSettingRequest(req)
	cacheResponse, httpResp, err := request.Execute()
	if err != nil {
		logger.Debug("Error while creating a Cache Setting", zap.Error(err))
		errBody, err := utils.LogAndRewindBodyV4(httpResp)
		if err != nil {
			return sdk.CacheSetting{}, err
		}

		return sdk.CacheSetting{}, utils.ErrorPerStatusCodeV4(errBody, httpResp, err)
	}

	return cacheResponse.Data, nil
}

func (c *ClientV4) Update(ctx context.Context, req *RequestUpdate, applicationID, cacheSettingID int64) (ResponseV4, error) {
	logger.Debug("Update Cache Settings")

	request := c.apiClient.ApplicationsCacheSettingsAPI.
		PartialUpdateCacheSetting(ctx, applicationID, cacheSettingID).
		PatchedCacheSettingRequest(req.PatchedCacheSettingRequest)
	cacheResponse, httpResp, err := request.Execute()
	if err != nil {
		logger.Debug("Error while updating a Cache Setting", zap.Any("ID", cacheSettingID), zap.Any("Name", req.PatchedCacheSettingRequest.Name), zap.Error(err))
		errBody, err := utils.LogAndRewindBodyV4(httpResp)
		if err != nil {
			return nil, err
		}

		return nil, utils.ErrorPerStatusCodeV4(errBody, httpResp, err)
	}

	return cacheResponse, nil
}

func (c *ClientV4) List(ctx context.Context, opts *contracts.ListOptions, edgeApplicationID int64,
) (GetResponseV4, error) {

	logger.Debug("List Cache Settings")
	if opts.OrderBy == "" {
		opts.OrderBy = "id"
	}

	resp, httpResp, err := c.apiClient.ApplicationsCacheSettingsAPI.
		ListCacheSettings(ctx, edgeApplicationID).
		Ordering(opts.OrderBy).
		Page(opts.Page).
		PageSize(opts.PageSize).
		Search(opts.Filter).
		Execute()
	if err != nil {
		logger.Debug("Error while listing Cache Settings", zap.Error(err))
		errBody, err := utils.LogAndRewindBodyV4(httpResp)
		if err != nil {
			return nil, err
		}
		return nil, utils.ErrorPerStatusCodeV4(errBody, httpResp, err)
	}

	return resp, nil
}
func (c *ClientV4) Get(ctx context.Context, edgeApplicationID, cacheSettingsID int64) (sdk.CacheSetting, error) {
	logger.Debug("Get Cache Settings")

	resp, httpResp, err := c.apiClient.ApplicationsCacheSettingsAPI.
		RetrieveCacheSetting(ctx, edgeApplicationID, cacheSettingsID).
		Execute()
	if err != nil {
		logger.Debug("Error while getting a Cache Setting", zap.Error(err))
		errBody, err := utils.LogAndRewindBodyV4(httpResp)
		if err != nil {
			return sdk.CacheSetting{}, err
		}

		return sdk.CacheSetting{}, utils.ErrorPerStatusCodeV4(errBody, httpResp, err)
	}
	return resp.Data, nil
}

func (c *ClientV4) Delete(ctx context.Context, edgeApplicationID, cacheSettingsID int64) (int, error) {
	logger.Debug("Delete Cache Settings")

	req := c.apiClient.ApplicationsCacheSettingsAPI.
		DeleteCacheSetting(ctx, edgeApplicationID, cacheSettingsID)
	_, httpResp, err := req.Execute()

	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while deleting a Cache Setting", zap.Error(err))
			errBody, err = utils.LogAndRewindBodyV4(httpResp)
			if err != nil {
				return httpResp.StatusCode, err
			}
		}
		return httpResp.StatusCode, utils.ErrorPerStatusCodeV4(errBody, httpResp, err)
	}

	return 0, nil
}
