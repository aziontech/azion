package wafexceptions

var (
	Usage               = "waf-exceptions"
	ShortDescription    = "Updates a WAF Exception"
	LongDescription     = "Updates a WAF Exception based on given attributes"
	FlagName            = "The WAF Exception's name"
	FlagRuleID          = "The WAF Rule ID to create an exception for"
	FlagPath            = "The path for the WAF Exception"
	FlagConditions      = "The conditions for the WAF Exception (JSON format)"
	FlagOperator        = "The operator for matching conditions: 'regex' or 'contains'"
	FlagActive          = "Whether the WAF Exception is active or not"
	FlagFile            = "Path to a JSON file containing the attributes that will be updated; you can use - for reading from stdin"
	FlagWafID           = "Unique identifier of the WAF"
	FlagExceptionID     = "Unique identifier of the WAF Exception"
	OutputSuccess       = "Updated WAF Exception with ID %d"
	HelpFlag            = "Displays more information about the update waf-exceptions command"
	AskInputName        = "Enter the new WAF Exception's name:"
	AskInputWafID       = "Enter the WAF's ID this Exception is associated with:"
	AskInputExceptionID = "Enter the WAF Exception's ID:"
	AskInputRuleID      = "Enter the Rule ID:"
)
