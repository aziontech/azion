package variables

var (
	// [ variables ]
	Usage            = "variables"
	ShortDescription = "Create variables for edges on Azion's platform"
	LongDescription  = "Build your Web applications in minutes without the need to manage infrastructure or security"
	FlagHelp         = "Displays more information about the variables command"

	// [ describe ]
	DescribeUsage            = "describe --variable-id <variable_id> [flags]"
	DescribeShortDescription = "Returns the variable data"
	DescribeLongDescription  = "Displays information about the variable via a given ID to show the application’s attributes in detail"
	DescribeFlagVariableID   = "Unique identifier for an variable. The '--variable-id' flag is mandatory"
	DescribeFlagOut          = "Exports the output to the given <file_path/file_name.ext>"
	DescribeFlagFormat       = "Changes the output format passing the json value to the flag"
	DescribeHelpFlag         = "Displays more information about the describe command"

	// [ general ]
	FileWritten = "File successfully written to: %s\n"
)
