package edgeconnector

var (
	// general
	Usage            = "connector"
	FileWritten      = "File successfully written to: %s\n"
	ShortDescription = "Manages your Azion account's Connectors"
	LongDescription  = "Manages serverless connectors on the Connectors library"

	//create cmd
	CreateShortDescription = "Creates a new serverless Connector"
	CreateLongDescription  = "Creates a Connector based on given attributes"
	CreateOutputSuccess    = "Created Connector with ID %d"

	//delete cmd
	DeleteShortDescription = "Deletes a Connector"
	DeleteLongDescription  = "Removes a Connector from the Connectors library based on its given ID"
	DeleteOutputSuccess    = "Connector %s was successfully deleted"
	DeleteHelpFlag         = "Displays more information about the delete connector command"

	//describe cmd
	DescribeShortDescription   = "Returns the Connector data"
	DescribeLongDescription    = "Displays information about the Connector via a given ID to show the connector’s attributes in detail"
	DescribeFlagOut            = "Exports the output to the given <file_path/file_name.ext>"
	DescribeFlagFormat         = "Changes the output format passing the json value to the flag"
	DescribeFlagWithCode       = "Displays the Connector's code; disabled by default"
	DescribeHelpFlag           = "Displays more information about the describe connector command"
	DescribeAskInputFunctionID = "Enter the ID of the Connector you wish to describe:"

	//list cmd
	ListShortDescription = "Displays your account's Connectors"
	ListLongDescription  = "Displays all connectors in the user account’s Connectors library"
	ListHelpFlag         = "Displays more information about the list connector command"

	//update cmd
	UpdateShortDescription           = "Updates a Connector"
	UpdateLongDescription            = "Modifies a Connector based on its ID to update its name, activity status, code path, and other attributes"
	UpdateFlagName                   = "The Connector's name"
	UpdateFlagCode                   = "Path and name to the file containing the Connector's code"
	UpdateFlagActive                 = "Whether the Connector should be active or not"
	UpdateFlagArgs                   = "Path and name of the JSON file containing the Connector's arguments"
	UpdateFlagFile                   = "Given path and JSON file to automatically update the Connector attributes; you can use - for reading from stdin"
	UpdateOutputSuccess              = "Updated Connector with ID %d"
	UpdateHelpFlag                   = "Displays more information about the update connector command"
	UpdateAskEdgeConnectorFunctionID = "Enter the ID of the Connector you wish to update:"
	UpdateAskEdgeConnectorType       = "Enter the type of the Connector you wish to update:"
	UpdateAskEdgeConnectorFile       = "Enter the path of the json to update the Connector:"

	// flags
	FlagID         = "Unique identifier of the Connector"
	FlagName       = "The Connector's name"
	FlagType       = "The Connector's type ('http', 'storage', 'live_ingest')"
	FlagAddresses  = "List of origin addresses"
	FlagActive     = "Whether the Connector is active or not"
	FlagIn         = "Given file path to create an Connector; you can use - for reading from stdin"
	CreateFlagHelp = "Displays more information about the create connector command"

	// ask
	AskName                  = "Enter the new Connector's name:"
	AskType                  = "Enter the Connector's type:"
	AskActive                = "Select whether the Connector is active or not"
	AskEdgeConnectorID       = "Enter the Connector's ID:"
	UpdateAskEdgeConnectorID = "Enter the ID of the Connector you wish to update:"
)
