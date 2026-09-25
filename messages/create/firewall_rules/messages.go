package firewallrules

// Used only by the v4 command tree.
var (
	Usage              = "firewall-rule"
	ShortDescription   = "Creates a new Firewall Rule"
	LongDescription    = "Creates a Firewall Rule based on given attributes"
	FlagFile           = "Path to a JSON file containing the attributes that will be updated; you can use - for reading from stdin"
	OutputSuccess      = "Created Firewall Rule with ID %d"
	HelpFlag           = "Displays more information about the create firewall-rule command"
	AskInputPathFile   = "Enter the path to the JSON file:"
	AskInputFirewallID = "Enter the Firewall's ID this Rule will be associated with:"
	FlagFirewallID     = "Unique identifier of the Firewall"
)
