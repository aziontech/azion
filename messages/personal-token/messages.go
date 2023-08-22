package personal_token

var (
	// [ personal_token ]
	Usage            = "personal_token"
	ShortDescription = "The \"Personal Token\" command is a security and authentication feature that allows users to generate unique individual tokens"
	LongDescription  = "The \"Personal Token\" command is a security and authentication feature that allows users to generate unique individual tokens. These tokens are used to authenticate and authorize actions in systems"
	FlagHelp         = "Displays more information about the origins command"

	// [ create ]
	CreateUsage            = "create [flags]"
	CreateShortDescription = "Creates a new personal token"
	CreateLongDescription  = "Creates an personal token"
	CreateFlagName         = "The flag name"
	CreateFlagExpiresAt    = "The flag expires at"
	CreateFlagDescription  = "The flag description"
	CreateFlagIn           = "Path to a JSON file containing the attributes of the personal token that will be created; you can use - for reading from stdin"
	CreateOutputSuccess    = "Created personal token with ID %s\n"
	CreateHelpFlag         = "Displays more information about the create subcommand"

	// [ list ]
	ListUsage            = "list [flags]"
	ListShortDescription = "Displays your personal token in a list"
	ListLongDescription  = "Displays all your personal token in a list"
	ListHelpFlag         = "Displays more information about the list subcommand"
)
