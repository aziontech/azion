package origins

import "errors"

var (
	ErrorGetOrigins                   = errors.New("Failed to describe the origins. Check your settings and try again. If the error persists, contact Azion support.")
	ErrorMissingApplicationIDArgument = errors.New("A required flag is missing. You must supply an application_id as an argument. Run 'azioncli origins list --help' command to display more information and try again")
)
