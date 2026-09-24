package firewall

// Used only by the v4 command tree.
var (
	Usage            = "firewall"
	ShortDescription = "Returns the Firewall data"
	LongDescription  = "Displays information about the Firewall via a given ID to show the firewall’s attributes in detail"
	HelpFlag         = "Displays more information about the describe command"

	FlagId             = "Unique identifier of the Firewall"
	AskInputFirewallID = "Enter the Firewall's ID:"
)
