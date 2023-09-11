package edge_application

import "errors"

var (
	ErrorCreate               = errors.New("Failed to create the Edge Application: %s. Check your settings and try again. If the error persists, contact Azion support")
	ErrorMandatoryCreateFlags = errors.New("Required inputs are missing. You must provide name as flags or input the json structure with the name and expiry field when using the --in flag example from json:\"{'name': 'One day token'}\". Run the command 'azion <command> <subcommand> --help' to display more information and try again.")
)
