package dev

import "errors"

// Used by both the v3 and the v4 command trees.
var (
	ErrorVulcanExecute       = errors.New("Error executing Bundler: %s")
	ErrFailedToRunDevCommand = errors.New("Failed to run dev command. Verify if the command is correct and check the output above for more details. Run the 'azion dev' command again or contact Azion's support")
)
