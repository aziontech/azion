package variables

import "errors"

var (
	ErrorGetItem                         = errors.New("Failed to describe the variable: %s. Check your settings and try again. If the error persists, contact Azion support.")
	ErrorMissingArguments                = errors.New("A required flag is missing. You must supply the --variable-id flag as an argument. Run 'azion <command> <subcommand> --help' command to display more information and try again")
	ErrorFailToDeleteVariable            = errors.New("Failed to delete the variable: %s. Check your settings and try again. If the error persists, contact Azion support")
	ErrorMissingVariableIdArgumentDelete = errors.New("A required flag is missing. You must provide the --variable_id flag as an argument. Run the command 'azion variables <subcommand> --help' to display more information and try again")
	ErrorMissingVariableIdArgument       = errors.New("Required lags are missing. You must provide the --variable-id, --key, --value, and --secret flags as arguments, or the --file flag informing the path to import the file. Run the command 'azion variables <subcommand> --help' to display more information and try again")
	ErrorMissingFieldUpdateVariables     = errors.New("Required flags are missing. You must provide the --key, --value, and --secret flags as arguments, or the --file flag informing the path to import the file. Run the command 'azion variables <subcommand> --help' to display more information and try again")
	ErrorSecretFlag                      = errors.New("Invalid --secret flag provided. The value must be 'true' or 'false'. Run the command 'azion variables <subcommand> --help' to display more information and try again")
	ErrorUpdateVariable                  = errors.New("Failed to update the variable: %s. Check your settings and try again. If the error persists, contact Azion support")
	ErrorCreateItem                      = errors.New("Failed to create the variable: %s. Check your settings and try again. If the error persists, contact Azion support.")
)
