package apply

import "errors"

// Used only by the v4 command tree.
var (
	ErrorReadingManifest     = errors.New("Failed to read manifest.json file")
	ErrorCreatingAzionJson   = errors.New("Failed to create azion.json file")
	ErrorAzionConfigNotFound = errors.New("azion.config file not found. Create an azion.config file to define your application configuration")
)
