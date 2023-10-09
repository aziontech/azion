package domains

import "errors"

var (
	ErrorUpdateDomain           = errors.New("Failed to update the Domain: %s. Check your settings and try again. If the error persists, contact Azion support.")
	ErrorActiveFlag             = errors.New("Invalid --active flag provided. The flag must have  'true' or 'false' values. Run the command 'azion domains <subcommand> --help' to display more information and try again.")
	ErrorDigitalCertificateFlag = errors.New("Invalid --digital-certificate-id flag provided. The flag must have an Integer or 'null' as a value. Run the command 'azion domains <subcommand> --help' to display more information and try again")
	ErrorConvertDomainID        = errors.New("The domain ID you provided is invalid. The value must be an integer. You may run the 'azion list domains' command to check your domain ID")
)
