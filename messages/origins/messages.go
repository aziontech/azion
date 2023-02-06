package origins

var (
  //origins cmd
	OriginsUsage                   = "origins"
	OriginsShortDescription        = "Origins is where data is fetched when the cache is not available."
  OriginsLongDescription         = "Origins is the original source of data in content delivery systems (CDN), where data is fetched when cache is not available. Data is stored at origin and can be retrieved by clients through cache servers distributed around the world. The fetch is done from the cache first, and if the data is not available, it is fetched from origin and saved in the cache for future use. This allows for fast data delivery."
	OriginsFlagHelp                = "Displays more information about the origins command"

  //list cmd
	OriginsListUsage            = "list [flags]"
	OriginsListShortDescription = "Displays yours origins"
	OriginsListLongDescription  = "Displays all your origins references to your edges"
	OriginsListHelpFlag         = "Displays more information about the list subcommand"
)
