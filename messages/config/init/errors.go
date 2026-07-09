package init

import "errors"

var (
	ErrorCreatingConfig = errors.New("Failed to create azion.config file")
	ErrorConfigExists   = errors.New("Configuration file already exists")
)
