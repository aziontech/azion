package edge_applications

import "errors"

var (
	ErrorMissingApplicationIdArgument = errors.New("A required flag is missing. You must provide an application_id as an argument or path to import the file. Run the command 'azioncli edge_applications <subcommand> --help' to display more information and try again")
	ErrorGetApplication               = errors.New("Failed to get the Edge Application: %s. Check your settings and try again. If the error persists, contact Azion support")
)
