package edgefunctions

var (
	// EDGE FUNCTIONS MESSAGES

	//used by more than one cmd
	EdgeFunctionFlagId      = "Unique identifier of the Edge Function"
	EdgeFunctionFileWritten = "File successfully written to: %s\n"

	//Edge Functions cmd
	EdgeFunctionUsage            = "edge_functions"
	EdgeFunctionShortDescription = "Manages your Azion account's Edge Functions"
	EdgeFunctionLongDescription  = "You can create, update, delete, list and describe your Azion account's Edge Functions"

	//create cmd
	EdgeFunctionCreateUsage            = "create [flags]"
	EdgeFunctionCreateShortDescription = "Create a new Edge Function"
	EdgeFunctionCreateLongDescription  = "Create a new Edge Function"
	EdgeFunctionCreateFlagName         = "Your Edge Function's name (Mandatory if --in is not sent)"
	EdgeFunctionCreateFlagCode         = "Path to the file containing your Edge Function's code (Mandatory if --in is not sent)"
	EdgeFunctionCreateFlagActive       = "Whether your Edge Function should be active or not: <true|false> (Mandatory if --in is not sent)"
	EdgeFunctionCreateFlagArgs         = "Path to the file containing your Edge Function's JSON arguments"
	EdgeFunctionCreateFlagIn           = "Uses provided file path to create an Edge Function. You can use - for reading from stdin"
	EdgeFunctionCreateOutputSuccess    = "Created Edge Function with ID %d\n"

	//delete cmd
	EdgeFunctionDeleteUsage            = "delete --function-id <function_id> [flags]"
	EdgeFunctionDeleteShortDescription = "Deletes an Edge Function"
	EdgeFunctionDeleteLongDescription  = "Deletes an Edge Function based on the id given"
	EdgeFunctionDeleteOutputSuccess    = "Edge Function %d was successfully deleted\n"

	//describe cmd
	EdgeFunctionDescribeUsage            = "describe --function-id <function_id> [flags]"
	EdgeFunctionDescribeShortDescription = "Describes an Edge Function"
	EdgeFunctionDescribeLongDescription  = "Details an Edge Function based on the id given"
	EdgeFunctionDescribeFlagOut          = "Exports the command result to the received file path"
	EdgeFunctionDescribeFlagFormat       = "You can change the results format by passing json value to this flag"
	EdgeFunctionDescribeFlagWithCode     = "Displays the Edge Function's code (disabled by default)"

	//list cmd
	EdgeFunctionListUsage            = "list [flags]"
	EdgeFunctionListShortDescription = "Lists your account's Edge Functions"
	EdgeFunctionListLongDescription  = "Lists your account's Edge Functions"

	//update cmd
	EdgeFunctionUpdateUsage            = "update --function-id <function_id> [flags]"
	EdgeFunctionUpdateShortDescription = "Updates an Edge Function"
	EdgeFunctionUpdateLongDescription  = "Updates an Edge Function based on the id given"
	EdgeFunctionUpdateFlagName         = "Your Edge Function's name"
	EdgeFunctionUpdateFlagCode         = "Path to the file containing your Edge Function's code"
	EdgeFunctionUpdateFlagActive       = "Whether your Edge Function should be active or not: <true|false>"
	EdgeFunctionUpdateFlagArgs         = "Path to the file containing your Edge Function's JSON arguments"
	EdgeFunctionUpdateFlagIn           = "Uses provided file path to update the fields. You can use - for reading from stdin"
	EdgeFunctionUpdateOutputSuccess    = "Updated Edge Function with ID %d\n"
)
