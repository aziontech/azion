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

	//describe cmd
	DomainsDescribeUsage            = "describe --domain-id <domain_id> [flags]"
	DomainsDescribeShortDescription = "Returns the domain data"
	DomainsDescribeLongDescription  = "Displays information about the domain via a given ID to show the application’s attributes in detail"
	DomainsDescribeFlagOut          = "Exports the output to the given <file_path/file_name.ext>"
	DomainsDescribeFlagFormat       = "Changes the output format passing the json value to the flag"
	DomainsDescribeHelpFlag         = "Displays more information about the describe command"
)
