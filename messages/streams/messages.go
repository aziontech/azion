package streams

var (
	// general
	Usage            = "streams"
	FileWritten      = "File successfully written to: %s\n"
	ShortDescription = "Manages your Data Stream streams"
	LongDescription  = "Manages the Data Stream streams used to send your data to external endpoints and analytics tools"

	// create cmd
	CreateShortDescription = "Creates a new stream"
	CreateLongDescription  = "Creates a Data Stream stream based on given attributes"
	CreateOutputSuccess    = "Created stream with ID %d"
	CreateFlagHelp         = "Displays more information about the create streams command"

	// delete cmd
	DeleteShortDescription = "Deletes a stream"
	DeleteLongDescription  = "Removes a stream based on its given ID"
	DeleteOutputSuccess    = "Stream %d was successfully deleted"
	DeleteHelpFlag         = "Displays more information about the delete streams command"

	// describe cmd
	DescribeShortDescription = "Returns the stream data"
	DescribeLongDescription  = "Displays information about the stream via a given ID to show its attributes in detail"
	DescribeFlagOut          = "Exports the output to the given <file_path/file_name.ext>"
	DescribeFlagFormat       = "Changes the output format passing the json value to the flag"
	DescribeHelpFlag         = "Displays more information about the describe streams command"

	// list cmd
	ListShortDescription = "Displays your account's streams"
	ListLongDescription  = "Displays all Data Stream streams in the user account"
	ListHelpFlag         = "Displays more information about the list streams command"

	// update cmd
	UpdateShortDescription = "Updates a stream"
	UpdateLongDescription  = "Modifies a stream based on its ID to update its name, activity status, inputs, transforms, and other attributes"
	UpdateFlagFile         = "Given path and JSON file to automatically update the stream attributes; you can use - for reading from stdin"
	UpdateOutputSuccess    = "Updated stream with ID %d"
	UpdateHelpFlag         = "Displays more information about the update streams command"
	UpdateAskStreamID      = "Enter the ID of the stream you wish to update:"
	UpdateAskStreamFile    = "Enter the path of the json to update the stream:"

	// flags
	FlagID = "Unique identifier of the stream"
	FlagIn = "Given file path to create a stream; you can use - for reading from stdin"

	// ask
	AskStreamID   = "Enter the stream's ID:"
	AskCreateFile = "Enter the path of the json to create the stream:"
)
