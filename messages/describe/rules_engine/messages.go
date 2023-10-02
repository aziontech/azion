package rulesengine

var (
	Usage            = "rules-engine"
	ShortDescription = "Returns the information related to the rule in Rules Engine"
	LongDescription  = "Returns the information related to the rule in Rules Engine, informed through the flag '--rule-id' in detail"
	FileWritten      = "File successfully written to: %s\n"

	FlagRuleID         = "Your Rule Engine ID"
	FlagAppID          = "Your Edge Application ID"
	FlagPhase          = "The phase of your Rule Engine (request/response)"
	HelpFlag           = "Displays more information about the describe rule-engine subcommand"
	DescribeFlagOut    = "Exports the output of the command to the given file path <file_path/file_name.ext>"
	DescribeFlagFormat = "Changes the output format passing the json value to the flag. Example '--format json'"

	AskInputRulesId       = "What is the id of the Rule Engine you wish to describe?"
	AskInputApplicationId = "What is the id of the Edge Application this Rule Engine is linked to?"
	AskInputPhase         = "What is the phase of your rule engine? (request/response)"
)
