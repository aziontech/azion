package origins

var (
	// [ origins ]
	Usage = "origin"

	// [create]
	CreateShortDescription = "Creates a new Origin"
	CreateLongDescription  = "Creates an Origin based on given attributes to be used in Edge Applications"
	CreateOutputSuccess    = "Created Origin with key %s\n"
	CreateFlagHelp         = "Displays more information about the create Origin command"

	// [delete]
	DeleteShortDescription = "Deletes an Origin"
	DeleteLongDescription  = "Deletes an Origin from the Edge Applications library based on its given ID"
	DeleteOutputSuccess    = "Origin %s was successfully deleted\n"
	DeleteHelpFlag         = "Displays more information about the delete Origin command"
	DeleteAskInputApp      = "Enter the ID of the Edge Application linked to this Origin:"
	DeleteAskInputOri      = "Enter the key of the Origin you wish to delete:"

	// [describe]
	DescribeShortDescription = "Returns information about a specific Origin"
	DescribeLongDescription  = "Returns information about a specific Origin, based on a given ID, in details"
	DescribeFlagOut          = "Exports the output to the given <file_path/file_name.ext>"
	DescribeFlagFormat       = "Changes the output format passing the json value to the flag"
	DescribeHelpFlag         = "Displays more information about the describe Origin command"

	// [list]
	ListShortDescription      = "Displays your origins"
	ListLongDescription       = "Displays all origins related to your applications"
	ListHelpFlag              = "Displays more information about the list Origin command"
	ListAskInputApplicationId = "Enter the ID of the Edge Application the Origins are linked to:"

	// [update]
	UpdateShortDescription      = "Updates an Origin"
	UpdateLongDescription       = "Updates an Origin based on its key and given attributes"
	UpdateFlagEdgeApplicationId = "Unique identifier for an Edge Application"
	UpdateFlagHelp              = "Displays more information about the update Origin command"
	UpdateOutputSuccess         = "Updated Origin with key %s\n"

	// [ ask ]
	AskAppID      = "Enter the ID of the Edge Application this Origin is linked to:"
	AskName       = "Enter the new Origin's Name:"
	AskAddresses  = "Enter the new Origin's Addresses:"
	AskHostHeader = "Enter the new Origin's Host Header:"
	AskOriginKey  = "Enter the Origin's Key:"

	// [ flags ]
	FlagEdgeApplicationID    = "Unique identifier for an Edge Application"
	FlagOriginKey            = "The Origin's key unique identifier"
	FlagName                 = "The Origin's name"
	FlagOriginType           = `Identifies the source of a record. I.e. "single_origin"`
	FlagAddresses            = "Passes a list of addresses linked to the Origin"
	FlagOriginProtocolPolicy = "Tells the protocol policy used in the Origin"
	FlagHostHeader           = "Specifies the hostname of the server being accessed"
	FlagOriginPath           = "Path to be appended to the URI when forwarding the request to the Origin. Leave it in blank to use only the URI"
	FlagHmacAuthentication   = "Whether Hmac Authentication is used or not"
	FlagHmacRegionName       = "Informs Hmac region name"
	FlagHmacAccessKey        = "Informs Hmac Access Key"
	FlagHmacSecretKey        = "Informs Hmac Secret Key"
	FlagFile                 = "Path to a JSON file containing the attributes of the Origin that will be created; you can use - for reading from stdin"
	OriginsFileWritten       = "File successfully written to: %s\n"
)
