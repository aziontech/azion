package firewallrules

var (
	Usage              = "firewall-rule"
	ShortDescription   = "Updates a Firewall Rule"
	LongDescription    = "Updates a Firewall Rule based on given attributes"
	FlagFile           = "Path to a JSON file containing the attributes that will be updated; you can use - for reading from stdin"
	OutputSuccess      = "Updated Firewall Rule with ID %d"
	HelpFlag           = "Displays more information about the update firewall-rule command"
	AskInputPathFile   = "Enter the path to the JSON file:"
	AskInputFirewallID = "Enter the Firewall's ID:"
	AskInputRuleID     = "Enter the Firewall Rule's ID:"
	FlagFirewallID     = "Unique identifier of the Firewall"
	FlagRuleID         = "Unique identifier of the Firewall Rule"
)
