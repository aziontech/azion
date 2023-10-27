package origin

import "errors"

var (
	ErrorMissingArguments = errors.New("Required flags are missing. You must supply application-id and origin-id as arguments. Run 'azion <command> <subcommand> --help' command to display more information and try again")
	ErrorGetOrigin        = errors.New("Failed to describe the origin: %s. Check your settings and try again. If the error persists, contact Azion support.")
	ErrorFormatOut        = errors.New("The server failed formatting data for display. Repeat the HTTP request and check the HTTP response's format")
	ErrorWriteFile        = errors.New("The file is read-only and/or isn't accessible. Change the attributes of the file to read and write and/or give access to it")
	OriginsFileWritten    = "File successfully written to: %s\n"
)
