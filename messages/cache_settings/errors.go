package cache_settings

import (
    "errors"
)

var (
    ErrorGetCaches          = errors.New("Failed to list your cache settings. Check your settings and try again. If the error persists, contact Azion support.")
    ErrorGetCache           = errors.New("Failed to get cache settings. Check your settings and try again. If the error persists, contact Azion support.")
    ErrorMandatoryListFlags = errors.New("Required flags are missing. You must provide application-id, name, addresses and host-header flags when the --application-id flag are not provided. Run the command 'azioncli <command> <subcommand> --help' to display more information and try again.")
)
