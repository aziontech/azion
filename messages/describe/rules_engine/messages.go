package rulesengine

// Used by both the v3 and the v4 command trees.
var (
	Usage                 = "rules-engine"
	ShortDescription      = "Returns the information related to the rule in Rules Engine"
	LongDescription       = "Returns the information related to the rule in Rules Engine, informed through the flag '--rule-id' in detail"
	FlagRuleID            = "Your Rule Engine ID"
	FlagAppID             = "Your Application ID"
	FlagPhase             = "The phase of your Rule Engine (request/response)"
	HelpFlag              = "Displays more information about the describe rule-engine subcommand"
	AskInputRulesId       = "Enter the ID of the Rules Engine you wish to describe:"
	AskInputApplicationId = "Enter the ID of the Application this Rules Engine is linked to:"
	AskInputPhase         = "Enter the phase of your Rules Engine (request/response):"
)
