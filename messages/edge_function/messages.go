package edgefunction

var (
	// general
	Usage            = "function"
	FileWritten      = "File successfully written to: %s\n"
	ShortDescription = "Manages your Azion account's Functions"
	LongDescription  = "Manages serverless functions on the Functions library"

	//create cmd
	CreateShortDescription = "Creates a new serverless function"
	CreateLongDescription  = "Creates a function based on given attributes to create a serverless code for Applications"
	CreateOutputSuccess    = "Created function with ID %d"

	//delete cmd
	DeleteShortDescription = "Deletes a function"
	DeleteLongDescription  = "Removes a function from the functions library based on its given ID"
	DeleteOutputSuccess    = "Function %s was successfully deleted"
	DeleteHelpFlag         = "Displays more information about the delete function command"

	//describe cmd
	DescribeShortDescription   = "Returns the function data"
	DescribeLongDescription    = "Displays information about the function via a given ID to show the function’s attributes in detail"
	DescribeFlagOut            = "Exports the output to the given <file_path/file_name.ext>"
	DescribeFlagFormat         = "Changes the output format passing the json value to the flag"
	DescribeFlagWithCode       = "Displays the function's code; disabled by default"
	DescribeHelpFlag           = "Displays more information about the describe function command"
	DescribeAskInputFunctionID = "Enter the ID of the function you wish to describe:"

	//list cmd
	ListShortDescription = "Displays your account's functions"
	ListLongDescription  = "Displays all functions in the user account’s functions library"
	ListHelpFlag         = "Displays more information about the list function command"

	//update cmd
	UpdateShortDescription = "Updates a function"
	UpdateLongDescription  = "Modifies a function based on its ID to update its name, activity status, code path, and other attributes"
	UpdateFlagName         = "The function's name"
	UpdateFlagCode         = "Path and name to the file containing the function's code"
	UpdateFlagActive       = "Whether the function should be active or not"
	UpdateFlagArgs         = "Path and name of the JSON file containing the function's arguments"
	UpdateFlagFile         = "Given path and JSON file to automatically update the function attributes; you can use - for reading from stdin"
	UpdateOutputSuccess    = "Updated function with ID %d"
	UpdateHelpFlag         = "Displays more information about the update function command"
	UpdateAskFunctionID    = "Enter the ID of the function you wish to update:"

	// flags
	FlagID                   = "Unique identifier of the function"
	FlagName                 = "The function's name"
	FlagCode                 = "Path to the function's code"
	FlagActive               = "Whether the function is active or not"
	FlagArgs                 = "Path to the function's arguments JSON file"
	FlagIn                   = "Given file path to create a function; you can use - for reading from stdin"
	FlagExecutionEnvironment = "Either 'edge_application' or 'edge_firewall'"
	CreateFlagHelp           = "Displays more information about the create function command"

	// ask
	AskName       = "Enter the new function's name:"
	AskInitiator  = "Enter the new function's initiator type:"
	AskCode       = "Enter the file path of the function's source code:"
	AskActive     = "Select whether the function is active or not"
	AskFunctionID = "Enter the function's ID:"
)
