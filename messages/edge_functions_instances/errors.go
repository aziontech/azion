package edge_functions_instances

import "errors"

var (
	ErrorGetFunctions           = errors.New("Failed to get the edge functions instances: %s. Check your settings and try again. If the error persists, contact Azion support")
	ErrorMissingArgumentsDelete = errors.New("Required flags are missing. You must supply application-id and instance-id as arguments. Run 'azioncli <command> <subcommand> --help' command to display more information and try again")
	ErrorFailToDeleteFuncInst   = errors.New("Failed to delete the edge functions instance: %s. Check your settings and try again. If the error persists, contact Azion support")
	ErrorMandatoryCreateFlags   = errors.New("Required flags are missing. You must provide the application-id, edge-function-id, and name flags when the --application-id and --in flag are not provided. Run the command 'azioncli <command> <subcommand> --help' to display more information and try again.")
	ErrorMandatoryListFlags     = errors.New("A required flag is missing. You must provide application-id. Run the command 'azioncli <command> <subcommand> --help' to display more information and try again.")
	ErrorCreate                 = errors.New("Failed to create the edge functions instances: %s. Check your settings and try again. If the error persists, contact Azion support.")
	ErrorMandatoryFlags         = errors.New("One or more required flags are missing. You must provide the --application-id and --instance-id flags. Run the command 'azioncli <command> <subcommand> --help' to display more information and try again.")
	ErrorGetEdgeFuncInstances   = errors.New("Failed to describe the edge functions instance: %s. Check your settings and try again. If the error persists, contact Azion support.")
	ErrorUpdateFuncInstance     = errors.New("Failed to update the edge functions instance: %s. Check your settings and try again. If the error persists, contact Azion support.")
	ErrorMandatoryUpdateFlags   = errors.New("Required flags are missing. You must provide the application-id, instance-id, and function-id flags when the --in flag is not provided. Run the command 'azioncli <command> <subcommand> --help' to display more information and try again.")
	ErrorMandatoryUpdateFlagsIn = errors.New("Required flags are missing. You must provide the application-id and instance-id flags when the --in flag is provided. Run the command 'azioncli <command> <subcommand> --help' to display more information and try again.")
)
