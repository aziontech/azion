package workloaddeployment

import "errors"

// Used by both the v3 and the v4 command trees.
var (
	ErrorUpdateWorkloadDeployment = errors.New("Failed to update the Workload Deployment")
	ErrorIsActiveFlag             = errors.New("Invalid --active flag provided. The value must be 'true' or 'false'. Run the command 'azion update workload-deployment --help' to display more information and try again")
	ErrorConvertWorkloadId        = errors.New("The Workload ID you provided is invalid. The value must be an integer. You may run the 'azion list workload' command to check your Workload ID")
	ErrorConvertDeploymentId      = errors.New("The Deployment ID you provided is invalid. The value must be an integer. You may run the 'azion list workload-deployment' command to check your Deployment ID")
	ErrorConvertEdgeApplication   = errors.New("Invalid --edge-application flag provided. The value must be an integer. Run the command 'azion update workload-deployment --help' to display more information and try again")
	ErrorConvertEdgeFirewall      = errors.New("Invalid --edge-firewall flag provided. The value must be an integer. Run the command 'azion update workload-deployment --help' to display more information and try again")
	ErrorConvertCustomPage        = errors.New("Invalid --custom-page flag provided. The value must be an integer. Run the command 'azion update workload-deployment --help' to display more information and try again")
	ErrorConvertCurrent           = errors.New("Invalid --current flag provided. The value must be 'true' or 'false'. Run the command 'azion update workload-deployment --help' to display more information and try again")
)
