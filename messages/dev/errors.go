package dev

import "errors"

var (
	ErrorVulcanExecute       = errors.New("Error executing Vulcan: %s")
	ErrFailedToRunDevCommand = errors.New("Failed to run dev command. Verify if the command is correct and check the output above for more details. Try the 'azion dev' command again or contact Azion's support")
)
