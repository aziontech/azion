package device_groups

import (
	"errors"
)

var (
	ErrorMissingApplicationIDArgument = errors.New("aaA mandatory flag is missing. You must provide a application-id as an argument or path to import the file. Run the command 'azioncli domains <subcommand> --help' to display more information and try again")
	ErrorGetDeviceGroups              = errors.New("Failed to describe the device groups: %s. Check your settings and try again. If the error persists, contact Azion support.")
	ErrorMandatoryFlags               = errors.New("One or more required flags are missing. You must provide the --application-id and --group-id flags. Run the command 'azioncli rules_engine <subcommand> --help' to display more information and try again.")
	ErrorFailToDelete                 = errors.New("Failed to delete the rule in Rules Engine: %s. Check your settings and try again. If the error persists, contact Azion support.")

	ErrorMandatoryFlagsUpdate         = errors.New("One or more required flags are missing. You must provide the --application-id and --group-id flags when --in flag is not sent. Run the command 'azioncli device_groups <subcommand> --help' to display more information and try again.")
	ErrorUpdateDeviceGroups           = errors.New("Failed to update the Device Group: %s. Check your settings and try again. If the error persists, contact Azion support")

	ErrorMandatoryCreateFlags         = errors.New("Required flags are missing. You must provide application-id, name, user-agent flags when the --application-id and --in flag are not provided. Run the command 'azioncli <command> <subcommand> --help' to display more information and try again.")
	ErrorCreateDeviceGroups           = errors.New("Failed to create the Device groups: %s. Check your settings and try again. If the error persists, contact Azion support.")
	ErrorListDeviceGroups             = errors.New("Failed to list your rules in Rules Engine: %s. Check your settings and try again. If the error persists, contact Azion support.")

)
