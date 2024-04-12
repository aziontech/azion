package domain

var (
	Usage                    = "domain"
	ShortDescription         = "Updates a Domain"
	LongDescription          = "Updates a Domain's name and other attributes based on a given ID"
	FlagDomainID             = "The '--domain-id'"
	FlagDigitalCertificateID = "Unique identifier of the Digital Certificate; this value is either an integer or null"
	FlagApplicationID        = "Unique identifier for an Edge Application used by this Domain."
	FlagName                 = "The Domain's name"
	FlagCnames               = "CNAMEs of your Domain"
	FlagCnameAccessOnly      = "Whether the Domain should be Accessed only through CNAMEs or not"
	FlagFile                 = "Given path and JSON file to automatically update the Domain attributes; you can use - for reading from stdin"
	OutputSuccess            = "Updated Domain with ID %d"
	FlagActive               = "Whether the Domain should be active or not"
	HelpFlag                 = "Displays more information about the update domains subcommand"
	AskInputDomainID         = "Enter the Domain's ID:"
)
