package edge_functions_instances

var (
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
)
