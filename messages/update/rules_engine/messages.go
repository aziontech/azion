package rules_engine

var (
	Usage            = "rules-engine"
	ShortDescription = "Updates a rule in Rules Engine"
	LongDescription  = "Updates a rule in Rules Engine based on given attributes to be used in edge applications"

	OutputSuccess    = "Updated rule engine with ID %d\n"
	RulesEnginePhase = "Rules Engine Phase <request|response>. The '--phase' flag is required"

	FlagApplicationID = "Unique identifier for the edge application that implements these rules. The '--application-id' flag is required"
	FlagRulesEngineID = "Unique identifier for a rule in Rules Engine. The '--rule-id' flag is required"
	FlagIn            = "Path to a JSON file containing the attributes of the rule that will be updated; you can use - for reading from stdin"
	FlagHelp          = "Displays more information about the Rules Engine command"

	AskInputApplicationID = "What is the ID of the Edge Application that the Rule Engine will be connected to?"
	AskInputRulesID       = "What is the ID of the rules engine to which it will be updated?"
	AskInputPhase         = "What is the phase of your rule engine? (request/response)"
	AskInputPathFile      = "What is the path of the json to update the rules engine?"
)
