package cachesetting

var (
	// [ cache_settings ]

	Usage = "cache-setting"

	CacheSettingsShortDescription = "Cache Settings allows you to manage existing cache configurations and create new ones"
	CacheSettingsLongDescription  = "Cache Settings allows you to check, remove or update existing cache configurations and create new ones"
	CacheSettingsFlagHelp         = "Displays more information about the cache_settings command"
	CacheSettingsId               = "Unique identifier for a Cache Settings configuration"
	ListAskInputApplicationID     = "What is the ID of the edge application the cache settings are linked to?"
	CreateAskInputApplicationID   = "What is the ID of the edge application this cache setting will be linked to?"
	DeleteAskInputCacheID         = "What is the ID of the cache setting you wish to delete?"
	DescribeAskInputCacheID       = "What is the ID of the cache setting you wish to describe?"
	DescibeAskInputApplicationID  = "What is the ID of the edge application the cache settings is linked to?"

	// [ list ]
	ListShortDescription = "Displays your Cache Settings configurations"
	ListLongDescription  = "Displays your Cache Settings configurations on the Azion platform"
	ListHelpFlag         = "Displays more information about the list subcommand"

	// [ create ]
	CreateShortDescription               = "Creates a new Cache Settings configuration"
	CreateLongDescription                = "Creates a Cache Settings configuration based on given attributes to be used in edge applications"
	CreateFlagEdgeApplicationId          = "Unique identifier for an edge application"
	CreateFlagName                       = "The Cache Settings configuration name"
	CreateFlagIn                         = "Path to a JSON file containing the attributes of the Cache Settings configuration that will be created; you can use - for reading from stdin"
	CreateOutputSuccess                  = "Created Cache Settings configuration with ID %d\n"
	CreateHelpFlag                       = "Displays more information about the create subcommand"
	CreateFlagBrowserCacheSettings       = "Configures the amount of time that the content is cached in the web browser"
	CreateFlagQueryStringFields          = "Gives a list of query strings parameters to be considered in the Cache Settings configuration, that will segregate the cache to the same URL"
	CreateFlagCookieNames                = "Distinguishes objects in the Azion cache by name/value of cookies"
	CreateFlagCacheByCookies             = "Whether cache by cookies is active or not"
	CreateFlagCacheByQueryString         = "Defines how you want the content to be cached according to variations of Query String in your URLs"
	CreateFlagCdnCacheSettingsEnabled    = "Configures the amount of time Azion's Edge Applications take to cache the content. It can either Honor Origin Cache Headers or Override Cache Settings"
	CreateFlagCachingForOptionsEnabled   = "Whether caching for options is active or not"
	CreateFlagCachingStringSortEnabled   = "Whether caching string sort is active or not"
	CreateFlagCachingForPostEnabled      = "Whether caching for post is active or not"
	CreateFlagSliceConfigurationEnabled  = "Whether slice configuration is active or not"
	CreateFlagSliceL2CachingEnabled      = "Whether slice L2 caching is active or not"
	CreateFlagSliceEdgeCachingEnabled    = "Whether slice edge caching is active or not"
	CreateFlagL2CachingEnabled           = "Whether L2 caching is active or not"
	CreateFlagSliceConfigurationRange    = "Informs slice configuration range"
	CreateFlagCdnCacheSettingsMaxTtl     = "Informs CDN Cache Settings configuration maximum TTL"
	CreateFlagBrowserCacheSettingsMaxTtl = "Informs Browser Cache Settings configuration maximum TTL"
	CreateFlagAdaptiveDeliveryAction     = "Informs the Cache Settings configuration adaptive delivery action."
	CreateAskInputName                   = "What is the Name of the cache setting?"

	// [ update ]
	UpdateUsage            = "update [flags]"
	UpdateShortDescription = "Updates a Cache Settings configuration"
	UpdateLongDescription  = "Updates a Cache Settings configuration based on given attributes to be used in edge applications"
	UpdateOutputSuccess    = "Updated a Cache Settings configuration with ID %d\n"

	// [ describe ]
	DescribeShortDescription    = "Returns information about a specific Cache Settings configuration"
	DescribeLongDescription     = "Returns information about a specific Cache Settings configuration, based on a given ID, in details"
	DescribeFlagApplicationID   = "Unique identifier for an edge application. The '--application-id' flag is required"
	DescribeFlagCacheSettingsID = "Unique identifier for a Cache Settings configuration. The '--cache-settings-id' flag is required"
	DescribeFlagOut             = "Exports the output to the given <file_path/file_name.ext>"
	DescribeFlagFormat          = "Changes the output format passing the json value to the flag"
	DescribeHelpFlag            = "Displays more information about the describe subcommand"

	// [ delete ]
	DeleteShortDescription    = "Deletes a Cache Settings configuration"
	DeleteLongDescription     = "Deletes a Caches Settings configuration from the Edge Applications library based on its given ID"
	DeleteOutputSuccess       = "Caches settings configuration %d was successfully deleted\n"
	DeleteFlagApplicationID   = "Unique identifier for an edge application"
	DeleteFlagCacheSettingsID = "The Cache Settings configuration key unique identifier"
	DeleteHelpFlag            = "Displays more information about the delete subcommand"

	FileWritten = "File successfully written to: %s\n"
)
