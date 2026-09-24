package deployremote

// Used by both the v3 and the v4 command trees.
var (
	// deploy cmd
	DeployUsage                          = "deploy-remote"
	DeployShortDescription               = "Deploys an Application"
	DeployLongDescription                = "Deploys an Application"
	DeploySuccessful                     = "Your Application was deployed successfully\n"
	DeployOutputWorkloadSuccess          = "\nTo visualize your application access the Domain: %s\n"
	EdgeApplicationDeployPathFlag        = "Path to where your static files are stored"
	EdgeApplicationDeployProjectConfFlag = "Relative path to where your custom azion.json and args.json files are stored"
	EdgeApplicationDeploySync            = "Synchronizes the local azion.json file with remote resources"
	EnvFlag                              = "Relative path to where your custom .env file is stored"
	AliasEnvFlag                         = "If sent, the --alias-env option is forwarded to Bundler during the build phase"
	DeployFlagHelp                       = "Displays more information about the deploy command"
	DeployFlagAuto                       = "If sent, the entire flow of the command will be run without interruptions"
	DeployFlagNoPrompt                   = "If sent, whenever the CLI would display an interactive prompt due to an error, it instead just returns the error"
	DeployPropagation                    = "Your application is being deployed to all Azion Edge Locations and it might take a few minutes.\n"
	SkipBucket                           = "Your project does not contain a '.edge/storage' folder. Skipping creation of bucket"
	SkipUpload                           = "Your project does not contain a '.edge/storage' folder. Skipping upload of static files"
	SkipUploadBuild                      = "Skipping upload of static files due to project not being built and no new static files being generated"
)

// Used only by the v3 command tree (bug-fix-only; see doc/plan.md).
var (
	DeployOutputDomainSuccess = "\nTo visualize your application access the Domain: %s\n"
	DeployFlagSkipBuild       = "If sent, the build command will not be called during the deploy process"
)
