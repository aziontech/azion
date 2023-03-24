package rules_engine

import "errors"

var (
	ErrorGetRulesEngines    = errors.New("Failed to list your rules engines: %s. Check your settings and try again. If the error persists, contact Azion support.")
	ErrorMandatoryListFlags = errors.New("One or more required flags are missing. You must provide --application-id and --phase flags. Run the command 'azioncli rules_engine list --help' to display more information and try again.")

	ErrorMandatoryFlags = errors.New("One or more required flags are missing. You must provide --application-id, --rules-id and --phase flags. Run the command 'azioncli rules_engine <subcommand> --help' to display more information and try again.")
	ErrorGetRulesEngine = errors.New("Failed to describe the rules engine: %s. Check your settings and try again. If the error persists, contact Azion support.")
)
