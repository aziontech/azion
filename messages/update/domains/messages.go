package domains

var (
	Usage                    = "domains --domain-id <domain_id> [flags]"
	ShortDescription         = "Updates a Domain"
	LongDescription          = "Updates a Domain based on its ID to update its name and other attributes"
	FlagDomainID             = "The '--domain-id'"
	FlagDigitalCertificateID = "Unique identifier of the Digital Certificate; this value is either an integer or null"
	FlagApplicationID        = "Unique identifier for an edge application used by this domain."
	FlagName                 = "The Domain's name"
	FlagCnames               = "Cnames of your Domain"
	FlagCnameAccessOnly      = "Whether the Domain should be Accessed through Cname only or not"
	FlagIn                   = "Given path and JSON file to automatically update the Domain attributes; you can use - for reading from stdin"
	OutputSuccess            = "Updated Domain with ID %d\n"
	FlagActive               = "Whether the Domain should be active or not"
	HelpFlag                 = "Displays more information about the update domains subcommand"
	AskInputDomainID         = "What is the ID of the domain?"
)
