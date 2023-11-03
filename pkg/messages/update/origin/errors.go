package origin

import "errors"

var (
	ErrorUpdateOrigin           = errors.New("Failed to update the Origin: %s. Check your settings and try again. If the error persists, contact Azion support.")
	ErrorHmacAuthenticationFlag = errors.New("Invalid --hmac-authentication flag provided. The flag must have  'true' or 'false' values. Run the command 'azion <command> <subcommand> --help' to display more information and try again.")
)
