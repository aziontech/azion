package cache_settings

var (
    // [ cache_settings ]
    CacheSettingsUsage            = "cache_settings"
    CacheSettingsShortDescription = "Cache Settings allows you to check, remove or update existing cache configurations and create new ones"
    CacheSettingsLongDescription  = "Cache Settings allows you to check, remove or update existing cache configurations and create new ones"
    CacheSettingsFlagHelp         = "Displays more information about the cache_settings command"
    CacheSettingsId               = "Unique identifier for a Cache Settings configuration"

    // [ list ]
    CacheSettingsListUsage            = "list [flags]"
    CacheSettingsListShortDescription = "Displays your Cache Settings configurations"
    CacheSettingsListLongDescription  = "Displays all Cache Settings configurations"
    CacheSettingsListHelpFlag         = "Displays more information about the list subcommand"

    // [ create ]
    CacheSettingsCreateUsage                          = "create [flags]"
    CacheSettingsCreateShortDescription               = "Creates a new Cache Settings configuration"
    CacheSettingsCreateLongDescription                = "Creates a Cache Settings configuration based on given attributes to be used in edge applications"
    CacheSettingsCreateFlagEdgeApplicationId          = "Unique identifier for an edge application"
    CacheSettingsCreateFlagName                       = "The Cache Settings configuration name"
    CacheSettingsCreateFlagIn                         = "Path to a JSON file containing the attributes of the Cache Settings configuration that will be created; you can use - for reading from stdin"
    CacheSettingsCreateOutputSuccess                  = "Created Cache Settings configuration with ID %d\n"
    CacheSettingsCreateHelpFlag                       = "Displays more information about the create subcommand"
    CacheSettingsCreateFlagBrowserCacheSettings       = "Configures the amount of time that the content is cached in the web browser"
    CacheSettingsCreateFlagQueryStringFields          = "Gives a list of query strings parameters to be considered in the Cache Settings configuration, that will segregate the cache to the same URL"
    CacheSettingsCreateFlagCookieNames                = "Distinguishes objects in the Azion cache by name/value of cookies"
    CacheSettingsCreateFlagCacheByCookies             = "Whether cache by cookies is active or not"
    CacheSettingsCreateFlagCacheByQueryString         = "Defines how you want the content to be cached according to variations of Query String in your URLs" 
    CacheSettingsCreateFlagCdnCacheSettings           = "Configures the amount of time Azion's Edge Applications take to cache the content. It can either Honor Origin Cache Headers or Override Cache Settings" 
    CacheSettingsCreateFlagCachingForOptions          = "Whether caching for options is active or not"
    CacheSettingsCreateFlagCachingStringSort          = "Whether caching string sort is active or not"
    CacheSettingsCreateFlagCachingForPost             = "Whether caching for post is active or not"
    CacheSettingsCreateFlagSliceConfigurationEnabled  = "Whether slice configuration is active or not"
    CacheSettingsCreateFlagSliceL2CachingEnabled      = "Whether slice L2 caching is active or not"
    CacheSettingsCreateFlagSliceEdgeCachingEnabled    = "Whether slice edge caching is active or not"
    CacheSettingsCreateFlagL2CachingEnabled           = "Whether L2 caching is active or not"
    CacheSettingsCreateFlagSliceConfigurationRange    = "Informs slice configuration range"
    CacheSettingsCreateFlagCdnCacheSettingsMaxTtl     = "Informs CDN Cache Settings configuration maximum TTL"
    CacheSettingsCreateFlagBrowserCacheSettingsMaxTtl = "Informs Browser Cache Settings configuration maximum TTL" 
    CacheSettingsCreateFlagAdaptiveDeliveryAction     = "Informs the Cache Settings configuration adaptive delivery action." 

    // [ update ]
    CacheSettingsUpdateUsage            = "update [flags]"
    CacheSettingsUpdateShortDescription = "Updates a Cache Settings configuration"
    CacheSettingsUpdateLongDescription  = "Updates a Cache Settings configuration based on given attributes to be used in edge applications"
    CacheSettingsUpdateOutputSuccess    = "Updated a Cache Settings configuration with ID %d\n"

    // [ describe ]
    CacheSettingsDescribeUsage               = "describe --application-id <application_id> --cache-settings-id <cache-settings-id> [flags]"
    CacheSettingsDescribeShortDescription    = "Returns information about a specific Cache Settings configuration"
    CacheSettingsDescribeLongDescription     = "Returns information about a specific Cache Settings configuration, based on a given ID, in details"
    CacheSettingsDescribeFlagApplicationID   = "Unique identifier for an edge application. The '--application-id' flag is required"
    CacheSettingsDescribeFlagCacheSettingsID = "Unique identifier for an origin. The '--cache-settings-id' flag is required"
    CacheSettingsDescribeFlagOut             = "Exports the output to the given <file_path/file_name.ext>"
    CacheSettingsDescribeFlagFormat          = "Changes the output format passing the json value to the flag"
    CacheSettingsDescribeHelpFlag            = "Displays more information about the describe subcommand"

    // [ delete ]
    CacheSettingsDeleteUsage               = "delete [flags]"
    CacheSettingsDeleteShortDescription    = "Deletes a Cache Settings configuration"
    CacheSettingsDeleteLongDescription     = "Deletes a Caches Settings configuration from the Edge Applications library based on its given ID"
    CacheSettingsDeleteOutputSuccess       = "Caches settings configuration %d was successfully deleted\n"
    CacheSettingsDeleteFlagApplicationID   = "Unique identifier for an edge application"
    CacheSettingsDeleteFlagCacheSettingsID = "The Cache Settings configuration key unique identifier"
    CacheSettingsDeleteHelpFlag            = "Displays more information about the delete subcommand"

    CacheSettingsFileWritten = "File successfully written to: %s\n"
)
