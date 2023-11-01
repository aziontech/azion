package domain

import "errors"

var (
	ErrorCreate               = errors.New("Failed to create the domain: %s. Check your settings and try again. If the error persists, contact Azion support")
	ErrorMissingCnames        = errors.New("Missing CNAMES. When the flag '--cname-access-only' is set as 'true', at least one CNAME must be provided through the flag '--cnames'. Add one or more CNAMES, or set '--cname-access-only' as false and try again.")
	ErrorConvertApplicationID = errors.New("The application ID you provided is invalid. The value must be an integer. You may run the 'azion list edge-application' command to check your application ID")
	ErrorIsActiveFlag         = errors.New("Invalid --active flag provided. The value must be 'true' or 'false'. Run the command 'azion create domains --help' to display more information and try again")
	ErrorCnameAccessOnlyFlag  = errors.New("Invalid --cname-access-only flag provided. The value must be 'true' or 'false'. Run the command 'azion create domains --help' to display more information and try again")
)
