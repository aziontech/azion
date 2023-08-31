package personal_token

var (
	// [ personal_token ]
	Usage            = "personal_token"
	ShortDescription = "The \"Personal Token\" command is a security and authentication feature that allows users to generate unique individual tokens"
	LongDescription  = "The \"Personal Token\" command is a security and authentication feature that allows users to generate unique individual tokens. These tokens are used to authenticate and authorize actions in systems"
	FlagHelp         = "Displays more information about the personal_token command"

	// [ create ]
	CreateUsage            = "create [flags]"
	CreateShortDescription = "Creates a new personal token"
	CreateLongDescription  = "Creates a personal token to be used for authentication and security"
	CreateFlagName         = "The personal token's name. It's required if the --in flag is not informed."
	CreateFlagExpiresAt    = "The personal token's expiration. It's required if the --in flag is not informed."
	CreateFlagDescription  = "The personal token's description"
	CreateFlagIn           = "Path to a JSON file containing the attributes of the personal token being created; you can use - for reading from stdin"
	CreateOutputSuccess    = "Created personal token with ID %s\n"
	CreateHelpFlag         = "Displays more information about the create subcommand"

	// [ list ]
	ListUsage            = "list [flags]"
	ListShortDescription = "Displays your personal tokens in a list"
	ListLongDescription  = "Displays all your personal token in a list"
	ListHelpFlag         = "Displays more information about the list subcommand"

	// [ delete ]
	DeleteOutputSuccess    = "Personal token %v was successfully deleted\n"
	DeleteHelpFlag         = "Displays more information about the delete subcommand"
	DeleteUsage            = "delete [flags]"
	DeleteShortDescription = "Deletes a personal token"
	DeleteLongDescription  = "Deletes a personal token based on its UUID"
	FlagID                 = "Unique identifier for a personal token. The '--id' flag is required"
)
