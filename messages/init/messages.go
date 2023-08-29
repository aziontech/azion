package init

var (

	//init cmd
	EdgeApplicationsInitUsage = `init [flags]
	--template string       The type of Edge application
	nextjs                  Create a serverless NextJS edge-runtime application on edge.`
	EdgeApplicationsInitShortDescription  = "Initializes an Edge Application"
	EdgeApplicationsInitLongDescription   = "Defines primary parameters based on a given name and application type to start an Edge Application on the Azion’s platform"
	EdgeApplicationsInitRunningCmd        = "Running init step command:\n\n"
	EdgeApplicationsInitFlagName          = "The Edge application's name"
	EdgeApplicationsInitFlagTemplate      = "The type of Edge Application"
	EdgeApplicationsInitFlagMode          = "The mode of Edge Application"
	EdgeApplicationsInitFlagYes           = "Forces the automatic response 'yes' to all user input"
	EdgeApplicationsInitFlagNo            = "Forces the automatic response 'no' to all user input"
	WebAppInitContentOverridden           = "This project was already configured. Do you want to override the previous configuration? <yes | no> (default: no) "
	WebAppInitCmdSuccess                  = "Template successfully fetched and configured\n\n"
	InitGettingTemplates                  = "Getting modes available"
	InitProjectQuestion                   = "(Hit enter to accept the suggested name in parenthesis) Your project's name: "
	EdgeApplicationsInitFlagHelp          = "Displays more information about the init subcommand"
	EdgeApplicationsInitSuccessful        = "Your project %s was initialized successfully\n"
	EdgeApplicationsInitNameNotSent       = "The Project Name was not sent through the --name flag; By default when --name is not informed the one found in your package.json file or working directory is used\n\n"
	EdgeApplicationsInitNameNotSentSimple = "The project name was not sent by the --name flag; By default, when --name is not given, the working directory is used\n"
	EdgeApplicationsInitNameNotSentStatic = "The project name was not sent by the --name flag; By default, when --name is not given, the working directory is used\n"
	EdgeApplicationsUpdateNamePackageJson = "Updating your package.json name field with the value informed through the --name flag\n"
	EdgeApplicationsInitTypeNotSent       = "The Project Type was not sent through the --template flag; By default when --template is not informed it is auto-detected based on the framework used by the user\n\n"
	InitDevCommand                        = "If you want to start a local development server later, run 'azion dev'\n"
	InitDeployCommand                     = "If you want to deploy your project later, run 'azion deploy'\n"
	InitInstallDeps                       = "Installing projct dependencies"
)
