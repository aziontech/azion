package workloads

import "errors"

var (
	ErrorCreate               = errors.New("Failed to create the Workload: %s. Check your settings and try again. If the error persists, contact Azion support")
	ErrorConvertApplicationID = errors.New("The application ID you provided is invalid. The value must be an integer. You may run the 'azion list application' command to check your application ID")
	ErrorIsActiveFlag         = errors.New("Invalid --active flag provided. The value must be 'true' or 'false'. Run the command 'azion create workload --help' to display more information and try again")
)
