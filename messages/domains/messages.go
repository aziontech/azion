package domains

var (
	//domains cmd
	DomainsUsage            = "domains"
	DomainsShortDescription = "Create domains for edges on Azion's platform"
	DomainsLongDescription  = "Build your Web applications in minutes without the need to manage infrastructure or security"
	DomainsFlagHelp         = "Displays more information about the domains command"
	DomainFlagId            = "Unique identifier of the Domain"

	//list cmd
	DomainsListUsage            = "list [flags]"
	DomainsListShortDescription = "Displays yours domains"
	DomainsListLongDescription  = "Displays all your domain references to your edges"
	DomainsListHelpFlag         = "Displays more information about the list subcommand"

	//create cmd
	DomainsCreateUsage                    = "create [flags]"
	DomainsCreateShortDescription         = "Makes a new domain"
	DomainsCreateLongDescription          = "Makes a Domain based on given attributes to be used in Edge Applications"
	DomainsCreateFlagName                 = "The Domain name"
	DomainsCreateFlagCnames               = "Domain Cnames"
	DomainsCreateFlagCnameAccessOnly      = "Whether the Domain should be accessed only through the Cnames or not"
	DomainsCreateFlagDigitalCertificateId = "Unique identifier of the Digital Certificate"
	DomainsCreateFlagEdgeApplicationId    = "Unique identifier of the Edge Application using this Domain"
	DomainsCreateFlagIsActive             = "Whether the Domain should be active or not"
	DomainsCreateFlagIn                   = "Given file path to create an Domain; you can use - for reading from stdin"
	DomainsCreateOutputSuccess            = "Created domain with ID %d\n"
	DomainsCreateHelpFlag                 = "Displays more information about the create subcommand"

	//describe cmd
	DomainsDescribeUsage            = "describe --domain-id <domain_id> [flags]"
	DomainsDescribeShortDescription = "Returns the domain data"
	DomainsDescribeLongDescription  = "Displays information about the domain via a given ID to show the application’s attributes in detail"
	DomainsDescribeFlagOut          = "Exports the output to the given <file_path/file_name.ext>"
	DomainsDescribeFlagFormat       = "Changes the output format passing the json value to the flag"
	DomainsDescribeHelpFlag         = "Displays more information about the describe command"

	//delete cmd
	DomainDeleteUsage            = "delete --domain-id <domain_id> [flags]"
	DomainDeleteShortDescription = "Removes a Domain"
	DomainDeleteLongDescription  = "Removes a Domain from the Domains library based on its given ID"
	DomainDeleteOutputSuccess    = "Domain %d was successfully deleted\n"
	DomainDeleteHelpFlag         = "Displays more information about the delete subcommand"
)
