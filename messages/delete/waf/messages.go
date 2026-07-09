package waf

var (
	Usage            = "waf"
	ShortDescription = "Deletes a WAF"
	LongDescription  = "Removes a Web Application Firewall (WAF) from the WAFs library based on a given ID"
	OutputSuccess    = "WAF %d was successfully deleted"
	HelpFlag         = "Displays more information about the delete waf command"
	FlagId           = "Unique identifier of the WAF"
	AskInput         = "Enter the ID of the WAF you wish to delete:"
)
