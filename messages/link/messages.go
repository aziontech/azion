package link

const (
	FLAG_PACKAGE_MANAGE = "Specify the package manager to use (e.g., npm, yarn, pnpm)"
	FLAG_SYNC           = "Synchronizes the local azion.json file with remote resources. Use this flag when deploying your project from this command"
	FLAG_LOCAL          = "Runs the entire build and deploy process locally. Use this flag when deploying your project from this command"
)

var (
	//link cmd
	EdgeApplicationsLinkUsage             = "link [flags]"
	EdgeApplicationsLinkShortDescription  = "Creates configuration used to build and deploy applications on Azion"
	EdgeApplicationsLinkLongDescription   = "Defines primary parameters based on a given name and application preset to link a Project to an Azion Application"
	EdgeApplicationsLinkRunningCmd        = "Running link step command:\n\n"
	EdgeApplicationsLinkFlagName          = "The Application's name"
	EdgeApplicationsLinkFlagTemplate      = "The Application's template"
	WebAppLinkCmdSuccess                  = "Project successfully configured\n"
	LinkGettingTemplates                  = "Getting templates available\n"
	LinkProjectQuestion                   = "(Hit enter to accept the suggested name in parenthesis) Your application's name: "
	EdgeApplicationsLinkFlagHelp          = "Displays more information about the link command"
	EdgeApplicationsLinkSuccessful        = "Your application %s was linked successfully\n"
	InstallDeps                           = "Installing application dependencies\n"
	EdgeApplicationsLinkNameNotSent       = "The application name was not sent through the --name flag; By default when --name is not informed the one found in your package.json file or working directory is used\n\n"
	EdgeApplicationsLinkNameNotSentSimple = "The application name was not sent through the --name flag; By default, when --name is not given, the working directory is used\n"
	EdgeApplicationsLinkNameNotSentStatic = "The application name was not sent through the --name flag; By default, when --name is not given, the working directory is used\n"
	LinkDevCommand                        = "If you want to start a local development server later, run 'azion dev'\n"
	LinkDeployCommand                     = "If you want to deploy your application later, run 'azion deploy'\n"
	LinkFlagAuto                          = "If sent, the entire flow of the command will be run without interruptions"
	AskDeploy                             = "Do you want to deploy your project? (y/N)"
	AskInstallDepsDev                     = "Do you want to install project dependencies? This may be required to start local development server (Y/n)"
	AskInstallDepsDeploy                  = "Do you want to install project dependencies? This may be required to deploy your project (Y/n)"
	AskLocalDev                           = "Do you want to start a local development server? (y/N)"
	AskGitignore                          = "Azion CLI creates some files during the build process for internal use. Would you like to add these to your .gitignore file? (Y/n)"
	WrittenGitignore                      = "Sucessfully written to your .gitignore file\n"
	SkipFrameworkBuild                    = "Indicates whether to bypass the framework build phase before executing 'azion build'"
)

const (
	FLAG_REMOTE  = "Clones a remote repository to be linked to an Azion Application"
	FLAGPATHCONF = "Relative path to where your custom azion.json and args.json files are stored"
	ASKPREBUILD  = "Do you allow Azion to build your project in order to generate configuration files? (Y/n)"
	BUILDLATER   = "Please, remember to run azion build --preset [preset-name], in order to generate the necessary configuration files\n"
)
