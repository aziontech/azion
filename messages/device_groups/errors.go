package device_groups

import "errors"

var (
	ErrorMandatoryFlags = errors.New("One or more required flags are missing. You must provide the --application-id and --group-id flags. Run the command 'azioncli rules_engine <subcommand> --help' to display more information and try again.")
	ErrorFailToDelete   = errors.New("Failed to delete the rule in Rules Engine: %s. Check your settings and try again. If the error persists, contact Azion support.")
)
