package deploy

import "errors"

var (
	ErrorOpeningAzionFile  = errors.New("Failed to open the azion.json file. The file doesn't exist, is corrupted, or has an invalid JSON format. Verify if the file format is JSON or fix its content according to the JSON format specification at https://www.json.org/json-en.html")
	ErrorCodeFlag          = errors.New("Failed to read the code file. Verify if the file name and its path are correct and the file content has a valid code format")
	ErrorArgsFlag          = errors.New("Failed to read the args file. Verify if the file name and its path are correct and the file's content has a valid JSON format")
	ErrorParseArgs         = errors.New("Failed to parse JSON args. Verify if the file's content has a valid JSON format")
	ErrorCreateFunction    = errors.New("Failed to create Edge Function: %s. Check your settings and try again. If the error persists, contact Azion support")
	ErrorUpdateFunction    = errors.New("Failed to update the Edge Function: %s. Check your settings and try again. If the error persists, contact Azion support")
	ErrorCreateApplication = errors.New("Failed to create the Edge Application: %s. Check your settings and try again. If the error persists, contact Azion support")
	ErrorUpdateApplication = errors.New("Failed to update the Edge Application: %s. Check your settings and try again. If the error persists, contact Azion support")
	ErrorCreateInstance    = errors.New("Failed to create the Edge Function Instance: %s. Check your settings and try again. If the error persists, contact Azion support")
	ErrorCreateDomain      = errors.New("Failed to create the Domain: %s. Check your settings and try again. If the error persists, contact Azion support")
	ErrorUpdateDomain      = errors.New("Failed to update the Domain: %s. Check your settings and try again. If the error persists, contact Azion support")
	ErrorInvalidToken      = errors.New("The configured token is invalid. You must create a new token and configure it to use with the CLI.")
)
