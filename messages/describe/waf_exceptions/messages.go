package wafexceptions

var (
	Usage               = "waf-exceptions"
	ShortDescription    = "Returns the WAF Exception data"
	LongDescription     = "Displays information in detail about the WAF Exception via a given ID"
	FlagOut             = "Exports the output to the given <file_path/file_name.ext>"
	FlagFormat          = "Changes the output format passing the json value to the flag"
	HelpFlag            = "Displays more information about the 'describe waf-exceptions' command"
	FlagExceptionID     = "Unique identifier of the WAF Exception"
	FlagWafID           = "Unique identifier of the WAF"
	FileWritten         = "File successfully written to: %s"
	AskInputExceptionID = "Enter the WAF Exception's ID:"
	AskInputWafID       = "Enter the WAF's ID:"
)
