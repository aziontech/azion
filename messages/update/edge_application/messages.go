package edgeapplication

var (
	Usage                       = "edge-application"
	ShortDescription            = "Updates an Edge Application"
	LongDescription             = "Modifies an Edge Application's name, activity status, and other attributes based on the given ID"
	FlagID                      = "The Edge Application's id"
	FlagName                    = "The Edge Application's name"
	FlagDebugRules              = "Allows you to check whether rules created using Rules Engine for Edge Application have been successfully executed in your application"
	FlagApplicationAcceleration = "Whether the Edge Application has Application Acceleration active or not"
	FlagCaching                 = "Whether the Edge Application has Caching active or not"
	FlagEdgeFunctions           = "Whether the Edge Application has Edge Functions active or not"
	FlagImageOptimization       = "Whether the Edge Application has Image Optimization active or not"
	FlagTieredCaching           = "Whether the Edge Application has Tiered Caching active or not"
	FlagFile                    = "Given path and JSON file to automatically update the Edge Application attributes; you can use - for reading from stdin"
	OutputSuccess               = "Updated Edge Application with ID %d"
	HelpFlag                    = "Displays more information about the update edge-application command"
	AskInputApplicationId       = "Enter the ID of the Edge Application you wish to update:"
)
