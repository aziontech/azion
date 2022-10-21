package edgeservices

var (

	// EDGE SERVICE MESSAGES

	//Edge Services cmd
	EdgeServiceUsage            = "edge_services <subcommand> [flags]"
	EdgeServiceShortDescription = "Manages your Azion account's Edge Services"
	EdgeServiceLongDescription  = "Manages your Edge Services of Edge Orchestrator"
	EdgeServiceHelpFlag         = "Displays more information about the edge_services command"

	//Edge Services Resources cmd
	EdgeServiceResourceUsage            = "resources <subcommand>"
	EdgeServiceResourceShortDescription = "Manages resources of a given Edge Service"
	EdgeServiceResourceLongDescription  = "Manages the Edge Services´s Resources"
	EdgeServiceResourceHelpFlag         = "Displays more information about the Resources subcommand"

	//used by more than one cmd
	EdgeServiceFlagId         = "Unique identifier of the Edge Service"
	EdgeServiceResourceFlagId = "Unique identifier of the Resource"
	EdgeServiceFlagOut        = "Exports the output to the given path <file_path/…/file_name.ext>"
	EdgeServiceFlagFormat     = "Changes the output format passing the json value to the flag"
	EdgeServiceFileWritten    = "File successfully written to: %s\n"

	//create cmd
	EdgeServiceCreateUsage            = "create [flags]"
	EdgeServiceCreateShortDescription = "Makes a new Edge Service"
	EdgeServiceCreateLongDescription  = "Makes a new Edge Service in the Azion Edge Orchestrator based on its name or configuration file"
	EdgeServiceCreateFlagName         = "The Edge Service's name"
	EdgeServiceCreateFlagIn           = "Path and file to create an Edge Service; you can use - for reading from stdin"
	EdgeServiceCreateOutputSuccess    = "Created Edge Service with ID %d\n"
	EdgeServiceCreateFlagHelp         = "Displays more information about the create subcommand"

	//delete cmd
	EdgeServiceDeleteUsage            = "delete --service-id <service_id> [flags]"
	EdgeServiceDeleteShortDescription = "Removes an Edge Service"
	EdgeServiceDeleteLongDescription  = "Removes an Edge Service based on its given ID"
	EdgeServiceDeleteOutputSuccess    = "Service %d was successfully deleted\n"
	EdgeServiceDeleteFlagHelp         = "Displays more information about the delete subcommand"

	//describe cmd
	EdgeServiceDescribeUsage            = "describe --service-id <service_id> [flags]"
	EdgeServiceDescribeShortDescription = "Returns the Edge Service data"
	EdgeServiceDescribeLongDescription  = "Displays information about the Edge Service via a given ID to show the service’s attributes in detail"
	EdgeServiceDescribeFlagWithVariable = "Displays the Edge Service's variables (disabled by default)"
	EdgeServiceDescribeOutputSuccess    = "Service %d was successfully deleted\n"
	EdgeServiceDescribeHelpFlag         = "Displays more information about the describe subcommand"

	//list cmd
	EdgeServiceListUsage            = "list [flags]"
	EdgeServiceListShortDescription = "Display your account’s Edge Services"
	EdgeServiceListLongDescription  = "Displays all Edge Services in the user’s Azion account"
	EdgeServiceListFlagHelp         = "Displays more information about the list subcommand"

	//update cmd
	EdgeServiceUpdateUsage            = "update --service-id <service_id> [flags]"
	EdgeServiceUpdateShortDescription = "Modifies an Edge Service"
	EdgeServiceUpdateLongDescription  = "Modifies the Edge Service attributes based on its ID"
	EdgeServiceUpdateFlagName         = "The Edge Service's name"
	EdgeServiceUpdateFlagActive       = "Whether the Edge Service should be active or not"
	EdgeServiceUpdateFlagVariables    = `Path to the file containing the Edge Service's variables; 
the accepted format to define the variables is one <KEY>=<VALUE> per line`
	EdgeServiceUpdateFlagIn        = "Given path and JSON file to automatically update the Edge Service attributes; you can use - for reading from stdin"
	EdgeServiceUpdateOutputSuccess = "Updated Edge Service with ID %d\n"
	EdgeServiceUpdateFlagHelp      = "Displays more information about the update subcommand"

	//EDGE SERVICE - RESOURCES MESSAGES

	//create cmd
	EdgeServiceResourceCreateUsage            = "create --service-id <service_id> [flags]"
	EdgeServiceResourceCreateShortDescription = "Makes a new Resource"
	EdgeServiceResourceCreateLongDescription  = "Makes a new resource in the Azion Platform based on its file’s path, name, and type"
	EdgeServiceResourceCreateFlagName         = "The Resource's path and name; mandatory"
	EdgeServiceResourceCreateFlagTrigger      = "The Resource's trigger; <Install|Reload|Uninstall>"
	EdgeServiceResourceCreateFlagContentType  = "The Resource's content-type; <shellscript|text>"
	EdgeServiceResourceCreateFlagContentFile  = "Path and name of the file with the Resource's content"
	EdgeServiceResourceCreateFlagIn           = "Path and file to create a Resource; you can use - for reading from stdin"
	EdgeServiceResourceCreateOutputSuccess    = "Created Resource with ID %d\n"
	EdgeServiceResourceCreateFlagHelp         = "Displays more information about the Resources create subcommand"

	//delete cmd
	EdgeServiceResourceDeleteUsage            = "delete --service-id <service_id> --resource-id <resource_id> [flags]"
	EdgeServiceResourceDeleteShortDescription = "Removes a Resource"
	EdgeServiceResourceDeleteLongDescription  = "Removes a Resource via given service ID and resource ID"
	EdgeServiceResourceDeleteOutputSuccess    = "Resource %d was successfully deleted\n"
	EdgeServiceResourceDeleteFlagHelp         = "Displays more information about the resources delete subcommand"

	//describe cmd
	EdgeServiceResourceDescribeUsage            = "describe --service-id <service_id> --resource-id <resource_id> [flags]"
	EdgeServiceResourceDescribeShortDescription = "Returns the Resource data"
	EdgeServiceResourceDescribeLongDescription  = "Displays information about the Resource via given service ID and resource ID to show the resources’ attributes in detail"
	EdgeServiceResourceDescribeOutputSuccess    = "Service %d was successfully deleted\n"
	EdgeServiceResourceDescribeFlagHelp         = "Displays more information about the resources describe subcommand"

	//list cmd
	EdgeServiceResourceListUsage            = "list --service-id <service_id> [flags]"
	EdgeServiceResourceListShortDescription = "Display the Resources of an Edge Service"
	EdgeServiceResourceListLongDescription  = "Displays all Resources of an Edge Service via the service ID"
	EdgeServiceResourceListFlagHelp         = "Displays more information about the resources list subcommand"

	//update cmd
	EdgeServiceResourceUpdateUsage            = "update --service-id <service_id> --resource-id <resource_id>[flags]"
	EdgeServiceResourceUpdateShortDescription = "Modifies a Resource"
	EdgeServiceResourceUpdateLongDescription  = "Modifies a Resource via a given service ID and resource ID to update its name, activity status, and other attributes"
	EdgeServiceResourceUpdateFlagName         = "The resource's path and name; <PATH>/<RESOURCE_NAME>"
	EdgeServiceResourceUpdateFlagTrigger      = "The resource's trigger; <Install|Reload|Uninstall>"
	EdgeServiceResourceUpdateFlagContentType  = "The resource's content-type; <shellscript | txt>"
	EdgeServiceResourceUpdateFlagContentFile  = "Path and name of the file with the resource's content"
	EdgeServiceResourceUpdateFlagIn           = "Path and file to update a resource; you can use - for reading from stdin"
	EdgeServiceResourceUpdateOutputSuccess    = "Updated Resource with ID %d\n"
	EdgeServiceResourceUpdateFlagHelp         = "Displays more information about the resources update subcommand"
)
