package origin

var (
	Usage                    = "origin"
	ShortDescription         = "Updates an Origin"
	LongDescription          = "Updates an Origin based on its ID and given attributes"
	FlagOriginKey            = "The Origin's key unique identifier"
	FlagEdgeApplicationId    = "Unique identifier for an edge application"
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
	FlagIn                   = "Path to a JSON file containing the attributes of the origin that will be updated; you can use - for reading from stdin"
	FlagHelp                 = "Displays more information about the update subcommand"
	OutputSuccess            = "Updated origin with ID %s\n"
)
