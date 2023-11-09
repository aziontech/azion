package edgefunction

import "errors"

var (
	ErrorMandatoryCreateFlags            = errors.New("One or more required flags are missing. You must provide --active, --code, and --name flags when the --in flag is not provided. Run the command 'azion edge_functions create --help' to display more information and try again")
	ErrorActiveFlag                      = errors.New("Invalid --active flag provided. The flag must have 'true' or 'false' values. Run the command 'azion edge_functions <subcommand> --help' to display more information and try again")
	ErrorCodeFlag                        = errors.New("Failed to read the code file. Verify if the file name and its path are correct and the file content has a valid code format. Run the command 'azion edge_functions <subcommand> --help' to display more information and try again")
	ErrorArgsFlag                        = errors.New("Failed to read the args file. Verify if the file name and its path are correct and the file's content has a valid JSON format. Run the command 'azion edge_functions <subcommand> --help' to display more information and try again")
	ErrorParseArgs                       = errors.New("Failed to parse JSON args. Verify if the file's content has a valid JSON format. Run the command 'azion edge_functions <subcommand> --help' to display more information and try again")
	ErrorCreateFunction                  = errors.New("Failed to create edge function: %s. Check your settings and try again. If the error persists, contact Azion support")
	ErrorMissingFunctionIdArgument       = errors.New("A required flag is missing. You must provide a function_id as an argument or path to import the file. Run the command 'azion edge_functions <subcommand> --help' to display more information and try again")
	ErrorMissingFunctionIdArgumentDelete = errors.New("A required flag is missing. You must provide a function_id as an argument. Run the command 'azion edge_functions <subcommand> --help' to display more information and try again")
	ErrorFailToDeleteFunction            = errors.New("Failed to delete the Edge Function: %s. Check your settings and try again. If the error persists, contact Azion support")
	ErrorGetFunction                     = errors.New("Failed to get the Edge Function: %s. Check your settings and try again. If the error persists, contact Azion support")
	ErrorGetFunctions                    = errors.New("Failed to list the Edge Functions: %s. Check your settings and try again. If the error persists, contact Azion support")
	ErrorUpdateFunction                  = errors.New("Failed to update the Edge Function: %s. Check your settings and try again. If the error persists, contact Azion support")
	ErrorConvertIdFunction               = errors.New("The function ID you provided is invalid. The value must be an integer. You may run the 'azion list edge-function' command to check your function ID")
)
