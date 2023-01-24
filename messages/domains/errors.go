package domains

import "errors"

var (
	ErrorGetDomains = errors.New("Failed to get the Domains: %s. Check your settings and try again. If the error persists, contact Azion support")

	ErrorMissingDomainIdArgumentDelete = errors.New("A mandatory flag is missing. You must provide a domain_id as an argument. Run the command 'azioncli domains <subcommand> --help' to display more information and try again")
	ErrorFailToDeleteDomain            = errors.New("Failed to delete the Domain: %s. Check your settings and try again. If the error persists, contact Azion support")
)
