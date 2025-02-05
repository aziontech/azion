package edgeapplication

var (
	Usage                       = "edge-application"
	ShortDescription            = "Updates an Edge Application"
	LongDescription             = "Modifies an Edge Application's name, activity status, and other attributes based on the given ID"
	FlagID                      = "The Edge Application's id"
	FlagName                    = "The Edge Application's name"
	FlagDeliveryProtocol        = "The Edge Application's Delivery Protocol"
	FlagHttpPort                = "The Edge Application's Http Port"
	FlagHttpsPort               = "The Edge Application's Https Port"
	FlagMinimumTlsVersion       = "The Edge Application's Minimum Tls Version"
	FlagApplicationAcceleration = "Whether the Edge Application has Application Acceleration active or not"
	FlagCaching                 = "Whether the Edge Application has Caching active or not"
	FlagDeviceDetection         = "Whether the Edge Application has Device Detection active or not"
	FlagEdgeFirewall            = "Whether the Edge Application has Edge Firewall active or not"
	FlagEdgeFunctions           = "Whether the Edge Application has Edge Functions active or not"
	FlagImageOptimization       = "Whether the Edge Application has Image Optimization active or not"
	FlagL2Caching               = "Whether the Edge Application has L2 Caching active or not"
	FlagLoadBalancer            = "Whether the Edge Application has Load Balancer active or not"
	RawLogs                     = "Whether the Edge Application has Raw Logs active or not"
	WebApplicationFirewall      = "Whether the Edge Application has Web Application Firewall active or not"
	FlagFile                    = "Given path and JSON file to automatically update the Edge Application attributes; you can use - for reading from stdin"
	OutputSuccess               = "Updated Edge Application with ID %d"
	HelpFlag                    = "Displays more information about the update edge-application command"
	AskInputApplicationId       = "Enter the ID of the Edge Application you wish to update:"
)
