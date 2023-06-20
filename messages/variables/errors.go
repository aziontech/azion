package variables

import "errors"

var (
	ErrorGetVariables                    = errors.New("Failed to describe the origins: %s. Check your settings and try again. If the error persists, contact Azion support.")
	ErrorGetItem                         = errors.New("Failed to describe the variable: %s. Check your settings and try again. If the error persists, contact Azion support.")
	ErrorMissingArguments                = errors.New("Required flags are missing. You must supply application-id and variable-id as arguments. Run 'azioncli <command> <subcommand> --help' command to display more information and try again")
	ErrorFailToDeleteVariable            = errors.New("Failed to delete the Variable: %s. Check your settings and try again. If the error persists, contact Azion support")
	ErrorMissingVariableIdArgumentDelete = errors.New("A mandatory flag is missing. You must provide a variable_id as an argument. Run the command 'azioncli variables <subcommand> --help' to display more information and try again")
	ErrorMissingVariableIdArgument       = errors.New("A required flag is missing. You must provide variable-id, key, value and secret flags as an argument or path to import the file. Run the command 'azioncli variables <subcommand> --help' to display more information and try again")
	ErrorSecretFlag                      = errors.New("Invalid --secret flag provided. The flag must have 'true' or 'false' values. Run the command 'azioncli variables <subcommand> --help' to display more information and try again")
	ErrorUpdateVariable                  = errors.New("Failed to update the Variable: %s. Check your settings and try again. If the error persists, contact Azion support")
)
