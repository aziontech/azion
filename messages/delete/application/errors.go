package application

import "errors"

// Used by both the v3 and the v4 command trees.
var (
	ErrorMissingAzionJson         = errors.New("Azion.json file is missing. Please initialize and deploy your project before using cascade delete")
	ErrorMissingApplicationIdJson = errors.New("Application ID is missing from azion.json. Please initialize and deploy your project before using cascade delete")
	ErrorFailToDeleteApplication  = errors.New("Failed to delete the Application: %s. Check your settings and try again. If the error persists, contact Azion support")
	ErrorConvertId                = errors.New("The application ID you provided is invalid. The value must be an integer. You may run the 'azion list application' command to check your application ID")
)

// Used only by the v3 command tree (bug-fix-only; see doc/plan.md).
var (
	ErrorFailedUpdateAzionJson = errors.New("Failed to update azion.json file to remove IDs of deleted resource")
)
