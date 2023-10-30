package origin

const (
	Usage            = "origin"
	ShortDescription = "Creates a new origin"
	LongDescription  = "Creates an origin based on given attributes to be used in edge applications"
	OutputSuccess    = "Created origin with ID %d\n"

	// [ flags ]
	FlagEdgeApplicationID    = "Unique identifier for an edge application"
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
	FlagIn                   = "Path to a JSON file containing the attributes of the origin that will be created; you can use - for reading from stdin"
	FlagHelp                 = "Displays more information about the create subcommand"

	// [ ask ]
	AskAppID      = "What is the ID of the Edge Application?"
	AskName       = "What is the Name of the Origin?"
	AskAddresses  = "What is the Addresses of the Origin?"
	AskHostHeader = "What is the Host Header of the Origin?"
)
