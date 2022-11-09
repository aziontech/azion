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

	ErrorNpmNotInstalled               = errors.New("Failed to open the NPM package Manager. Visit the website 'https://nodejs.org/en/download/' and follow the instructions to install the Node.js JavaScript runtime environment in your operating system. Node.js installation includes the NPM package manager")
	FailedUpdatingScriptsDeployField   = errors.New("Failed to update scripts.deploy field in package.json file. Make sure you have the needed permissions and try again. If the error persists, contact the Azion support")
	FailedUpdatingScriptsBuildField    = errors.New("Failed to update scripts.build field in package.json file. Make sure you have the needed permissions and try again. If the error persists, contact the Azion support")
	ErrorMandatoryEnvs                 = errors.New("You must provide the following enviroment variables: AWS_SECRET_ACCESS_KEY and AWS_ACCESS_KEY_ID. Please edit the following file 'azion/webdev.env' and add your credentials")
	ErrorFailedCreatingWorkerDirectory = errors.New("Failed to create the worker directory. The worker's parent directory is read-only and/or isn't accessible. Change the permissions of the parent directory to read and write and/or give access to it")
)
