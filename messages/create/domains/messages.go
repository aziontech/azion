package domains

var (
	Usage                    = "domains [flags]"
	ShortDescription         = "Creates a new domain"
	LongDescription          = "Creates a Domain based on given attributes to be used in Edge Applications"
	FlagName                 = "The Domain's name"
	FlagCnames               = "List of domains' names"
	FlagCnameAccessOnly      = "Whether the domain is accessed only through the CNAMES or not"
	FlagDigitalCertificateID = "The digital certificate's unique identifier. It can be an integer or null."
	FlagEdgeApplicationId    = "The Edge Application's unique identifier"
	FlagIsActive             = "Whether the Domain is active or not"
	FlagIn                   = "Path to a JSON file containing the attributes that will be updated; you can use - for reading from stdin"
	OutputSuccess            = "Created domain with ID %d\n"
	HelpFlag                 = "Displays more information about the create subcommand"
	AskInputApplicationID    = "What is the ID of the Edge Application that the Rule Engine will be connected to?"
	AskInputName             = "What will the domain name be?"
)
