package webapp

var (

	//used by more than one cmd
	WebappOutput = "\nCommand exited with code %d\n"

	//webapp cmd
	WebappUsage            = "webapp"
	WebappShortDescription = "Creates Web Applications on Azion's platform"
	WebappLongDescription  = "Build your Web applications in minutes without the need to manage infrastructure or security"
	WebappFlagHelp         = "Displays more information about the webapp command"

	//build cmd
	WebappBuildUsage            = "build [flags]"
	WebappBuildShortDescription = "Builds a Web Application"
	WebappBuildLongDescription  = "Builds your Web Application to run on Azion’s Edge Computing Platform"
	WebappBuildCmdNotSpecified  = "Build step command not specified. No action will be taken\n"
	WebappBuildRunningCmd       = "Running build step command:\n\n"
	WebappBuildFlagHelp         = "Displays more information about the build subcommand"

	//init cmd
	WebappInitUsage             = "init [flags]"
	WebappInitShortDescription  = "Initializes a Web Application"
	WebappInitLongDescription   = "Defines primary parameters based on a given name and application type to start a Web Application on Azion’s platform"
	WebappInitCmdNotSpecified   = "Init step command not specified. No action will be taken\n"
	WebappInitRunningCmd        = "Running init step command:\n\n"
	WebappInitFlagName          = "The Web application's name"
	WebappInitFlagType          = "The type of  Web application <javascript | flareact | nextjs>"
	WebappInitFlagYes           = "Forces the automatic response 'yes' to all user input"
	WebappInitFlagNo            = "Forces the automatic response 'no' to all user input"
	WebAppInitContentOverridden = "This project was already configured. Do you want to override the previous configuration? <yes | no> (default: no) "
	WebAppInitCmdSuccess        = "Template successfully fetched and configured"
	WebappInitFlagHelp          = "Displays more information about the init subcommand"

	//publish cmd
	WebappPublishUsage                       = "publish"
	WebappPublishShortDescription            = "Publishes a Web Application on the Azion platform"
	WebappPublishLongDescription             = "Publishes a Web Application based on the Azion’s Platform"
	WebappPublishCmdNotSpecified             = "Publish pre command not specified. No action will be taken\n"
	WebappPublishRunningCmd                  = "Running publish pre command:\n\n"
	WebappPublishOutputDomainSuccess         = "\nTo visualize your application access the domain: %s\n"
	WebappPublishOutputCachePurge            = "Domain cache was purged"
	WebappPublishOutputEdgeFunctionCreate    = "Created Edge Function %s with ID %d\n"
	WebappPublishOutputEdgeFunctionUpdate    = "Updated Edge Function %s with ID %d\n"
	WebappPublishOutputEdgeApplicationCreate = "Created Edge Application %s with ID %d\n"
	WebappPublishOutputEdgeApplicationUpdate = "Updated Edge Application %s with ID %d\n"
	WebappPublishOutputDomainCreate          = "Created Domain %s with ID %d\n"
	WebappPublishOutputDomainUpdate          = "Updated Domain %s with ID %d\n"
	WebappPublishFlagHelp                    = "Displays more information about the publish subcommand"
	WebappPublishPropagation                 = "Content is being propagated to all Azion POPs and it might take a few minutes for all edges to be up to date\n"
)
