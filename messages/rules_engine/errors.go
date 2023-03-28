package rules_engine

import "errors"

var (
	ErrorGetRulesEngines    = errors.New("Failed to list your rules engines: %s. Check your settings and try again. If the error persists, contact Azion support.")
	ErrorMandatoryListFlags = errors.New("One or more required flags are missing. You must provide --application-id and --phase flags. Run the command 'azioncli rules_engine list --help' to display more information and try again.")

	ErrorMandatoryFlags = errors.New("One or more required flags are missing. You must provide --application-id, --rules-id and --phase flags. Run the command 'azioncli rules_engine <subcommand> --help' to display more information and try again.")
	ErrorGetRulesEngine = errors.New("Failed to describe the rules engine: %s. Check your settings and try again. If the error persists, contact Azion support.")

	ErrorMissingArgumentsDelete = errors.New("Required flags are missing. You must supply application-id and phase and rule-id as arguments. Run 'azioncli <command> <subcommand> --help' command to display more information and try again")
	ErrorFailToDelete           = errors.New("Failed to delete the Rules Engine: %s. Check your settings and try again. If the error persists, contact Azion support.")

	ErrorMandatoryCreateFlags = errors.New("Required flags are missing. You must provide application-id and phase flags when the --application-id and --in flag are not provided. Run the command 'azioncli <command> <subcommand> --help' to display more information and try again.")

	ErrorCreateRulesEngine = errors.New("Failed to create the Rules Engine: %s. Check your settings and try again. If the error persists, contact Azion support.")
)
