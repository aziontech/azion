package origin

import "errors"

var (
	ErrorFailToDelete = errors.New("Failed to delete the Origin: %s. Check your settings and try again. If the error persists, contact Azion support.")
	ErrorConvertIdApp = errors.New("The application ID you provided is invalid. The value must be an integer. You may run the 'azion list edge-application' command to check your application ID")
)
