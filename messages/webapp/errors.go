package webapp

import "errors"

var (
	ErrOpeningConfigFile   = errors.New("Failed to open the config.json file. The file doesn't exist, is corrupted, or has an invalid JSON format. Verify if the file was deleted or changed or run the 'azioncli webapp init' command again")
	ErrUnmarshalConfigFile = errors.New("Failed to parse the config.json file. Verify if the file format is JSON or fix its content according to the JSON format specification at https://www.json.org/json-en.html")
	ErrReadEnvFile         = errors.New("Failed to read the build.env file. Verify if the file was deleted or changed or run the 'azioncli webapp init' command again")
	ErrFailedToRunCommand  = errors.New("Failed to run the build step command. Verify if the command is correct and check its output for more details. Try the 'azioncli webapp build' command again or contact Azion's support")

	ErrorOpeningConfigFile   = errors.New("Failed to open the config.json file. The file doesn't exist, is corrupted, or has an invalid JSON format. Verify if the file format is JSON or fix its content according to the JSON format specification at https://www.json.org/json-en.html")
	ErrorUnmarshalConfigFile = errors.New("Failed to unmarshal the config.json file. Verify if the file format is JSON or fix its content according to the JSON format specification at https://www.json.org/json-en.html")
	ErrorOpeningAzionFile    = errors.New("Failed to open the azion.json file. The file doesn't exist, is corrupted, or has an invalid JSON format. Verify if the file format is JSON or fix its content according to the JSON format specification at https://www.json.org/json-en.html")
	ErrorUnmarshalAzionFile  = errors.New("Failed to unmarshal the azion.json file. Verify if the file format is JSON or fix its content according to the JSON format specification at https://www.json.org/json-en.html")
	ErrorPackageJsonNotFound = errors.New("Failed to find the package.json file in the current directory. Verify if you are currently in your project's directory and provide an existing path and file with a valid JSON format")
	ErrorYesAndNoOptions     = errors.New("You can use only one option at a time. Choose either --yes or --no options. Run the command 'azioncli webapp <subcommand> --help' to display more information and try again")

	ErrorCreateApplication = errors.New("Failed to create the Edge Application: %s. Check your settings and try again. If the error persists, contact Azion support")
	ErrorUpdateApplication = errors.New("Failed to update the Edge Application: %s. Check your settings and try again. If the error persists, contact Azion support")
	ErrorCreateInstance    = errors.New("Failed to create the Edge Function Instance: %s. Check your settings and try again. If the error persists, contact Azion support")
	ErrorCreateDomain      = errors.New("Failed to create the Domain: %s. Check your settings and try again. If the error persists, contact Azion support")
	ErrorUpdateDomain      = errors.New("Failed to update the Domain: %s. Check your settings and try again. If the error persists, contact Azion support")

	ErrorWebappInitCmdNotSpecified = errors.New("Init step command not specified. No action will be taken")
)
