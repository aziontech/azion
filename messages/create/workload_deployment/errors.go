package workloaddeployment

import "errors"

var (
	ErrorCreate                  = errors.New("Failed to create the Workload Deployment: %s. Check your settings and try again. If the error persists, contact Azion support")
	ErrorConvertWorkloadId       = errors.New("The workload ID you provided is invalid. The value must be an integer. You may run the 'azion list workload' command to check your workload ID")
	ErrorConvertEdgeApplication  = errors.New("The edge application ID you provided is invalid. The value must be an integer. You may run the 'azion list edge-application' command to check your edge application ID")
	ErrorConvertEdgeFirewall     = errors.New("The edge firewall ID you provided is invalid. The value must be an integer")
)
