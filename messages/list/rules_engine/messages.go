package ruleengine

var (
	//list cmd
	RulesEngineListUsage            = "rules-engine"
	RulesEngineListShortDescription = "Displays the rules related to a specific Edge Application."
	RulesEngineListLongDescription  = "Displays the rules related to a specific Edge Application, informed through the '--application-id' flag"
	RulesEngineListHelpFlag         = "Displays more information about the list rule-engine command"
	ApplicationFlagId               = "Unique identifier for the Edge Application that implements these rules"
	RulesEnginePhase                = "Rules Engine Phase (request/response)"
	AskInputApplicationId           = "What is the id of the Edge Application the rule engines are linked to?"
	AskInputPhase                   = "What is the phase of your rule engine? (request/response)"
)
