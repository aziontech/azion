package root

var (
	RootUsage       = "azion <command> <subcommand> [flags]"
	RootDescription = "The Azion Command Line Interface is a unified tool to manage your Azion projects and resources"
	RootHelpFlag    = "Displays more information about the Azion CLI"
	RootDoNotUpdate = "Do not receive update notification"
	RootLogDebug    = "Displays log at a debug level"
	RootLogLevel    = "Set the logging level, \"debug\", \"info\", or \"error\"."
	RootLogSilent   = "Silences log completely; mostly used for automation purposes"
	RootTokenFlag   = "Saves a given Personal Token locally to authorize CLI commands"
	RootConfigFlag  = "Sets the Azion configuration folder for the current command only, without changing persistent settings."
	RootYesFlag     = "Answers all yes/no interactions automatically with yes"
	TokenSavedIn    = "Token saved in %s\n"
	TokenUsedIn     = "This token will be used by default with all commands"

	// update messages
	NewVersion        = "There is a new version of Azion CLI available\n"
	BrewUpdate        = "Please run: 'brew upgrade azion' to update to the latest version\n"
	ReleasePage       = "Please visit our Releases page and download the appropriate file\n"
	CouldNotGetUser   = "Sadly, we could not get information on your system to indicate the correct update form\n"
	DownloadRelease   = "Visit https://github.com/aziontech/azion/releases to download the correct package"
	RpmUpdate         = "Please run: 'sudo rpm -i <downloaded_file>' to update to the latest version\n"
	DpkgUpdate        = "Please run: 'sudo dpkg -i <downloaded_file>' to update to the latest version\n"
	ApkUpdate         = "Please run: 'sudo apk add <downloaded_file>' to update to the latest version\n"
	AskCollectMetrics = "To better understand user needs and enhance our application, we gather anonymous data. Do you agree to participate? (Y/n)"
	UnsupportedOS     = "Unsupported Operating System\n"
)
