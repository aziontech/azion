package rulesengine

var (
	AskInputRulesId       = "What's the id of the rules engine you wish to delete?"
	AskInputApplicationId = "What's the id of the edge application this rule is linked to?"
	AskInputPhase         = "What's the phase of your rule engine? (request/response)"
	DeleteOutputSuccess   = "Rule Engine %d was successfully deleted\n"
	FlagRuleID            = "Your rules engine's ID"
	FlagAppID             = "Your edge application's ID"
	FlagPhase             = "The phase of your Rule Engine (request/response)"
	HelpFlag              = "Displays more information about the delete rules-engine subcommand"

	Usage            = "rules-engine"
	ShortDescription = "Deletes a rule in Rules Engine"
	LongDescription  = "Deletes a rule in Rules Engine based on the given '--rule-id', '--application-id', and '--phase'"
)
