package edge_functions_instances

var (
	// [ EdgeFunctionsInstances ]
	EdgeFunctionsInstancesUsage            = "edge_functions_instances"
	EdgeFunctionsInstancesShortDescription = "edge functions instances is the original source of data."
	EdgeFunctionsInstancesLongDescription  = "EdgeFunctionsInstances is the original source of data on edge platforms, where data is fetched when cache is not available."
	EdgeFunctionsInstancesFlagHelp         = "Displays more information about the edge functions instances command"

	// [ list ]
	EdgeFunctionsInstancesListUsage                 = "list [flags]"
	EdgeFunctionsInstancesListShortDescription      = "Displays your edge functions instances"
	EdgeFunctionsInstancesListLongDescription       = "Displays all edge functions instances related to your applications"
	EdgeFunctionsInstancesListHelpFlag              = "Displays more information about the list subcommand"
	EdgeFunctionsInstancesListFlagEdgeApplicationID = "Unique identifier for an edge application."

	// [ flags ]
	EdgeApplicationFlagId        = "Unique identifier of the Edge Application"
	EdgeFunctionsInstancesFlagId = "Unique identifier of the Edge Functions Instances"

	//domains cmd
	EdgeFuncInstanceUsage            = "edge_functions_instances"
	EdgeFuncInstanceShortDescription = "Create Edge Functions Instances for edges on Azion's platform"
	EdgeFuncInstanceLongDescription  = "Create Edge Functions Instances for edges on Azion's platform"
	EdgeFuncInstanceFlagHelp         = "Displays more information about the edge_functions_instances command"
	EdgeFuncInstanceFlagId           = "Unique identifier of the Edge Function Instance"
	ApplicationFlagId                = "Unique identifier for an edge application used by this domain.. The '--application-id' flag is mandatory"

	//delete cmd
	EdgeFuncInstanceDeleteUsage            = "delete --application-id <application_id> --function-id <function_id> [flags]"
	EdgeFuncInstanceDeleteShortDescription = "Removes an Edge Function Instance"
	EdgeFuncInstanceDeleteLongDescription  = "Removes an Edge Function Instance from the Domains library based on its given ID"
	EdgeFuncInstanceDeleteOutputSuccess    = "Edge Function Instance %s was successfully deleted\n"
	EdgeFuncInstanceDeleteHelpFlag         = "Displays more information about the delete subcommand"

	// [ create ]
	EdgeFuncInstanceCreateUsage                 = "create [flags]"
	EdgeFuncInstanceCreateShortDescription      = "Creates a new Function Instances"
	EdgeFuncInstanceCreateLongDescription       = "Creates an Function Instances based on given attributes to be used in edge applications"
	EdgeFuncInstanceCreateFlagEdgeApplicationId = "Unique identifier for an edge application"
	EdgeFuncInstanceCreateFlagEdgeFunctionID    = "Unique identifier for an Edge Function Instances"
	EdgeFuncInstanceCreateFlagName              = "The Function Instances name"
	EdgeFuncInstanceCreateFlagArgs              = "The args name"
	EdgeFuncInstanceCreateFlagIn                = "Path to a JSON file containing the attributes of the origin that will be created; you can use - for reading from stdin"
	EdgeFuncInstanceCreateOutputSuccess         = "Created Function Instances with ID %d\n"
	EdgeFuncInstanceCreateHelpFlag              = "Displays more information about the create subcommand"
)
