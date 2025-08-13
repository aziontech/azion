package cachesetting

import (
	sdkOld "github.com/aziontech/azionapi-go-sdk/edgeapplications"
	sdk "github.com/aziontech/azionapi-v4-go-sdk/edge-api"
)

type Request struct {
	sdk.CacheSettingRequest
}

type RequestUpdate struct {
	sdk.PatchedCacheSettingRequest
}

type ResponseV4 interface {
	GetState() string
	GetData() sdk.CacheSetting
}

type GetResponseV4 interface {
	GetCount() int64
	GetResults() []sdk.ResponseListCacheSetting
}

type UpdateRequest struct {
	sdkOld.ApplicationCachePatchRequest
}

type CreateRequest struct {
	sdkOld.ApplicationCacheCreateRequest
}

type Response interface {
	GetId() int64
	GetName() string
	GetBrowserCacheSettings() string
	GetBrowserCacheSettingsMaximumTtl() int64
	GetCdnCacheSettingsMaximumTtl() int64
	GetCdnCacheSettings() string
	GetCacheByQueryString() string
	GetQueryStringFields() []string
	GetEnableQueryStringSort() bool
	GetCacheByCookies() string
	GetCookieNames() []string
	GetEnableCachingForPost() bool
	GetL2CachingEnabled() bool
	GetAdaptiveDeliveryAction() string
	GetDeviceGroup() []int32
}

type GetResponse interface {
	GetId() int64
	GetName() string
	GetBrowserCacheSettings() string
	GetBrowserCacheSettingsMaximumTtl() int64
	GetCdnCacheSettingsMaximumTtl() int64
	GetCdnCacheSettings() string
	GetCacheByQueryString() string
	GetQueryStringFields() []string
	GetEnableQueryStringSort() bool
	GetCacheByCookies() string
	GetCookieNames() []*string
	GetEnableCachingForPost() bool
	GetL2CachingEnabled() bool
	GetAdaptiveDeliveryAction() string
	GetDeviceGroup() []int32
}
