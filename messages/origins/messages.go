package origins

var (
	// [ origins ]
	OriginsUsage            = "origins"
	OriginsShortDescription = "Origins is the original source of data."
	OriginsLongDescription  = "Origins is the original source of data on edge platforms, where data is fetched when cache is not available."
	OriginsFlagHelp         = "Displays more information about the origins command"

	// [ list ]
	OriginsListUsage                 = "list [flags]"
	OriginsListShortDescription      = "Displays your origins"
	OriginsListLongDescription       = "Displays all origins related to your applications"
	OriginsListHelpFlag              = "Displays more information about the list subcommand"
	OriginsListFlagEdgeApplicationID = "Unique identifier for an edge application."

	// [ describe ]
	OriginsDescribeUsage             = "describe --application-id <application_id> --origin-id <origin_id> [flags]"
	OriginsDescribeShortDescription  = "Returns information about a specific origin"
	OriginsDescribeLongDescription   = "Returns information about a specific origin, based on a given ID, in details"
	OriginsDescribeFlagApplicationID = "Unique identifier for an edge application. The '--application-id' flag is mandatory"
	OriginsDescribeFlagOriginID      = "Unique identifier for an origin. The '--origin-id' flag is mandatory"
	OriginsDescribeFlagOut           = "Exports the output to the given <file_path/file_name.ext>"
	OriginsDescribeFlagFormat        = "Changes the output format passing the json value to the flag"
	OriginsDescribeHelpFlag          = "Displays more information about the describe subcommand"

	// [ create ]
	OriginsCreateUsage                    = "create [flags]"
	OriginsCreateShortDescription         = "Creates a new origin"
	OriginsCreateLongDescription          = "Creates an origin based on given attributes to be used in edge applications"
	OriginsCreateFlagEdgeApplicationId    = "Unique identifier for an edge application"
	OriginsCreateFlagName                 = "The origin's name"
	OriginsCreateFlagOriginType           = `Identifies the source of a record. I.e. "single_origin"`
	OriginsCreateFlagAddresses            = "Passes a list of addresses linked to the origin"
	OriginsCreateFlagOriginProtocolPolicy = "Tells the protocol policy used in the origin"
	OriginsCreateFlagHostHeader           = "Specifies the hostname of the server being accessed"
	OriginsCreateFlagOriginPath           = "Path to be appended to the URI when forwarding the request to the origin. Leave it in blank to use only the URI"
	OriginsCreateFlagHmacAuthentication   = "Whether Hmac Authentication is used or not"
	OriginsCreateFlagHmacRegionName       = "Informs Hmac region name"
	OriginsCreateFlagHmacAccessKey        = "Informs Hmac Access Key"
	OriginsCreateFlagHmacSecretKey        = "Informs Hmac Secret Key"
	OriginsCreateFlagIn                   = "Path to a JSON file containing the attributes of the origin that will be created; you can use - for reading from stdin"
	OriginsCreateOutputSuccess            = "Created origin with ID %d\n"
	OriginsCreateHelpFlag                 = "Displays more information about the create subcommand"

	// [ update ]
	OriginsUpdateUsage                    = "update [flags]"
	OriginsUpdateShortDescription         = "Updates an Origin"
	OriginsUpdateLongDescription          = "Updates an Origin based on its ID and given attributes"
	OriginsUpdateFlagOriginKey            = "The Origin's key unique identifier"
	OriginsUpdateFlagEdgeApplicationId    = "Unique identifier for an edge application"
	OriginsUpdateFlagName                 = "The origin's name"
	OriginsUpdateFlagOriginType           = `Identifies the source of a record. I.e. "single_origin"`
	OriginsUpdateFlagAddresses            = "Passes a list of addresses linked to the origin"
	OriginsUpdateFlagOriginProtocolPolicy = "Tells the protocol policy used in the origin"
	OriginsUpdateFlagHostHeader           = "Specifies the hostname of the server being accessed"
	OriginsUpdateFlagOriginPath           = "Path to be appended to the URI when forwarding the request to the origin. Leave it in blank to use only the URI"
	OriginsUpdateFlagHmacAuthentication   = "Whether Hmac Authentication is used or not"
	OriginsUpdateFlagHmacRegionName       = "Informs Hmac region name"
	OriginsUpdateFlagHmacAccessKey        = "Informs Hmac Access Key"
	OriginsUpdateFlagHmacSecretKey        = "Informs Hmac Secret Key"
	OriginsUpdateFlagIn                   = "Path to a JSON file containing the attributes of the origin that will be updated; you can use - for reading from stdin"
	OriginsUpdateOutputSuccess            = "Updated origin with ID %s\n"
	OriginsUpdateHelpFlag                 = "Displays more information about the update subcommand"

	// [ delete ]
	OriginsDeleteUsage             = "delete [flags]"
	OriginsDeleteShortDescription  = "Deletes an Origin"
	OriginsDeleteLongDescription   = "Deletes an Origin from the Edge Applications library based on its given ID"
	OriginsDeleteOutputSuccess     = "Origin %s was successfully deleted\n"
	OriginsDeleteFlagApplicationID = "Unique identifier for an edge application"
	OriginsDeleteFlagOriginKey     = "The Origin's key unique identifier"
	OriginsDeleteHelpFlag          = "Displays more information about the delete subcommand"

	// [ general ]
	OriginsFileWritten = "File successfully written to: %s\n"
)
