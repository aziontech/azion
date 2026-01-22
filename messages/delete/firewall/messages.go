package firewall

var (
	Usage            = "firewall"
	ShortDescription = "Deletes a Firewall"
	LongDescription  = "Removes a Firewall from the Firewalls library based on a given ID"
	OutputSuccess    = "Firewall %d was successfully deleted"
	HelpFlag         = "Displays more information about the delete firewall command"
	FlagId           = "Unique identifier of the Firewall"
	AskInput         = "Enter the ID of the Firewall you wish to delete:"
)
