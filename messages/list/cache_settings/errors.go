package cache_settings

import "errors"

var (
	ErrorConvertIdApplication = errors.New("The application ID you provided is invalid. The value must be an integer. You may run the 'azion list edge-application' command to check your application ID")
	ErrorGetCaches            = errors.New("Failed to list your Cache Settings configurations. Check your settings and try again. If the error persists, contact Azion support.")
)
