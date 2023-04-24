package edge_functions_instances

var (
	// [ EdgeFunctionsInstances ]
	EdgeFunctionsInstancesUsage            = "edge_functions_instances"
	EdgeFunctionsInstancesShortDescription = "Edge Functions Instances allows you to instantiate serverless functions in your edge applications at Azion."
	EdgeFunctionsInstancesLongDescription  = "Edge Functions Instances allows you to instantiate serverless functions in your edge applications at Azion, as well as set up conditions for their execution."
	EdgeFunctionsInstancesFlagHelp         = "Displays more information about the edge functions instances command"

	// [ list ]
	EdgeFunctionsInstancesListUsage                 = "list [flags]"
	EdgeFunctionsInstancesListShortDescription      = "Displays your edge functions instances."
	EdgeFunctionsInstancesListLongDescription       = "Displays all edge functions instances related to a specific edge application."
	EdgeFunctionsInstancesListHelpFlag              = "Displays more information about the list subcommand"
	EdgeFunctionsInstancesListFlagEdgeApplicationID = "Unique identifier for an edge application."

	// [ flags ]
	EdgeApplicationFlagId        = "Unique identifier for an edge application"
	EdgeFunctionsInstancesFlagId = "Unique identifier for an edge functions instance"

	//Edge Functions Instances cmd
	EdgeFuncInstanceUsage            = "edge_functions_instances"
	EdgeFuncInstanceShortDescription = "Edge Functions Instances allows you to instantiate serverless functions in your edge applications at Azion."
	EdgeFuncInstanceLongDescription  = "Edge Functions Instances allows you to instantiate serverless functions in your edge applications at Azion, as well as set up conditions for their execution."
	EdgeFuncInstanceFlagHelp         = "Displays more information about the edge functions instances command"
	EdgeFuncInstanceFlagId           = "Unique identifier for an edge functions instance"
	ApplicationFlagId                = "Unique identifier for the edge application related to an edge functions instance. The '--application-id' flag is required"

	//delete cmd
	EdgeFuncInstanceDeleteUsage            = "delete --application-id <application_id> --instance-id <instance-id>"
	EdgeFuncInstanceDeleteShortDescription = "Removes an edge functions instance"
	EdgeFuncInstanceDeleteLongDescription  = "Removes an edge functions instance, instantiated in a specific edge application, based on the given flags."
	EdgeFuncInstanceDeleteOutputSuccess    = "Edge functions instance %s was successfully deleted\n"
	EdgeFuncInstanceDeleteHelpFlag         = "Displays more information about the delete subcommand"

	// [ create ]
	EdgeFuncInstanceCreateUsage                 = "create [flags]"
	EdgeFuncInstanceCreateShortDescription      = "Creates a new edge functions instance"
	EdgeFuncInstanceCreateLongDescription       = "Creates a new edge functions instance based on given attributes to be used in an edge application"
	EdgeFuncInstanceCreateFlagEdgeApplicationId = "Unique identifier for an edge application"
	EdgeFuncInstanceCreateFlagEdgeFunctionID    = "Unique identifier for an edge functions instance"
	EdgeFuncInstanceCreateFlagName              = "The edge functions instance name"
	EdgeFuncInstanceCreateFlagArgs              = "The JSON args related to the edge functions instance being created"
	EdgeFuncInstanceCreateFlagIn                = "Path to a JSON file containing the attributes of the edge functions instance being created; you can use - for reading from stdin"
	EdgeFuncInstanceCreateOutputSuccess         = "Created Function Instances with ID %d\n"
	EdgeFuncInstanceCreateHelpFlag              = "Displays more information about the create subcommand"

	//describe cmd
	EdgeFuncInstanceDescribeUsage            = "describe --application-id <application_id> --instance-id <instance_id> [flags]"
	EdgeFuncInstanceDescribeShortDescription = "Returns the information related to the edge functions instance"
	EdgeFuncInstanceDescribeLongDescription  = "Returns the information related to the edge functions instance, informed through the flag '--instance-id' in detail"
	EdgeFuncInstanceDescribeFlagOut          = "Exports the output of the subcommand 'describe' to the given file path <file_path/file_name.ext>"
	EdgeFuncInstanceDescribeFlagFormat       = "Changes the output format passing the json value to the flag. Example '--format json'"
	EdgeFuncInstanceDescribeHelpFlag         = "Displays more information about the describe subcommand"
	EdgeFuncInstanceFileWritten              = "File successfully written to: %s\n"


	// [ Update ]
	EdgeFuncInstanceUpdateUsage                 = "update --application-id <application_id> --instance-id <instance_id> [flags]"
	EdgeFuncInstanceUpdateShortDescription      = "Updates an edge functions instance"
	EdgeFuncInstanceUpdateLongDescription       = "Updates an edge functions instance, based on given attributes, to be used in edge applications"
	EdgeFuncInstanceUpdateFlagEdgeApplicationId = "Unique identifier for an edge application"
	EdgeFuncInstanceUpdateFlagIn                = "Path to a JSON file containing the attributes of the edge functions instance that will be updated; you can use - for reading from stdin"
	EdgeFuncInstanceUpdateFlagName              = "The edge functions instance name"
	EdgeFuncInstanceUpdateFlagArgs              = "The JSON args related to the edge functions instance being created"
	EdgeFuncInstanceUpdateFlagInstanceID        = "Unique identifier for an edge functions instance"
	EdgeFuncInstanceUpdateFlagFunctionID        = "Edge function ID to be used in the edge functions instance"
	EdgeFuncInstanceUpdateOutputSuccess         = "Updated edge functions instance with ID %d\n"
	EdgeFuncInstanceUpdateHelpFlag              = "Displays more information about the update subcommand"
)