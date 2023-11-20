package origins

var (
	// [ origins ]
	Usage = "origin"

	// [create]
	CreateShortDescription = "Creates a new origin"
	CreateLongDescription  = "Creates an origin based on given attributes to be used in edge applications"
	CreateOutputSuccess    = "Created origin with key %s\n"
	CreateFlagHelp         = "Displays more information about the create origin command"

	// [delete]
	DeleteShortDescription = "Deletes an Origin"
	DeleteLongDescription  = "Deletes an Origin from the Edge Applications library based on its given ID"
	DeleteOutputSuccess    = "Origin %s was successfully deleted\n"
	DeleteHelpFlag         = "Displays more information about the delete origin command"
	DeleteAskInputApp      = "What is the id of the edge application linked to this origin?"
	DeleteAskInputOri      = "What is the key of the origin you wish to delete?"

	// [describe]
	DescribeShortDescription = "Returns information about a specific origin"
	DescribeLongDescription  = "Returns information about a specific origin, based on a given ID, in details"
	DescribeFlagOut          = "Exports the output to the given <file_path/file_name.ext>"
	DescribeFlagFormat       = "Changes the output format passing the json value to the flag"
	DescribeHelpFlag         = "Displays more information about the describe origin command"

	// [list]
	ListShortDescription      = "Displays your origins"
	ListLongDescription       = "Displays all origins related to your applications"
	ListHelpFlag              = "Displays more information about the list origin command"
	ListAskInputApplicationId = "What is the id of the Edge Application the origins are linked to?"

	// [update]
	UpdateShortDescription      = "Updates an Origin"
	UpdateLongDescription       = "Updates an Origin based on its key and given attributes"
	UpdateFlagEdgeApplicationId = "Unique identifier for an edge application"
	UpdateFlagHelp              = "Displays more information about the update origin command"
	UpdateOutputSuccess         = "Updated origin with key %s\n"

	// [ ask ]
	AskAppID      = "What is the ID of the Edge Application this origin is linked to?"
	AskName       = "What is the Name of the Origin?"
	AskAddresses  = "What is the Addresses of the Origin?"
	AskHostHeader = "What is the Host Header of the Origin?"
	AskOriginKey  = "What is the Key of the Origin?"

	// [ flags ]
	FlagEdgeApplicationID    = "Unique identifier for an edge application"
	FlagOriginKey            = "The Origin's key unique identifier"
	FlagName                 = "The origin's name"
	FlagOriginType           = `Identifies the source of a record. I.e. "single_origin"`
	FlagAddresses            = "Passes a list of addresses linked to the origin"
	FlagOriginProtocolPolicy = "Tells the protocol policy used in the origin"
	FlagHostHeader           = "Specifies the hostname of the server being accessed"
	FlagOriginPath           = "Path to be appended to the URI when forwarding the request to the origin. Leave it in blank to use only the URI"
	FlagHmacAuthentication   = "Whether Hmac Authentication is used or not"
	FlagHmacRegionName       = "Informs Hmac region name"
	FlagHmacAccessKey        = "Informs Hmac Access Key"
	FlagHmacSecretKey        = "Informs Hmac Secret Key"
	FlagFile                 = "Path to a JSON file containing the attributes of the origin that will be created; you can use - for reading from stdin"
	OriginsFileWritten       = "File successfully written to: %s\n"
)
