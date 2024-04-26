package cachesetting

var (
	// [ general ]
	Usage       = "cache-setting"
	FileWritten = "File successfully written to: %s\n"

	CacheSettingsShortDescription = "Cache Settings allows you to manage existing cache configurations and create new ones"
	CacheSettingsLongDescription  = "Cache Settings allows you to check, remove or update existing cache configurations and create new ones"
	CreateFlagHelp                = "Displays more information about the create cache-setting command"
	CacheSettingsId               = "Unique identifier for a Cache Settings configuration"
	ListAskInputApplicationID     = "Enter the ID of the Edge Application the Cache Setting is linked to:"
	CreateAskInputApplicationID   = "Enter the ID of the Edge Application the Cache Setting will be linked to:"
	UpdateAskInputCacheSettingID  = "Enter the ID of the Cache Setting you wish to update:"
	DeleteAskInputCacheID         = "Enter the ID of the Cache Setting you wish to delete:"
	DescribeAskInputCacheID       = "Enter the ID of the Cache Setting you wish to describe:"
	DescibeAskInputApplicationID  = "Enter the ID of the Edge Application the Cache Settings is linked to:"

	// [ list ]
	ListShortDescription = "Displays your Cache Settings configurations"
	ListLongDescription  = "Displays your Cache Settings configurations to be used with an Edge Application"
	ListHelpFlag         = "Displays more information about the list cache-setting command"

	// [ create ]
	CreateShortDescription = "Creates a new Cache Settings configuration"
	CreateLongDescription  = "Creates a Cache Settings configuration based on given attributes to be used in Edge Applications"
	CreateOutputSuccess    = "Created Cache Settings configuration with ID %d"
	CreateAskInputName     = "Enter the new Cache Setting's name:"

	// [ update ]
	UpdateUsage            = "update [flags]"
	UpdateShortDescription = "Updates a Cache Settings configuration"
	UpdateLongDescription  = "Updates a Cache Settings configuration based on given attributes to be used in Edge Applications"
	UpdateOutputSuccess    = "Updated a Cache Settings configuration with ID %d"
	UpdateFlagHelp         = "Displays more information about the update cache-setting command"

	// [ describe ]
	DescribeShortDescription    = "Returns information about a specific Cache Settings configuration"
	DescribeLongDescription     = "Returns information about a specific Cache Settings configuration, based on a given ID, in details"
	DescribeFlagApplicationID   = "Unique identifier for an Edge Application. The '--application-id' flag is required"
	DescribeFlagCacheSettingsID = "Unique identifier for a Cache Settings configuration. The '--cache-settings-id' flag is required"
	DescribeFlagOut             = "Exports the output to the given <file_path/file_name.ext>"
	DescribeFlagFormat          = "Changes the output format passing the json value to the flag"
	DescribeHelpFlag            = "Displays more information about the describe cache-setting command"

	// [ delete ]
	DeleteShortDescription    = "Deletes a Cache Settings configuration"
	DeleteLongDescription     = "Deletes a Caches Settings configuration from the Edge Applications library based on its given ID"
	DeleteOutputSuccess       = "Caches settings configuration %d was successfully deleted"
	DeleteFlagApplicationID   = "Unique identifier for an Edge Application"
	DeleteFlagCacheSettingsID = "The Cache Settings configuration key unique identifier"
	DeleteHelpFlag            = "Displays more information about the delete cache-setting command"

	// [ flags ]
	FlagEdgeApplicationID          = "Unique identifier for an Edge Application"
	FlagCacheSettingID             = "Unique identifier for an Cache Setting"
	FlagName                       = "The Cache Settings configuration name"
	FlagFile                       = "Path to a JSON file containing the attributes of the Cache Settings configuration that will be created; you can use - for reading from stdin"
	FlagBrowserCacheSettings       = "Configures the amount of time that the content is cached in the web browser"
	FlagQueryStringFields          = "Gives a list of query strings parameters to be considered in the Cache Settings configuration, that will segregate the cache to the same URL"
	FlagCookieNames                = "Distinguishes objects in the Azion cache by name/value of cookies"
	FlagCacheByCookiesEnabled      = "Whether cache by cookies is active or not"
	FlagCacheByQueryString         = "Defines how you want the content to be cached according to variations of Query String in your URLs"
	FlagCdnCacheSettingsEnabled    = "Configures the amount of time Azion's Edge Applications take to cache the content. It can either 'honor' Origin Cache Headers or 'override' Cache Settings"
	FlagCachingForOptionsEnabled   = "Whether caching for options is active or not"
	FlagCachingStringSortEnabled   = "Whether caching string sort is active or not"
	FlagCachingForPostEnabled      = "Whether caching for post is active or not"
	FlagSliceConfigurationEnabled  = "Whether slice configuration is active or not"
	FlagSliceL2CachingEnabled      = "Whether slice L2 caching is active or not"
	FlagL2CachingEnabled           = "Whether L2 caching is active or not"
	FlagSliceConfigurationRange    = "Informs slice configuration range"
	FlagCdnCacheSettingsMaxTtl     = "Informs CDN Cache Settings configuration maximum TTL"
	FlagBrowserCacheSettingsMaxTtl = "Informs Browser Cache Settings configuration maximum TTL"
	FlagAdaptiveDeliveryAction     = "Informs the Cache Settings configuration adaptive delivery action."
)
