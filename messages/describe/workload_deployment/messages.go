package workloaddeployment

var (
	Usage                = "workload-deployment"
	ShortDescription     = "Returns the Workload Deployment data"
	LongDescription      = "Displays information in detail about the Workload Deployment via a given ID"
	FlagOut              = "Exports the output to the given <file_path/file_name.ext>"
	FlagFormat           = "Changes the output format passing the json value to the flag"
	HelpFlag             = "Displays more information about the 'describe workload-deployment' command"
	FlagDeploymentID     = "Unique identifier of the Workload-Deployment"
	FlagWorkloadID       = "Unique identifier of the Workload"
	FileWritten          = "File successfully written to: %s"
	AskInputDeploymentID = "Enter the Workload Deployment's ID:"
	AskInputWorkloadID   = "Enter the Workload's ID:"
)
