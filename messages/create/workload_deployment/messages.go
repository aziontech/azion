package workloaddeployment

var (
	Usage                 = "workload-deployment"
	ShortDescription      = "Creates a new Workload Deployment"
	LongDescription       = "Creates a Workload Deployment based on given attributes"
	FlagName              = "The Workload Deployment's name"
	FlagIsActive          = "Whether the Workload Deployment is active or not"
	FlagIsCurrent         = "Whether the Workload Deployment is current or not"
	FlagStrategyType      = "The type of deployment strategy"
	FlagStrategyAttrs     = "JSON string with strategy attributes"
	FlagFile              = "Path to a JSON file containing the attributes that will be updated; you can use - for reading from stdin"
	OutputSuccess         = "Created Workload Deployment with ID %d"
	HelpFlag              = "Displays more information about the create workload-deployment command"
	AskInputName          = "Enter the new Workload Deployment's name:"
	AskInputActive        = "Enter the new Workload Deployment's active:"
	AskInputCurrent       = "Enter the new Workload Deployment's current status:"
	AskInputWorkloadID    = "Enter the Workload's ID:"
	AskInputDeploymentID  = "Enter the Workload Deployment's ID:"
	AskInputApplicationID = "Enter the Application's ID:"
)
