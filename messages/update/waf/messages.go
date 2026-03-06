package waf

var (
	Usage                  = "waf"
	UpdateShortDescription = "Updates a WAF"
	UpdateLongDescription  = "Modifies a WAF's name, activity status, engine settings, and other attributes based on the given ID"
	FlagID                 = "The WAF's id"
	UpdateFlagName         = "The WAF's name"
	UpdateFlagActive       = "Whether the WAF is active or not"
	UpdateFlagFile         = "Given path and JSON file to automatically update the WAF attributes; you can use - for reading from stdin"
	UpdateOutputSuccess    = "Updated WAF with ID %d"
	UpdateHelpFlag         = "Displays more information about the update waf command"
	UpdateAskWafID         = "Enter the ID of the WAF you wish to update:"

	UpdateFlagEngineVersion = "WAF engine version"
	UpdateFlagType          = "WAF engine type"
	UpdateFlagRulesets      = "Comma-separated list of ruleset IDs to enable"
	UpdateFlagThresholds    = "Comma-separated list of threat=sensitivity pairs"
)
