package rulesengine

import "errors"

// Used by both the v3 and the v4 command trees.
var (
	ErrorGetRulesEngine = errors.New("Failed to describe the rule in Rules Engine: %s. Check your settings and try again. If the error persists, contact Azion support.")
	ErrorConvertIdRule  = errors.New("The Rules Engine ID you provided is invalid. The value must be an integer. You may run the 'azion list rule-engine' command to check your Rules Engine ID")
)

// Used only by the v4 command tree.
var (
	ErrorConvertIdApplication = errors.New("The application ID you provided is invalid. The value must be an integer. You may run the 'azion list application' command to check your application ID")
	ErrorInvalidPhase         = errors.New("The phase you provided is invalid. The value must be 'request' or 'response'")
)
