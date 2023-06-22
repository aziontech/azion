package variables

import "errors"

var (
	ErrorGetItem                         = errors.New("Failed to describe the variable: %s. Check your settings and try again. If the error persists, contact Azion support.")
	ErrorMissingArguments                = errors.New("Required flags are missing. You must supply the --application-id and --variable-id flags as arguments. Run 'azioncli <command> <subcommand> --help' command to display more information and try again")
	ErrorFailToDeleteVariable            = errors.New("Failed to delete the variable: %s. Check your settings and try again. If the error persists, contact Azion support")
	ErrorMissingVariableIdArgumentDelete = errors.New("A mandatory flag is missing. You must provide the --variable_id flag as an argument. Run the command 'azioncli variables <subcommand> --help' to display more information and try again")
	ErrorMissingVariableIdArgument       = errors.New("Required flags are missing. You must provide the --variable-id, --key, --value, and --secret flags as arguments, or the --in flag informing the path to import the file. Run the command 'azioncli variables <subcommand> --help' to display more information and try again")
	ErrorSecretFlag                      = errors.New("Invalid --secret flag provided. The value must be 'true' or 'false'. Run the command 'azioncli variables <subcommand> --help' to display more information and try again")
	ErrorUpdateVariable                  = errors.New("Failed to update the variable: %s. Check your settings and try again. If the error persists, contact Azion support")
	ErrorMandatoryCreateFlags            = errors.New("One or more required flags are missing. You must provide the --key, --name, and --secret flags when the --in flag is not provided. Run the command 'azioncli variables create --help' to display more information and try again")
	ErrorCreateItem                      = errors.New("Failed to create the variable: %s. Check your settings and try again. If the error persists, contact Azion support.")
)
