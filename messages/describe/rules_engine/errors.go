package rulesengine

import "errors"

var (
	ErrorGetRulesEngine       = errors.New("Failed to describe the rule in Rules Engine: %s. Check your settings and try again. If the error persists, contact Azion support.")
	ErrorConvertIdRule        = errors.New("The Rules Engine ID you provided is invalid. The value must be an integer. You may run the 'azion list rule-engine' command to check your Rules Engine ID")
	ErrorConvertIdApplication = errors.New("The application ID you provided is invalid. The value must be an integer. You may run the 'azion list edge-application' command to check your application ID")
)
