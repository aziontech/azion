package connector

// Used only by the v4 command tree.
var (
	// general
	Usage = "connector"

	//create cmd
	CreateShortDescription = "Creates a new serverless Connector"
	CreateLongDescription  = "Creates a Connector based on given attributes"
	CreateOutputSuccess    = "Created Connector with ID %d"

	//delete cmd
	DeleteShortDescription = "Deletes a Connector"
	DeleteLongDescription  = "Removes a Connector from the Connectors library based on its given ID"
	DeleteOutputSuccess    = "Connector %d was successfully deleted"
	DeleteHelpFlag         = "Displays more information about the delete connector command"

	//describe cmd
	DescribeShortDescription = "Returns the Connector data"
	DescribeLongDescription  = "Displays information about the Connector via a given ID to show the connector’s attributes in detail"
	DescribeHelpFlag         = "Displays more information about the describe connector command"

	//list cmd
	ListShortDescription = "Displays your account's Connectors"
	ListLongDescription  = "Displays all connectors in the user account’s Connectors library"
	ListHelpFlag         = "Displays more information about the list connector command"

	//update cmd
	UpdateShortDescription = "Updates a Connector"
	UpdateLongDescription  = "Modifies a Connector based on its ID to update its name, activity status, code path, and other attributes"
	UpdateFlagFile         = "Given path and JSON file to automatically update the Connector attributes; you can use - for reading from stdin"
	UpdateOutputSuccess    = "Updated Connector with ID %d"
	UpdateHelpFlag         = "Displays more information about the update connector command"
	UpdateAskConnectorType = "Enter the type of the Connector you wish to update:"
	UpdateAskConnectorFile = "Enter the path of the json to update the Connector:"

	// flags
	FlagID         = "Unique identifier of the Connector"
	FlagType       = "The Connector's type ('http', 'storage', 'live_ingest')"
	FlagIn         = "Given file path to create an Connector; you can use - for reading from stdin"
	CreateFlagHelp = "Displays more information about the create connector command"

	AskConnectorID       = "Enter the Connector's ID:"
	UpdateAskConnectorID = "Enter the ID of the Connector you wish to update:"
)
