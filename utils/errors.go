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
	ErrorOpeningAzionFile           = errors.New("Failed to open azion.json file")
	ErrorUnmarshalAzionFile         = errors.New("Failed to unmarshal azion.json file")
	ErrorCleaningDirectory          = errors.New("Failed to clean the directory's contents")
	ErrorPackageJsonNotFound        = errors.New("Failed to find package.json in current directory. Verify if your are currently in your project's directory")
	ErrorUnsupportedType            = errors.New("Unsupported type. Use -h or --help for more information")
	ErrorInvalidOption              = errors.New("Invalid option")
)
