package workloads

var (
	Usage                 = "workload"
	ShortDescription      = "Creates a new Workload"
	LongDescription       = "Creates a Workload based on given attributes"
	FlagName              = "The Workload's name"
	FlagEdgeFirewall      = "ID of the Edge Firewall connected to this Workload"
	FlagAlternateDomains  = "List of alternate domains"
	FlagEdgeApplicationId = "The Edge Application's unique identifier"
	FlagIsActive          = "Whether the Workload is active or not"
	FlagFile              = "Path to a JSON file containing the attributes that will be updated; you can use - for reading from stdin"
	OutputSuccess         = "Created Workload with ID %d"
	HelpFlag              = "Displays more information about the create workload command"
	AskInputApplicationID = "Enter the ID of the Edge Application that the Workload will be connected to:"
	AskInputName          = "Enter the new Workload's name:"
)
