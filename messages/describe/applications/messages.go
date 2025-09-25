package applications

var (
	Usage            = "application"
	ShortDescription = "Returns the Application data"
	LongDescription  = "Displays information about the Application via a given ID to show the application’s attributes in detail"
	FlagOut          = "Exports the output to the given <file_path/file_name.ext>"
	FlagFormat       = "Changes the output format passing the json value to the flag"
	HelpFlag         = "Displays more information about the describe command"

	FlagId                = "Unique identifier of the Application"
	FileWritten           = "File successfully written to: %s\n"
	AskInputApplicationID = "Enter the Application's ID:"
)
