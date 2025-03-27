package workloads

import "errors"

var (
	ErrorUpdateDomain      = errors.New("Failed to update the Workload: %s. Check your settings and try again. If the error persists, contact Azion support.")
	ErrorActiveFlag        = errors.New("Invalid --active flag provided. The flag must have  'true' or 'false' values. Run the command 'azion update workload --help' to display more information and try again.")
	ErrorConvertWorkloadID = errors.New("The workload ID you provided is invalid. The value must be an integer. You may run the 'azion list workload' command to check your workload ID")
)
