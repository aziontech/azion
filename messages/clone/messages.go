package clone

var (
	Usage            = "clone <subcommand> [flags]"
	ShortDescription = "Clones a resource"
	LongDescription  = "Clones a resource based on the given name"
	FlagHelp         = "Displays more information about the clone command"

	// Edge Application
	UsageEdgeApplication            = "edge-application"
	ShortDescriptionEdgeApplication = "Clones an Edge Application"
	LongDescriptionEdgeApplication  = "Clones an Edge Application based on the given name"
	FlagNameEdgeApplication         = "Name that will be used by the new Edge Application"
	FlagIdEdgeApplication           = "Identifier of which Edge Application to clone"
	FlagHelpEdgeApplication         = "Displays more information about the 'clone edge-application' command"
	CloneSuccessful                 = "Edge Application %s cloned successfully"
	AskApplicationIdClone           = "Enter the Edge Applicatin ID you wish to clone"
	AskApplicationNameClone         = "Enter the name that will be used by the new Edge Application"
)
