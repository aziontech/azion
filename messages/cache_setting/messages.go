package cachesetting

var (
	// [ general ]
	Usage       = "cache-setting"
	FileWritten = "File successfully written to: %s\n"

	CacheSettingsShortDescription = "Cache Settings allows you to manage existing cache configurations and create new ones"
	CacheSettingsLongDescription  = "Cache Settings allows you to check, remove or update existing cache configurations and create new ones"
	CreateFlagHelp                = "Displays more information about the create cache-setting command"
	CacheSettingsId               = "Unique identifier for a Cache Settings configuration"
	ListAskInputApplicationID     = "What is the ID of the edge application the cache settings are linked to?"
	CreateAskInputApplicationID   = "What is the ID of the edge application this cache setting will be linked to?"
	UpdateAskInputCacheSettingID  = "What is the ID of the Cache Setting you wish to update?"
	AskInputCacheID               = "What is the ID of the cache setting you wish to delete?"
	DeleteAskInputCacheID         = "What is the ID of the cache setting you wish to delete?"
	DescribeAskInputCacheID       = "What is the ID of the cache setting you wish to describe?"
	DescibeAskInputApplicationID  = "What is the ID of the edge application the cache settings is linked to?"

	// [ list ]
	ListShortDescription = "Displays your Cache Settings configurations"
	ListLongDescription  = "Displays your Cache Settings configurations on the Azion platform"
	ListHelpFlag         = "Displays more information about the list cache-setting command"

	// [ create ]
	CreateShortDescription = "Creates a new Cache Settings configuration"
	CreateLongDescription  = "Creates a Cache Settings configuration based on given attributes to be used in edge applications"
	CreateOutputSuccess    = "Created Cache Settings configuration with ID %d\n"
	CreateAskInputName     = "What is the Name of the cache setting?"

	// [ update ]
	UpdateUsage            = "update [flags]"
	UpdateShortDescription = "Updates a Cache Settings configuration"
	UpdateLongDescription  = "Updates a Cache Settings configuration based on given attributes to be used in edge applications"
	UpdateOutputSuccess    = "Updated a Cache Settings configuration with ID %d\n"
	UpdateFlagHelp         = "Displays more information about the update cache-setting command"

	// [ describe ]
	DescribeShortDescription    = "Returns information about a specific Cache Settings configuration"
	DescribeLongDescription     = "Returns information about a specific Cache Settings configuration, based on a given ID, in details"
	DescribeFlagApplicationID   = "Unique identifier for an edge application. The '--application-id' flag is required"
	DescribeFlagCacheSettingsID = "Unique identifier for a Cache Settings configuration. The '--cache-settings-id' flag is required"
	DescribeFlagOut             = "Exports the output to the given <file_path/file_name.ext>"
	DescribeFlagFormat          = "Changes the output format passing the json value to the flag"
	DescribeHelpFlag            = "Displays more information about the describe cache-setting command"

	// [ delete ]
	DeleteShortDescription    = "Deletes a Cache Settings configuration"
	DeleteLongDescription     = "Deletes a Caches Settings configuration from the Edge Applications library based on its given ID"
	DeleteOutputSuccess       = "Caches settings configuration %d was successfully deleted\n"
	DeleteFlagApplicationID   = "Unique identifier for an edge application"
	DeleteFlagCacheSettingsID = "The Cache Settings configuration key unique identifier"
	DeleteHelpFlag            = "Displays more information about the delete cache-setting command"

	// [ flags ]
	FlagEdgeApplicationID          = "Unique identifier for an edge application"
	FlagCacheSettingID             = "Unique identifier for an cache setting"
	FlagName                       = "The Cache Settings configuration name"
	FlagIn                         = "Path to a JSON file containing the attributes of the Cache Settings configuration that will be created; you can use - for reading from stdin"
	FlagBrowserCacheSettings       = "Configures the amount of time that the content is cached in the web browser"
	FlagQueryStringFields          = "Gives a list of query strings parameters to be considered in the Cache Settings configuration, that will segregate the cache to the same URL"
	FlagCookieNames                = "Distinguishes objects in the Azion cache by name/value of cookies"
	FlagCacheByCookiesEnabled      = "Whether cache by cookies is active or not"
	FlagCacheByQueryString         = "Defines how you want the content to be cached according to variations of Query String in your URLs"
	FlagCdnCacheSettingsEnabled    = "Configures the amount of time Azion's Edge Applications take to cache the content. It can either Honor Origin Cache Headers or Override Cache Settings"
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
