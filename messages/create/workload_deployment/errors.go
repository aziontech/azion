package workloaddeployment

import "errors"

var (
	ErrorCreateWorkloadDeployment = errors.New("Failed to create the Workload Deployment: %s. Check your settings and try again. If the error persists, contact Azion support.")
	ErrorIsActiveFlag             = errors.New("Invalid --active flag provided. The value must be 'true' or 'false'. Run the command 'azion create workload-deployment --help' to display more information and try again")
	ErrorConvertWorkloadId        = errors.New("The Workload ID you provided is invalid. The value must be an integer. You may run the 'azion list workload' command to check your Workload ID")
	ErrorConvertDeploymentId      = errors.New("The Deployment ID you provided is invalid. The value must be an integer. You may run the 'azion list workload-deployment' command to check your Deployment ID")
	ErrorConvertStrategyType      = errors.New("Invalid --strategy-type flag provided. The value must be 'blue-green' or 'canary'. Run the command 'azion create workload-deployment --help' to display more information and try again")
	ErrorConvertEdgeApplication   = errors.New("Invalid --edge-application flag provided. The value must be an integer. Run the command 'azion create workload-deployment --help' to display more information and try again")
	ErrorConvertEdgeFirewall      = errors.New("Invalid --edge-firewall flag provided. The value must be an integer. Run the command 'azion create workload-deployment --help' to display more information and try again")
	ErrorConvertCustomPage        = errors.New("Invalid --custom-page flag provided. The value must be an integer. Run the command 'azion create workload-deployment --help' to display more information and try again")
)
