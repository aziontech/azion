package rulesengine

import "errors"

var (
	ErrorConvertIdRule        = errors.New("The Rules Engine ID you provided is invalid. The value must be an integer. You may run the 'azion list rules-engine' command to check your Rules Engine ID")
	ErrorConvertIdApplication = errors.New("The application ID you provided is invalid. The value must be an integer. You may run the 'azion application' command to check your application ID")
	ErrorFailToDelete         = errors.New("Failed to delete the rule in Rules Engine: %s. Check your settings and try again. If the error persists, contact Azion support.")
	ErrorInvalidPhase         = errors.New("Invalid phase. Accepted values are 'request' or 'response'")
)
