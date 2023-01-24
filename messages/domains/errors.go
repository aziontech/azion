package domains

import "errors"

var (
	//used by more than one cmd
	DomainsFlagId      = "Unique identifier of the Edge Application"
	DomainsFileWritten = "File successfully written to: %s\n"

	ErrorGetDomains                   = errors.New("Failed to get the Domains: %s. Check your settings and try again. If the error persists, contact Azion support")
	ErrorMissingApplicationIdArgument = errors.New("A mandatory flag is missing. You must provide an domain_id as an argument or path to import the file. Run the command 'azioncli domains <subcommand> --help' to display more information and try again")
	ErrorGetDomain                    = errors.New("Failed to get the domain: %s. Check your settings and try again. If the error persists, contact Azion support")
)
