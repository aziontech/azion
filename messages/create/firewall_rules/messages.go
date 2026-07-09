package firewallrules

var (
	Usage              = "firewall-rule"
	ShortDescription   = "Creates a new Firewall Rule"
	LongDescription    = "Creates a Firewall Rule based on given attributes"
	FlagName           = "The Firewall Rule's name"
	FlagDescription    = "The Firewall Rule's description"
	FlagIsActive       = "Whether the Firewall Rule is active or not"
	FlagFile           = "Path to a JSON file containing the attributes that will be updated; you can use - for reading from stdin"
	OutputSuccess      = "Created Firewall Rule with ID %d"
	HelpFlag           = "Displays more information about the create firewall-rule command"
	AskInputName       = "Enter the new Firewall Rule's name:"
	AskInputPathFile   = "Enter the path to the JSON file:"
	AskInputFirewallID = "Enter the Firewall's ID this Rule will be associated with:"
	FlagFirewallID     = "Unique identifier of the Firewall"
)
