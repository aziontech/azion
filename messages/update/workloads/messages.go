package workloads

var (
	Usage                 = "workload"
	ShortDescription      = "Updates a Workload"
	LongDescription       = "Updates a Workload's name and other attributes based on a given ID"
	FlagWorkloadID        = "Unique identifier of the Workload"
	FlagEdgeApplicationId = "The Edge Application's unique identifier"
	FlagName              = "The Domain's name"
	FlagCnames            = "CNAMEs of your Domain"
	FlagCnameAccessOnly   = "Whether the Domain should be Accessed only through CNAMEs or not"
	FlagFile              = "Given path and JSON file to automatically update the Domain attributes; you can use - for reading from stdin"
	OutputSuccess         = "Updated Domain with ID %d"
	FlagAlternateDomains  = "List of alternate domains"
	FlagActive            = "Whether the Domain should be active or not"
	HelpFlag              = "Displays more information about the 'update workload' command"
	AskInputWorkloadID    = "Enter the Workload's ID:"
)
