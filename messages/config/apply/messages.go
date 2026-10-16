package apply

// Used only by the v4 command tree.
const (
	Usage               = "apply"
	ShortDescription    = "Apply configuration from configuration file to Azion Platform"
	LongDescription     = "Apply configuration resources defined in azion.config file to Azion Platform"
	FlagHelp            = "Displays more information about the apply command"
	FlagConfigFile      = "Path to the azion.config file to apply. If not sent, the azion.config.* file in the current directory is used"
	FlagConfigDir       = "Path to the configuration directory containing azion.json and azion.config (default: current directory)"
	ApplySuccessful     = "Configuration applied successfully: %d resource(s) applied\n"
	NoResourcesToApply  = "No resources found in manifest to apply\n"
	CreatingAzionJson   = "Creating azion.json file\n"
	AzionConfigNotFound = "azion.config file not found. Please create an azion.config file to define your application configuration before running 'azion config apply'\n"
)
