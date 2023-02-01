package domains

var (
	//domains cmd
	DomainsUsage                   = "domains"
	DomainsShortDescription        = "Create domains for edges on Azion's platform"
	DomainsLongDescription         = "Build your Web applications in minutes without the need to manage infrastructure or security"
	DomainsFlagHelp                = "Displays more information about the domains command"
	DomainFlagId                   = "Unique identifier of the Domain"
	DomainFlagDigitalCertificateId = "Unique identifier of the Digital Certificate; this value is either an Integer or null"
	ApplicationFlagId              = "Unique identifier of the Edge Application used by this domain"

	//list cmd
	DomainsListUsage            = "list [flags]"
	DomainsListShortDescription = "Displays yours domains"
	DomainsListLongDescription  = "Displays all your domain references to your edges"
	DomainsListHelpFlag         = "Displays more information about the list subcommand"

	//create cmd
	DomainsCreateUsage                    = "create [flags]"
	DomainsCreateShortDescription         = "Makes a new domain"
	DomainsCreateLongDescription          = "Makes a Domain based on given attributes to be used in Edge Applications"
	DomainsCreateFlagName                 = "The Domain's name"
	DomainsCreateFlagCnames               = "List of domains names"
	DomainsCreateFlagCnameAccessOnly      = "Whether the domain is accessed only through the CNAMES or not"
	DomainsCreateFlagDigitalCertificateId = "The digital certificate's unique identifier. It can be an integer or null."
	DomainsCreateFlagEdgeApplicationId    = "The Edge Application's unique identifier"
	DomainsCreateFlagIsActive             = "Whether the Domain is active or not"
	DomainsCreateFlagIn                   = " Path to a JSON file containing the attributes that will be updated; you can use - for reading from stdin"
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

	//update cnd
	DomainUpdateUsage               = "update --domain-id <domain_id> [flags]"
	DomainUpdateShortDescription    = "Modifies a Domain"
	DomainUpdateLongDescription     = "Modifies a Domain based on its ID to update its name and other attributes"
	DomainUpdateFlagName            = "The Domain's name"
	DomainUpdateFlagCnames          = "Cnames of your Domain"
	DomainUpdateFlagCnameAccessOnly = "Whether the Domain should be Accessed through Cname only or not"
	DomainUpdateFlagIn              = "Given path and JSON file to automatically update the Domain attributes; you can use - for reading from stdin"
	DomainUpdateOutputSuccess       = "Updated Domain with ID %d\n"
	DomainUpdateFlagActive          = "Whether the Domain should be active or not"
	DomainUpdateHelpFlag            = "Displays more information about the update subcommand"
)
