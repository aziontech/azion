package rulesengine

var (
	AskInputRulesId       = "What is the id of the Rule Engine you wish to delete?"
	AskInputApplicationId = "What is the id of the Edge Application this Rule Engine is linked to?"
	AskInputPhase         = "What is the phase of your rule engine? (request/response)"
	DeleteOutputSuccess   = "Rule Engine %d was successfully deleted\n"
	FlagRuleID            = "Your Rule Engine ID"
	FlagAppID             = "Your Edge Application ID"
	FlagPhase             = "The phase of your Rule Engine (request/response)"
	HelpFlag              = "Displays more information about the delete rule-engine subcommand"

	Usage            = "rule-engine"
	ShortDescription = "Deletes a rule in Rules Engine"
	LongDescription  = "Deletes a rule in Rules Engine based on the given '--rule-id', '--application-id', and '--phase'"
)
