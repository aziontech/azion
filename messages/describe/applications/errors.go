package applications

import "errors"

// Used by both the v3 and the v4 command trees.
var (
	ErrorGetApplication = errors.New("Failed to get the Application: %s. Check your settings and try again. If the error persists, contact Azion support")
)

// Used only by the v4 command tree.
var (
	ErrorConvertIdApplication = errors.New("The application ID you provided is invalid. The value must be an integer. You may run the 'azion list application' command to check your application ID")
)
