package variables

import "errors"

var (
	ErrorGetItem          = errors.New("Failed to describe the variable: %s. Check your settings and try again. If the error persists, contact Azion support.")
	ErrorMissingArguments = errors.New("Required flags are missing. You must supply application-id and origin-id as arguments. Run 'azioncli <command> <subcommand> --help' command to display more information and try again")
)
