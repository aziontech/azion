package webapp

import "errors"

var (
	ErrOpeningConfigFile   = errors.New("Failed to open config.json file")
	ErrUnmarshalConfigFile = errors.New("Failed to parse config.json file")
	ErrReadEnvFile         = errors.New("Failed to read build.env file")
	ErrFailedToRunCommand  = errors.New("Failed to run build step command. Verify if the command is correct and/or check its output for more details")

	ErrorOpeningConfigFile   = errors.New("Failed to open config.json file")
	ErrorUnmarshalConfigFile = errors.New("Failed to unmarshal config.json file")
	ErrorOpeningAzionFile    = errors.New("Failed to open azion.json file")
	ErrorUnmarshalAzionFile  = errors.New("Failed to unmarshal azion.json file")
	ErrorPackageJsonNotFound = errors.New("Failed to find package.json in current directory. Verify if your are currently in your project's directory")
	ErrorYesAndNoOptions     = errors.New("You can only use one option at a time. Please use either --yes or --no")

	ErrorCreateApplication = errors.New("Failed to create Edge Application")
	ErrorUpdateApplication = errors.New("Failed to update Edge Application")
	ErrorCreateInstance    = errors.New("Failed to create Edge Function Instance")
	ErrorCreateDomain      = errors.New("Failed to create Domain")
	ErrorUpdateDomain      = errors.New("Failed to update Domain")
)
