package domains

import "errors"

var (
	ErrorMandatoryCreateFlags    = errors.New("One or more required flags are missing. You must provide --name and --application-id flags when the --in flag is not provided. Run the command 'azion domains create --help' to display more information and try again.")
	ErrorActiveFlag              = errors.New("Invalid --active flag provided. The flag must have  'true' or 'false' values. Run the command 'azion domains <subcommand> --help' to display more information and try again.")
	ErrorMissingDomainIdArgument = errors.New("The flag '--domain-id' must be informed. Please, inform the correct id and try again or run the command ‘azion domains <subcommand> --help’ to display more information and try again. ")
	ErrorCreateDomain            = errors.New("Failed to create the Domain: %s. Check your settings and try again. If the error persists, contact Azion support.")
	ErrorUpdateDomain            = errors.New("Failed to update the Domain: %s. Check your settings and try again. If the error persists, contact Azion support.")

	//used by more than one cmd
	DomainsFlagId      = "Unique identifier of the Domain"
	DomainsFileWritten = "File successfully written to: %s\n"

	ErrorMissingApplicationIdArgument  = errors.New("A mandatory flag is missing. You must provide a domain_id as an argument or path to import the file. Run the command 'azion domains <subcommand> --help' to display more information and try again")
	ErrorGetDomain                     = errors.New("Failed to describe the domain: %s. Check your settings and try again. If the error persists, contact Azion support.")
	ErrorMissingDomainIdArgumentDelete = errors.New("A mandatory flag is missing. You must provide a domain_id as an argument. Run the command 'azion domains <subcommand> --help' to display more information and try again")

	ErrorFailToDeleteDomain     = errors.New("Failed to delete the Domain: %s. Check your settings and try again. If the error persists, contact Azion support")
	ErrorMissingCnames          = errors.New("Missing Cnames. When the flag '--cname-access-only`is set as 'true', at least one CNAME must be provided through the flag '--cnames'. Add one or more CNAMES, or set '--cname-access-only' as false and try again.")
	ErrorDigitalCertificateFlag = errors.New("Invalid --digital-certificate-id flag provided. The flag must have an Integer or 'null' as a value. Run the command 'azion domains <subcommand> --help' to display more information and try again")
)
