package init

import "errors"

var (
	ErrorOpeningAzionFile    = errors.New("Failed to open azion.json file")
	ErrorUnmarshalAzionFile  = errors.New("Failed to unmarshal azion.json file")
	ErrorPackageJsonNotFound = errors.New("Failed to find package.json in current directory. Verify if your are currently in your project's directory")
)
