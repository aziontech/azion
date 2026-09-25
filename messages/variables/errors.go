package variables

import "errors"

// Used by both the v3 and the v4 command trees.
var (
	ErrorGetItem              = errors.New("Failed to describe the variable: %s. Check your settings and try again. If the error persists, contact Azion support.")
	ErrorFailToDeleteVariable = errors.New("Failed to delete the variable: %s. Check your settings and try again. If the error persists, contact Azion support")
	ErrorSecretFlag           = errors.New("Invalid --secret flag provided. The value must be 'true' or 'false'. Run the command 'azion variables <subcommand> --help' to display more information and try again")
	ErrorUpdateVariable       = errors.New("Failed to update the variable: %s. Check your settings and try again. If the error persists, contact Azion support")
	ErrorCreateItem           = errors.New("Failed to create the variable: %s. Check your settings and try again. If the error persists, contact Azion support.")
)
