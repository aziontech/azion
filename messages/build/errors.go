package build

import "errors"

var (
	ErrorBuilding              = errors.New("Failed to build your resource. Azion configuration not found. Make sure you are in the root directory of your local repository and have already initialized or linked your resource with the commands 'azion init' or 'azion link'")
	ErrorVulcanExecute         = errors.New("Error executing Bundler: %s")
	EdgeApplicationsOutputErr  = errors.New("This output-ctrl option is not available. Read the readme files found in the repository https://github.com/aziontech/azion-template and try again")
	ErrFailedToRunBuildCommand = errors.New("Failed to run the build command. Verify if the command is correct and check the output above for more details. Run the 'azion build' command again or contact Azion's support")
	ErrorUnmarshalConfigFile   = errors.New("Failed to unmarshal the config.json file. Verify if the file format is JSON or fix its content according to the JSON format specification at https://www.json.org/json-en.html")
	ErrorOpeningConfigFile     = errors.New("Failed to open the config.json file. The file doesn't exist, is corrupted, or has an invalid JSON format. Verify if the file format is JSON or fix its content according to the JSON format specification at https://www.json.org/json-en.html")
	ErrorOpeningAzionFile      = errors.New("Failed to open the azion.json file. The file doesn't exist, is corrupted, or has an invalid JSON format. Verify if the file format is JSON or fix its content according to the JSON format specification at https://www.json.org/json-en.html")
	ErrorPolyfills             = errors.New("Invalid --use-node-polyfills flag provided. The flag must have  'true' or 'false' values. Run the command 'azion build --help' to display more information and try again.")
	ErrorWorker                = errors.New("Invalid --use-own-worker flag provided. The flag must have  'true' or 'false' values. Run the command 'azion build --help' to display more information and try again.")
)
