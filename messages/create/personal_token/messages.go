package personaltoken

var (
	CreateUsage            = "personal-token"
	CreateShortDescription = "Creates a Personal Token"
	CreateLongDescription  = "Creates a Personal Token to be used for authentication and security"
	CreateFlagName         = "The Personal Token's name"
	CreateFlagExpiresAt    = "The Personal Token's expiration"
	CreateFlagDescription  = "The Personal Token's description"
	CreateFlagFile         = "Path to a JSON file containing the attributes of the Personal Token being created; you can use - for reading from stdin"
	CreateOutputSuccess    = "Created Personal Token: %s"
	CreateHelpFlag         = "Displays more information about the 'create personal-token' subcommand"
	AskInputName           = "Enter the new Personal Token's name:"
	AskInputExpiration     = "Enter the new Personal Token's expiration:"
)
