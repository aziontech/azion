package edge_application

import "errors"

var (
	ErrorCreate               = errors.New("Failed to create the edge application: %s. Check your settings and try again. If the error persists, contact Azion support")
	ErrorMandatoryCreateFlags = errors.New("Required inputs are missing. You must provide a name or the --in flag followed by the filepath with the settings. Run the command 'azion create edge-application --help' to display more information and try again.")
)
