package origin

import "errors"

var (
	ErrorGetOrigins           = errors.New("Failed to list your origins: %s. Check your settings and try again. If the error persists, contact Azion support.")
	ErrorConvertIdApplication = errors.New("The application ID you provided is invalid. The value must be an integer. You may run the 'azion list edge-application' command to check your application ID")
)
