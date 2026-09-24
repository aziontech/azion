package build

import "errors"

// Used by both the v3 and the v4 command trees.
var (
	ErrorBuilding              = errors.New("Failed to build your resource. Azion configuration not found. Make sure you are in the root directory of your local repository and have already initialized or linked your resource with the commands 'azion init' or 'azion link'")
	ErrorVulcanExecute         = errors.New("Error executing Bundler: %s")
	ErrFailedToRunBuildCommand = errors.New("Failed to run the build command. Verify if the command is correct and check the output above for more details. Run the 'azion build' command again or contact Azion's support")
	ErrorPolyfills             = errors.New("Invalid --use-node-polyfills flag provided. The flag must have  'true' or 'false' values. Run the command 'azion build --help' to display more information and try again.")
	ErrorWorker                = errors.New("Invalid --use-own-worker flag provided. The flag must have  'true' or 'false' values. Run the command 'azion build --help' to display more information and try again.")
)
