package origins

// NOTE: @PatrickMenoti Why do we leave the messages of the global commands if they are specific to each command, I would like to know what you think about leaving them in the root of each command, and if there are messages that are shared, which is something that already happens, we add them to the ultils.

const (
	usage            = "origins [flags]"
	shortDescription = "Creates a new origin"
	longDescription  = "Creates an origin based on given attributes to be used in edge applications"

	outputSuccess = "Created origin with ID %d\n"

	example = `
        $ azion create origins --application-id 1673635839 --name "drink coffe" --addresses "asdfg.asd" --host-header "host"
        $ azion create origins --application-id 1673635839 --in "create.json"
        `

	// [ flags ]
	flagEdgeApplicationID    = "Unique identifier for an edge application"
	flagName                 = "The origin's name"
	flagOriginType           = `Identifies the source of a record. I.e. "single_origin"`
	flagAddresses            = "Passes a list of addresses linked to the origin"
	flagOriginProtocolPolicy = "Tells the protocol policy used in the origin"
	flagHostHeader           = "Specifies the hostname of the server being accessed"
	flagOriginPath           = "Path to be appended to the URI when forwarding the request to the origin. Leave it in blank to use only the URI"
	flagHmacAuthentication   = "Whether Hmac Authentication is used or not"
	flagHmacRegionName       = "Informs Hmac region name"
	flagHmacAccessKey        = "Informs Hmac Access Key"
	flagHmacSecretKey        = "Informs Hmac Secret Key"
	flagIn                   = "Path to a JSON file containing the attributes of the origin that will be created; you can use - for reading from stdin"
	flagHelp                 = "Displays more information about the create subcommand"

	// [ errors ]
	errorCreateOrigins          = "Failed to create the Origin: %s. Check your settings and try again. If the error persists, contact Azion support."
	errorHmacAuthenticationFlag = "Invalid --hmac-authentication flag provided. The flag must have  'true' or 'false' values. Run the command 'azion <command> <subcommand> --help' to display more information and try again."
)
