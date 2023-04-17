package device_groups

import (
	"errors"
)

var (
	ErrorMissingApplicationIDArgument = errors.New("A mandatory flag is missing. You must provide a application-id as an argument or path to import the file. Run the command 'azioncli domains <subcommand> --help' to display more information and try again")
	ErrorGetDeviceGroups              = errors.New("Failed to describe the device groups: %s. Check your settings and try again. If the error persists, contact Azion support.")
)
