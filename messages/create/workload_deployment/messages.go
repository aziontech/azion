package workloaddeployment

// Used only by the v4 command tree.
var (
	Usage                 = "workload-deployment"
	ShortDescription      = "Creates a new Workload Deployment"
	LongDescription       = "Creates a Workload Deployment based on given attributes"
	FlagName              = "The Workload Deployment's name"
	FlagIsActive          = "Whether the Workload Deployment is active or not"
	FlagIsCurrent         = "Whether the Workload Deployment is current or not"
	FlagStrategyType      = "The type of deployment strategy"
	FlagFile              = "Path to a JSON file containing the attributes that will be updated; you can use - for reading from stdin"
	OutputSuccess         = "Created Workload Deployment with ID %d"
	HelpFlag              = "Displays more information about the create workload-deployment command"
	AskInputName          = "Enter the new Workload Deployment's name:"
	AskInputWorkloadID    = "Enter the Workload's ID:"
	AskInputApplicationID = "Enter the Application's ID:"
)
