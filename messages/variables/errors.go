package variables

import "errors"

var (
	ErrorGetItem                         = errors.New("Failed to describe the variable: %s. Check your settings and try again. If the error persists, contact Azion support.")
	ErrorMissingArguments                = errors.New("Required flags are missing. You must supply application-id and variable-id as arguments. Run 'azioncli <command> <subcommand> --help' command to display more information and try again")
	ErrorGetVariables                    = errors.New("Failed to describe the variables: %s. Check your settings and try again. If the error persists, contact Azion support.")
	ErrorFailToDeleteVariable            = errors.New("Failed to delete the Variable: %s. Check your settings and try again. If the error persists, contact Azion support")
	ErrorMissingVariableIdArgumentDelete = errors.New("A mandatory flag is missing. You must provide a variable_id as an argument. Run the command 'azioncli variables <subcommand> --help' to display more information and try again")
	ErrorMandatoryCreateFlags            = errors.New("One or more required flags are missing. You must provide --key, --name, and --secret flags when the --in flag is not provided. Run the command 'azioncli variables create --help' to display more information and try again")
	ErrorCreateItem                      = errors.New("Failed to create the variable: %s. Check your settings and try again. If the error persists, contact Azion support.")
)
