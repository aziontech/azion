package edgefunctions

import "errors"

var (
	ErrorMandatoryCreateFlags      = errors.New("Absence of mandatory flags. You must provide --active, --code, and --name flags when the --in flag is not provided. Run the command 'azion cli edge_functions create --help' to display more information and try again")
	ErrorActiveFlag                = errors.New("Invalid --active flag provided. The flag must have  'true' or 'false' values. Run the command 'azioncli edge_functions <subcommand> --help' to display more information and try again")
	ErrorCodeFlag                  = errors.New("Failed to read the code file. Verify if the file name and its path are correct and the file content has a valid code format. Run the command 'azioncli edge_functions <subcommand> --help' to display more information and try again")
	ErrorArgsFlag                  = errors.New("Failed to read the args file. Verify if the file name and its path are correct and the file's content has a valid JSON format. Run the command 'azioncli edge_functions <subcommand> --help' to display more information and try again")
	ErrorParseArgs                 = errors.New("Failed to parse JSON args. Verify if the file's content has a valid JSON format. Run the command 'azioncli edge_functions <subcommand> --help' to display more information and try again")
	ErrorCreateFunction            = errors.New("Failed to create edge function. Run the command 'azioncli edge_functions create --help' to display more information and try again")
	ErrorMissingFunctionIdArgument = errors.New("Absence of mandatory --function_id flag. You must provide a function_id to identify the function. Run the command 'azioncli edge_functions <subcommand> --help' to display more information and try again")
	ErrorMissingArgumentUpdate     = errors.New("Absence of mandatory --function_id flag or a path to import the file. You must provide a function_id to identify the function or a valid path and file to import. Run the command 'azioncli edge_functions <subcommand> --help' to display more information and try again")
	ErrorFailToDeleteFunction      = errors.New("Failed to delete an Edge Function based on its function_id. You must provide a valid --function_id flag of an existing function of the command. Run the command 'azioncli edge_functions delete --help' to display more information and try again")
	ErrorGetFunction               = errors.New("Failed to get the Edge Function based on its function_id. You must provide a valid --function_id flag of an existing function as an argument of the command. Run the command 'azioncli edge_functions <subcommand> --help' to display more information and try again")
	ErrorGetFunctions              = errors.New("Failed to get the list of Edge Functions. After a while, try again or contact Azion's support if the error persists")
	ErrorUpdateFunction            = errors.New("Failed to update the Edge Function based on its function_id. You must provide a valid --function_id flag of an existing function. Run the command 'azioncli edge_functions update --help' to display more information and try again")
)
