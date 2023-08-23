package link

var (

	//link cmd
	EdgeApplicationsLinkUsage             = "link [flags]"
	EdgeApplicationsLinkShortDescription  = "Links a project to Azion"
	EdgeApplicationsLinkLongDescription   = "Defines primary parameters based on a given name and application type to link a Project on the Azion’s platform"
	EdgeApplicationsLinkRunningCmd        = "Running link step command:\n\n"
	EdgeApplicationsLinkFlagName          = "The Edge application's name"
	EdgeApplicationsLinkFlagTemplate      = "The type of Edge Application"
	EdgeApplicationsLinkFlagMode          = "The mode of Edge Application"
	WebAppLinkCmdSuccess                  = "Template successfully fetched and configured\n\n"
	LinkGettingTemplates                  = "Getting templates available"
	LinkProjectQuestion                   = "(Hit enter to accept the suggested name in parenthesis) Your project's name: "
	EdgeApplicationsLinkFlagHelp          = "Displays more information about the link command"
	EdgeApplicationsLinkSuccessful        = "Your project %s was linked successfully\n"
	EdgeApplicationsLinkNameNotSent       = "The Project Name was not sent through the --name flag; By default when --name is not informed the one found in your package.json file or working directory is used\n\n"
	EdgeApplicationsLinkNameNotSentSimple = "The project name was not sent by the --name flag; By default, when --name is not given, the working directory is used\n"
	EdgeApplicationsLinkNameNotSentStatic = "The project name was not sent by the --name flag; By default, when --name is not given, the working directory is used\n"
)
