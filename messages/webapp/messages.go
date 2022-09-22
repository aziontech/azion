package webapp

var (

	//used by more than one cmd
	WebappOutput = "\nCommand exited with code %d\n"

	//webapp cmd
	WebappUsage            = "webapp"
	WebappShortDescription = "Create Web Applications on Azion's platform"
	WebappLongDescription  = "Create Web Applications on Azion's platform"

	//build cmd
	WebappBuildUsage            = "build [flags]"
	WebappBuildShortDescription = "Build your Web application"
	WebappBuildLongDescription  = "Build your Web application"
	WebappBuildCmdNotSpecified  = "Build step command not specified. No action will be taken\n"
	WebappBuildRunningCmd       = "Running build step command:\n\n"

	//init cmd
	WebappInitUsage             = "init [flags]"
	WebappInitShortDescription  = "Use Azion templates along with your Web applications"
	WebappInitLongDescription   = "Use Azion templates along with your Web applications"
	WebappInitCmdNotSpecified   = "Init step command not specified. No action will be taken\n"
	WebappInitRunningCmd        = "Running init step command:\n\n"
	WebappInitFlagName          = "Your Web Application's name"
	WebappInitFlagType          = "Your Web Application's type <javascript|flareact|nextjs>"
	WebappInitFlagYes           = "Force yes to all user input"
	WebappInitFlagNo            = "Force no to all user input"
	WebAppInitContentOverridden = "This project was already configured. Do you want to override the previous configuration? <yes | no> (default: no) "
	WebAppInitCmdSuccess        = "Template successfully fetched and configured"

	//publish cmd
	WebappPublishUsage                       = "publish"
	WebappPublishShortDescription            = "Publish your Web Application to Azion"
	WebappPublishLongDescription             = "Publish your Web Application to Azion"
	WebappPublishCmdNotSpecified             = "Publish pre command not specified. No action will be taken\n"
	WebappPublishRunningCmd                  = "Running publish pre command:\n\n"
	WebappPublishOutputDomainSuccess         = "\nYour Domain name: %s\n"
	WebappPublishOutputCachePurge            = "Domain cache was purged"
	WebappPublishOutputEdgeFunctionCreate    = "Created Edge Function with ID %d\n"
	WebappPublishOutputEdgeFunctionUpdate    = "Updated Edge Function with ID %d\n"
	WebappPublishOutputEdgeApplicationCreate = "Created Edge Application with ID %d\n"
	WebappPublishOutputEdgeApplicationUpdate = "Updated Edge Application with ID %d\n"
	WebappPublishOutputDomainCreate          = "Created Domain with ID %d\n"
	WebappPublishOutputDomainUpdate          = "Updated Domain with ID %d\n"
	WebappPublishOutputRulesEngineUpdate     = "Updated Rules Engine with ID %d\n"
)
