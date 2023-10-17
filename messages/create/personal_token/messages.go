package personaltoken

var (
	CreateUsage            = "personal-token"
	CreateShortDescription = "Creates a new personal token"
	CreateLongDescription  = "Creates a personal token to be used for authentication and security"
	CreateFlagName         = "The personal token's name"
	CreateFlagExpiresAt    = "The personal token's expiration"
	CreateFlagDescription  = "The personal token's description"
	CreateFlagIn           = "Path to a JSON file containing the attributes of the personal token being created; you can use - for reading from stdin"
	CreateOutputSuccess    = "Created personal token: %s\n"
	CreateHelpFlag         = "Displays more information about the 'create personal-token' subcommand"
	AskInputName           = "What is the name of the Personal Token?"
	AskInputExpiration     = "What is the expiration of the Personal Token?"
)
