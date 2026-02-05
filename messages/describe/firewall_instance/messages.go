package firewallinstance

var (
	Usage                              = "firewall-instance"
	ShortDescription                   = "Returns the Firewall Function Instance data"
	LongDescription                    = "Displays information in detail about the Firewall Function Instance via a given ID"
	FlagOut                            = "Exports the output to the given <file_path/file_name.ext>"
	FlagFormat                         = "Changes the output format passing the json value to the flag"
	HelpFlag                           = "Displays more information about the 'describe firewall-instance' command"
	FlagFirewallFunctionInstanceID     = "Unique identifier of the Firewall Function Instance"
	FlagFirewallID                     = "Unique identifier of the Firewall"
	FileWritten                        = "File successfully written to: %s"
	AskInputFirewallFunctionInstanceID = "Enter the Firewall Function Instance's ID:"
	AskInputFirewallID                 = "Enter the Firewall's ID:"
)
