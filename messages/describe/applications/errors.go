package applications

import "errors"

// Used by the v3, v4 and v6 command trees.
var (
	ErrorGetApplication = errors.New("Failed to get the Application: %s. Check your settings and try again. If the error persists, contact Azion support")
)

// Used by the v4 and v6 command trees.
var (
	ErrorConvertIdApplication = errors.New("The application ID you provided is invalid. The value must be an integer. You may run the 'azion list application' command to check your application ID")
)
