package domains

var (
	Usage            = "domains --domain-id <domain_id> [flags]"
	ShortDescription = "Returns the domain data"
	LongDescription  = "Displays information about the domain via a given ID to show the application’s attributes in detail"
	FlagOut          = "Exports the output to the given <file_path/file_name.ext>"
	FlagFormat       = "Changes the output format passing the json value to the flag"
	HelpFlag         = "Displays more information about the describe command"
	FlagDomainID     = "Unique identifier of the Domain"
	FileWritten      = "File successfully written to: %s\n"
	AskInputDomainID = "What is the id of the Domain?"
)
