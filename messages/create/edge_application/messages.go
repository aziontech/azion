package edge_application

var (
	// [ edge_applications ]
	Usage            = "edge-application"
	ShortDescription = "Creates an Edge Application"
	LongDescription  = "Creates an Edge Application without the need to manage infrastructure or security"
	FlagFile         = "Path to a JSON file containing the attributes of the Edge Application being created; you can use - for reading from stdin"
	FlagHelp         = "Displays more information about the create edge-application command"
	OutputSuccess    = "Created Edge Application with ID %d"

	FlagName                    = "Edge Application's name"
	FlagActive                  = "Whether the Edge Application is active or not"
	FlagDebugRules              = "Allows you to check whether rules created using Rules Engine for Edge Application have been successfully executed in your application"
	FlagApplicationAcceleration = "Whether the Edge Application has Application Acceleration active or not"
	FlagCaching                 = "Whether the Edge Application has Caching active or not"
	FlagEdgeFunctions           = "Whether the Edge Application has Edge Functions active or not"
	FlagImageOptimization       = "Whether the Edge Application has Image Optimization active or not"
	FlagTieredCaching           = "Whether the Edge Application has Tiered Caching active or not"
)
