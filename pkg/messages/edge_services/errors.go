package edgeservices

import "errors"

var (
	ErrorMissingServiceIdArgument      = errors.New("A required --service_id flag is missing. You must provide a valid service_id. Run the command 'azion edge_services <subcommand> --help' to display more information and try again")
	ErrorMissingResourceIdArgument     = errors.New("One or more required flags are missing. You must provide a valid service_id and resource_id. Run the command 'azion edge_services <subcommand> --help' to display more information and try again")
	ErrorInvalidResourceTrigger        = errors.New("The trigger is invalid. You must provide a valid trigger. Run the command 'azion edge_services <subcommand> --help' to display more information and try again")
	ErrorUpdateNoFlagsSent             = errors.New("No values/flags sent during update. You must provide at least one valid value in the update. Run the command 'azion edge_services resources update --help' to display more information and try again")
	ErrorDeleteResource                = errors.New("Failed to delete the Resource: %s. Check your settings and try again. If the error persists, contact Azion support")
	ErrorGetResource                   = errors.New("Failed to get the Resource: %s. Check your settings and try again. If the error persists, contact Azion support")
	ErrorGetResources                  = errors.New("Failed to get the Resources: %s. Check your settings and try again. If the error persists, contact Azion support")
	ErrorInvalidNameFlag               = errors.New("Invalid Edge Service name. You must provide a valid Edge Service name with the flag --name. Run the command 'azion edge_services <subcommand> --help' to display more information and try again")
	ErrorInvalidTriggerFlag            = errors.New("The trigger flag is invalid. You must provide a valid flag --trigger value. Run the command 'azion edge_services resources update --help' to display more information and try again")
	ErrorInvalidContentTypeFlag        = errors.New("The resource content type is invalid. You must provide a valid flag --content-type with value <shellscript|text>. Run the command 'azion edge_services resources <subcommand> --help' to display more information and try again")
	ErrorUpdateResource                = errors.New("Failed to update the Resource: %s. Check your settings and try again. If the error persists, contact Azion support")
	ErrorCreateResource                = errors.New("Failed to create the Resource: %s. Check your settings and try again. If the error persists, contact Azion support")
	ErrorGetServices                   = errors.New("Failed to get the Edge Services: %s. Check your settings and try again. If the error persists, contact Azion support")
	ErrorGetSerivce                    = errors.New("Failed to get the Edge Service: %s. Check your settings and try again. If the error persists, contact Azion support")
	ErrorWithVariablesFlag             = errors.New("Failed to process --with-variables flag. Try again and, if the error persists, please contact Azion support")
	ErrorDeleteService                 = errors.New("Failed to delete Edge Service: %s. Check your settings and try again. If the error persists, contact Azion support")
	ErrorCreateService                 = errors.New("Failed to create the Edge Service: %s. Check your settings and try again. If the error persists, contact Azion support")
	ErrorUpdateService                 = errors.New("Failed to update the Edge Service: %s. Check your settings and try again. If the error persists, contact Azion support")
	ErrorMandatoryName                 = errors.New("A required flag is missing. You must provide --name flag when --in flag is not sent. Run the command 'azion edge_services <subcommand> --help' to display more information and try again")
	ErrorMandatoryFlagsResource        = errors.New("One or more required flags are missing. You must provide --name, --content-type, and --content-file flags when the --in flag is not sent. Run the command 'azion edge_services <subcommand> --help' to display more information and try again")
	ErrorMissingArgumentUpdate         = errors.New("A mandatory flag or file is missing. You must provide a service_id as an argument or path to import the file. Run the command 'azion edge_services <subcommand> --help' to display more information and try again")
	ErrorMissingArgumentUpdateResource = errors.New("One or more required flags or a file is missing. You must provide a service_id and a resource_id as an argument or path to import the file. Run the command 'azion edge_services <subcommand> --help' to display more information and try againx")
)
