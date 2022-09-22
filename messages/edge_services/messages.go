package edgeservices

var (

	// EDGE SERVICE MESSAGES

	//Edge Services cmd
	EdgeServiceUsage            = "edge_services"
	EdgeServiceShortDescription = "Manages your Azion account's Edge Services"
	EdgeServiceLongDescription  = "You can create, update, delete, list and describe your Azion account's Edge Services"

	//Edge Services Resources cmd
	EdgeServiceResourceUsage            = "resources"
	EdgeServiceResourceShortDescription = "Manages resources in a given Edge Service"
	EdgeServiceResourceLongDescription  = "You can create, update, delete, list and describe Resources in a given Edge Service"

	//used by more than one cmd
	EdgeServiceFlagId         = "Unique identifier of the Edge Service"
	EdgeServiceResourceFlagId = "Unique identifier of the Resource"
	EdgeServiceFlagOut        = "Exports the command result to the received file path"
	EdgeServiceFlagFormat     = "You can change the results format by passing json value to this flag"
	EdgeServiceFileWritten    = "File successfully written to: %s\n"

	//create cmd
	EdgeServiceCreateUsage            = "create [flags]"
	EdgeServiceCreateShortDescription = "Creates a new Edge Service"
	EdgeServiceCreateLongDescription  = "Creates a new Edge Service"
	EdgeServiceCreateFlagName         = "Your Edge Service's name (Mandatory)"
	EdgeServiceCreateFlagIn           = "Uses provided file path to create an Edge Service. You can use - for reading from stdin"
	EdgeServiceCreateOutputSuccess    = "Created Edge Service with ID %d\n"

	//delete cmd
	EdgeServiceDeleteUsage            = "delete --service-id <service_id> [flags]"
	EdgeServiceDeleteShortDescription = "Deletes an Edge Service"
	EdgeServiceDeleteLongDescription  = "Deletes an Edge Service based on the id given"
	EdgeServiceDeleteOutputSuccess    = "Service %d was successfully deleted\n"

	//describe cmd
	EdgeServiceDescribeUsage            = "describe --service-id <service_id> [flags]"
	EdgeServiceDescribeShortDescription = "Describes an Edge Service"
	EdgeServiceDescribeLongDescription  = "Details an Edge Service based on the id given"
	EdgeServiceDescribeFlagWithVariable = "Displays the Edge Service's variables (disabled by default)"
	EdgeServiceDescribeOutputSuccess    = "Service %d was successfully deleted\n"

	//list cmd
	EdgeServiceListUsage            = "list [flags]"
	EdgeServiceListShortDescription = "Lists your account's Edge Services"
	EdgeServiceListLongDescription  = "Lists your account's Edge Services"

	//update cmd
	EdgeServiceUpdateUsage            = "update --service-id <service_id> [flags]"
	EdgeServiceUpdateShortDescription = "Updates an Edge Service"
	EdgeServiceUpdateLongDescription  = "Updates an Edge Service"
	EdgeServiceUpdateFlagName         = "Your Edge Service's name"
	EdgeServiceUpdateFlagActive       = "Whether your Edge Service should be active or not: <true|false>"
	EdgeServiceUpdateFlagVariables    = `Path to the file containing your Edge Service's Variables.
The accepted format for defining variables is one <KEY>=<VALUE> per line`
	EdgeServiceUpdateFlagIn        = "Uses provided file path to update an Edge Service. You can use - for reading from stdin"
	EdgeServiceUpdateOutputSuccess = "Updated Edge Service with ID %d\n"

	//EDGE SERVICE - RESOURCES MESSAGES

	//create cmd
	EdgeServiceResourceCreateUsage            = "create --service-id <service_id> [flags]"
	EdgeServiceResourceCreateShortDescription = "Creates a new Resource"
	EdgeServiceResourceCreateLongDescription  = "Creates a new Resource in an Edge Service based on the service_id given`"
	EdgeServiceResourceCreateFlagName         = "Your Resource's name: <PATH>/<RESOURCE_NAME> (Mandatory)"
	EdgeServiceResourceCreateFlagTrigger      = "Your Resource's trigger: <Install|Reload|Uninstall>"
	EdgeServiceResourceCreateFlagContentType  = "Your Resource's content-type: <shellscript|text> (Mandatory)"
	EdgeServiceResourceCreateFlagContentFile  = "Path to the file with your Resource's content (Mandatory)"
	EdgeServiceResourceCreateFlagIn           = "Uses provided file path to create a Resource. You can use - for reading from stdin"
	EdgeServiceResourceCreateOutputSuccess    = "Created Resource with ID %d\n"

	//delete cmd
	EdgeServiceResourceDeleteUsage            = "delete --service-id <service_id> --resource-id <resource_id> [flags]"
	EdgeServiceResourceDeleteShortDescription = "Deletes a Resource"
	EdgeServiceResourceDeleteLongDescription  = "Deletes a Resource based on the service_id and resource_id given"
	EdgeServiceResourceDeleteOutputSuccess    = "Resource %d was successfully deleted\n"

	//describe cmd
	EdgeServiceResourceDescribeUsage            = "describe --service-id <service_id> --resource-id <resource_id> [flags]"
	EdgeServiceResourceDescribeShortDescription = "Describes a Resource"
	EdgeServiceResourceDescribeLongDescription  = "Provides a long description of a Resource based on a service_id and a resource_id given"
	EdgeServiceResourceDescribeOutputSuccess    = "Service %d was successfully deleted\n"

	//list cmd
	EdgeServiceResourceListUsage            = "list --service-id <service_id> [flags]"
	EdgeServiceResourceListShortDescription = "Lists the Resources in a given Edge Service"
	EdgeServiceResourceListLongDescription  = "Lists the Resources in a given Edge Service"

	//update cmd
	EdgeServiceResourceUpdateUsage            = "update --service-id <service_id> --resource-id <resource_id>[flags]"
	EdgeServiceResourceUpdateShortDescription = "Updates a Resource"
	EdgeServiceResourceUpdateLongDescription  = "Updates a Resource based on a resource_id"
	EdgeServiceResourceUpdateFlagName         = "Your Resource's name: <PATH>/<RESOURCE_NAME>"
	EdgeServiceResourceUpdateFlagTrigger      = "Your Resource's trigger: <Install|Reload|Uninstall>"
	EdgeServiceResourceUpdateFlagContentType  = "Your Resource's content-type: <shellscript|text>"
	EdgeServiceResourceUpdateFlagContentFile  = "Path to the file with your Resource's content"
	EdgeServiceResourceUpdateFlagIn           = "Uses provided file path to update a Resource. You can use - for reading from stdin"
	EdgeServiceResourceUpdateOutputSuccess    = "Updated Resource with ID %d\n"
)
