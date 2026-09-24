package application

// Used by both the v3 and the v4 command trees.
var (
	ShortDescription            = "Updates an Application"
	LongDescription             = "Modifies an Application's name, activity status, and other attributes based on the given ID"
	FlagID                      = "The Application's id"
	FlagName                    = "The Application's name"
	FlagApplicationAcceleration = "Whether the Application has Application Acceleration active or not"
	FlagFunctions               = "Whether the Application has Functions active or not"
	FlagImageOptimization       = "Whether the Application has Image Optimization active or not"
	FlagFile                    = "Given path and JSON file to automatically update the Application attributes; you can use - for reading from stdin"
	OutputSuccess               = "Updated Application with ID %d"
	HelpFlag                    = "Displays more information about the update application command"
	AskInputApplicationId       = "Enter the ID of the Application you wish to update:"
)

// Used only by the v4 command tree.
var (
	Usage          = "application"
	FlagDebugRules = "Allows you to check whether rules created using Rules Engine for Application have been successfully executed in your application"
	FlagCaching    = "Whether the Application has Caching active or not"
)

// Used only by the v3 command tree (bug-fix-only; see doc/plan.md).
var (
	FlagDeliveryProtocol   = "Specify whether the data should be delivered via HTTP or HTTPS."
	FlagHttpPort           = "Flag to set the HTTP port or ports your application will use. 80 as default."
	FlagHttpsPort          = "Flag to set the HTTPs port or ports your application will use. 443 as default."
	FlagMinimumTlsVersion  = "The Application's Minimum Tls Version"
	FlagDeviceDetection    = "Whether the Application has Device Detection active or not"
	FlagFirewall           = "Whether the Application has Firewall active or not"
	FlagL2Caching          = "Whether the Application has L2 Caching active or not"
	FlagLoadBalancer       = "Whether the Application has Load Balancer active or not"
	RawLogs                = "Whether the Application has Raw Logs active or not"
	WebApplicationFirewall = "Whether the Application has Web Application Firewall active or not"
)
