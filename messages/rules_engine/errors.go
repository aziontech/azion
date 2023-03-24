package domains

import "errors"

var (
	ErrorGetRulesEngine         = errors.New("Failed to list your rules engines: %s. Check your settings and try again. If the error persists, contact Azion support.")
	ErrorMandatoryListFlags     = errors.New("One or more required flags are missing. You must provide --application-id and --phase flags. Run the command 'azioncli rules_engine list --help' to display more information and try again.")
	ErrorMissingArgumentsDelete = errors.New("Required flags are missing. You must supply application-id and phase and rule-id as arguments. Run 'azioncli <command> <subcommand> --help' command to display more information and try again")
	ErrorFailToDelete           = errors.New("Failed to delete the Rules Engine: %s. Check your settings and try again. If the error persists, contact Azion support.")
)
