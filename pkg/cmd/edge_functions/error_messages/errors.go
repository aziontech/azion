package errormessages

import "errors"

var (
	ErrorMandatoryCreateFlags      = errors.New("You must provide --active, --code and --name flags when --in flag is not sent")
	ErrorActiveFlag                = errors.New("Invalid --active flag")
	ErrorCodeFlag                  = errors.New("Failed to read code file")
	ErrorArgsFlag                  = errors.New("Failed to read args file")
	ErrorParseArgs                 = errors.New("Failed to parse json args")
	ErrorCreateFunction            = errors.New("Failed to create edge function")
	ErrorMissingFunctionIdArgument = errors.New("You must provide a function_id. Use -h or --help for more information")
	ErrorMissingArgumentUpdate     = errors.New("You must provide a function_id as an argument or path to import file")
	ErrorFailToDeleteFunction      = errors.New("Failed to delete Edge Function")
	ErrorGetFunction               = errors.New("Failed to get Edge Function")
	ErrorGetFunctions              = errors.New("Failed to get Edge Functions")
	ErrorUpdateFunction            = errors.New("Failed to update Edge Function")
	ErrorPurgeDomainCache          = errors.New("Could not purge domain cache")
)
