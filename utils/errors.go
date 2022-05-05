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
	ErrorFormatOut                  = errors.New("Failed to format response")
	ErrorWriteFile                  = errors.New("Failed to write to file")
	ErrorTokenManager               = errors.New("Failed to create token manager")
	ErrorTokenNotProvided           = errors.New("Token not provided, loading the saved one")
	ErrorInvalidToken               = errors.New("Invalid token")
	ErrorMissingGitBinary           = errors.New("You must have git binary installed")
	ErrorFetchingTemplates          = errors.New("Failed to fetch azioncli templates from Github")
	ErrorMovingFiles                = errors.New("Failed to move files to destination directory")
	ErrorUnsupportedType            = errors.New("Unsupported type. Use -h or --help for more information")
	ErrorInvalidOption              = errors.New("Invalid option")
	ErrorCleaningDirectory          = errors.New("Failed to clean the directory's contents")
	ErrorRunningCommand             = errors.New("Failed to run specified command")
	ErrorLoadingEnvVars             = errors.New("Failed to load environment variables")
	ErrorOpeningAzionJsonFile       = errors.New("Failed to open azion.json")
	ErrorUnmarshalAzionJsonFile     = errors.New("Failed to parse azion.json. Verify the file format.")
	ErrorMarshalAzionJsonFile       = errors.New("Failed to encode azion.json. Verify the file format.")
	ErrorWritingAzionJsonFile       = errors.New("Failed to write azion.json. Verify the file format.")
)
