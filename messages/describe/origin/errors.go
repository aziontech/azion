package origin

import "errors"

var (
	ErrorGetOrigin     = errors.New("Failed to describe the origin: %s. Check your settings and try again. If the error persists, contact Azion support.")
	ErrorFormatOut     = errors.New("The server failed formatting data for display. Repeat the HTTP request and check the HTTP response's format")
	ErrorWriteFile     = errors.New("The file is read-only and/or isn't accessible. Change the attributes of the file to read and write and/or give access to it")
	OriginsFileWritten = "File successfully written to: %s\n"
)
