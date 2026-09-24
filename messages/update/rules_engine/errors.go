package rules_engine

import "errors"

// Used by both the v3 and the v4 command trees.
var (
	ErrorUpdate               = errors.New("Failed to update the rule in Rules Engine: %s. Check your settings and try again. If the error persists, contact Azion support")
	ErrorConvertApplicationID = errors.New("The application ID you provided is invalid. The value must be an integer. You may run the 'azion list application' command to check your application ID")
	ErrorConvertRulesID       = errors.New("The Rules Engine ID you provided is invalid. The value must be an integer. You can run the 'azion list rules-engine' command to check your ID.")
)

// Used only by the v4 command tree.
var (
	ErrorInvalidPhase = errors.New("Invalid phase. Accepted values are 'request' or 'response'")
)

// Used only by the v3 command tree (bug-fix-only; see doc/plan.md).
var (
	ErrorConditionalEmpty   = errors.New("The conditional field shouldn't be empty")
	ErrorVariableEmpty      = errors.New("The variable field shouldn't be empty")
	ErrorOperatorEmpty      = errors.New("The operator field shouldn't be empty")
	ErrorInputValueEmpty    = errors.New("The input value field shouldn't be empty")
	ErrorNameBehaviorsEmpty = errors.New("The behavior name field cannot be empty")
)
