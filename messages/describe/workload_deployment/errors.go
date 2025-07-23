package workloaddeployment

import "errors"

var (
	ErrorGetDeployment       = errors.New("Failed to describe the Workload Deployment: %s. Check your settings and try again. If the error persists, contact Azion support.")
	ErrorConvertWorkloadId   = errors.New("The Workload ID you provided is invalid. The value must be an integer. You may run the 'azion list workload' command to check your Workload ID")
	ErrorConvertDeploymentId = errors.New("The Deployment ID you provided is invalid. The value must be an integer. You may run the 'azion list workload-deployment' command to check your Deployment ID")
)
