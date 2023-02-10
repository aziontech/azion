package origins

var (
	// [ origins ]
	OriginsUsage            = "origins"
	OriginsShortDescription = "Origins is the original source of data."
	OriginsLongDescription  = "Origins is the original source of data in edge platforms, where data is fetched when cache is not available. Data is stored at origin and can be retrieved by clients through cache servers distributed around the world. The fetch is done from the cache first, and if the data is not available, it is fetched from origin and saved in the cache for future use. This allows for fast data delivery."
	OriginsFlagHelp         = "Displays more information about the origins command"

	// [ list ]
	OriginsListUsage                 = "list [flags]"
	OriginsListShortDescription      = "Displays your origins"
	OriginsListLongDescription       = "Displays all origins of your applications"
	OriginsListHelpFlag              = "Displays more information about the list subcommand"
	OriginsListFlagEdgeApplicationID = "Unique identifier for an edge application."

	// [ describe ]
	OriginsDescribeUsage             = "describe --application-id <application_id> --origin-id <origin_id> [flags]"
	OriginsDescribeShortDescription  = "Returns information about a specific origin"
	OriginsDescribeLongDescription   = "Displays information about a specific origin based on a given ID origin's attributes in detail"
	OriginsDescribeFlagApplicationID = "Unique identifier for an edge application."
	OriginsDescribeFlagOriginID      = "Unique identifier for an origin."
	OriginsDescribeFlagOut           = "Exports the output to the given <file_path/file_name.ext>"
	OriginsDescribeFlagFormat        = "Changes the output format passing the json value to the flag"
	OriginsDescribeHelpFlag          = "Displays more information about the describe subcommand"

	// [ create ]
	OriginsCreateUsage                    = "create [flags]"
	OriginsCreateShortDescription         = "Creates a new origin"
	OriginsCreateLongDescription          = "Creates an origin based on given attributes to be used in edge applications"
	OriginsCreateFlagEdgeApplicationId    = "Unique identifier for an edge application."
	OriginsCreateFlagName                 = "The origin's name"
	OriginsCreateFlagOriginType           = "Identifies the source of a record" // Research to list possible types
	OriginsCreateFlagAddresses            = "Passes a list of addresses linked to the origin"
	OriginsCreateFlagOriginProtocolPolicy = "Tells the protocol policy used in the origin" //Is an origin protocol policy that specifies how Amazon CloudFront should respond to requests for content. This policy specifies whether CloudFront should use the origin protocol (HTTP or HTTPS) to get content from an origin server, or whether it should use the origin protocol regardless of the protocol used for the request."
	OriginsCreateFlagHostHeader           = "Specifies the hostname of the server being accessed"//Informs the HostHeaderThe HostHeader is an HTTP header field that specifies the hostname of the server being accessed. It is used to identify which website or application is being accessed. The HostHeader is sent by the browser to the web server and is used to determine which website or application to load."
	OriginsCreateFlagOriginPath           = "Gives a file path to the origin"//OriginPath is a file path that is used to identify the source of a file or directory. It is typically used to track the original location of a file or directory before it was moved or copied to a new location. OriginPath is often used in backup and restore operations to ensure that the original file or directory is not overwritten or lost."
	OriginsCreateFlagHmacAuthentication   = "Whether Hmac Authentication is used or not" 
	OriginsCreateFlagHmacRegionName       = "Informs Hmac region name" //"HmacRegionName is a field used in the Amazon Web Services (AWS) API to identify the region in which a particular request is being made. It is used to ensure that requests are routed to the correct region and that the correct authentication credentials are used."
	OriginsCreateFlagHmacAccessKey        = "Informs Hmac Access Key"
	OriginsCreateFlagHmacSecretKey        = "Informs Hmac Secret Key"
	OriginsCreateFlagIn                   = "Path to a JSON file containing the attributes of the origin that will be created; you can use - for reading from stdin"
	OriginsCreateOutputSuccess            = "Created origin with ID %d\n"
	OriginsCreateHelpFlag                 = "Displays more information about the create subcommand"

	// [ general ]
	OriginsFileWritten = "File successfully written to: %s\n"
)
