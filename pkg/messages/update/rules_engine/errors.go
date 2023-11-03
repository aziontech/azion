package rules_engine

import "errors"

var (
	ErrorUpdate               = errors.New("Failed to update the rule in Rules Engine: %s. Check your settings and try again. If the error persists, contact Azion support")
	ErrorNameEmpty            = errors.New("The name field shouldn't be empty")
	ErrorConditionalEmpty     = errors.New("The conditional field shouldn't be empty")
	ErrorVariableEmpty        = errors.New("The variable field shouldn't be empty")
	ErrorOperatorEmpty        = errors.New("The operator field shouldn't be empty")
	ErrorInputValueEmpty      = errors.New("The input value field shouldn't be empty")
	ErrorNameBehaviorsEmpty   = errors.New("The behavior name field cannot be empty")
	ErrorStructCriteriaNil    = errors.New("You must inform a criteria")
	ErrorStructBehaviorsNil   = errors.New("You must inform a behavior")
	ErrorConvertApplicationID = errors.New("The application ID you provided is invalid. The value must be an integer. You may run the 'azion list edge-application' command to check your application ID")
	ErrorConvertRulesID       = errors.New("The rules engine ID you provided is invalid. The value must be an integer. You can run the 'azion list rules-engine' command to check your ID.")
)
