package rules_engine

import "errors"

// Used by both the v3 and the v4 command trees.
var (
	ErrorCreateRulesEngine    = errors.New("Failed to create the rule in Rules Engine: %s. Check your settings and try again. If the error persists, contact Azion support.")
	ErrorConvertApplicationId = errors.New("Invalid --application-id flag provided. The value must be an integer. Run the command 'azion create rules-engine --help' to display more information and try again")
)

// Used only by the v4 command tree.
var (
	ErrorInvalidPhase = errors.New("Invalid phase value provided. The value must be 'request' or 'response'. Run the command 'azion create rules-engine --help' to display more information and try again")
)

// Used only by the v3 command tree (bug-fix-only; see doc/plan.md).
var (
	ErrorNameEmpty          = errors.New("The name field shouldn't be empty")
	ErrorConditionalEmpty   = errors.New("The conditional field shouldn't be empty")
	ErrorVariableEmpty      = errors.New("The variable field shouldn't be empty")
	ErrorOperatorEmpty      = errors.New("The operator field shouldn't be empty")
	ErrorInputValueEmpty    = errors.New("The input value field shouldn't be empty")
	ErrorNameBehaviorsEmpty = errors.New("The behavior name field cannot be empty")
	ErrorStructCriteriaNil  = errors.New("You must inform a criteria")
	ErrorStructBehaviorsNil = errors.New("You must inform a behavior")
)
