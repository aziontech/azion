package function

import "errors"

// Used by both the v3 and the v4 command trees.
var (
	ErrorActiveFlag           = errors.New("Invalid --active flag provided. The flag must have 'true' or 'false' values. Run the command again using the --help flag to display more information and try again")
	ErrorCodeFlag             = errors.New("Failed to read the code file. Verify if the file name and its path are correct and the file content has a valid code format. Run the command again using the --help flag to display more information and try again")
	ErrorArgsFlag             = errors.New("Failed to read the args file. Verify if the file name and its path are correct and the file's content has a valid JSON format. Run the command again using the --help flag to display more information and try again")
	ErrorParseArgs            = errors.New("Failed to parse JSON args. Verify if the file's content has a valid JSON format. Run the command again using the --help flag to display more information and try again")
	ErrorCreateFunction       = errors.New("Failed to create function: %s. Check your settings and try again. If the error persists, contact Azion support")
	ErrorFailToDeleteFunction = errors.New("Failed to delete the Edge Function: %s. Check your settings and try again. If the error persists, contact Azion support")
	ErrorGetFunction          = errors.New("Failed to get the Edge Function: %s. Check your settings and try again. If the error persists, contact Azion support")
	ErrorGetFunctions         = errors.New("Failed to list the Edge Functions: %s. Check your settings and try again. If the error persists, contact Azion support")
	ErrorUpdateFunction       = errors.New("Failed to update the Edge Function: %s. Check your settings and try again. If the error persists, contact Azion support")
	ErrorConvertIdFunction    = errors.New("The function ID you provided is invalid. The value must be an integer. You may run the 'azion list function' command to check your function ID")
)

// Used only by the v4 command tree.
var (
	ErrorConvertFunctionId = errors.New("Invalid --function-id flag provided. The value must be an integer. Run the command 'azion delete function --help' to display more information and try again")
)
