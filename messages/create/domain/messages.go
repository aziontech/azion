package domain

var (
	Usage                    = "domain"
	ShortDescription         = "Creates a new domain"
	LongDescription          = "Creates a domain based on given attributes to be used with edge applications"
	FlagName                 = "The domain's name"
	FlagCnames               = "List of domains' CNAMES"
	FlagCnameAccessOnly      = "Whether the domain is accessed only through the CNAMES or not"
	FlagDigitalCertificateID = "The digital certificate's unique identifier. It can be an integer or null."
	FlagEdgeApplicationId    = "The edge application's unique identifier"
	FlagIsActive             = "Whether the Domain is active or not"
	FlagFile                 = "Path to a JSON file containing the attributes that will be updated; you can use - for reading from stdin"
	OutputSuccess            = "Created domain with ID %d\n"
	HelpFlag                 = "Displays more information about the create domain command"
	AskInputApplicationID    = "What is the ID of the edge application that the domain will be connected to?"
	AskInputName             = "What will the domain name be?"
)
