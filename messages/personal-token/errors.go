package personal_token

import "errors"

var (
	ErrorMandatoryCreateFlags = errors.New("Required flags are missing. You must provide name and expiration flags when the --in flag is not provided. Run the command 'azion <command> <subcommand> --help' to display more information and try again.")
	ErrorCreate               = errors.New("Failed to create the Personal Token: %s. Check your settings and try again. If the error persists, contact Azion support.")
	ErrorGet                  = errors.New("Failed to describe the personal token: %s. Check your settings and try again. If the error persists, contact Azion support.")
)
