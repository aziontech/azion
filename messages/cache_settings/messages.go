package cache_settings

var (
	// [ cache_settings ]
	CacheSettingsUsage            = "cache_settings"
	CacheSettingsShortDescription = "Cache Settings lets you check, remove or update existing settings, besides creating new ones"
	CacheSettingsLongDescription  = "Cache Settings lets you check, remove or update existing settings, besides creating new ones"
	CacheSettingsFlagHelp         = "Displays more information about the cache_settings command"
	CacheSettingsId               = "Unique identifier for a cache setting"

	// [ list ]
	CacheSettingsListUsage            = "list [flags]"
	CacheSettingsListShortDescription = "Displays yours cache settings"
	CacheSettingsListLongDescription  = "Displays all cache settings"
	CacheSettingsListHelpFlag         = "Displays more information about the list subcommand"

	// [ create ]
	CacheSettingsCreateUsage                          = "create [flags]"
	CacheSettingsCreateShortDescription               = "Creates a new Cache Setting"
	CacheSettingsCreateLongDescription                = "Creates a Cache Setting based on given attributes to be used in edge applications"
	CacheSettingsCreateFlagEdgeApplicationId          = "Unique identifier for an edge application"
	CacheSettingsCreateFlagName                       = "The Cache Settings' name"
	CacheSettingsCreateFlagIn                         = "Path to a JSON file containing the attributes of the Cache Setting that will be created; you can use - for reading from stdin"
	CacheSettingsCreateOutputSuccess                  = "Created cache setting with ID %d\n"
	CacheSettingsCreateHelpFlag                       = "Displays more information about the create subcommand"
	CacheSettingsCreateFlagBrowserCacheSettings       = "Browser Cache Settings"
	CacheSettingsCreateFlagQueryStringFields          = "Cache Settings' query string fields"
	CacheSettingsCreateFlagCookieNames                = "Cache Settings' cookie names"
	CacheSettingsCreateFlagCacheByCookies             = "Whether cache by cookies is active or not"
	CacheSettingsCreateFlagCacheByQueryString         = "Cache Settings' cache by query string"
	CacheSettingsCreateFlagCdnCacheSettings           = "CDN cache settings"
	CacheSettingsCreateFlagCachingForOptions          = "Whether caching for options is active or not"
	CacheSettingsCreateFlagCachingStringSort          = "Whether caching string sort is active or not"
	CacheSettingsCreateFlagCachingForPost             = "Whether caching for post is active or not"
	CacheSettingsCreateFlagSliceConfigurationEnabled  = "Whether slice configuration is active or not"
	CacheSettingsCreateFlagSliceL2CachingEnabled      = "Whether slice L2 caching is active or not"
	CacheSettingsCreateFlagSliceEdgeCachingEnabled    = "Whether slice edge caching is active or not"
	CacheSettingsCreateFlagL2CachingEnabled           = "Whether slice edge caching is active or not"
	CacheSettingsCreateFlagSliceConfigurationRange    = "Cache Settings' slice configuration range"
	CacheSettingsCreateFlagCdnCacheSettingsMaxTtl     = "CDN cache settings' maximum TTL"
	CacheSettingsCreateFlagBrowserCacheSettingsMaxTtl = "Browser cache settings' maximum TTL"
	CacheSettingsCreateFlagAdaptiveDeliveryAction     = "Cache Settings' adaptive delivery action"

	// [ update ]
	CacheSettingsUpdateUsage            = "update [flags]"
	CacheSettingsUpdateShortDescription = "Modifies a Cache Setting"
	CacheSettingsUpdateLongDescription  = "Modifies a Cache Setting based on given attributes to be used in edge applications"
	CacheSettingsUpdateOutputSuccess    = "Updated cache setting with ID %d\n"
)
