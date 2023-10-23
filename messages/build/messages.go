package build

var (
	BuildUsage            = "build [flags]"
	BuildShortDescription = "Builds an edge application locally"
	BuildLongDescription  = "Builds an edge application locally to run on Azion’s Edge Computing Platform"
	BuildRunningCmd       = "Running build step command:\n\n"
	BuildStart            = "Building your edge application. This process may take a few minutes\n"
	BuildSuccessful       = "Your edge application was built successfully\n"
	BuildFlagHelp         = "Displays more information about the build command"
	BuildSimple           = "Skipping build step. Build isn't applied to the type 'simple'\n"
	BuildStatic           = "Skipping build step. Build isn't applied to the type 'static'\n"
	BuildNotNecessary     = "Skipping build step. There were no changes detected in your project"
	FlagTemplate          = "The edge application's preset; Inform this flag if you wish to change the project's preset during build"
	FlagMode              = "The edge application's mode; Inform this flag if you wish to change the project's mode during build"
)
