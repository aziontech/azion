package edge_functions_instances

var (
	// [ EdgeFunctionsInstances ]
	EdgeFunctionsInstancesUsage            = "edge_functions_instances"
	EdgeFunctionsInstancesShortDescription = "Functions Instances allows you to instantiate serverless functions in your Applications at Azion."
	EdgeFunctionsInstancesLongDescription  = "Functions Instances allows you to instantiate serverless functions in your Applications at Azion, as well as set up conditions for their execution."
	EdgeFunctionsInstancesFlagHelp         = "Displays more information about the Functions instances command"

	// [ list ]
	EdgeFunctionsInstancesListUsage                 = "list [flags]"
	EdgeFunctionsInstancesListShortDescription      = "Displays your Functions instances."
	EdgeFunctionsInstancesListLongDescription       = "Displays all Functions instances related to a specific Application."
	EdgeFunctionsInstancesListHelpFlag              = "Displays more information about the list subcommand"
	EdgeFunctionsInstancesListFlagEdgeApplicationID = "Unique identifier for an Application."

	// [ flags ]
	EdgeApplicationFlagId        = "Unique identifier for an Application"
	EdgeFunctionsInstancesFlagId = "Unique identifier for an Functions instance"

	//Functions Instances cmd
	EdgeFuncInstanceUsage            = "edge_functions_instances"
	EdgeFuncInstanceShortDescription = "Functions Instances allows you to instantiate serverless functions in your Applications at Azion."
	EdgeFuncInstanceLongDescription  = "Functions Instances allows you to instantiate serverless functions in your Applications at Azion, as well as set up conditions for their execution."
	EdgeFuncInstanceFlagHelp         = "Displays more information about the Functions instances command"
	EdgeFuncInstanceFlagId           = "Unique identifier for an Functions instance"
	ApplicationFlagId                = "Unique identifier for the Application related to an Functions instance. The '--application-id' flag is required"

	//delete cmd
	EdgeFuncInstanceDeleteUsage            = "delete --application-id <application_id> --instance-id <instance-id>"
	EdgeFuncInstanceDeleteShortDescription = "Removes an Functions instance"
	EdgeFuncInstanceDeleteLongDescription  = "Removes an Functions instance, instantiated in a specific Application, based on the given flags."
	EdgeFuncInstanceDeleteOutputSuccess    = "functions instance %s was successfully deleted"
	EdgeFuncInstanceDeleteHelpFlag         = "Displays more information about the delete subcommand"

	// [ create ]
	EdgeFuncInstanceCreateUsage                 = "create [flags]"
	EdgeFuncInstanceCreateShortDescription      = "Creates a new Functions instance"
	EdgeFuncInstanceCreateLongDescription       = "Creates a new Functions instance based on given attributes to be used in an Application"
	EdgeFuncInstanceCreateFlagEdgeApplicationId = "Unique identifier for an Application"
	EdgeFuncInstanceCreateFlagEdgeFunctionID    = "Unique identifier for an Functions instance"
	EdgeFuncInstanceCreateFlagName              = "The Functions instance name"
	EdgeFuncInstanceCreateFlagArgs              = "The JSON args related to the Functions instance being created"
	EdgeFuncInstanceCreateFlagIn                = "Path to a JSON file containing the attributes of the Functions instance being created; you can use - for reading from stdin"
	EdgeFuncInstanceCreateOutputSuccess         = "Created Function Instances with ID %d\n"
	EdgeFuncInstanceCreateHelpFlag              = "Displays more information about the create subcommand"

	//describe cmd
	EdgeFuncInstanceDescribeUsage            = "describe --application-id <application_id> --instance-id <instance_id> [flags]"
	EdgeFuncInstanceDescribeShortDescription = "Returns the information related to the Functions instance"
	EdgeFuncInstanceDescribeLongDescription  = "Returns the information related to the Functions instance, informed through the flag '--instance-id' in detail"
	EdgeFuncInstanceDescribeFlagOut          = "Exports the output of the subcommand 'describe' to the given file path <file_path/file_name.ext>"
	EdgeFuncInstanceDescribeFlagFormat       = "Changes the output format passing the json value to the flag. Example '--format json'"
	EdgeFuncInstanceDescribeHelpFlag         = "Displays more information about the describe subcommand"
	EdgeFuncInstanceFileWritten              = "File successfully written to: %s\n"

	// [ Update ]
	EdgeFuncInstanceUpdateUsage                 = "update --application-id <application_id> --instance-id <instance_id> [flags]"
	EdgeFuncInstanceUpdateShortDescription      = "Updates an Functions instance"
	EdgeFuncInstanceUpdateLongDescription       = "Updates an Functions instance, based on given attributes, to be used in Applications"
	EdgeFuncInstanceUpdateFlagEdgeApplicationId = "Unique identifier for an Application"
	EdgeFuncInstanceUpdateFlagIn                = "Path to a JSON file containing the attributes of the Functions instance that will be updated; you can use - for reading from stdin"
	EdgeFuncInstanceUpdateFlagName              = "The Functions instance name"
	EdgeFuncInstanceUpdateFlagArgs              = "The JSON args related to the Functions instance being updated"
	EdgeFuncInstanceUpdateFlagInstanceID        = "Unique identifier for an Functions instance"
	EdgeFuncInstanceUpdateFlagFunctionID        = "Function ID to be used in the Functions instance"
	EdgeFuncInstanceUpdateOutputSuccess         = "Updated Functions instance with ID %d\n"
	EdgeFuncInstanceUpdateHelpFlag              = "Displays more information about the update subcommand"
)
