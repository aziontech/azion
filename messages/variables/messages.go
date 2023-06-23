package variables

var (
	// [ variables ]
	Usage            = "variables"
	ShortDescription = "Manages your environment variables and secrets on the Azion's platform"
	LongDescription  = "Manages your environment variables and secrets to be used inside edge functions on the Azion's platform"
	FlagHelp         = "Displays more information about the variables command"
	FlagVariableID   = "Unique identifier for a variable. The '--variable-id' flag is mandatory"

	// [ describe ]
	DescribeUsage            = "describe --variable-id <variable_id> [flags]"
	DescribeShortDescription = "Returns the specific variable's key and value"
	DescribeLongDescription  = "Displays information about a variable based on a given UUID to show the variable's attributes in detail"
	DescribeFlagOut          = "Exports the output to the given filepath, such as: <file_path/file_name.ext>"
	DescribeFlagFormat       = "Changes the output format passing the json value to the flag"
	DescribeHelpFlag         = "Displays more information about the describe subcommand"

	// [ general ]
	FileWritten = "File successfully written to: %s\n"

	// [ list ]
	VariablesListUsage            = "list [flags]"
	VariablesListShortDescription = "Displays your variables in a list"
	VariablesListLongDescription  = "Displays all your environment variables and secrets in a list"
	VariablesListHelpFlag         = "Displays more information about the list subcommand"

	// [ delete ]
	DeleteOutputSuccess    = "Variable %v was successfully deleted\n"
	DeleteHelpFlag         = "Displays more information about the delete subcommand"
	DeleteUsage            = "delete [flags]"
	DeleteShortDescription = "Deletes a variable"
	DeleteLongDescription  = "Deletes a variable based on its UUID"

	//update cmd
	UpdateUsage            = "update --variable-id <variable_id> [flags]"
	UpdateShortDescription = "Modifies a variable's attributes"
	UpdateLongDescription  = "Modifies a variable's attributes based on its UUID"
	UpdateFlagKey          = "The variable's key"
	UpdateFlagValue        = "The variable's value"
	UpdateFlagSecret       = "Indicates whether the value is meant to be confidential."
	UpdateFlagIn           = "Given path and JSON file to automatically update the variable attributes; you can use - for reading from stdin"
	UpdateOutputSuccess    = "Updated variable with UUID %d\n"
	UpdateHelpFlag         = "Displays more information about the update subcommand"

	// [ create ]
	CreateUsage            = "create [flags]"
	CreateShortDescription = "Creates a new environment variable or secret on the Azion's platform"
	CreateLongDescription  = "Creates a new environment variable or secret to be used inside edge functions on the Azion's platform"
	CreateFlagKey          = "Informs the variable's key"
	CreateFlagValue        = "Informs the variable's value"
	CreateFlagSecret       = "Indicates whether the value is meant to be confidential."
	CreateFlagIn           = "Path to a JSON file containing the attributes of the variable that will be created; you can use - for reading from stdin"
	CreateOutputSuccess    = "Created variable with UUID %s\n"
	CreateHelpFlag         = "Displays more information about the create subcommand"
)
