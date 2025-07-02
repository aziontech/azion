package build

var (
	BuildUsage            = "build [flags]"
	BuildShortDescription = "Builds an Edge Application locally"
	BuildLongDescription  = "Builds an Edge Application locally"
	BuildRunningCmd       = "Running build step command:\n\n"
	BuildStart            = "Building your Edge Application. This process may take a few minutes\n"
	BuildSuccessful       = "Your Edge Application was built successfully\n"
	BuildFlagHelp         = "Displays more information about the build command"
	BuildSimple           = "Skipping build step. Build isn't applied to this type\n"
	BuildStatic           = "Skipping build step. Build isn't applied to the type 'static'\n"
	BuildNotNecessary     = "Skipping build step. There were no changes detected in your project"
	FlagTemplate          = "The Edge Application's preset; Inform this flag if you wish to change the project's preset during build"
	FlagWorker            = "Indicates that the constructed code inserts its own worker expression, such as addEventListener(\"fetch\") or similar, without the need to inject a provider"
	FlagPolyfill          = "Use node polyfills in build"
	FlagEntry             = "Code entrypoint; (default: ./main.js)"
	ProjectConfFlag       = "Relative path to where your custom azion.json and args.json files are stored"
	SkipFrameworkBuild    = "Indicates whether to bypass the framework build phase before executing 'azion build'."
)
