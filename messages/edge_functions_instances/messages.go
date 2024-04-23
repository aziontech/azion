package edge_functions_instances

var (
	// [ EdgeFunctionsInstances ]
	EdgeFunctionsInstancesUsage            = "edge_functions_instances"
	EdgeFunctionsInstancesShortDescription = "Edge Functions Instances allows you to instantiate serverless functions in your Edge Applications at Azion."
	EdgeFunctionsInstancesLongDescription  = "Edge Functions Instances allows you to instantiate serverless functions in your Edge Applications at Azion, as well as set up conditions for their execution."
	EdgeFunctionsInstancesFlagHelp         = "Displays more information about the Edge Functions instances command"

	// [ list ]
	EdgeFunctionsInstancesListUsage                 = "list [flags]"
	EdgeFunctionsInstancesListShortDescription      = "Displays your Edge Functions instances."
	EdgeFunctionsInstancesListLongDescription       = "Displays all Edge Functions instances related to a specific Edge Application."
	EdgeFunctionsInstancesListHelpFlag              = "Displays more information about the list subcommand"
	EdgeFunctionsInstancesListFlagEdgeApplicationID = "Unique identifier for an Edge Application."

	// [ flags ]
	EdgeApplicationFlagId        = "Unique identifier for an Edge Application"
	EdgeFunctionsInstancesFlagId = "Unique identifier for an Edge Functions instance"

	//Edge Functions Instances cmd
	EdgeFuncInstanceUsage            = "edge_functions_instances"
	EdgeFuncInstanceShortDescription = "Edge Functions Instances allows you to instantiate serverless functions in your Edge Applications at Azion."
	EdgeFuncInstanceLongDescription  = "Edge Functions Instances allows you to instantiate serverless functions in your Edge Applications at Azion, as well as set up conditions for their execution."
	EdgeFuncInstanceFlagHelp         = "Displays more information about the Edge Functions instances command"
	EdgeFuncInstanceFlagId           = "Unique identifier for an Edge Functions instance"
	ApplicationFlagId                = "Unique identifier for the Edge Application related to an Edge Functions instance. The '--application-id' flag is required"

	//delete cmd
	EdgeFuncInstanceDeleteUsage            = "delete --application-id <application_id> --instance-id <instance-id>"
	EdgeFuncInstanceDeleteShortDescription = "Removes an Edge Functions instance"
	EdgeFuncInstanceDeleteLongDescription  = "Removes an Edge Functions instance, instantiated in a specific Edge Application, based on the given flags."
	EdgeFuncInstanceDeleteOutputSuccess    = "Edge functions instance %s was successfully deleted"
	EdgeFuncInstanceDeleteHelpFlag         = "Displays more information about the delete subcommand"

	// [ create ]
	EdgeFuncInstanceCreateUsage                 = "create [flags]"
	EdgeFuncInstanceCreateShortDescription      = "Creates a new Edge Functions instance"
	EdgeFuncInstanceCreateLongDescription       = "Creates a new Edge Functions instance based on given attributes to be used in an Edge Application"
	EdgeFuncInstanceCreateFlagEdgeApplicationId = "Unique identifier for an Edge Application"
	EdgeFuncInstanceCreateFlagEdgeFunctionID    = "Unique identifier for an Edge Functions instance"
	EdgeFuncInstanceCreateFlagName              = "The Edge Functions instance name"
	EdgeFuncInstanceCreateFlagArgs              = "The JSON args related to the Edge Functions instance being created"
	EdgeFuncInstanceCreateFlagIn                = "Path to a JSON file containing the attributes of the Edge Functions instance being created; you can use - for reading from stdin"
	EdgeFuncInstanceCreateOutputSuccess         = "Created Function Instances with ID %d\n"
	EdgeFuncInstanceCreateHelpFlag              = "Displays more information about the create subcommand"

	//describe cmd
	EdgeFuncInstanceDescribeUsage            = "describe --application-id <application_id> --instance-id <instance_id> [flags]"
	EdgeFuncInstanceDescribeShortDescription = "Returns the information related to the Edge Functions instance"
	EdgeFuncInstanceDescribeLongDescription  = "Returns the information related to the Edge Functions instance, informed through the flag '--instance-id' in detail"
	EdgeFuncInstanceDescribeFlagOut          = "Exports the output of the subcommand 'describe' to the given file path <file_path/file_name.ext>"
	EdgeFuncInstanceDescribeFlagFormat       = "Changes the output format passing the json value to the flag. Example '--format json'"
	EdgeFuncInstanceDescribeHelpFlag         = "Displays more information about the describe subcommand"
	EdgeFuncInstanceFileWritten              = "File successfully written to: %s\n"

	// [ Update ]
	EdgeFuncInstanceUpdateUsage                 = "update --application-id <application_id> --instance-id <instance_id> [flags]"
	EdgeFuncInstanceUpdateShortDescription      = "Updates an Edge Functions instance"
	EdgeFuncInstanceUpdateLongDescription       = "Updates an Edge Functions instance, based on given attributes, to be used in Edge Applications"
	EdgeFuncInstanceUpdateFlagEdgeApplicationId = "Unique identifier for an Edge Application"
	EdgeFuncInstanceUpdateFlagIn                = "Path to a JSON file containing the attributes of the Edge Functions instance that will be updated; you can use - for reading from stdin"
	EdgeFuncInstanceUpdateFlagName              = "The Edge Functions instance name"
	EdgeFuncInstanceUpdateFlagArgs              = "The JSON args related to the Edge Functions instance being updated"
	EdgeFuncInstanceUpdateFlagInstanceID        = "Unique identifier for an Edge Functions instance"
	EdgeFuncInstanceUpdateFlagFunctionID        = "Edge Function ID to be used in the Edge Functions instance"
	EdgeFuncInstanceUpdateOutputSuccess         = "Updated Edge Functions instance with ID %d\n"
	EdgeFuncInstanceUpdateHelpFlag              = "Displays more information about the update subcommand"
)
