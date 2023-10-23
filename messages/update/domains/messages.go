package domains

var (
	Usage                    = "domains --domain-id <domain_id> [flags]"
	ShortDescription         = "Updates a domain"
	LongDescription          = "Updates a domain's name and other attributes based on a given ID"
	FlagDomainID             = "The '--domain-id'"
	FlagDigitalCertificateID = "Unique identifier of the Digital Certificate; this value is either an integer or null"
	FlagApplicationID        = "Unique identifier for an edge application used by this domain."
	FlagName                 = "The domain's name"
	FlagCnames               = "CNAMEs of your domain"
	FlagCnameAccessOnly      = "Whether the domain should be Accessed only through CNAMEs or not"
	FlagIn                   = "Given path and JSON file to automatically update the domain attributes; you can use - for reading from stdin"
	OutputSuccess            = "Updated domain with ID %d\n"
	FlagActive               = "Whether the domain should be active or not"
	HelpFlag                 = "Displays more information about the update domains subcommand"
	AskInputDomainID         = "What is the ID of the domain?"
)
