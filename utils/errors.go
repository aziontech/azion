package utils

import (
	"errors"
)

var (
	//Generic errors that can be used by any package
	ErrorConvertingStringToBool     = errors.New("The given data isn’t a boolean type value. Provide a valid boolean type in the command’s flag value and try again. Use the flags -h or --help with a command or subcommand to display more information and try again")
	ErrorConvertingStringToInt      = errors.New("The given data isn’t a integer type value. Provide a valid integer type in the command’s flag value and try again. Use the flags -h or --help with a command or subcommand to display more information and try again")
	ErrorHandlingFile               = errors.New("The file name doesn’t exist or is invalid. Provide a valid and/or existing path and file name. Use the flags -h or --help with a command or subcommand to display more information and try again")
	ErrorEmptyFile                  = errors.New("The file’s content is empty. Provide a path and file with valid content and try the command again. Use the flags -h or --help with a command or subcommand to display more information and try again")
	ErrorOpeningFile                = errors.New("Failed to open a file. Verify if the path and file exists and/or the file is corrupted and try the command again")
	ErrorReadingFile                = errors.New("Failed to read a file. Verify if the path and file exists and/or the file is corrupted and try the command again")
	ErrorParsingModel               = errors.New("Failed in parsing the model")
	ErrorExecTemplate               = errors.New("Failed to apply template to given data and store result in buffer")
	ErrorInvalidVariablesFileFormat = errors.New("The format of the variables in the file is invalid. You must provide a file with valid variable formats. Use the flags -h or --help with a command or subcommand to display more information")
	ErrorInternalServerError        = errors.New("The server could not process the request because an internal and unexpected problem occurred. Wait a few seconds and try again. For more information run the command again using the '--debug' flag. If the problem persists, contact Azion’s support")
	ErrorUpdateNoFlagsSent          = errors.New("The subcommand update needs at least one flag with a valid value. Run the command `azion <command> update --help` to display more information and try again")
	ErrorUnmarshalReader            = errors.New("Failed to decode the given 'json' file. Verify if the file format is JSON or fix its content according to the JSON format specification at https://www.json.org/json-en.html")
	ErrorFormatOut                  = errors.New("The server failed formatting data for display. Repeat the HTTP request and check the HTTP response's format")
	ErrorWriteFile                  = errors.New("The file is read-only and/or isn't accessible. Change the attributes of the file to read and write and/or give access to it")
	ErrorTokenManager               = errors.New("Internal token handling failure. Run 'azion configure --help' command to display more information and try again")
	ErrorTokenNotProvided           = errors.New("Token was not provided; the CLI uses a previous stored token if it was configured. You must provide a valid token, or create a new one, and configure it to use with the CLI. Manage your personal tokens on RTM using the Account Menu > Personal Tokens and configure the token again with the command 'azion -t <token>'")
	ErrorInvalidToken               = errors.New("The provided token is invalid. You must create a new token and configure it to use with the CLI. Manage your personal tokens on RTM using the Account Menu > Personal Tokens and configure the new token with the command 'azion -t <new_token>'")
	ErrorToken401                   = errors.New("The token doesn't exist or has expired. Manage your personal tokens on RTM using the Account Menu > Personal Tokens and configure a valid token with the command 'azion -t <my_token>'")
	ErrorForbidden403               = errors.New("You do not have the permissions to access the API. Make sure the feature is enabled in your profile")
	ErrorNotFound404                = errors.New("The given ID or API's endpoint doesn't exist or isn't available. Check that the identifying information is correct")
	ErrorFetchingTemplates          = errors.New("Failed to fetch templates from the Azion's GitHub remote repository. Verify the connectivity to the repository https://github.com/aziontech/azioncli-template and try again")
	ErrorMovingFiles                = errors.New("Failed to initialize your project with the Azion template. Please verify if you have write permissions to this directory")
	ErrorInvalidOption              = errors.New("You must inform 'yes' or 'no' as input, or force --yes or --no by using the flags")
	ErrorCleaningDirectory          = errors.New("Failed to clean the directory's contents because the directory is read-only and/or isn't accessible. Change the attributes of the directory to read/write and/or give access to it")
	ErrorRunningCommand             = errors.New("Failed to run the command specified in the template (config.json)")
	ErrorRunningCommandStream       = errors.New("Failed to run the command specified in the template (config.json): %s")
	ErrorLoadingEnvVars             = errors.New("Failed to load the Applications's environment variables. Verify if the environment variables exist and/or if their values are valid and try again")
	ErrorOpeningAzionJsonFile       = errors.New("Failed to open the azion.json file. The file doesn't exist, is corrupted, or has an invalid JSON format. Verify if you have initialized your project, if the file format is JSON or fix its content according to the JSON format specification at https://www.json.org/json-en.html")
	ErrorUnmarshalAzionJsonFile     = errors.New("Failed to parse the given 'azion.json' file. Verify if the file format is JSON or fix its content according to the JSON format specification at https://www.json.org/json-en.html")
	ErrorMarshalAzionJsonFile       = errors.New("Failed to encode the given 'azion.json' file. Verify if the file format is JSON or fix its content according to the JSON format specification at https://www.json.org/json-en.html")
	ErrorWritingAzionJsonFile       = errors.New("Failed to write in the given 'azion.json' file. Verify if the file is writable and/or you have access to it, if the data format is JSON, or fix the content according to the JSON format specification at https://www.json.org/json-en.html")
	ErrorTimeoutAPICall             = errors.New("CLI's request has timed out during communication with Azion. Verify if it has completed successfully or wait some time and try the command again")
	ErrorCreateFile                 = errors.New("Failed to create %s file")
	ErrorProductNotOwned            = errors.New("This account does not own the following product")
	ErrorUnknownSystem              = errors.New("Unknown system")
	ErrorCommandNotFound            = errors.New("Command '%s' not found")
	ErrorGetAssetsNamesAzioncli     = errors.New("Failed to fetch the assets names from azion")
	ErrorArgumentIsEmpty            = errors.New("Argument is empty")
	ErrorParseResponse              = errors.New("Failed to parse your response. Check your response and try again. If the error persists, contact Azion support")
	ErrorMinTlsVersion              = errors.New("This is not a valid TLS Version. Run azion edge_applications <subcommand> --help for more information")
	ErrorNameInUse                  = errors.New("The name you've selected is already in use by another resource. Please choose a different name. Run 'azion list [resource]' to see all your resources")
	ErrorCancelledContextInput      = errors.New("Execution interrupted by the user. All interactions of this flow were lost.")
	ErrorWriteProfiles              = errors.New("Failed to write profiles.toml file: %w")
	ErrorReadProfiles               = errors.New("Failed to read profiles.toml file: %w")
	ErrorWriteSettings              = errors.New("Failed to write settings.toml file: %w")
	ErrorCheckingProfilesFile       = errors.New("Failed to check profiles file: %w")
	ErrorCreatingConfigDirectory    = errors.New("Failed to create config directory: %w")
	ErrorCreatingDefaultProfiles    = errors.New("Failed to create default profiles.json: %w")
)

const (
	ERROR_CLONE        = "Error cloning the repository: %v"
	ERROR_CDDIR        = "Error entering the repository folder: %v"
	ERROR_INVALID_REPO = "Invalid repository URL: %s"
)
