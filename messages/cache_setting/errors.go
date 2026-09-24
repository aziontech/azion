package cachesetting

import (
	"errors"
)

// Used by both the v3 and the v4 command trees.
var (
	ErrorGetCaches                = errors.New("Failed to list your Cache Settings configurations. Check your settings and try again. If the error persists, contact Azion support.")
	ErrorGetCache                 = errors.New("Failed to get Cache Settings configuration: %s. Check your settings and try again. If the error persists, contact Azion support.")
	ErrorCachingForOptionsFlag    = errors.New("Invalid --enable-caching-for-options flag provided. The value must be either 'true' or 'false'. Run the command 'azion <command> <subcommand> --help' to display more information and try again.")
	ErrorCachingForPostFlag       = errors.New("Invalid --enable-caching-for-post flag provided. The value must be either 'true' or 'false'. Run the command 'azion <command> <subcommand> --help' to display more information and try again.")
	ErrorCachingStringSortFlag    = errors.New("Invalid --enable-caching-string-sort flag provided. The value must be either 'true' or 'false'. Run the command 'azion <command> <subcommand> --help' to display more information and try again.")
	ErrorCreateCacheSettings      = errors.New("Failed to create the Cache Settings configuration: %s. Check your settings and try again. If the error persists, contact Azion support.")
	ErrorBrowserMaximumTtlNotSent = errors.New("When browser Cache Settings is 'override' you must inform the --browser-cache-max-age flag.")
	ErrorFailToDelete             = errors.New("Failed to delete the Cache Settings configuration: %s. Check your settings and try again. If the error persists, contact Azion support.")
	ErrorConvertIdApplication     = errors.New("The application ID you provided is invalid. The value must be an integer. You may run the 'azion list application' command to check your application ID")
)

// Used only by the v4 command tree.
var (
	ErrorTieredCachingFlag   = errors.New("Invalid --tiered-caching-enabled flag provided. The value must be either 'true' or 'false'. Run the command 'azion <command> <subcommand> --help' to display more information and try again.")
	ErrorUpdateCacheSettings = errors.New("Failed to update the Cache Settings configuration: %s. Check your settings and try again. If the error persists, contact Azion support.")
)

// Used only by the v3 command tree (bug-fix-only; see doc/plan.md).
var (
	ErrorSliceConfigurationFlag            = errors.New("Invalid --slice-configuration-enable flag provided. The value must be either 'true' or 'false'. Run the command 'azion <command> <subcommand> --help' to display more information and try again.")
	ErrorSliceL2CachingFlag                = errors.New("Invalid --slice-l2-caching-enabled flag provided. The value must be either 'true' or 'false'. Run the command 'azion <command> <subcommand> --help' to display more information and try again.")
	ErrorL2CachingEnabledFlag              = errors.New("Invalid --l2-caching-enabled flag provided. The value must be either 'true' or 'false'. Run the command 'azion <command> <subcommand> --help' to display more information and try again.")
	ErrorApplicationAccelerationNotEnabled = errors.New("When --enable-caching-string-sort, --enable-caching-for-post or --enable-caching-for-options is sent, application acceleration must be enabled.")
)
