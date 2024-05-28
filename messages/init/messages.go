package init

const (
	USAGE             = "init"
	SHORT_DESCRIPTION = "Initializes an Edge Application from a starter template"
	LONG_DESCRIPTION  = "Defines primary parameters based on a given name and application preset to start an Edge Application"
	EXAMPLE           = "$ azion init\n$ azion init --help\n$ azion init --name testproject"
	FLAG_NAME         = "The Edge Application's name"
	FLAG_PRESET       = "The Preset's name"
	FLAG_TEMPLATE     = "The Template's name"
)

var (
	EdgeApplicationsInitRunningCmd        = "Running init step command:\n\n"
	EdgeApplicationsInitFlagYes           = "Answers all yes/no interactions automatically with yes"
	EdgeApplicationsInitFlagNo            = "Answers all yes/no interactions automatically with no"
	WebAppInitContentOverridden           = "This application was already configured. Do you want to override the previous configuration? <yes | no> (default: no) "
	WebAppInitCmdSuccess                  = "Template successfully fetched and configured\n"
	InitGettingTemplates                  = "\nGetting modes available (Some dependencies may need to be installed)\n"
	InitGettingVulcan                     = "Getting templates available\n"
	InitProjectQuestion                   = "Your application's name: "
	EdgeApplicationsInitFlagHelp          = "Displays more information about the init command"
	EdgeApplicationsInitSuccessful        = "Your application %s was initialized successfully\n"
	EdgeApplicationsInitNameNotSent       = "The application name was not sent through the --name flag; By default when --name is not informed the one found in your package.json file or working directory is used\n\n"
	EdgeApplicationsInitNameNotSentSimple = "The application name was not sent through the --name flag; By default, when --name is not given, the working directory is used\n"
	EdgeApplicationsInitNameNotSentStatic = "The application name was not sent through the --name flag; By default, when --name is not given, the working directory is used\n"
	EdgeApplicationsUpdateNamePackageJson = "Updating your package.json name field with the value informed through the --name flag\n"
	EdgeApplicationsInitTypeNotSent       = "The application preset was not sent through the --template flag; By default when --template is not informed it is auto-detected based on the framework used by the user\n\n"
	InitDevCommand                        = "If you want to start a local development server later, run 'azion dev'\n"
	InitDeployCommand                     = "If you want to deploy your application later, run 'azion deploy'\n"
	InitInstallDeps                       = "Installing application dependencies"
	AskDeploy                             = "Do you want to deploy your project? (y/N)"
	AskInstallDepsDev                     = "Do you want to install project dependencies? This may be required to start local development server (y/N)"
	AskInstallDepsDeploy                  = "Do you want to install project dependencies? This may be required to deploy your project (y/N)"
	AskLocalDev                           = "Do you want to start a local development server? (y/N)"
)
