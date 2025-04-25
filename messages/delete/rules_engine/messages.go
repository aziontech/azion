package rulesengine

var (
	AskInputRulesId       = "Enter the ID of the Rules Engine you wish to delete:"
	AskInputApplicationId = "Enter the ID of the Edge Application this rule is linked to:"
	AskInputPhase         = "Enter the phase of the Rules Engine (request/response):"
	DeleteOutputSuccess   = "Rule Engine %s was successfully deleted"
	FlagRuleID            = "Your Rules Engine's ID"
	FlagAppID             = "Your Edge Application's ID"
	FlagPhase             = "The phase of your Rule Engine (request/response)"
	HelpFlag              = "Displays more information about the delete rules-engine subcommand"

	Usage            = "rules-engine"
	ShortDescription = "Deletes a rule in Rules Engine"
	LongDescription  = "Deletes a rule in Rules Engine based on the given '--rule-id', '--application-id', and '--phase'"
)
