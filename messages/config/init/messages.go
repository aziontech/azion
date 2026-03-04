package init

const (
	Usage            = "init"
	ShortDescription = "Initialize Azion configuration file"
	LongDescription  = "Create an azion.json file in the current directory to store your application state"
	FlagHelp         = "Displays more information about the init command"
	FlagConfigDir    = "Path to the directory where azion.json will be created (default: current directory)"
	InitSuccessful   = "Configuration file created successfully: %s\n"
	ConfigExists     = "Configuration file already exists: %s\n"
	CreatingConfig   = "Creating azion.json file\n"
	DocsURL          = "For more information on azion.config, visit: https://github.com/aziontech/lib?tab=readme-ov-file#config/\n"
)
