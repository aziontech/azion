package workloads

import "errors"

// Used only by the v4 command tree.
var (
	ErrorCreate       = errors.New("Failed to create the Workload: %s. Check your settings and try again. If the error persists, contact Azion support")
	ErrorIsActiveFlag = errors.New("Invalid --active flag provided. The value must be 'true' or 'false'. Run the command 'azion create workload --help' to display more information and try again")
)
