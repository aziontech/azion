package waf

var (
	Usage                  = "waf"
	CreateShortDescription = "Creates a WAF"
	CreateLongDescription  = "Creates a Web Application Firewall (WAF) to protect your applications from threats and attacks"
	FlagIn                 = "Path to a JSON file containing the attributes of the WAF being created; you can use - for reading from stdin"
	CreateFlagHelp         = "Displays more information about the create waf command"
	CreateOutputSuccess    = "Created WAF with ID %d"

	FlagName           = "WAF's name"
	FlagActive         = "Whether the WAF is active or not"
	FlagProductVersion = "WAF's product version"

	FlagEngineVersion = "WAF engine version"
	FlagType          = "WAF engine type"
	FlagRulesets      = "Comma-separated list of ruleset IDs to enable"
	FlagThresholds    = "Comma-separated list of threat=sensitivity pairs"

	AskName   = "Enter the WAF's name:"
	AskActive = "Is the WAF active?"
)
