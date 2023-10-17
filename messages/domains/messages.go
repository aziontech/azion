package domains

var (
	//domains cmd
	DomainsUsage                   = "domains"
	DomainsShortDescription        = "Create domains for edges on Azion's platform"
	DomainsLongDescription         = "Build your Web applications in minutes without the need to manage infrastructure or security"
	DomainsFlagHelp                = "Displays more information about the domains command"
	DomainFlagId                   = "Unique identifier of the Domain"
	DomainFlagDigitalCertificateId = "Unique identifier of the Digital Certificate; this value is either an integer or null"
	ApplicationFlagId              = "Unique identifier for an edge application used by this domain.. The '--application-id' flag is mandatory"

	//describe cmd
	DomainsDescribeUsage            = "describe --domain-id <domain_id> [flags]"
	DomainsDescribeShortDescription = "Returns the domain data"
	DomainsDescribeLongDescription  = "Displays information about the domain via a given ID to show the application’s attributes in detail"
	DomainsDescribeFlagOut          = "Exports the output to the given <file_path/file_name.ext>"
	DomainsDescribeFlagFormat       = "Changes the output format passing the json value to the flag"
	DomainsDescribeHelpFlag         = "Displays more information about the describe command"

	//update cnd
	DomainUpdateUsage               = "update --domain-id <domain_id> [flags]"
	DomainUpdateShortDescription    = "Updates a Domain"
	DomainUpdateLongDescription     = "Updates a Domain based on its ID to update its name and other attributes"
	DomainUpdateFlagName            = "The Domain's name"
	DomainUpdateFlagCnames          = "Cnames of your Domain"
	DomainUpdateFlagCnameAccessOnly = "Whether the Domain should be Accessed through Cname only or not"
	DomainUpdateFlagIn              = "Given path and JSON file to automatically update the Domain attributes; you can use - for reading from stdin"
	DomainUpdateOutputSuccess       = "Updated Domain with ID %d\n"
	DomainUpdateFlagActive          = "Whether the Domain should be active or not"
	DomainUpdateHelpFlag            = "Displays more information about the update subcommand"
)
