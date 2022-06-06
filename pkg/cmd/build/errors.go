package build

import "errors"

var (
	ErrOpeningConfigFile   = errors.New("Failed to open config.json file")
	ErrUnmarshalConfigFile = errors.New("Failed to parse config.json file")
	ErrReadEnvFile         = errors.New("Failed to read build.env file")
	ErrFailedToRunCommand  = errors.New("Failed to run build step command. Verify if the command is correct and/or check its output for more details")
)
