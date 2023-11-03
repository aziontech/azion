package cachesetting

var (
	// [ cache_settings ]
	CacheSettingsUsage            = "cache_settings"
	CacheSettingsShortDescription = "Cache Settings allows you to manage existing cache configurations and create new ones"
	CacheSettingsLongDescription  = "Cache Settings allows you to check, remove or update existing cache configurations and create new ones"
	CacheSettingsFlagHelp         = "Displays more information about the cache_settings command"
	CacheSettingsId               = "Unique identifier for a Cache Settings configuration"

	// [ list ]
	CacheSettingsListUsage            = "list [flags]"
	CacheSettingsListShortDescription = "Displays your Cache Settings configurations"
	CacheSettingsListLongDescription  = "Displays your Cache Settings configurations on the Azion platform"
	CacheSettingsListHelpFlag         = "Displays more information about the list subcommand"

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
	CacheSettingsDescribeFlagCacheSettingsID = "Unique identifier for a Cache Settings configuration. The '--cache-settings-id' flag is required"
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
