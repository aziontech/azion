package ruleengine

import "errors"

var (
	ErrorGetRulesEngines      = errors.New("Failed to list your rules in Rules Engine: %s. Check your settings and try again. If the error persists, contact Azion support.")
	ErrorConvertIdApplication = errors.New("The application ID you provided is invalid. The value must be an integer. You may run the 'azion list edge-application' command to check your application ID")
	ErrorInvalidPhase         = errors.New("Invalid phase. Accepted values are 'request' or 'response'")
)
