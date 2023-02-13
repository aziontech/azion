package origins

import "errors"

var (
	ErrorGetOrigins                   = errors.New("Failed to list your origins. Check your settings and try again. If the error persists, contact Azion support.")
	ErrorGetOrigin                    = errors.New("Failed to retrieve the origin's data. Check your settings and try again. If the error persists, contact Azion support.")
	ErrorMissingApplicationIDArgument = errors.New("A required flag is missing. You must supply an 'application-id' as argument. Run 'azioncli origins list --help' command to display more information and try again")
	ErrorMissingArguments             = errors.New("One or more required flags are missing. Run 'azioncli <command> <subcommand> --help' to display more information and try again")
)
