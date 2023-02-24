package cachesettings

import "errors"

var (
	ErrorMandatoryCreateFlags   = errors.New("Required flags are missing. You must provide the application-id and name flags when --in flag is not provided. Run the command 'azioncli <command> <subcommand> --help' to display more information and try again.")
	ErrorMandatoryCreateFlagsIn = errors.New("Required flags are missing. You must provide the application-id flag when the --in flag is provided. Run the command 'azioncli <command> <subcommand> --help' to display more information and try again.")
	ErrorCachingForOptionsFlag  = errors.New("Invalid --enable-caching-for-options flag provided. The flag must have  'true' or 'false' values. Run the command 'azioncli <command> <subcommand> --help' to display more information and try again.")
	ErrorCachingForPostFlag     = errors.New("Invalid --enable-caching-for-post flag provided. The flag must have  'true' or 'false' values. Run the command 'azioncli <command> <subcommand> --help' to display more information and try again.")
	ErrorCachingStringSortFlag  = errors.New("Invalid --enable-caching-string-sort flag provided. The flag must have  'true' or 'false' values. Run the command 'azioncli <command> <subcommand> --help' to display more information and try again.")
	ErrorSliceConfigurationFlag = errors.New("Invalid --slice-configuration-enable flag provided. The flag must have  'true' or 'false' values. Run the command 'azioncli <command> <subcommand> --help' to display more information and try again.")
	ErrorSliceL2CachingFlag     = errors.New("Invalid --slice-l2-caching-enabled flag provided. The flag must have  'true' or 'false' values. Run the command 'azioncli <command> <subcommand> --help' to display more information and try again.")
	ErrorL2CachingEnabledFlag   = errors.New("Invalid --l2-caching-enabled flag provided. The flag must have  'true' or 'false' values. Run the command 'azioncli <command> <subcommand> --help' to display more information and try again.")

	ErrorCreateCacheSettings               = errors.New("Failed to create the Cache Setting: %s. Check your settings and try again. If the error persists, contact Azion support.")
	ErrorBrowserMaximumTtlNotSent          = errors.New("When browser cache settings is 'override' you must send the --browser-cache-settings-maximum-ttl flag")
	ErrorApplicationAccelerationNotEnabled = errors.New("When --enable-caching-string-sort, --enable-caching-for-post or --enable-caching-for-options is sent, application acceleration must be enabled")
)
