package build

// Used by both the v3 and the v4 command trees.
var (
	BuildUsage            = "build [flags]"
	BuildShortDescription = "Builds an Application locally"
	BuildLongDescription  = "Builds an Application locally"
	BuildRunningCmd       = "Running build step command:\n\n"
	BuildStart            = "Building your Application. This process may take a few minutes\n"
	BuildSuccessful       = "Your Application was built successfully\n"
	BuildFlagHelp         = "Displays more information about the build command"
	FlagTemplate          = "The Application's preset; Inform this flag if you wish to change the project's preset during build"
	FlagWorker            = "Indicates that the constructed code inserts its own worker expression, such as addEventListener(\"fetch\") or similar, without the need to inject a provider"
	FlagPolyfill          = "Use node polyfills in build"
	FlagEntry             = "Code entrypoint; (default: ./main.js)"
	ProjectConfFlag       = "Relative path to where your custom azion.json and args.json files are stored"
	SkipFrameworkBuild    = "Indicates whether to bypass the framework build phase before executing 'azion build'."
	FlagAliasEnv          = "If sent, the --alias-env option is forwarded to Bundler during the build phase"
)

// Used only by the v3 command tree
var (
	IsFirewall = "Indicates whether the function to be run is intended for the Firewall"
)
