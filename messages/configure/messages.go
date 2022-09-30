package configure

var (
	ConfigureUsage            = "configure"
	ConfigureShortDescription = "Authorizes connections with Azion platform’s services"
	ConfigureLongDescription  = "Sets up CLI parameters and an authentication token to enable connections with Azion platform’s services and to run CLI commands"
	ConfigureFlagToken        = "Saves a given personal token locally to authorize CLI commands"
	TokenSavedIn              = "Token saved in %v\n"
	TokenUsedIn               = "This token will be used by default when calling any command"
	ConfigureHelpFlag         = "Displays more information about the configure command"
)
