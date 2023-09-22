package edge_application

var (
	// [ edge_applications ]
	Usage            = "edge-application"
	ShortDescription = "Creates an edge application on Azion's platform"
	LongDescription  = "Build your Edge applications in minutes without the need to manage infrastructure or security"
	FlagIn           = "Path to a JSON file containing the attributes of the edge application being created; you can use - for reading from stdin"
	FlagHelp         = "Displays more information about the edge_application command"
	OutputSuccess    = "Created edge application with ID %d\n"

	FlagName                           = "Edge application name created"
	FlagApplicationAcceleration        = "Used for application acceleration, enable or disable"
	FlagDeliveryProtocol               = "specify whether the data should be delivered via HTTP, HTTPS, FTP or another communication protocol."
	FlagOriginType                     = "Type of data source. It can be used to identify whether the data source is a web server, a data repository"
	FlagAddress                        = "specify the location of a resource or server."
	FlagOriginProtocolPolicy           = "define rules governing how data is handled when communicating with the origin"
	FlagBrowserCacheSettings           = "Browser cache settings. Can be used to control the behavior of the browser cache in relation to application or website resources"
	FlagCdnCacheSettings               = "Cache settings of a Content Distribution Network (CDN). It can be used to specify how resources are cached on the CDN servers."
	FlagBrowserCacheSettingsMaximumTtl = "The maximum time to live (TTL) of cached resources in the browser. It can be used to set a time limit for how long resources can be cached in the browser."
	FlagCdnCacheSettingsMaximumTtl     = "Maximum time to live (TTL) of cached resources in the CDN. It can be used to set a time limit for how long resources can be cached on the CDN servers."
)
