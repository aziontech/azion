package cache_settings

var (
    // [ cache_settings ]
    CacheSettingsUsage            = "cache_settings"
    CacheSettingsShortDescription = "Cache Settings lets you check, remove or update existing settings, besides creating new ones"
    CacheSettingsLongDescription  = "Cache Settings lets you check, remove or update existing settings, besides creating new ones"
    CacheSettingsFlagHelp         = "Displays more information about the cache_settings command"

    // [ list ]
    CacheSettingsListUsage            = "list [flags]"
    CacheSettingsListShortDescription = "Displays yours cache settings"
    CacheSettingsListLongDescription  = "Displays all cache settings"
    CacheSettingsListHelpFlag         = "Displays more information about the list subcommand"

    // [ create ]
    CacheSettingsCreateUsage                          = "create [flags]"
    CacheSettingsCreateShortDescription               = "Creates a new Cache Settings"
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

    // [ describe ]
    CacheSettingsDescribeUsage               = "describe --application-id <application_id> --cache-settings-id <cache-settings-id> [flags]"
    CacheSettingsDescribeShortDescription    = "Returns information about a specific cache settings"
    CacheSettingsDescribeLongDescription     = "Returns information about a specific cache settings, based on a given ID, in details"
    CacheSettingsDescribeFlagApplicationID   = "Unique identifier for an edge application. The '--application-id' flag is mandatory"
    CacheSettingsDescribeFlagCacheSettingsID = "Unique identifier for an origin. The '--cache-settings-id' flag is mandatory"
    CacheSettingsDescribeFlagOut             = "Exports the output to the given <file_path/file_name.ext>"
    CacheSettingsDescribeFlagFormat          = "Changes the output format passing the json value to the flag"
    CacheSettingsDescribeHelpFlag            = "Displays more information about the describe subcommand"

    // [ delete ]
    CacheSettingsDeleteUsage               = "delete [flags]"
    CacheSettingsDeleteShortDescription    = "Deletes an Cache Settings"
    CacheSettingsDeleteLongDescription     = "Deletes an Caches Settings from the Edge Applications library based on its given ID"
    CacheSettingsDeleteOutputSuccess       = "Caches settings %d was successfully deleted\n"
    CacheSettingsDeleteFlagApplicationID   = "Unique identifier for an edge application"
    CacheSettingsDeleteFlagCacheSettingsID = "The Cache Settings key unique identifier"
    CacheSettingsDeleteHelpFlag            = "Displays more information about the delete subcommand"

    CacheSettingsFileWritten = "File successfully written to: %s\n"
)
