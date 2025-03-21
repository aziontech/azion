package workloads

import "errors"

var (
	ErrorConvertId            = errors.New("The Workload ID you provided is invalid. The value must be an integer. You may run the 'azion list workloads' command to check your Workload ID")
	ErrorFailToDeleteWorkload = errors.New("Failed to delete the Workload: %s. Check your settings and try again. If the error persists, contact Azion support")
)
