package apply

import "errors"

var (
	ErrorReadingManifest     = errors.New("Failed to read manifest.json file")
	ErrorApplyingResources   = errors.New("Failed to apply resources from manifest")
	ErrorCreatingAzionJson   = errors.New("Failed to create azion.json file")
	ErrorGeneratingManifest  = errors.New("Failed to generate manifest from azion.config")
	ErrorAzionConfigNotFound = errors.New("azion.config file not found. Create an azion.config file to define your application configuration")
)
