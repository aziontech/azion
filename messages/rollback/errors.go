package rollback

import "errors"

var (
	ERRORROLLBACK    = errors.New("Failed to roll back to previous static files")
	ERRORNEEDSDEPLOY = errors.New("You cannot use the rollback command unless you have already deployed this project. Please check if you are in the correct working directory")
	ERRORAZION       = errors.New("Failed to open the azion.json file. The file doesn't exist, is corrupted, or has an invalid JSON format")
	ERRORNOPREVIOUS  = errors.New("No previous deployment found to roll back to. You need at least two deployments to use the rollback command")
)
