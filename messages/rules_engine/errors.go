package rules_engine

import "errors"

var (
	ErrorGetRulesEngines        = errors.New("Failed to list your rules in Rules Engine: %s. Check your settings and try again. If the error persists, contact Azion support.")
	ErrorMandatoryListFlags     = errors.New("One or more required flags are missing. You must provide --application-id and --phase flags. Run the command 'azioncli rules_engine list --help' to display more information and try again.")
	ErrorMandatoryFlags         = errors.New("One or more required flags are missing. You must provide --application-id, --rules-id, and --phase flags. Run the command 'azioncli rules_engine <subcommand> --help' to display more information and try again.")
	ErrorGetRulesEngine         = errors.New("Failed to describe the rule in Rules Engine: %s. Check your settings and try again. If the error persists, contact Azion support.")
	ErrorMissingArgumentsDelete = errors.New("Required flags are missing. You must supply --application-id, --phase, and --rule-id as arguments. Run 'azioncli <command> <subcommand> --help' command to display more information and try again.")
	ErrorFailToDelete           = errors.New("Failed to delete the rule in Rules Engine: %s. Check your settings and try again. If the error persists, contact Azion support.")

	ErrorUpdateRulesengine    = errors.New("Failed to update the rule in Rules Engine: %s. Check your settings and try again. If the error persists, contact Azion support")
	ErrorMandatoryFlagsUpdate = errors.New("One or more required flags are missing. You must provide --application-id, --rule-id, --phase, and --in flags. Run the command 'azioncli rules_engine <subcommand> --help' to display more information and try again.")

	ErrorMandatoryCreateFlags = errors.New("Required flags are missing. You must provide the --application-id and --phase flags when the --application-id and --in flags are not provided. Run the command 'azioncli <command> <subcommand> --help' to display more information and try again.")
	ErrorCreateRulesEngine    = errors.New("Failed to create the rule in Rules Engine: %s. Check your settings and try again. If the error persists, contact Azion support.")
	ErrorNameEmpty            = errors.New("The name field shouldn't be empty")
	ErrorConditionalEmpty     = errors.New("The conditional field shouldn't be empty")
	ErrorVariableEmpty        = errors.New("The variable field shouldn't be empty")
	ErrorOperatorEmpty        = errors.New("The operator field shouldn't be empty")
	ErrorInputValueEmpty      = errors.New("The input value field shouldn't be empty")
	ErrorNameBehaviorsEmpty   = errors.New("The behaviors name field shouldn't be empty")
	ErrorStructCriteriaNil    = errors.New("You must inform a criteria")
	ErrorStructBehaviorsNil   = errors.New("You must inform a behavior")

	ErrorWriteTemplate = errors.New("Failed to create the template file. Verify if you have permission to write to this directory and/or you have access to it")
)
