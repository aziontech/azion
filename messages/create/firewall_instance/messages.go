package firewallinstance

var (
	Usage              = "firewall-instance"
	ShortDescription   = "Creates a new Firewall Function Instance"
	LongDescription    = "Creates a Firewall Function Instance based on given attributes"
	FlagName           = "The Firewall Function Instance's name"
	FlagIsActive       = "Whether the Firewall Function Instance is active or not"
	FlagFile           = "Path to a JSON file containing the attributes that will be updated; you can use - for reading from stdin"
	OutputSuccess      = "Created Firewall Function Instance with ID %d"
	HelpFlag           = "Displays more information about the create firewall-instance command"
	AskInputName       = "Enter the new Firewall Function Instance's name:"
	AskInputFirewallID = "Enter the Firewall's ID this Function Instance will be associated with:"
	AskInputFunctionID = "Enter the Function's ID:"
	FlagFirewallID     = "Unique identifier of the Firewall"
	FlagFunctionID     = "Unique identifier of the Function"
	FlagArgs           = "The Firewall Function Instance's arguments"
)
