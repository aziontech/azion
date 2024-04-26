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
	CreateOutputSuccess    = "Created Edge Function with ID %d"

	//delete cmd
	DeleteShortDescription = "Removes an Edge Function"
	DeleteLongDescription  = "Removes an Edge Function from the Edge Functions library based on its given ID"
	DeleteOutputSuccess    = "Edge Function %d was successfully deleted"
	DeleteHelpFlag         = "Displays more information about the delete edge-function command"

	//describe cmd
	DescribeShortDescription   = "Returns the Edge Function data"
	DescribeLongDescription    = "Displays information about the Edge Function via a given ID to show the function’s attributes in detail"
	DescribeFlagOut            = "Exports the output to the given <file_path/file_name.ext>"
	DescribeFlagFormat         = "Changes the output format passing the json value to the flag"
	DescribeFlagWithCode       = "Displays the Edge Function's code; disabled by default"
	DescribeHelpFlag           = "Displays more information about the describe edge-function command"
	DescribeAskInputFunctionID = "Enter the ID of the Edge Function you wish to describe:"

	//list cmd
	ListShortDescription = "Displays your account's Edge Functions"
	ListLongDescription  = "Displays all functions in the user account’s Edge Functions library"
	ListHelpFlag         = "Displays more information about the list edge-function command"

	//update cmd
	UpdateShortDescription  = "Modifies an Edge Function"
	UpdateLongDescription   = "Modifies an Edge Function based on its ID to update its name, activity status, code path, and other attributes"
	UpdateFlagName          = "The Edge Function's name"
	UpdateFlagCode          = "Path and name to the file containing the Edge Function's code"
	UpdateFlagActive        = "Whether the Edge Function should be active or not"
	UpdateFlagArgs          = "Path and name of the JSON file containing the Edge Function's arguments"
	UpdateFlagFile          = "Given path and JSON file to automatically update the Edge Function attributes; you can use - for reading from stdin"
	UpdateOutputSuccess     = "Updated Edge Function with ID %d"
	UpdateHelpFlag          = "Displays more information about the update edge-function command"
	UpdateAskEdgeFunctionID = "Enter the ID of the Edge Function you wish to update:"

	// flags
	FlagID         = "Unique identifier of the Edge Function"
	FlagName       = "The Edge Function's name"
	FlagCode       = "Path to the Edge Function's code"
	FlagActive     = "Whether the Edge Function is active or not"
	FlagArgs       = "Path to the Edge Function's arguments JSON file"
	FlagIn         = "Given file path to create an Edge Function; you can use - for reading from stdin"
	CreateFlagHelp = "Displays more information about the create edge-function command"

	// ask
	AskName           = "Enter the new Edge Function's name:"
	AskCode           = "Enter the file path of the Edge Function's source code:"
	AskActive         = "Select whether the Edge Function is active or not"
	AskEdgeFunctionID = "Enter the Edge Function's ID:"
)
