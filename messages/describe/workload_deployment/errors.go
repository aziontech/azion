package workloaddeployment

import "errors"

var ErrorGetDeployment = errors.New("Failed to describe the Workload Deployment: %s. Check your settings and try again. If the error persists, contact Azion support.")
