package application

var (
	// [ edge_applications ]
	Usage            = "application"
	ShortDescription = "Creates an Application"
	LongDescription  = "Creates an Application without the need to manage infrastructure or security"
	FlagFile         = "Path to a JSON file containing the attributes of the Application being created; you can use - for reading from stdin"
	FlagHelp         = "Displays more information about the create application command"
	OutputSuccess    = "Created Application with ID %d"

	FlagName                    = "Application's name"
	FlagActive                  = "Whether the Application is active or not"
	FlagDebugRules              = "Allows you to check whether rules created using Rules Engine for Application have been successfully executed in your application"
	FlagApplicationAcceleration = "Whether the Application has Application Acceleration active or not"
	FlagCaching                 = "Whether the Application has Caching active or not"
	FlagEdgeFunctions           = "Whether the Application has Functions active or not"
	FlagImageOptimization       = "Whether the Application has Image Optimization active or not"
	FlagTieredCaching           = "Whether the Application has Tiered Caching active or not"

	//V3 flags
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
	FlagSupportedCiphers               = "Determines which cryptographic algorithms will be used in the TLS connections of your Application"
	FlagWebsocket                      = "Allows you to establish the WebSocket communication protocol between your Origin and your users under the reverse proxy architecture."
	FlagBrowserCacheSettingsMaximumTtl = "Defines the maximum time to live (TTL) of cached resources in the browser. It can be used to set a time limit for how long resources can be cached in the browser."
	FlagCdnCacheSettingsMaximumTtl     = "Defines the maximum time to live (TTL) of cached resources in the Application. It can be used to set a time limit for how long resources can be cached on the Application servers."
)
