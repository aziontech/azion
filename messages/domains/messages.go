package domains

var (
	//domains cmd
	DomainsUsage            = "domains"
	DomainsShortDescription = "Create domains for edges on Azion's platform"
	DomainsLongDescription  = "Build your Web applications in minutes without the need to manage infrastructure or security"
	DomainsFlagHelp         = "Displays more information about the domains command"

	//list cmd
	DomainsListUsage            = "list [flags]"
	DomainsListShortDescription = "Displays yours domains"
	DomainsListLongDescription  = "Displays all your domain references to your edges"
	DomainsListHelpFlag         = "Displays more information about the list subcommand"

	//create cmd
	DomainsCreateUsage                    = "create [flags]"
	DomainsCreateShortDescription         = "create new domain"
	DomainsCreateLongDescription          = "creates a new domain to be used in edges_applications"
	DomainsCreateFlagName                 = "The Domain name"
	DomainsCreateFlagCnames               = "multiples domains name"
	DomainsCreateFlagCnameAccessOnly      = "only access to the cnames"
	DomainsCreateFlagDigitalCertificateId = "field where to put the id of the digital certificate"
	DomainsCreateFlagEdgeApplicationId    = "edge applications reference field"
	DomainsCreateFlagIsActive             = "boolean field to check if it is active"
	DomainsCreateFlagIn                   = "Given file path to create an Domain; you can use - for reading from stdin"
	DomainsCreateOutputSuccess            = "Created domain with ID %d\n"
	DomainsCreateHelpFlag                 = "Displays more information about the create subcommand"
)
