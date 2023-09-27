package rulesengine

import "errors"

var (
	ErrorConvertIdRule        = errors.New("The rule engine ID you provided is invalid. The value must be an integer. You may run the 'azion list rule-engine' command to check your rule engine ID")
	ErrorConvertIdApplication = errors.New("The application ID you provided is invalid. The value must be an integer. You may run the 'azion list edge-application' command to check your application ID")
	ErrorFailToDelete         = errors.New("Failed to delete the rule in Rules Engine: %s. Check your settings and try again. If the error persists, contact Azion support.")
)
