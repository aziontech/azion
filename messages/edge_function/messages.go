package edgefunction

var (
	// general
	Usage            = "edge-function"
	FileWritten      = "File successfully written to: %s\n"
	ShortDescription = "Manages your Azion account's Edge Functions"
	LongDescription  = "Manages serverless functions on the Edge Functions library"

	//create cmd
	CreateShortDescription = "Creates a new serverless Edge Function"
	CreateLongDescription  = "Creates an Edge Function based on given attributes to create a serverless code for Edge Applications"
	CreateOutputSuccess    = "Created Edge Function with ID %d\n"

	//delete cmd
	DeleteUsage            = "delete --function-id <function_id> [flags]"
	DeleteShortDescription = "Removes an Edge Function"
	DeleteLongDescription  = "Removes an Edge Function from the Edge Functions library based on its given ID"
	DeleteOutputSuccess    = "Edge Function %d was successfully deleted\n"
	DeleteHelpFlag         = "Displays more information about the delete subcommand"

	//describe cmd
	DescribeUsage            = "describe --function-id <function_id> [flags]"
	DescribeShortDescription = "Returns the Edge Function data"
	DescribeLongDescription  = "Displays information about the Edge Function via a given ID to show the function’s attributes in detail"
	DescribeFlagOut          = "Exports the output to the given <file_path/file_name.ext>"
	DescribeFlagFormat       = "Changes the output format passing the json value to the flag"
	DescribeFlagWithCode     = "Displays the Edge Function's code; disabled by default"
	DescribeHelpFlag         = "Displays more information about the describe command"

	//list cmd
	ListUsage            = "list [flags]"
	ListShortDescription = "Displays your account's Edge Functions"
	ListLongDescription  = "Displays all functions in the user account’s Edge Functions library"
	ListHelpFlag         = "Displays more information about the list subcommand"

	//update cmd
	UpdateUsage            = "update --function-id <function_id> [flags]"
	UpdateShortDescription = "Modifies an Edge Function"
	UpdateLongDescription  = "Modifies an Edge Function based on its ID to update its name, activity status, code path, and other attributes"
	UpdateFlagName         = "The Edge Function's name"
	UpdateFlagCode         = "Path and name to the file containing the Edge Function's code"
	UpdateFlagActive       = "Whether the Edge Function should be active or not"
	UpdateFlagArgs         = "Path and name of the JSON file containing the Edge Function's arguments"
	UpdateFlagIn           = "Given path and JSON file to automatically update the Edge Function attributes; you can use - for reading from stdin"
	UpdateOutputSuccess    = "Updated Edge Function with ID %d\n"

	// flags
	FlagID     = "Unique identifier of the Edge Function"
	FlagName   = "The Edge Function's name"
	FlagCode   = "Path to the Edge Function's code"
	FlagActive = "Whether the Edge Function is active or not"
	FlagArgs   = "Path to the Edge Function's arguments JSON file"
	FlagIn     = "Given file path to create an Edge Function; you can use - for reading from stdin"
	FlagHelp   = "Displays more information about the edge function subcommand"

	// ask
	AskName           = "What is the Name of the edge function?"
	AskCode           = "What is the Code Path  of the edge function?"
	AskActive         = "Does activating the edge function?"
	AskEdgeFunctionID = "What is the ID of the edge function you wish to delete?"
)
