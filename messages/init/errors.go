package init

import "errors"

var (
	EdgeApplicationsOutputErr = errors.New("This output-ctrl option is not available. Read the readme files found in the repository https://github.com/aziontech/azion-template and try again")

	ErrorVulcanExecute                 = errors.New("Error executing Bundler: %s")
	ErrorUpdatingVulcan                = errors.New("Failed to update Bundler: %s")
	ErrorInstallVulcan                 = errors.New("Failed to install Bundler: %s")
	ErrorOpeningConfigFile             = errors.New("Failed to open the config.json file. The file doesn't exist, is corrupted, or has an invalid JSON format. Verify if the file format is JSON or fix its content according to the JSON format specification at https://www.json.org/json-en.html")
	ErrorUnmarshalConfigFile           = errors.New("Failed to unmarshal the config.json file. Verify if the file format is JSON or fix its content according to the JSON format specification at https://www.json.org/json-en.html")
	ErrFailedToRunInitCommand          = errors.New("Failed to run the init step command. Verify if the command is correct and check the output above for more details. Try the 'azion init' command again or contact Azion's support")
	ErrorGetAllTags                    = errors.New("Failed to return all reference tags in a repository. Verify your repository tags and try again. If the error persists, contact Azion support.")
	ErrorIterateAllTags                = errors.New("Failed to iterate over Git reference. Verify the credentials to access your Git repository and try again. If the error persists, contact Azion support.")
	ErrorOpeningAzionFile              = errors.New("Failed to open the azion.json file. The file doesn't exist, is corrupted, or has an invalid JSON format. Verify if the file format is JSON or fix its content according to the JSON format specification at https://www.json.org/json-en.html")
	ErrorUnmarshalAzionFile            = errors.New("Failed to unmarshal the azion.json file. Verify if the file format is JSON or fix its content according to the JSON format specification at https://www.json.org/json-en.html")
	ErrorNpmNotInstalled               = errors.New("Failed to open the NPM package Manager. Visit the website 'https://nodejs.org/en/download/' and follow the instructions to install the Node.js JavaScript runtime environment in your operating system. Node.js installation includes the NPM package manager")
	ErrorFailedCreatingWorkerDirectory = errors.New("Failed to create the worker directory. The worker's parent directory is read-only and/or isn't accessible. Change the permissions of the parent directory to read and write and/or give access to it")
	ErrorFailedCreatingAzionDirectory  = errors.New("Failed to create the azion directory. The public's parent directory is read-only and/or isn't accessible. Change the permissions of the parent directory to read and write and/or give access to it")
	ErrorDeps                          = errors.New("Failed to install project dependencies")
	ErrorWorkingDir                    = errors.New("Failed to change current working directory")
	ErrorReadingGitignore              = errors.New("Failed to read your .gitignore file")
	ErrorWritingGitignore              = errors.New("Failed to write to your .gitignore file")
	ErrorGetProjectInfo                = errors.New("Failed to get project preset")
)
