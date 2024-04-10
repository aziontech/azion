package purge

var (
	Usage            = "purge"
	ShortDescription = "Removes cache object before time-out"
	LongDescription  = "Deletes an object from the Edge Cache or Tiered Cache layers before time-out"
	FlagHelp         = "Displays more information about the purge command"
	FlagWildcard     = "Specifies the Wildcard URL or Cache Key for the objects you want to purge. Only one Wildcard expression can be used per request."
	FlagCacheKeys    = "Provides a list of URLs that must be purged from Azion Edge Cache"
	FlagUrls         = "Provides a list of URLs that must be purged from Azion Edge Cache"
	FlagLayer        = "Specifies the layer the purge will be executed. Possible values: 'edge_caching' or 'l2_caching'"
	PurgeSuccessful  = "Purge carried out successfully"
	AskForInput      = "Enter the URLs you wish to purge, separated by commas"
)
