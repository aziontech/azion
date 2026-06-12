package device_groups

import (
	"errors"
)

var (
	ErrorMissingApplicationIDArgument = errors.New("A mandatory flag is missing. You must provide a application-id as an argument or path to import the file. Run the command 'azion device-group <subcommand> --help' to display more information and try again")
	ErrorGetDeviceGroups              = errors.New("Failed to describe the device groups: %s. Check your settings and try again. If the error persists, contact Azion support.")
	ErrorMandatoryFlags               = errors.New("One or more required flags are missing. You must provide the --application-id and --group-id flags. Run the command 'azion device-group <subcommand> --help' to display more information and try again.")
	ErrorFailToDelete                 = errors.New("Failed to delete the device group: %s. Check your settings and try again. If the error persists, contact Azion support.")

	ErrorMandatoryFlagsUpdate = errors.New("One or more required flags are missing. You must provide the --application-id and --group-id flags when --in flag is not sent. Run the command 'azion device-group <subcommand> --help' to display more information and try again.")
	ErrorUpdateDeviceGroups   = errors.New("Failed to update the device group: %s. Check your settings and try again. If the error persists, contact Azion support")

	ErrorMandatoryCreateFlags = errors.New("Required flags are missing. You must provide the application-id, name, and user-agent flags when the --application-id and --in flags are not provided. Run the command 'azion device-group <subcommand> --help' to display more information and try again.")
	ErrorCreateDeviceGroups   = errors.New("Failed to create the device group: %s. Check your settings and try again. If the error persists, contact Azion support.")
	ErrorListDeviceGroups     = errors.New("Failed to list your device groups: %s. Check your settings and try again. If the error persists, contact Azion support.")

	ErrorConvertIdApplication = errors.New("The Application ID you provided is invalid. The value must be an integer. You may run the 'azion list application' command to check your Applications' IDs")
	ErrorConvertIdDeviceGroup = errors.New("The device group ID you provided is invalid. The value must be an integer. You may run the 'azion list device-group' command to check your device groups' IDs")
)
