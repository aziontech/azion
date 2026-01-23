package firewallinstance

var (
	Usage                  = "firewall-instance"
	ShortDescription       = "Deletes a Firewall Function Instance"
	LongDescription        = "Removes a Firewall Function Instance from the Firewall based on a given ID"
	OutputSuccess          = "Firewall Function Instance %d was successfully deleted"
	HelpFlag               = "Displays more information about the delete firewall-instance subcommand"
	FlagId                 = "Unique identifier of the Firewall Function Instance"
	AskDeleteInput         = "Enter the ID of the Firewall Function Instance you wish to delete:"
	AskDeleteFirewallInput = "Enter the ID of the Firewall:"
)
