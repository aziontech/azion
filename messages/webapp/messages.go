package webapp

var (
	//webapp cmd
	WebappUsage            = "webapp"
	WebappShortDescription = "Creates Web Applications on Azion's platform"
	WebappLongDescription  = "Build your Web applications in minutes without the need to manage infrastructure or security"
	WebappFlagHelp         = "Displays more information about the webapp command"
	WebappAutoDetectec     = "Auto-detected Project Settings (%s)\n"

	//build cmd
	WebappBuildUsage            = "build [flags]"
	WebappBuildShortDescription = "Builds a Web Application"
	WebappBuildLongDescription  = "Builds your Web Application to run on Azion’s Edge Computing Platform"
	WebappBuildRunningCmd       = "Running build step command:\n\n"
	WebappBuildStart            = "Building your Web Application\n"
	WebappBuildSuccessful       = "Your Web Application was built successfully\n"
	WebappBuildFlagHelp         = "Displays more information about the build subcommand"

	//init cmd
	WebappInitUsage             = "init [flags]"
	WebappInitShortDescription  = "Initializes a Web Application"
	WebappInitLongDescription   = "Defines primary parameters based on a given name and application type to start a Web Application on Azion’s platform"
	WebappInitRunningCmd        = "Running init step command:\n\n"
	WebappInitFlagName          = "The Web application's name"
	WebappInitFlagType          = "The type of  Web application <javascript | flareact | nextjs>"
	WebappInitFlagYes           = "Forces the automatic response 'yes' to all user input"
	WebappInitFlagNo            = "Forces the automatic response 'no' to all user input"
	WebAppInitContentOverridden = "This project was already configured. Do you want to override the previous configuration? <yes | no> (default: no) "
	WebAppInitCmdSuccess        = "Template successfully fetched and configured\n"
	WebappInitFlagHelp          = "Displays more information about the init subcommand"
	WebappInitSuccessful        = "Your project %s was initialized successfully"

	//publish cmd
	WebappPublishUsage                       = "publish"
	WebappPublishShortDescription            = "Publishes a Web Application on the Azion platform"
	WebappPublishLongDescription             = "Publishes a Web Application based on the Azion’s Platform"
	WebappPublishRunningCmd                  = "Running pre publish command:\n\n"
	WebappPublishSuccessful                  = "Your Web Application was published successfully\n"
	WebappPublishOutputDomainSuccess         = "\nTo visualize your application access the domain: %s\n"
	WebappPublishOutputCachePurge            = "Domain cache was purged\n"
	WebappPublishOutputEdgeFunctionCreate    = "Created Edge Function %s with ID %d\n"
	WebappPublishOutputEdgeFunctionUpdate    = "Updated Edge Function %s with ID %d\n"
	WebappPublishOutputEdgeApplicationCreate = "Created Edge Application %s with ID %d\n"
	WebappPublishOutputEdgeApplicationUpdate = "Updated Edge Application %s with ID %d\n"
	WebappPublishOutputDomainCreate          = "Created Domain %s with ID %d\n"
	WebappPublishOutputDomainUpdate          = "Updated Domain %s with ID %d\n"
	WebappPublishFlagHelp                    = "Displays more information about the publish subcommand"
	WebappPublishPropagation                 = "Content is being propagated to all Azion POPs and it might take a few minutes for all edges to be up to date\n"
)
