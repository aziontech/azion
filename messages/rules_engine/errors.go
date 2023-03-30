package rules_engine

import "errors"

var (
	ErrorGetRulesEngines    = errors.New("Failed to list your rules engines: %s. Check your settings and try again. If the error persists, contact Azion support.")
	ErrorMandatoryListFlags = errors.New("One or more required flags are missing. You must provide --application-id and --phase flags. Run the command 'azioncli rules_engine list --help' to display more information and try again.")

	ErrorMandatoryFlags = errors.New("One or more required flags are missing. You must provide --application-id, --rules-id and --phase flags. Run the command 'azioncli rules_engine <subcommand> --help' to display more information and try again.")
	ErrorGetRulesEngine = errors.New("Failed to describe the rules engine: %s. Check your settings and try again. If the error persists, contact Azion support.")

	ErrorMissingArgumentsDelete = errors.New("Required flags are missing. You must supply application-id and phase and rule-id as arguments. Run 'azioncli <command> <subcommand> --help' command to display more information and try again")
	ErrorFailToDelete           = errors.New("Failed to delete the Rules Engine: %s. Check your settings and try again. If the error persists, contact Azion support.")

	ErrorUpdateRulesengine    = errors.New("Failed to update the Rule Engine: %s. Check your settings and try again. If the error persists, contact Azion support")
	ErrorMandatoryFlagsUpdate = errors.New("One or more required flags are missing. You must provide --application-id, --phase and --in flags. Run the command 'azioncli rules_engine <subcommand> --help' to display more information and try again.")

	ErrorConditionalEmpty   = errors.New("field conditional empty")
	ErrorVariableEmpty      = errors.New("field variable empty")
	ErrorOperatorEmpty      = errors.New("field operator empty")
	ErrorInputValueEmpty    = errors.New("field input value empty")
	ErrorNameBehaviorsEmpty = errors.New("field name from behaviors empty")
	ErrorStructCriteriaNil  = errors.New("all struct criteria nil")
	ErrorStructBehaviorsNil = errors.New("all struct behaviors nil")
)
