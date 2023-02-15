package origins

import "errors"

var (
	ErrorGetOrigins                   = errors.New("Failed to describe the origins. Check your settings and try again. If the error persists, contact Azion support.")
	ErrorGetOrigin                    = errors.New("Failed to describe the origin. Check your settings and try again. If the error persists, contact Azion support.")
	ErrorMissingApplicationIDArgument = errors.New("A required flag is missing. You must supply an application_id as an argument. Run 'azioncli origins list --help' command to display more information and try again")
	ErrorMissingArguments             = errors.New("A required flag is missing. You must supply an application_id and origin_id as an argument. Run 'azioncli origins list --help' command to display more information and try again")
	ErrorMandatoryCreateFlags         = errors.New("A mandatory flag is missing. You must provide --application-id, --name, --addresses and --host-header flags when the --application-id and --in flag is not provided. Run the command 'azioncli domains create --help' to display more information and try again.")
	ErrorMandatoryUpdateFlags         = errors.New("A mandatory flag is missing. You must provide --application-id and --origin-id flags when the --application-id and --in flag is not provided. Run the command 'azioncli domains create --help' to display more information and try again.")
	ErrorHmacAuthenticationFlag       = errors.New("Invalid --hmac-authentication flag provided. The flag must have  'true' or 'false' values. Run the command 'azioncli domains <subcommand> --help' to display more information and try again.")
	ErrorCreateOrigin                 = errors.New("Failed to create the Origin: %s. Check your settings and try again. If the error persists, contact Azion support.")
	ErrorUpdateOrigin                 = errors.New("Failed to update the Origin: %s. Check your settings and try again. If the error persists, contact Azion support.")
  ErrorFailToDelete                 = errors.New("Failed to delete the Origin: %s. Check your settings and try again. If the error persists, contact Azion support")
)

