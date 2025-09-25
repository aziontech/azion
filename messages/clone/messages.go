package clone

var (
	Usage            = "clone <subcommand> [flags]"
	ShortDescription = "Clones a resource"
	LongDescription  = "Clones a resource based on the given name"
	FlagHelp         = "Displays more information about the clone command"

	// Application
	UsageApplication            = "application"
	ShortDescriptionApplication = "Clones an Application"
	LongDescriptionApplication  = "Clones an Application based on the given name"
	FlagNameApplication         = "Name that will be used by the new Application"
	FlagIdApplication           = "Identifier of which Application to clone"
	FlagHelpApplication         = "Displays more information about the 'clone application' command"
	CloneSuccessful             = "Application %s cloned successfully"
	AskApplicationIdClone       = "Enter the Application ID you wish to clone"
	AskApplicationNameClone     = "Enter the name that will be used by the new Application"
)
