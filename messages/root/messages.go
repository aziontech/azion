package root

var (
	RootUsage       = "azion <command> <subcommand> [flags]\n azion [flags]\n\t"
	RootDescription = "Azion CLI %s"
	RootHelpFlag    = "Displays more information about the Azion CLI"
	RootDoNotUpdate = "Do not receive update notification"
	RootLogDebug    = "Displays log at a debug level"
	RootLogLevel    = "Set the logging level, \"debug\", \"info\", or \"error\". Default \"info\"."
	RootLogSilent   = "Silences log completely; mostly used for automation purposes"
	RootTokenFlag   = "Saves a given personal token locally to authorize CLI commands"
	RootConfigFlag  = "Sets the Azion configuration folder for the current command only, without changing persistent settings."
	RootYesFlag     = "Answers all yes/no interactions automatically with yes"
	TokenSavedIn    = "Token saved in %v\n"
	TokenUsedIn     = "This token will be used by default when calling any command"
)
