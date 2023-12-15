package ruleengine

var (
	//list cmd
	RulesEngineListUsage            = "rules-engine"
	RulesEngineListShortDescription = "Displays the rules related to a specific Edge Application."
	RulesEngineListLongDescription  = "Displays the rules related to a specific Edge Application, informed through the '--application-id' flag"
	RulesEngineListHelpFlag         = "Displays more information about the list rule-engine command"
	ApplicationFlagId               = "Unique identifier for the Edge Application that implements these rules"
	RulesEnginePhase                = "Rules Engine Phase (request/response)"
	AskInputApplicationId           = "Enter the ID of the Edge Application the Rules Engines are linked to:"
	AskInputPhase                   = "Enter the Rules Engines' phase (request/response):"
)
