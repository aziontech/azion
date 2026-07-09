package wafexceptions

var (
	Usage            = "waf-exceptions"
	ShortDescription = "Deletes a WAF Exception"
	LongDescription  = "Removes a WAF Exception from a WAF based on a given ID"
	OutputSuccess    = "WAF Exception %d was successfully deleted"
	HelpFlag         = "Displays more information about the delete waf-exceptions subcommand"
	FlagWafID        = "Unique identifier of the WAF"
	FlagExceptionID  = "Unique identifier of the WAF Exception"
	AskDeleteWafID   = "Enter the ID of the WAF:"
	AskDeleteInput   = "Enter the ID of the WAF Exception you wish to delete:"
)
