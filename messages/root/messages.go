package root

var (
	RootUsage       = "azion <command> <subcommand> [flags]"
	RootDescription = "The Azion Command Line Interface is a unified tool to manage your Azion projects and resources"
	RootHelpFlag    = "Displays more information about the Azion CLI"
	RootDoNotUpdate = "Do not receive update notification"
	RootLogDebug    = "Displays log at a debug level"
	RootLogLevel    = "Set the logging level, \"debug\", \"info\", or \"error\"."
	RootLogSilent   = "Silences log completely; mostly used for automation purposes"
	RootTokenFlag   = "Saves a given personal token locally to authorize CLI commands"
	RootConfigFlag  = "Sets the Azion configuration folder for the current command only, without changing persistent settings."
	RootYesFlag     = "Answers all yes/no interactions automatically with yes"
	TokenSavedIn    = "Token saved in %v\n"
	TokenUsedIn     = "This token will be used by default with all commands"
)
