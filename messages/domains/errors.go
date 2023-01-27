package domains

import "errors"

var (
	//used by more than one cmd
	DomainsFlagId                      = "Unique identifier of the Domain"
	DomainsFileWritten                 = "File successfully written to: %s\n"
	ErrorMissingApplicationIdArgument  = errors.New("A mandatory flag is missing. You must provide a domain_id as an argument or path to import the file. Run the command 'azioncli domains <subcommand> --help' to display more information and try again")
	ErrorGetDomain                     = errors.New("Failed to get the domain: %s. Check your settings and try again. If the error persists, contact Azion support")
	ErrorMissingDomainIdArgumentDelete = errors.New("A mandatory flag is missing. You must provide a domain_id as an argument. Run the command 'azioncli domains <subcommand> --help' to display more information and try again")
	ErrorFailToDeleteDomain            = errors.New("Failed to delete the Domain: %s. Check your settings and try again. If the error persists, contact Azion support")
	ErrorActiveFlag                    = errors.New("Invalid --active flag provided. The flag must have 'true' or 'false' values. Run the command 'azioncli domains <subcommand> --help' to display more information and try again")
	ErrorUpdateDomain                  = errors.New("Failed to update the Domain: %s. Check your settings and try again. If the error persists, contact Azion support")
)
