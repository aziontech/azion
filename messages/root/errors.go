package root

import "errors"

var (
	ErrorCurrentUser          = errors.New("Failed to get current user's information.")
	ErrorMarshalUserInfo      = errors.New("Failed to marshal current user information.")
	ErrorUnmarshalUserInfo    = errors.New("Failed to unmarshal current user information.")
	ErrorReadFileSettingsToml = errors.New("Provide the correct path of the configuration file. Make sure the file is in .toml format, access the document for more information https://www.azion.com/en/documentation/devtools/cli/globals/#config")
	ErrorPrefix               = errors.New("A configuration path is expected for your location, not a flag")
	ErrorParseTimeout         = errors.New("Failed to parse timeout flag")
)
