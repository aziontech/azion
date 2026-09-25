package waf

// Used only by the v4 command tree.
var (
	Usage            = "waf"
	ShortDescription = "Returns the WAF data"
	LongDescription  = "Displays information about the Web Application Firewall (WAF) via a given ID to show the WAF's attributes in detail"
	HelpFlag         = "Displays more information about the describe command"

	FlagId        = "Unique identifier of the WAF"
	AskInputWafID = "Enter the WAF's ID:"
)
