package domains

import "errors"

var (
	ErrorGetRulesEngine     = errors.New("Failed to list your rules engines: %s. Check your settings and try again. If the error persists, contact Azion support.")
	ErrorMandatoryListFlags = errors.New("One or more required flags are missing. You must provide --application-id and --phase flags. Run the command 'azioncli rules_engine list --help' to display more information and try again.")
)
