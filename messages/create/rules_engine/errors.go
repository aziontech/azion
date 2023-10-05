package rules_engine

import "errors"

var (
	ErrorCreateRulesEngine    = errors.New("Failed to create the rule in Rules Engine: %s. Check your settings and try again. If the error persists, contact Azion support.")
	ErrorNameEmpty            = errors.New("The name field shouldn't be empty")
	ErrorConditionalEmpty     = errors.New("The conditional field shouldn't be empty")
	ErrorVariableEmpty        = errors.New("The variable field shouldn't be empty")
	ErrorOperatorEmpty        = errors.New("The operator field shouldn't be empty")
	ErrorInputValueEmpty      = errors.New("The input value field shouldn't be empty")
	ErrorNameBehaviorsEmpty   = errors.New("The behaviors name field shouldn't be empty")
	ErrorStructCriteriaNil    = errors.New("You must inform a criteria")
	ErrorStructBehaviorsNil   = errors.New("You must inform a behavior")
	ErrorMandatoryCreateFlags = errors.New("Required flags are missing. You must provide the --application-id and --phase flags when the --application-id and --in flags are not provided. Run the command 'azion <command> <subcommand> --help' to display more information and try again.")
	ErrorConvertIdApplication = errors.New("The application ID you provided is invalid. The value must be an integer. You may run the 'azion list edge-application' command to check your application ID")
)
