package workloads

var (
	Usage                = "workload"
	ShortDescription     = "Creates a new Workload"
	LongDescription      = "Creates a Workload based on given attributes"
	FlagName             = "The Workload's name"
	FlagAlternateDomains = "List of alternate domains"
	FlagIsActive         = "Whether the Workload is active or not"
	FlagFile             = "Path to a JSON file containing the attributes that will be updated; you can use - for reading from stdin"
	OutputSuccess        = "Created Workload with ID %d"
	HelpFlag             = "Displays more information about the create workload command"
	AskInputName         = "Enter the new Workload's name:"
)
