package init

import "errors"

// Used only by the v4 command tree.
var (
	ErrorCreatingConfig = errors.New("Failed to create azion.config file")
	ErrorConfigExists   = errors.New("Configuration file already exists")
)
