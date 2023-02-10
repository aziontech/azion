package origins

var (
  // [ origins ]
	OriginsUsage                     = "origins"
	OriginsShortDescription          = "Origins is where data is fetched when the cache is not available."
  OriginsLongDescription           = "Origins is the original source of data in content delivery systems (CDN), where data is fetched when cache is not available. Data is stored at origin and can be retrieved by clients through cache servers distributed around the world. The fetch is done from the cache first, and if the data is not available, it is fetched from origin and saved in the cache for future use. This allows for fast data delivery."
	OriginsFlagHelp                  = "Displays more information about the origins command"

  // [ list ]
	OriginsListUsage                 = "list [flags]"
	OriginsListShortDescription      = "Displays yours origins"
	OriginsListLongDescription       = "Displays all your origins references to your edges"
	OriginsListHelpFlag              = "Displays more information about the list subcommand"
	OriginsListFlagEdgeApplicationID = "Is a unique identifier for the edge application that references the origins to direct data requests correctly."

	// [ describe ]
	OriginsDescribeUsage             = "describe --application-id <domain_id> --origin-id <origin_id> [flags]"
	OriginsDescribeShortDescription  = "Returns the origin data"
	OriginsDescribeLongDescription   = "Displays information about the origin via a given ID to show the application’s attributes in detail"
	OriginsDescribeFlagApplicationID = "Is a unique identifier for the edge application that references the origins to direct data requests correctly."
	OriginsDescribeFlagOriginID      = `is a unique identifier that identifies an "origins" in a list of results returned by the API. The "GetOrigin" function uses the "Origin Id" to search for the desired "origins" and returns an error if it is not found.`
	OriginsDescribeFlagOut           = "Exports the output to the given <file_path/file_name.ext>"
	OriginsDescribeFlagFormat        = "Changes the output format passing the json value to the flag"
	OriginsDescribeHelpFlag          = "Displays more information about the describe command"

  // [ general ]
	OriginsFileWritten                 = "File successfully written to: %s\n"
)
