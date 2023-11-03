package edge_applications

import "errors"

var (
	ErrorGetApplication       = errors.New("Failed to get the Edge Application: %s. Check your settings and try again. If the error persists, contact Azion support")
	ErrorConvertIdApplication = errors.New("The application ID you provided is invalid. The value must be an integer. You may run the 'azion list edge-application' command to check your application ID")
)
