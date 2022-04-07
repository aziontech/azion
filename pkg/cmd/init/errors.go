package init

import "errors"

var (
	ErrorOpeningConfigFile   = errors.New("Failed to open config.json file")
	ErrorUnmarshalConfigFile = errors.New("Failed to unmarshal config.json file")
	ErrorOpeningAzionFile    = errors.New("Failed to open azion.json file")
	ErrorUnmarshalAzionFile  = errors.New("Failed to unmarshal azion.json file")
	ErrorPackageJsonNotFound = errors.New("Failed to find package.json in current directory. Verify if your are currently in your project's directory")
	ErrorYesAndNoOptions     = errors.New("You can only use one option at a time. Please use either --yes or --no")
)
