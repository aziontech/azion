package workloads

import "errors"

var (
	ErrorGetDomain         = errors.New("Failed to describe the Domain: %s. Check your settings and try again. If the error persists, contact Azion support.")
	ErrorConvertWorkloadId = errors.New("The Workload ID you provided is invalid. The value must be an integer. You may run the 'azion list workloads' command to check your Workload ID")
)
