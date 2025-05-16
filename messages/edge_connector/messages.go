package edgeconnector

var (
	// general
	Usage            = "edge-connector"
	FileWritten      = "File successfully written to: %s\n"
	ShortDescription = "Manages your Azion account's Edge Connectors"
	LongDescription  = "Manages serverless connectors on the Edge Connectors library"

	//create cmd
	CreateShortDescription = "Creates a new serverless Edge Connector"
	CreateLongDescription  = "Creates an Edge Connector based on given attributes to create a serverless code for Edge Applications"
	CreateOutputSuccess    = "Created Edge Connector with ID %d"

	//delete cmd
	DeleteShortDescription = "Deletes an Edge Connector"
	DeleteLongDescription  = "Removes an Edge Connector from the Edge Connectors library based on its given ID"
	DeleteOutputSuccess    = "Edge Connector %s was successfully deleted"
	DeleteHelpFlag         = "Displays more information about the delete edge-connector command"

	//describe cmd
	DescribeShortDescription   = "Returns the Edge Connector data"
	DescribeLongDescription    = "Displays information about the Edge Connector via a given ID to show the connector’s attributes in detail"
	DescribeFlagOut            = "Exports the output to the given <file_path/file_name.ext>"
	DescribeFlagFormat         = "Changes the output format passing the json value to the flag"
	DescribeFlagWithCode       = "Displays the Edge Connector's code; disabled by default"
	DescribeHelpFlag           = "Displays more information about the describe edge-connector command"
	DescribeAskInputFunctionID = "Enter the ID of the Edge Connector you wish to describe:"

	//list cmd
	ListShortDescription = "Displays your account's Edge Connectors"
	ListLongDescription  = "Displays all connectors in the user account’s Edge Connectors library"
	ListHelpFlag         = "Displays more information about the list edge-connector command"

	//update cmd
	UpdateShortDescription           = "Updates an Edge Connector"
	UpdateLongDescription            = "Modifies an Edge Connector based on its ID to update its name, activity status, code path, and other attributes"
	UpdateFlagName                   = "The Edge Connector's name"
	UpdateFlagCode                   = "Path and name to the file containing the Edge Connector's code"
	UpdateFlagActive                 = "Whether the Edge Connector should be active or not"
	UpdateFlagArgs                   = "Path and name of the JSON file containing the Edge Connector's arguments"
	UpdateFlagFile                   = "Given path and JSON file to automatically update the Edge Connector attributes; you can use - for reading from stdin"
	UpdateOutputSuccess              = "Updated Edge Connector with ID %d"
	UpdateHelpFlag                   = "Displays more information about the update edge-connector command"
	UpdateAskEdgeConnectorFunctionID = "Enter the ID of the Edge Connector you wish to update:"

	// flags
	FlagID                   = "Unique identifier of the Edge Connector"
	FlagName                 = "The Edge Connector's name"
	FlagType                 = "The Edge Connector's type ('http', 's3', 'edge_storage', 'live_ingest')"
	FlagAddresses            = "List of origin addresses"
	FlagLoadBalancer         = ""
	FlagOriginShield         = ""
	FlagTlsPolicy            = ""
	FlagLoadBalanceMethod    = ""
	FlagConnectionPreference = ""
	FlagConnectionTimeout    = ""
	FlagReadWriteTimeout     = ""
	FlagMaxRetries           = ""
	FlagVersions             = ""
	FlagHost                 = ""
	FlagPath                 = ""
	FlagFollowingRedirect    = ""
	FlagRealIpHeader         = ""
	FlagRealPortHeader       = ""
	FlagBucket               = ""
	FlagPrefix               = ""
	FlagActive               = "Whether the Edge Connector is active or not"
	FlagIn                   = "Given file path to create an Edge Connector; you can use - for reading from stdin"
	CreateFlagHelp           = "Displays more information about the create edge-connector command"

	// ask
	AskName                  = "Enter the new Edge Connector's name:"
	AskType                  = "Enter the Edge Connector's type:"
	AskActive                = "Select whether the Edge Connector is active or not"
	AskEdgeConnectorID       = "Enter the Edge Connector's ID:"
	UpdateAskEdgeConnectorID = "Enter the ID of the Edge Connector you wish to update:"
)
