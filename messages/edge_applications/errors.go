package edge_applications

import "errors"

var (
	EdgeApplicationsOutputErr = errors.New("This output-ctrl option is not available. Read the readme files found in the repository https://github.com/aziontech/azioncli-template and try again")

	ErrOpeningConfigFile         = errors.New("Failed to open the config.json file. The file doesn't exist, is corrupted, or has an invalid JSON format. Verify if the file was deleted or changed or run the 'azioncli edge_applications init' command again")
	ErrUnmarshalConfigFile       = errors.New("Failed to parse the config.json file. Verify if the file format is JSON or fix its content according to the JSON format specification at https://www.json.org/json-en.html")
	ErrReadEnvFile               = errors.New("Failed to read the webdev.env file. Verify if the file is corrupted or changed or run the 'azioncli edge_applications publish' command again")
	ErrFailedToRunBuildCommand   = errors.New("Failed to run the build step command. Verify if the command is correct and check the output above for more details. Try the 'azioncli edge_applications build' command again or contact Azion's support")
	ErrFailedToRunInitCommand    = errors.New("Failed to run the init step command. Verify if the command is correct and check the output above for more details. Try the 'azioncli edge_applications build' command again or contact Azion's support")
	ErrFailedToRunPublishCommand = errors.New("Failed to run the publish step command. Verify if the command is correct and check the output above for more details. Try the 'azioncli edge_applications build' command again or contact Azion's support")
	ErrorOpeningConfigFile       = errors.New("Failed to open the config.json file. The file doesn't exist, is corrupted, or has an invalid JSON format. Verify if the file format is JSON or fix its content according to the JSON format specification at https://www.json.org/json-en.html")
	ErrorUnmarshalConfigFile     = errors.New("Failed to unmarshal the config.json file. Verify if the file format is JSON or fix its content according to the JSON format specification at https://www.json.org/json-en.html")
	ErrorOpeningAzionFile        = errors.New("Failed to open the azion.json file. The file doesn't exist, is corrupted, or has an invalid JSON format. Verify if the file format is JSON or fix its content according to the JSON format specification at https://www.json.org/json-en.html")
	ErrorUnmarshalAzionFile      = errors.New("Failed to unmarshal the azion.json file. Verify if the file format is JSON or fix its content according to the JSON format specification at https://www.json.org/json-en.html")
	ErrorPackageJsonNotFound     = errors.New("Failed to find the package.json file in the current directory. Verify if you are currently in your project's directory and provide an existing path and file with a valid JSON format")
	ErrorYesAndNoOptions         = errors.New("You can use only one option at a time. Choose either --yes or --no options. Run the command 'azioncli edge_applications <subcommand> --help' to display more information and try again")

	ErrorCreateApplication = errors.New("Failed to create the Edge Application: %s. Check your settings and try again. If the error persists, contact Azion support")
	ErrorUpdateApplication = errors.New("Failed to update the Edge Application: %s. Check your settings and try again. If the error persists, contact Azion support")
	ErrorCreateInstance    = errors.New("Failed to create the Edge Function Instance: %s. Check your settings and try again. If the error persists, contact Azion support")
	ErrorCreateDomain      = errors.New("Failed to create the Domain: %s. Check your settings and try again. If the error persists, contact Azion support")
	ErrorUpdateDomain      = errors.New("Failed to update the Domain: %s. Check your settings and try again. If the error persists, contact Azion support")

	ErrorNpmNotInstalled               = errors.New("Failed to open the NPM package Manager. Visit the website 'https://nodejs.org/en/download/' and follow the instructions to install the Node.js JavaScript runtime environment in your operating system. Node.js installation includes the NPM package manager")
	FailedUpdatingScriptsDeployField   = errors.New("Failed to update scripts.deploy field in package.json file. Make sure you have the needed permissions and try again. If the error persists, contact the Azion support")
	FailedUpdatingScriptsBuildField    = errors.New("Failed to update scripts.build field in package.json file. Make sure you have the needed permissions and try again. If the error persists, contact the Azion support")
	FailedUpdatingNameField            = errors.New("Failed to update name field in package.json file. Make sure you have the needed permissions and try again. If the error persists, contact the Azion support")
	ErrorMandatoryEnvs                 = errors.New("You must provide the following enviroment variables: AWS_SECRET_ACCESS_KEY and AWS_ACCESS_KEY_ID. Please edit the following file 'azion/webdev.env' and add your credentials")
	ErrorFailedCreatingWorkerDirectory = errors.New("Failed to create the worker directory. The worker's parent directory is read-only and/or isn't accessible. Change the permissions of the parent directory to read and write and/or give access to it")
	ErrorFailedCreatingPublicDirectory = errors.New("Failed to create the public directory. The public's parent directory is read-only and/or isn't accessible. Change the permissions of the parent directory to read and write and/or give access to it")
	ErrorFailedCreatingAzionDirectory  = errors.New("Failed to create the azion directory. The public's parent directory is read-only and/or isn't accessible. Change the permissions of the parent directory to read and write and/or give access to it")

	ErrorCreateFunction = errors.New("Failed to create edge function: %s. Check your settings and try again. If the error persists, contact Azion support")
	ErrorUpdateFunction = errors.New("Failed to update the Edge Function: %s. Check your settings and try again. If the error persists, contact Azion support")
	ErrorCodeFlag       = errors.New("Failed to read the code file. Verify if the file name and its path are correct and the file content has a valid code format")
	ErrorArgsFlag       = errors.New("Failed to read the args file. Verify if the file name and its path are correct and the file's content has a valid JSON format")
	ErrorParseArgs      = errors.New("Failed to parse JSON args. Verify if the file's content has a valid JSON format")

	ErrorGetAllTags       = errors.New("Failed returning all Reference tags in a repository. Verify your repository tags and try again. If the error persists, contact Azion support.")
	ErrorIterateAllTags   = errors.New("Failed to iterate over Git reference. Verify the credentials to access your Git repository and try again. If the error persists, contact Azion support.")
	ErrorWritingWebdevEnv = errors.New("Failed to write 'webdev.env' file. Verify if the file is writable and/or you have access to it, if the data format is JSON, or fix the content according to the JSON format specification at https://www.json.org/json-en.html")

	ErrorMissingApplicationIdArgument = errors.New("A mandatory flag is missing. You must provide an application_id as an argument or path to import the file. Run the command 'azioncli edge_applications <subcommand> --help' to display more information and try again")
	ErrorGetApplication               = errors.New("Failed to get the Edge Application: %s. Check your settings and try again. If the error persists, contact Azion support")
	ErrorFailToDeleteApplication      = errors.New("Failed to delete the Edge Application: %s. Check your settings and try again. If the error persists, contact Azion support")

	ErrorMandatoryCreateFlags        = errors.New("A mandatory flag is missing. You must provide --active, --code, and --name flags when the --in flag is not provided. Run the command 'azioncli edge_applications create --help' to display more information and try again")
	ErrorActiveFlag                  = errors.New("Invalid --active flag provided. The flag must have  'true' or 'false' values. Run the command 'azioncli edge_applications <subcommand> --help' to display more information and try again")
	ErrorApplicationAccelerationFlag = errors.New("Invalid --application-acceleration flag provided. The flag must have  'true' or 'false' values. Run the command 'azioncli edge_applications <subcommand> --help' to display more information and try again")
	ErrorCachingFlag                 = errors.New("Invalid --caching flag provided. The flag must have  'true' or 'false' values. Run the command 'azioncli edge_applications <subcommand> --help' to display more information and try again")
	ErrorDeviceDetectionFlag         = errors.New("Invalid --device-detection flag provided. The flag must have  'true' or 'false' values. Run the command 'azioncli edge_applications <subcommand> --help' to display more information and try again")
	ErrorEdgeFirewallFlag            = errors.New("Invalid --edge-direwall flag provided. The flag must have  'true' or 'false' values. Run the command 'azioncli edge_applications <subcommand> --help' to display more information and try again")
	ErrorEdgeFunctionsFlag           = errors.New("Invalid --edge-functions flag provided. The flag must have  'true' or 'false' values. Run the command 'azioncli edge_applications <subcommand> --help' to display more information and try again")
	ErrorImageOptimizationFlag       = errors.New("Invalid --image-optimization flag provided. The flag must have  'true' or 'false' values. Run the command 'azioncli edge_applications <subcommand> --help' to display more information and try again")
	ErrorL2CachingFlag               = errors.New("Invalid --l2-caching flag provided. The flag must have  'true' or 'false' values. Run the command 'azioncli edge_applications <subcommand> --help' to display more information and try again")
	ErrorLoadBalancerFlag            = errors.New("Invalid --load-balancer flag provided. The flag must have  'true' or 'false' values. Run the command 'azioncli edge_applications <subcommand> --help' to display more information and try again")
	ErrorRawLogsFlag                 = errors.New("Invalid --raw-logs flag provided. The flag must have  'true' or 'false' values. Run the command 'azioncli edge_applications <subcommand> --help' to display more information and try again")
	ErrorWebApplicationFirewallFlag  = errors.New("Invalid --webapp-firewall flag provided. The flag must have  'true' or 'false' values. Run the command 'azioncli edge_applications <subcommand> --help' to display more information and try again")
	ErrorMinTlsVersion               = errors.New("This is not a valid TLS Version. Run azioncli edge_applications <subcommand> --help for more information")
	ErrorMissingApplicationIdJson    = errors.New("Application ID is missing from azion.json. Please initialize and publish your project first before using cascade delete")
	ErrorMissingAzionJson            = errors.New("Azion.json file is missing. Please initialize and publish your project first before using cascade delete")
	ErrorFailedUpdateAzionJson       = errors.New("Failed to update azion.json file to remove IDs of deleted resource")

	ErrorGetVersionId = errors.New("Failed to get Version Id: %s")
)
