package edgeapplication

import "errors"

var (
	ErrorMissingAzionJson             = errors.New("Azion.json file is missing. Please initialize and deploy your project before using cascade delete")
	ErrorMissingApplicationIdJson     = errors.New("Application ID is missing from azion.json. Please initialize and publish your project first before using cascade delete")
	ErrorFailToDeleteApplication      = errors.New("Failed to delete the Edge Application: %s. Check your settings and try again. If the error persists, contact Azion support")
	ErrorMissingApplicationIdArgument = errors.New("A required flag is missing. You must provide an application_id as an argument or path to import the file. Run the command 'azion list edge-application' to retrieve the specific ID and try again")
	ErrorFailedUpdateAzionJson        = errors.New("Failed to update azion.json file to remove IDs of deleted resource")
	ErrorConvertId                    = errors.New("The application ID you provided is invalid. The value must be an integer. You may run the 'azion list edge-application' command to check your application ID")
)
