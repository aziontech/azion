package workloaddeployment

import "errors"

// Used only by the v4 command tree.
var (
	ErrorCreateWorkloadDeployment = errors.New("Failed to create the Workload Deployment: %s. Check your settings and try again. If the error persists, contact Azion support.")
	ErrorIsActiveFlag             = errors.New("Invalid --active flag provided. The value must be 'true' or 'false'. Run the command 'azion create workload-deployment --help' to display more information and try again")
	ErrorConvertCustomPage        = errors.New("Invalid --custom-page flag provided. The value must be an integer. Run the command 'azion create workload-deployment --help' to display more information and try again")
)
