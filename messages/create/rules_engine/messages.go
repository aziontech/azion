package rules_engine

var (
	Usage                 = "rules-engine [flags]"
	ShortDescription      = "Creates a new rule in Rules Engine"
	LongDescription       = "Creates a new rule in Rules Engine based on given attributes to be used in edge applications"
	FlagEdgeApplicationID = "Unique identifier for an edge application"
	FlagName              = "The rule name"
	FlagPhase             = "The phase is either 'request' or 'response'"
	FlagIn                = "Path to a JSON file containing the attributes of the rule that will be created; you can use - for reading from stdin"
	OutputSuccess         = "Created Rules Engine with ID %d\n"
	HelpFlag              = "Displays more information about the create subcommand"
	AskInputApplicationId = "What is the ID of the Edge Application that the Rule Engine will be connected to?"
	AskInputPhase         = "What is the phase of your rule engine? (request/response)"
	AskInputPathFile      = "What is the path of the json to create the rules engine?"
)
