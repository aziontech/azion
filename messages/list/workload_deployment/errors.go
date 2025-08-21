package workloaddeployments

import "errors"

var (
	ErrorGetWorkloadDeployments = errors.New("Failed to list your workload deployments. Check your settings and try again. If the error persists, contact Azion support.")
	ErrorConvertId              = errors.New("The Workload ID you provided is invalid. The value must be an integer. You may run the 'azion list workload' command to check your Workload ID")
)
