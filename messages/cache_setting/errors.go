package cachesetting

import (
	"errors"
)

var (
	ErrorGetCaches          = errors.New("Failed to list your Cache Settings configurations. Check your settings and try again. If the error persists, contact Azion support.")
	ErrorGetCache           = errors.New("Failed to get Cache Settings configuration: %s. Check your settings and try again. If the error persists, contact Azion support.")
	ErrorMandatoryListFlags = errors.New("A Required flag is missing. You must provide the application-id flag. Run the command 'azion <command> <subcommand> --help' to display more information and try again.")

	ErrorMandatoryCreateFlags   = errors.New("Required flags are missing. You must provide the application-id and name flags when --in flag is not provided. Run the command 'azion <command> <subcommand> --help' to display more information and try again.")
	ErrorMandatoryCreateFlagsIn = errors.New("A required flag is missing. You must provide the application-id flag when the --in flag is provided. Run the command 'azion <command> <subcommand> --help' to display more information and try again.")
	ErrorCachingForOptionsFlag  = errors.New("Invalid --enable-caching-for-options flag provided. The value must be either 'true' or 'false'. Run the command 'azion <command> <subcommand> --help' to display more information and try again.")
	ErrorCachingForPostFlag     = errors.New("Invalid --enable-caching-for-post flag provided. The value must be either 'true' or 'false'. Run the command 'azion <command> <subcommand> --help' to display more information and try again.")
	ErrorCachingStringSortFlag  = errors.New("Invalid --enable-caching-string-sort flag provided. The value must be either 'true' or 'false'. Run the command 'azion <command> <subcommand> --help' to display more information and try again.")
	ErrorSliceConfigurationFlag = errors.New("Invalid --slice-configuration-enable flag provided. The value must be either 'true' or 'false'. Run the command 'azion <command> <subcommand> --help' to display more information and try again.")
	ErrorSliceL2CachingFlag     = errors.New("Invalid --slice-l2-caching-enabled flag provided. The value must be either 'true' or 'false'. Run the command 'azion <command> <subcommand> --help' to display more information and try again.")
	ErrorL2CachingEnabledFlag   = errors.New("Invalid --l2-caching-enabled flag provided. The value must be either 'true' or 'false'. Run the command 'azion <command> <subcommand> --help' to display more information and try again.")

	ErrorCreateCacheSettings               = errors.New("Failed to create the Cache Settings configuration: %s. Check your settings and try again. If the error persists, contact Azion support.")
	ErrorBrowserMaximumTtlNotSent          = errors.New("When browser Cache Settings is 'override' you must inform the --browser-cache-settings-maximum-ttl flag.")
	ErrorApplicationAccelerationNotEnabled = errors.New("When --enable-caching-string-sort, --enable-caching-for-post or --enable-caching-for-options is sent, application acceleration must be enabled.")

	ErrorMissingArguments = errors.New("Required flags are missing. You must supply application-id and cache-settings-id as arguments. Run 'azion <command> <subcommand> --help' command to display more information and try again.")

	ErrorFailToDelete = errors.New("Failed to delete the Cache Settings configuration: %s. Check your settings and try again. If the error persists, contact Azion support.")

	ErrorMandatoryUpdateFlags   = errors.New("Required flags are missing. You must provide the application-id and cache-settings-id flags when --in flag is not provided. Run the command 'azion <command> <subcommand> --help' to display more information and try again.")
	ErrorMandatoryUpdateInFlags = errors.New("Required flags are missing. You must provide the application-id flag when --in flag is not provided. Run the command 'azion <command> <subcommand> --help' to display more information and try again.")

	ErrorConvertIdApplication = errors.New("The application ID you provided is invalid. The value must be an integer. You may run the 'azion list edge-application' command to check your application ID")
)
