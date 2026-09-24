package init

// Used by both the v3 and the v4 command trees.
var (
	USAGE                          = "init"
	SHORT_DESCRIPTION              = "Initializes an Application from a starter template"
	LONG_DESCRIPTION               = "Defines primary parameters based on a given name and application preset to start an Application"
	EXAMPLE                        = "$ azion init\n$ azion init --help\n$ azion init --name testproject"
	FLAG_NAME                      = "The Application's name"
	FLAG_PACKAGE_MANAGE            = "Specify the package manager to use (e.g., npm, yarn, pnpm)"
	FLAG_AUTO                      = "If sent, the entire flow of the command will be run without interruptions"
	WebAppInitCmdSuccess           = "Template successfully configured\n"
	InitGettingTemplates           = "\nGetting presets available (Some dependencies may need to be installed)\n"
	InitProjectQuestion            = "Your application's name: "
	EdgeApplicationsInitSuccessful = "Your application %s was initialized successfully\n"
	InitDevCommand                 = "If you want to start a local development server later, run 'azion dev'\n"
	InitDeployCommand              = "If you want to deploy your application later, run 'azion deploy'\n"
	InstallDeps                    = "Installing application dependencies\n"
	AskDeploy                      = "Do you want to deploy your project? (y/N)"
	AskInstallDepsDev              = "Do you want to install project dependencies? This may be required to start local development server (Y/n)"
	AskLocalDev                    = "Do you want to start a local development server? (y/N)"
	ChangeWorkingDir               = "Make sure to change to the new working directory before running building or deploying your project\n"
)

// Used only by the v4 command tree.
var (
	FLAG_SYNC            = "Synchronizes the local azion.json file with remote resources. Use this flag when deploying your project from this command"
	FLAG_LOCAL           = "Runs the entire build and deploy process locally. Use this flag when deploying your project from this command"
	FLAG_CONFIG_DIR      = "Relative path to where your custom azion.json and args.json files are stored"
	AskInstallDepsBuild  = "Do you want to install project dependencies? This may be required to generate initial configuration file (Y/n)"
	AskInstallDepsDeploy = "Do you want to install project dependencies? This may be required to deploy your project (Y/n)"
	SkipFrameworkBuild   = "Indicates whether to bypass the framework build phase before executing 'azion build'"
	AliasEnvFlag         = "If sent, the --alias-env option is forwarded to Bundler during the build phase"
)
