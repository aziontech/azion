package edge_application

var (
	// [ edge_applications ]
	Usage            = "edge-application"
	ShortDescription = "Creates an edge application on Azion's platform"
	LongDescription  = "Creates an edge application without the need to manage infrastructure or security"
	FlagIn           = "Path to a JSON file containing the attributes of the edge application being created; you can use - for reading from stdin"
	FlagHelp         = "Displays more information about the edge_application command"
	OutputSuccess    = "Created edge application with ID %d\n"

	FlagName                           = "Edge application's name"
	FlagApplicationAcceleration        = "Used for application acceleration, enable or disable"
	FlagDeliveryProtocol               = "Specify whether the data should be delivered via HTTP, HTTPS, FTP or another communication protocol."
	FlagHttp3                          = "Flag to enable HTTP3"
	FlagHttpPort                       = "Flag to set your HTTP port custom"
	FlagOriginType                     = "Type of the origin. Possible values: 'single_origin'(default value), 'load_balancer' or 'live_ingest'."
	FlagAddress                        = "Specify the address of a resource or server."
	FlagHostHeader                     = "Flag to customize your host headers"
	FlagOriginProtocolPolicy           = "Type of connection between the edge nodes and your origin. Possible values: 'preserve', 'http' or 'https'"
	FlagBrowserCacheSettings           = "Configures the amount of time that content is cached in the user’s browser. Possible values: 'honor' or 'override'"
	FlagCdnCacheSettings               = "Configures how Azion caches the content at the edge. Possible values: 'honor' or 'override'"
	FlagDebugRules                     = "Allows you to check whether rules or rule sets created using the Rules Engine module for Edge Application and Edge Firewall have been successfully executed in your application"
	FlagSupportedCiphers               = "Determines which cryptographic algorithms will be used in the TLS connections of your edge application"
	FlagWebsocket                      = "Allows you to establish the WebSocket communication protocol between your origin and your users under the reverse proxy architecture."
	FlagBrowserCacheSettingsMaximumTtl = "Defines the maximum time to live (TTL) of cached resources in the browser. It can be used to set a time limit for how long resources can be cached in the browser."
	FlagCdnCacheSettingsMaximumTtl     = "Defines the maximum time to live (TTL) of cached resources in the CDN. It can be used to set a time limit for how long resources can be cached on the CDN servers."
)
