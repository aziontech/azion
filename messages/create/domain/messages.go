package domain

var (
	Usage                    = "domain"
	ShortDescription         = "Creates a new Domain"
	LongDescription          = "Creates a Domain based on given attributes to be used with Edge Applications"
	FlagName                 = "The Domain's name"
	FlagCnames               = "List of a Domain's CNAMES"
	FlagCnameAccessOnly      = "Whether the Domain is accessed only through CNAMES or not"
	FlagDigitalCertificateID = "The Digital Certificate's unique identifier. It can be an integer or null."
	FlagEdgeApplicationId    = "The Edge Application's unique identifier"
	FlagIsActive             = "Whether the Domain is active or not"
	FlagFile                 = "Path to a JSON file containing the attributes that will be updated; you can use - for reading from stdin"
	OutputSuccess            = "Created Domain with ID %d\n"
	HelpFlag                 = "Displays more information about the create domain command"
	AskInputApplicationID    = "Enter the ID of the Edge Application that the Domain will be connected to:"
	AskInputName             = "Enter the new Domain's name:"
)
