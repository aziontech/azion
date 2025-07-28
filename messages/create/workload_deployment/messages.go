package workloaddeployment

var (
	Usage                     = "workload-deployment"
	ShortDescription          = "Creates a new Workload Deployment"
	LongDescription           = "Creates a Workload Deployment based on given attributes"
	FlagWorkloadId            = "The Workload's ID"
	FlagEdgeApplication       = "The Edge Application ID to associate with the deployment"
	FlagEdgeFirewall          = "The Edge Firewall ID to associate with the deployment (optional)"
	FlagFile                  = "Path to a JSON file containing the attributes that will be created; you can use - for reading from stdin"
	OutputSuccess             = "Created Workload Deployment with ID %d"
	HelpFlag                  = "Displays more information about the create workload-deployment command"
	AskInputWorkloadId        = "Enter the Workload ID:"
	AskInputEdgeApplication   = "Enter the Edge Application ID:"
)
