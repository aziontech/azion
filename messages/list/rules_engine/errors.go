package ruleengine

import "errors"

// Used by both the v3 and the v4 command trees.
var (
	ErrorGetRulesEngines      = errors.New("Failed to list your rules in Rules Engine: %s. Check your settings and try again. If the error persists, contact Azion support.")
	ErrorConvertIdApplication = errors.New("The application ID you provided is invalid. The value must be an integer. You may run the 'azion list application' command to check your application ID")
)

// Used only by the v4 command tree.
var (
	ErrorInvalidPhase = errors.New("Invalid phase. Accepted values are 'request' or 'response'")
)
