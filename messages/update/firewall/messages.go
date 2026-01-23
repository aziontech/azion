package firewall

var (
	Usage                       = "firewall"
	UpdateShortDescription      = "Updates a Firewall"
	UpdateLongDescription       = "Modifies a Firewall's name, activity status, and other attributes based on the given ID"
	FlagID                      = "The Firewall's id"
	UpdateFlagName              = "The Firewall's name"
	UpdateFlagDebug             = "Allows you to check whether rules created using Rules Engine for Firewall have been successfully executed in your firewall"
	UpdateFlagActive            = "Whether the Firewall is active or not"
	UpdateFlagFunctionsEnabled  = "Whether the Firewall has Functions module enabled or not"
	UpdateFlagNetworkProtection = "Whether the Firewall has Network Layer Protection module enabled or not"
	UpdateFlagWafEnabled        = "Whether the Firewall has Web Application Firewall (WAF) module enabled or not"
	UpdateFlagFile              = "Given path and JSON file to automatically update the Firewall attributes; you can use - for reading from stdin"
	UpdateOutputSuccess         = "Updated Firewall with ID %d"
	UpdateHelpFlag              = "Displays more information about the update firewall command"
	UpdateAskFirewallID         = "Enter the ID of the Firewall you wish to update:"
)
