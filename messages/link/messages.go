package link

var (

	//link cmd
	EdgeApplicationsLinkUsage             = "link [flags]"
	EdgeApplicationsLinkShortDescription  = "Links a local repo or project folder to an existing application on Azion"
	EdgeApplicationsLinkLongDescription   = "Defines primary parameters based on a given name and application preset to link a Project to an Azion Edge Application"
	EdgeApplicationsLinkRunningCmd        = "Running link step command:\n\n"
	EdgeApplicationsLinkFlagName          = "The Edge Application's name"
	EdgeApplicationsLinkFlagTemplate      = "The Edge Application's template"
	EdgeApplicationsLinkFlagMode          = "The Edge Application's mode"
	WebAppLinkCmdSuccess                  = "Template successfully fetched and configured\n\n"
	LinkGettingTemplates                  = "Getting templates available\n"
	LinkProjectQuestion                   = "(Hit enter to accept the suggested name in parenthesis) Your application's name: "
	EdgeApplicationsLinkFlagHelp          = "Displays more information about the link command"
	EdgeApplicationsLinkSuccessful        = "Your application %s was linked successfully\n"
	EdgeApplicationsLinkNameNotSent       = "The application name was not sent through the --name flag; By default when --name is not informed the one found in your package.json file or working directory is used\n\n"
	EdgeApplicationsLinkNameNotSentSimple = "The application name was not sent through the --name flag; By default, when --name is not given, the working directory is used\n"
	EdgeApplicationsLinkNameNotSentStatic = "The application name was not sent through the --name flag; By default, when --name is not given, the working directory is used\n"
	LinkDevCommand                        = "If you want to start a local development server later, run 'azion dev'\n"
	LinkDeployCommand                     = "If you want to deploy your application later, run 'azion deploy'\n"
	LinkFlagAuto                          = "If sent, the entire flow of the command will be run without interruptions"
)
