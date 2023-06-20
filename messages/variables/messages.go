package variables

var (
	// [ variables ]

	Usage            = "variables"
	ShortDescription = "Manage your variables on the Azion Edge platform"
	LongDescription  = "Manage your variables' varaibles on the Azion Edge platform"
	FlagHelp         = "Displays more information about the Rules Engine command"
	FlagId           = "Unique identifier of the Variable"

	// [ list ]
	VariablesListUsage            = "list [flags]"
	VariablesListShortDescription = "Displays your variables"
	VariablesListLongDescription  = "Displays all variables related to your applications"
	VariablesListHelpFlag         = "Displays more information about the list subcommand"

	//delete cmd
	DeleteOutputSuccess    = "Variable %v was successfully deleted\n"
	DeleteHelpFlag         = "Displays more information about the delete subcommand"
	DeleteUsage            = "delete [flags]"
	DeleteShortDescription = "Delete a Variable"
	DeleteLongDescription  = "Delete a Variable using UUID"

	//update cmd
	UpdateUsage            = "update --variable-id <variable_id> [flags]"
	UpdateShortDescription = "Modifies a Variable"
	UpdateLongDescription  = "Modifies a Variable based on its ID to update its fields"
	UpdateFlagKey          = "The Variable's key"
	UpdateFlagValue        = "The value for the key"
	UpdateFlagSecret       = "Whether the key and value should be secret or not"
	UpdateFlagIn           = "Given path and JSON file to automatically update the Edge Function attributes; you can use - for reading from stdin"
	UpdateOutputSuccess    = "Updated Variable with ID %d\n"
	UpdateHelpFlag         = "Displays more information about the update subcommand"
)
