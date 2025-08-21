package workloads

var (
	Usage                 = "workload"
	ShortDescription      = "Updates a Workload"
	LongDescription       = "Updates a Workload's name and other attributes based on a given ID"
	FlagWorkloadID        = "Unique identifier of the Workload"
	FlagEdgeApplicationId = "The Edge Application's unique identifier"
	FlagName              = "The Workload's name"
	FlagDomains           = "List of domains"
	FlagFile              = "Given path and JSON file to automatically update the Workload attributes; you can use - for reading from stdin"
	OutputSuccess         = "Updated Workload with ID %d"
	FlagActive            = "Whether the Workload should be active or not"
	HelpFlag              = "Displays more information about the 'update workload' command"
	AskInputWorkloadID    = "Enter the Workload's ID:"
)
