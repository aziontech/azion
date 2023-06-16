package variables

var (
	// [ variables ]
  
	Usage 				= "variables"
	ShortDescription   = "Manage your variables on the Azion Edge platform"
	LongDescription    = "Manage your variables' varaibles on the Azion Edge platform"
	FlagHelp           = "Displays more information about the Rules Engine command"
	FlagId              = "Unique identifier of the Variable"

	// [ list ]
	VariablesListUsage            = "list [flags]"
	VariablesListShortDescription = "Displays your variables"
	VariablesListLongDescription  = "Displays all variables related to your applications"
	VariablesListHelpFlag         = "Displays more information about the list subcommand"

	
	//delete cmd
	DeleteOutputSuccess            = "Variable %v was successfully deleted\n"
	DeleteHelpFlag                 = "Displays more information about the delete subcommand"
	DeleteUsage               	   = "delete [flags]"
	DeleteShortDescription         = "Delete a Variable"
	DeleteLongDescription		   = "Delete a Variable using UUID"

)
