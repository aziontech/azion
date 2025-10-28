package profile

var (
	UsageCreate = "profile"

	CreateShortDescription = "Create a new profile"
	CreateLongDescription  = "Create a new profile setting up all the required fields"
	CreateFlagHelp         = "Displays more information about the create profile subcommand"
	FlagName               = "Name of the new profile"
	FlagToken              = "Token for the new profile"
	FlagFile               = "Path to the toml file containing the settings for the new profile"
	QuestionToken          = "Would you like to set a token for the new profile? (Y/n)"
	QuestionProvideToken   = "Would you like to provide a token for the new profile? If answer is no, the CLI will create a new token for you (Y/n)"
	FieldToken             = "Please inform a token for the new profile"
	FieldProfileName       = "Please inform a name for the new profile"
	QuestionCollectMetrics = "To better understand user needs and enhance our application, we gather anonymous data. Do you agree to participate? (Y/n)"
	CreateOutputSuccess    = "Profile '%s' created successfully"

	UsageProfiles            = "profiles"
	ProfilesShortDescription = "Manage profiles"
	ProfilesLongDescription  = "Manage profiles that you have configured"
	ProfilesFlagHelp         = "Displays more information about the profiles command"
	SwitchSuccessful         = "Profile switched successfully"

	UsageDelete            = "profile"
	DeleteShortDescription = "Delete a profile"
	DeleteLongDescription  = "Delete a profile and all its associated data"
	DeleteFlagHelp         = "Displays more information about the delete profile subcommand"
	DeleteOutputSuccess    = "Profile '%s' deleted successfully"
	QuestionDeleteProfile  = "Choose a profile to delete:"
	ConfirmDeleteProfile   = "Are you sure you want to delete profile '%s'? This action cannot be undone (Y/n)"
	WarningDeleteToken     = "Warning: Failed to delete token from server: %v"
	WarningSetActiveProfile = "Warning: Failed to set active profile to default: %v"
)
