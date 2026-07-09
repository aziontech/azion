package waf

var (
	Usage            = "waf"
	ShortDescription = "Returns the WAF data"
	LongDescription  = "Displays information about the Web Application Firewall (WAF) via a given ID to show the WAF's attributes in detail"
	FlagOut          = "Exports the output to the given <file_path/file_name.ext>"
	FlagFormat       = "Changes the output format passing the json value to the flag"
	HelpFlag         = "Displays more information about the describe command"

	FlagId        = "Unique identifier of the WAF"
	FileWritten   = "File successfully written to: %s\n"
	AskInputWafID = "Enter the WAF's ID:"
)
