package utils

import "errors"

var (
	GenericUseHelp                 = errors.New("Use -h or --help for more information")
	ErrorMissingServiceIdArgument  = errors.New("You must provide a service_id as an argument. Use -h or --help for more information")
	ErrorMissingResourceIdArgument = errors.New("You must provide a service_id and a resource_id as arguments. Use -h or --help for more information")
	ErrorConvertingIdArgumentToInt = errors.New("You must provide a valid id")
	ErrorConvertingStringToBool    = errors.New("You must provide a valid value. Use -h or --help for more information")
	ErrorInternalServerError       = errors.New("Something went wrong, please try again")
	ErrorInvalidResourceTrigger    = errors.New("You musr provide a velid trigger. Use -h or --help for more information")
)
