package firewall

var (
	Usage                  = "firewall"
	CreateShortDescription = "Creates a Firewall"
	CreateLongDescription  = "Creates a Firewall to protect your applications from threats and attacks"
	FlagIn                 = "Path to a JSON file containing the attributes of the Firewall being created; you can use - for reading from stdin"
	CreateFlagHelp         = "Displays more information about the create firewall command"
	CreateOutputSuccess    = "Created Firewall with ID %d"

	FlagName              = "Firewall's name"
	FlagDebug             = "Allows you to check whether rules created using Rules Engine for Firewall have been successfully executed in your firewall"
	FlagActive            = "Whether the Firewall is active or not"
	FlagFunctionsEnabled  = "Whether the Firewall has Functions module enabled or not"
	FlagNetworkProtection = "Whether the Firewall has Network Layer Protection module enabled or not"
	FlagWafEnabled        = "Whether the Firewall has Web Application Firewall (WAF) module enabled or not"

	AskName   = "Enter the Firewall's name:"
	AskActive = "Is the Firewall active?"
)
