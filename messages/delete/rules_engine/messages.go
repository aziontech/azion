package rulesengine

var (
	AskInputRulesId       = "What's the id of the Rules Engine you wish to delete?"
	AskInputApplicationId = "What's the id of the Edge Application this rule is linked to?"
	AskInputPhase         = "What's the phase of your Rules Engine? (request/response)"
	DeleteOutputSuccess   = "Rule Engine %d was successfully deleted\n"
	FlagRuleID            = "Your Rules Engine's ID"
	FlagAppID             = "Your Edge Application's ID"
	FlagPhase             = "The phase of your Rule Engine (request/response)"
	HelpFlag              = "Displays more information about the delete rules-engine subcommand"

	Usage            = "rules-engine"
	ShortDescription = "Deletes a rule in Rules Engine"
	LongDescription  = "Deletes a rule in Rules Engine based on the given '--rule-id', '--application-id', and '--phase'"
)
