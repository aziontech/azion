package firewall

var (
	Usage            = "firewall"
	ShortDescription = "Returns the Firewall data"
	LongDescription  = "Displays information about the Firewall via a given ID to show the firewall’s attributes in detail"
	FlagOut          = "Exports the output to the given <file_path/file_name.ext>"
	FlagFormat       = "Changes the output format passing the json value to the flag"
	HelpFlag         = "Displays more information about the describe command"

	FlagId             = "Unique identifier of the Firewall"
	FileWritten        = "File successfully written to: %s\n"
	AskInputFirewallID = "Enter the Firewall's ID:"
)
