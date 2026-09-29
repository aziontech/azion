package cachesetting

import sdk "github.com/aziontech/azionapi-go-sdk/edgeapplications"

type CreateRequest struct {
	sdk.ApplicationCacheCreateRequest
}

type UpdateRequest struct {
	sdk.ApplicationCachePatchRequest
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
