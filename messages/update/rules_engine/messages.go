package rules_engine

var (
	Usage            = "rules-engine"
	ShortDescription = "Updates a rule in Rules Engine"
	LongDescription  = "Updates a rule in Rules Engine based on given attributes to be used in Edge Applications"

	OutputSuccess    = "Updated Rules Engine with ID %d\n"
	RulesEnginePhase = "Rules Engine Phase <request|response>. The '--phase' flag is required"

	FlagApplicationID = "Unique identifier for the Edge Application that implements these rules. The '--application-id' flag is required"
	FlagRulesEngineID = "Unique identifier for a rule in Rules Engine. The '--rule-id' flag is required"
	FlagFile          = "Path to a JSON file containing the attributes of the rule that will be updated; you can use - for reading from stdin"
	FlagHelp          = "Displays more information about the Rules Engine command"

	AskInputApplicationID = "What's the ID of the Edge Application that the Rule Engine will be connected to?"
	AskInputRulesID       = "What's the ID of the Rules Engine to which it will be updated?"
	AskInputPhase         = "What's the phase of your Rules Engine? (request/response)"
	AskInputPathFile      = "What's the path of the json to update the Rules Engine?"
)
