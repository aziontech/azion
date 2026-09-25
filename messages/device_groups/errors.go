package device_groups

import (
	"errors"
)

// Used only by the v4 command tree.
var (
	ErrorGetDeviceGroups = errors.New("Failed to describe the device groups: %s. Check your settings and try again. If the error persists, contact Azion support.")
	ErrorFailToDelete    = errors.New("Failed to delete the device group: %s. Check your settings and try again. If the error persists, contact Azion support.")

	ErrorUpdateDeviceGroups = errors.New("Failed to update the device group: %s. Check your settings and try again. If the error persists, contact Azion support")

	ErrorCreateDeviceGroups = errors.New("Failed to create the device group: %s. Check your settings and try again. If the error persists, contact Azion support.")
	ErrorListDeviceGroups   = errors.New("Failed to list your device groups: %s. Check your settings and try again. If the error persists, contact Azion support.")

	ErrorConvertIdApplication = errors.New("The Application ID you provided is invalid. The value must be an integer. You may run the 'azion list application' command to check your Applications' IDs")
	ErrorConvertIdDeviceGroup = errors.New("The device group ID you provided is invalid. The value must be an integer. You may run the 'azion list device-group' command to check your device groups' IDs")
)
