package vulcan

import "errors"

// Used by both the v3 and the v4 command trees.
var (
	ErrConfigFileNotFound = errors.New("The azion.config file informed in --config-file was not found. Verify the path and try again")
)
