package application

import "errors"

// Used by the v3, v4 and v6 command trees. The cascade messages are v3 and v4
// only: v6 has no cascade delete until its deploy model is known.
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
