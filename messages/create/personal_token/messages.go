package personaltoken

var (
	CreateUsage            = "personal-token"
	CreateShortDescription = "Creates a Personal Token"
	CreateLongDescription  = "Creates a Personal Token to be used for authentication and security"
	CreateFlagName         = "The Personal Token's name"
	CreateFlagExpiresAt    = "The Personal Token's expiration"
	CreateFlagDescription  = "The Personal Token's description"
	CreateFlagFile         = "Path to a JSON file containing the attributes of the Personal Token being created; you can use - for reading from stdin"
	CreateOutputSuccess    = "Created Personal Token: %s\n"
	CreateHelpFlag         = "Displays more information about the 'create personal-token' subcommand"
	AskInputName           = "What is the name of the Personal Token?"
	AskInputExpiration     = "What is the expiration of the Personal Token?"
)
