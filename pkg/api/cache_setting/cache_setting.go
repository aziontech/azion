package cachesetting

import (
	"context"
	"strconv"

	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/utils"
	sdk "github.com/aziontech/azionapi-v4-go-sdk/edge"
	"go.uber.org/zap"
)

func (c *ClientV4) Create(ctx context.Context, req *Request, applicationID int64) (sdk.CacheSetting, error) {
	logger.Debug("Create Cache Settings")

	applicationIDStr := strconv.Itoa(int(applicationID))

	request := c.apiClient.EdgeApplicationsCacheSettingsAPI.
		CreateCacheSetting(ctx, applicationIDStr).
		CacheSettingRequest(req.CacheSettingRequest)
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

	applicationIDStr := strconv.Itoa(int(applicationID))
	cacheSettingIDStr := strconv.Itoa(int(cacheSettingID))

	request := c.apiClient.EdgeApplicationsCacheSettingsAPI.
		PartialUpdateCacheSetting(ctx, applicationIDStr, cacheSettingIDStr).
		PatchedCacheSettingRequest(req.PatchedCacheSettingRequest)
	cacheResponse, httpResp, err := request.Execute()
	if err != nil {
		logger.Debug("Error while updating a Cache Setting", zap.Error(err))
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

	applicationIDStr := strconv.Itoa(int(edgeApplicationID))

	resp, httpResp, err := c.apiClient.EdgeApplicationsCacheSettingsAPI.
		ListCacheSettings(ctx, applicationIDStr).
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

	edgeApplicationIDStr := strconv.Itoa(int(edgeApplicationID))
	cacheSettingIDStr := strconv.Itoa(int(cacheSettingsID))

	resp, httpResp, err := c.apiClient.EdgeApplicationsCacheSettingsAPI.
		RetrieveCacheSetting(ctx, edgeApplicationIDStr, cacheSettingIDStr).
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

	edgeApplicationIDStr := strconv.Itoa(int(edgeApplicationID))
	cacheSettingIDStr := strconv.Itoa(int(cacheSettingsID))

	_, httpResp, err := c.apiClient.EdgeApplicationsCacheSettingsAPI.
		DestroyCacheSetting(ctx, edgeApplicationIDStr, cacheSettingIDStr).
		Execute()
	if err != nil {
		logger.Debug("Error while deleting a Cache Setting", zap.Error(err))

		return httpResp.StatusCode, utils.ErrorPerStatusCode(httpResp, err)
	}

	return 0, nil
}
