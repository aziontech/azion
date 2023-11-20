package origins

import "errors"

var (
	ErrorCreateOrigins          = errors.New("Failed to create the Origin: %s. Check your settings and try again. If the error persists, contact Azion support.")
	ErrorHmacAuthenticationFlag = errors.New("Invalid --hmac-authentication flag provided. The flag must have  'true' or 'false' values. Run the command 'azion <command> <subcommand> --help' to display more information and try again.")
	ErrorFailToDelete           = errors.New("Failed to delete the Origin: %s. Check your settings and try again. If the error persists, contact Azion support.")
	ErrorConvertIdApp           = errors.New("The application ID you provided is invalid. The value must be an integer. You may run the 'azion list edge-application' command to check your application ID")
	ErrorGetOrigin              = errors.New("Failed to describe the origin: %s. Check your settings and try again. If the error persists, contact Azion support.")
	ErrorFormatOut              = errors.New("The server failed formatting data for display. Repeat the HTTP request and check the HTTP response's format")
	ErrorWriteFile              = errors.New("The file is read-only and/or isn't accessible. Change the attributes of the file to read and write and/or give access to it")
	ErrorGetOrigins             = errors.New("Failed to list your origins: %s. Check your settings and try again. If the error persists, contact Azion support.")
	ErrorConvertIdApplication   = errors.New("The application ID you provided is invalid. The value must be an integer. You may run the 'azion list edge-application' command to check your application ID")
	ErrorUpdateOrigin           = errors.New("Failed to update the Origin: %s. Check your settings and try again. If the error persists, contact Azion support.")
)
