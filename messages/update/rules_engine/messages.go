package rules_engine

var (
	Usage            = "rules-engine"
	ShortDescription = "Updates a rule in Rules Engine"
	LongDescription  = "Updates a rule in Rules Engine based on given attributes to be used in Applications"

	OutputSuccess    = "Updated Rules Engine with ID %d"
	RulesEnginePhase = "Rules Engine Phase <request|response>. The '--phase' flag is required"

	FlagApplicationID = "Unique identifier for the Application that implements these rules. The '--application-id' flag is required"
	FlagRulesEngineID = "Unique identifier for a rule in Rules Engine. The '--rule-id' flag is required"
	FlagFile          = "Path to a JSON file containing the attributes of the rule that will be updated; you can use - for reading from stdin"
	FlagHelp          = "Displays more information about the Rules Engine command"

	AskInputApplicationID = "Enter the ID of the Application the Rules Engine will be connected to:"
	AskInputRulesID       = "Enter the ID of the Rules Engine you wish to update:"
	AskInputPhase         = "Enter the phase of your Rules Engine (request/response):"
	AskInputPathFile      = "Enter the path of the json to update the Rules Engine:"
)
