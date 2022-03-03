package utils

import "errors"

var (
	//Generic errors that can be used by any package
	GenericUseHelp                  = errors.New("Use -h or --help for more information")
	ErrorConvertingIdArgumentToInt  = errors.New("You must provide a valid id")
	ErrorConvertingStringToBool     = errors.New("You must provide a valid value. Use -h or --help for more information")
	ErrorHandlingFile               = errors.New("You must provide a valid file name. Use -h or --help for more information")
	ErrorOpeningFile                = errors.New("Failed to open file")
	ErrorInvalidVariablesFileFormat = errors.New("You must provide a valid variables file content. Use -h or --help for more information")
	ErrorInternalServerError        = errors.New("Something went wrong, please try again")
	ErrorUpdateNoFlagsSent          = errors.New("You must provide at least one value in update. Use -h or --help for more information")
	ErrorUnmarshalReader            = errors.New("Failed to unmarshal from reader")
	ErrorGetHttpClient              = errors.New("Failed to get http client")
	ErrorFormatOut                  = errors.New("Failed to format response")
	ErrorWriteFile                  = errors.New("Failed to write to file")
	ErrorTokenManager               = errors.New("Failed to create token manager")
	ErrorTokenNotProvided           = errors.New("Token not provided, loading the saved one")
	ErrorInvalidToken               = errors.New("Invalid token")
)
