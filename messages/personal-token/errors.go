package personal_token

import "errors"

var (
	ErrorMandatoryCreateFlags    = errors.New("Required inputs are missing. You must provide name and expiry as flags or input the json structure with the name and expiry field when using the --in flag example from json:\"{'name': 'One day token', 'expires_at': '9m'}\". Run the command 'azion <command> <subcommand> --help' to display more information and try again.")
	ErrorCreate                  = errors.New("Failed to create the Personal Token: %s. Check your settings and try again. If the error persists, contact Azion support.")
	ErrorGet                     = errors.New("Failed to describe the personal token: %s. Check your settings and try again. If the error persists, contact Azion support.")
	ErrorMissingIDArgumentDelete = errors.New("A required flag is missing. You must provide the --id flag as an argument. Run the command 'azion personal_token <subcommand> --help' to display more information and try again")
	ErrorList                    = errors.New("Failed to list your personal tokens: %s. Check your settings and try again. If the error persists, contact Azion support.")
	ErrorFailToDelete            = errors.New("Failed to delete the personal token: %s. Check your settings and try again. If the error persists, contact Azion support")
)
