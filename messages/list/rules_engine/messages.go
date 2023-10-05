package ruleengine

var (
	//list cmd
	RulesEngineListUsage            = "rules-engine"
	RulesEngineListShortDescription = "Displays the rules related to a specific edge application."
	RulesEngineListLongDescription  = "Displays the rules related to a specific edge application, informed through the '--application-id' flag"
	RulesEngineListHelpFlag         = "Displays more information about the list rule-engine command"
	ApplicationFlagId               = "Unique identifier for the edge application that implements these rules"
	RulesEnginePhase                = "Rules Engine Phase (request/response)"
	AskInputApplicationId           = "What is the id of the Edge Application this Rule Engine is linked to?"
	AskInputPhase                   = "What is the phase of your rule engine? (request/response)"
)
