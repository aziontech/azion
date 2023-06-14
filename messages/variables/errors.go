package variables

import "errors"

var (
	ErrorFailToDeleteVariable = errors.New("Failed to delete the Variable: %s. Check your settings and try again. If the error persists, contact Azion support")
	ErrorMissingVariableIdArgumentDelete = errors.New("A mandatory flag is missing. You must provide a variable_id as an argument. Run the command 'azioncli domains <subcommand> --help' to display more information and try again")

)