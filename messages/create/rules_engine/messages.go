package rules_engine

var (
	Usage                 = "rules-engine"
	ShortDescription      = "Creates a rule in Rules Engine"
	LongDescription       = "Creates a rule in Rules Engine based on given attributes to be used in Applications"
	FlagEdgeApplicationID = "Unique identifier for an Application"
	FlagName              = "The rule name"
	FlagPhase             = "The phase is either 'request' or 'response'"
	FlagFile              = "Path to a JSON file containing the attributes of the rule that will be created; you can use - for reading from stdin"
	OutputSuccess         = "Created Rules Engine with ID %d"
	HelpFlag              = "Displays more information about the azion create rules-engine subcommand"
	AskInputApplicationId = "Enter the ID of the Application that the Rules Engine will be connected to:"
	AskInputPhase         = "Enter the new Rule Engine's phase (request/response):"
	AskInputPathFile      = "Enter the path of the json to create the Rules Engine:"
)
