package workloaddeployment

var (
	Usage                 = "workload-deployment"
	ShortDescription      = "Updates a Workload Deployment"
	LongDescription       = "Updates a Workload Deployment's attributes based on a given ID"
	FlagWorkloadID        = "Unique identifier of the Workload"
	FlagDeploymentID      = "Unique identifier of the Workload Deployment"
	FlagName              = "The Workload Deployment's name"
	FlagIsActive          = "Whether the Workload Deployment is active or not"
	FlagIsCurrent         = "Whether the Workload Deployment is current or not"
	FlagStrategyType      = "The type of deployment strategy"
	FlagStrategyAttrs     = "JSON string with strategy attributes"
	FlagEdgeApplicationId = "Unique identifier of the Edge Application"
	FlagEdgeFirewallId    = "Unique identifier of the Edge Firewall"
	FlagCustomPageId      = "Unique identifier of the Custom Page"
	FlagTag               = "The Workload Deployment's Tag"
	FlagCurrent           = "Whether this Deployment is the current one"
	FlagFile              = "Given path and JSON file to automatically update the Workload Deployment attributes; you can use - for reading from stdin"
	OutputSuccess         = "Updated Workload Deployment with ID %d"
	HelpFlag              = "Displays more information about the 'update workload-deployment' command"
	AskInputWorkloadID    = "Enter the Workload's ID:"
	AskInputDeploymentID  = "Enter the Workload Deployment's ID:"
)
