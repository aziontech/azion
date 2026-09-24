package deploy

import "errors"

// Used by both the v3 and the v4 command trees.
var (
	ErrorCreateApplication = errors.New("Failed to create the Application: %s. Check your settings and try again. If the error persists, contact Azion support")
	ErrorUpdateApplication = errors.New("Failed to update the Application: %s. Check your settings and try again. If the error persists, contact Azion support")
	ErrorCreateDomain      = errors.New("Failed to create the Domain: %s. Check your settings and try again. If the error persists, contact Azion support")
	ErrorUpdateDomain      = errors.New("Failed to update the Domain: %s. Check your settings and try again. If the error persists, contact Azion support")
	ErrorInvalidToken      = errors.New("The configured token is invalid. You must create a new token and configure it to use with the CLI.")
	ErrorDeployRemote      = errors.New("Failed to read the response from remote deploy process. Please verify if your deploy finished successfully, and update your azion.json file, if necessary.")
	ErrorUnableSDKConfig   = "Unable to load SDK config, "
	ErrorUploadFileBucket  = "Failed to upload file to bucket %s: %w"
	ErrorCreateZip         = "Failed to create zip file %s: %w"
	ErrorRelPath           = "Error determining relative path for %s: %v"
	ErrorResetPointFile    = "Error resetting file pointer %s: %v"
	ErrorCopyContentFile   = "Error copying contents of file %s to ZIP: %v"
	ERRORMARSHALMANIFEST   = errors.New("Failed to marshal manifest structure.")
	ERRORWRITEMANIFEST     = errors.New("Failed to write manifest.json file.")
	ERRORCAPTURELOGS       = "Failed to capture deploy logs: %s. Please check your account to verify if the resources were successfully created."
)

// Used only by the v3 command tree (bug-fix-only; see doc/plan.md).
var (
	ErrorCodeFlag       = errors.New("Failed to read the code file. Verify if the file name and its path are correct and the file content has a valid code format")
	ErrorArgsFlag       = errors.New("Failed to read the args file. Verify if the file name and its path are correct and the file's content has a valid JSON format")
	ErrorParseArgs      = errors.New("Failed to parse JSON args. Verify if the file's content has a valid JSON format")
	ErrorCreateFunction = errors.New("Failed to create Function: %s. Check your settings and try again. If the error persists, contact Azion support")
	ErrorUpdateFunction = errors.New("Failed to update the Function: %s. Check your settings and try again. If the error persists, contact Azion support")
)
