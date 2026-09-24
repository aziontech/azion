package wafexceptions

// Used only by the v4 command tree.
var (
	Usage               = "waf-exceptions"
	ShortDescription    = "Returns the WAF Exception data"
	LongDescription     = "Displays information in detail about the WAF Exception via a given ID"
	HelpFlag            = "Displays more information about the 'describe waf-exceptions' command"
	FlagExceptionID     = "Unique identifier of the WAF Exception"
	FlagWafID           = "Unique identifier of the WAF"
	AskInputExceptionID = "Enter the WAF Exception's ID:"
	AskInputWafID       = "Enter the WAF's ID:"
)
