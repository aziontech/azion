package edge_applications

var (
	Usage            = "edge_applications --application-id <application_id> [flags]"
	ShortDescription = "Returns the Edge Application data"
	LongDescription  = "Displays information about the Edge Application via a given ID to show the application’s attributes in detail"
	FlagOut          = "Exports the output to the given <file_path/file_name.ext>"
	FlagFormat       = "Changes the output format passing the json value to the flag"
	HelpFlag         = "Displays more information about the describe command"

	FlagId      = "Unique identifier of the Edge Application"
	FileWritten = "File successfully written to: %s\n"

	DescribeFlagOut    = "Exports the output to the given <file_path/file_name.ext>"
	DescribeFlagFormat = "Changes the output format passing the json value to the flag"
	DescribeHelpFlag   = "Displays more information about the describe command"
)
