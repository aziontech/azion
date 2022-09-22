package edgeservices

import "errors"

var (
	ErrorMissingServiceIdArgument      = errors.New("You must provide a service_id. Use -h or --help for more information")
	ErrorMissingResourceIdArgument     = errors.New("You must provide a service_id and a resource_id. Use -h or --help for more information")
	ErrorInvalidVariablesFileFormat    = errors.New("You must provide a valid variables file content")
	ErrorInvalidResourceTrigger        = errors.New("You must provide a valid trigger")
	ErrorUpdateNoFlagsSent             = errors.New("You must provide at least one value in update")
	ErrorDeleteResource                = errors.New("Failed to delete Resource")
	ErrorGetResource                   = errors.New("Failed to get Resource")
	ErrorGetResources                  = errors.New("Failed to get Resources")
	ErrorInvalidNameFlag               = errors.New("Invalid --name flag")
	ErrorInvalidTriggerFlag            = errors.New("Invalid --trigger flag")
	ErrorInvalidContentTypeFlag        = errors.New("Invalid --content-type flag")
	ErrorUpdateResource                = errors.New("Failed to update Resource")
	ErrorCreateResource                = errors.New("Failed to create Resource")
	ErrorGetServices                   = errors.New("Failed to get Edge Services")
	ErrorGetSerivce                    = errors.New("Failed to get Edge Service")
	ErrorWithVariablesFlag             = errors.New("Invalid --with-variables flag")
	ErrorDeleteService                 = errors.New("Failed to delete Edge Service")
	ErrorCreateService                 = errors.New("Failed to create Edge Service")
	ErrorUpdateService                 = errors.New("Failed to update Edge Service")
	ErrorMandatoryName                 = errors.New("You must provide --name flag when --in flag is not sent")
	ErrorMandatoryFlagsResource        = errors.New("You must provide --name, --content-type and --content-file flags when --in flag is not sent")
	ErrorMissingArgumentUpdate         = errors.New("You must provide a service_id as an argument or path to import file")
	ErrorMissingArgumentUpdateResource = errors.New("You must provide a service_id and a resource_id as an argument or path to import file")
)
