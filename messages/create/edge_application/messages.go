package edge_application

var (
	// [ edge_applications ]
	Usage            = "edge-application"
	ShortDescription = "Creates an Edge Application"
	LongDescription  = "Creates an Edge Application without the need to manage infrastructure or security"
	FlagFile         = "Path to a JSON file containing the attributes of the Edge Application being created; you can use - for reading from stdin"
	FlagHelp         = "Displays more information about the create edge-application command"
	OutputSuccess    = "Created Edge Application with ID %d"

	FlagName                           = "Edge Application's name"
	FlagApplicationAcceleration        = "Used for Application Acceleration, enable or disable"
	FlagDeliveryProtocol               = "Specify whether the data should be delivered via HTTP or HTTPS."
	FlagHttp3                          = "Flag to enable HTTP3"
	FlagOriginType                     = "Type of the Origin. Possible values: 'single_origin'(default value), 'load_balancer' or 'live_ingest'."
	FlagHttpPort                       = "Flag to set the HTTP port or ports your application will use. 80 as default."
	FlagHttpsPort                      = "Flag to set the HTTPs port or ports your application will use. 443 as default."
	FlagAddress                        = "Specify the address of a resource or server."
	FlagHostHeader                     = "Flag to customize your host headers"
	FlagOriginProtocolPolicy           = "Type of connection between the edge nodes and your Origin. Possible values: 'preserve', 'http' or 'https'"
	FlagBrowserCacheSettings           = "Configures the amount of time that content is cached in the user’s browser. Possible values: 'honor' or 'override'"
	FlagCdnCacheSettings               = "Configures how Azion caches the content at the edge. Possible values: 'honor' or 'override'"
	FlagDebugRules                     = "Allows you to check whether rules created using Rules Engine for Edge Application have been successfully executed in your application"
	FlagSupportedCiphers               = "Determines which cryptographic algorithms will be used in the TLS connections of your Edge Application"
	FlagWebsocket                      = "Allows you to establish the WebSocket communication protocol between your Origin and your users under the reverse proxy architecture."
	FlagBrowserCacheSettingsMaximumTtl = "Defines the maximum time to live (TTL) of cached resources in the browser. It can be used to set a time limit for how long resources can be cached in the browser."
	FlagCdnCacheSettingsMaximumTtl     = "Defines the maximum time to live (TTL) of cached resources in the CDN. It can be used to set a time limit for how long resources can be cached on the CDN servers."
)
