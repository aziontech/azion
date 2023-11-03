package edgefunctions

var (
	// EDGE FUNCTIONS MESSAGES

	//used by more than one cmd
	EdgeFunctionFlagId      = "Unique identifier of the Edge Function"
	EdgeFunctionFileWritten = "File successfully written to: %s\n"

	//Edge Functions cmd
	EdgeFunctionUsage            = "edge_functions <subcommand>"
	EdgeFunctionShortDescription = "Manages your Azion account's Edge Functions"
	EdgeFunctionLongDescription  = "Manages serverless functions on the Edge Functions library"
	EdgeFunctionHelpFlag         = "Displays more information about the edge_functions command"

	//create cmd
	EdgeFunctionCreateUsage            = "create [flags]"
	EdgeFunctionCreateShortDescription = "Creates a new serverless Edge Function"
	EdgeFunctionCreateLongDescription  = "Creates an Edge Function based on given attributes to create a serverless code for Edge Applications"
	EdgeFunctionCreateFlagName         = "The Edge Function's name"
	EdgeFunctionCreateFlagCode         = "Path to the Edge Function's code"
	EdgeFunctionCreateFlagActive       = "Whether the Edge Function is active or not"
	EdgeFunctionCreateFlagArgs         = "Path to the Edge Function's arguments JSON file"
	EdgeFunctionCreateFlagIn           = "Given file path to create an Edge Function; you can use - for reading from stdin"
	EdgeFunctionCreateOutputSuccess    = "Created Edge Function with ID %d\n"
	EdgeFunctionCreateHelpFlag         = "Displays more information about the create subcommand"

	//delete cmd
	EdgeFunctionDeleteUsage            = "delete --function-id <function_id> [flags]"
	EdgeFunctionDeleteShortDescription = "Removes an Edge Function"
	EdgeFunctionDeleteLongDescription  = "Removes an Edge Function from the Edge Functions library based on its given ID"
	EdgeFunctionDeleteOutputSuccess    = "Edge Function %d was successfully deleted\n"
	EdgeFunctionDeleteHelpFlag         = "Displays more information about the delete subcommand"

	//describe cmd
	EdgeFunctionDescribeUsage            = "describe --function-id <function_id> [flags]"
	EdgeFunctionDescribeShortDescription = "Returns the Edge Function data"
	EdgeFunctionDescribeLongDescription  = "Displays information about the Edge Function via a given ID to show the function’s attributes in detail"
	EdgeFunctionDescribeFlagOut          = "Exports the output to the given <file_path/file_name.ext>"
	EdgeFunctionDescribeFlagFormat       = "Changes the output format passing the json value to the flag"
	EdgeFunctionDescribeFlagWithCode     = "Displays the Edge Function's code; disabled by default"
	EdgeFunctionDescribeHelpFlag         = "Displays more information about the describe command"

	//list cmd
	EdgeFunctionListUsage            = "list [flags]"
	EdgeFunctionListShortDescription = "Displays your account's Edge Functions"
	EdgeFunctionListLongDescription  = "Displays all functions in the user account’s Edge Functions library"
	EdgeFunctionListHelpFlag         = "Displays more information about the list subcommand"

	//update cmd
	EdgeFunctionUpdateUsage            = "update --function-id <function_id> [flags]"
	EdgeFunctionUpdateShortDescription = "Modifies an Edge Function"
	EdgeFunctionUpdateLongDescription  = "Modifies an Edge Function based on its ID to update its name, activity status, code path, and other attributes"
	EdgeFunctionUpdateFlagName         = "The Edge Function's name"
	EdgeFunctionUpdateFlagCode         = "Path and name to the file containing the Edge Function's code"
	EdgeFunctionUpdateFlagActive       = "Whether the Edge Function should be active or not"
	EdgeFunctionUpdateFlagArgs         = "Path and name of the JSON file containing the Edge Function's arguments"
	EdgeFunctionUpdateFlagIn           = "Given path and JSON file to automatically update the Edge Function attributes; you can use - for reading from stdin"
	EdgeFunctionUpdateOutputSuccess    = "Updated Edge Function with ID %d\n"
	EdgeFunctionUpdateHelpFlag         = "Displays more information about the update subcommand"
)
