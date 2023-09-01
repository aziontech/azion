package build

var (
	BuildUsage            = "build [flags]"
	BuildShortDescription = "Builds an Edge Application"
	BuildLongDescription  = "Builds your Edge Application to run on Azion’s Edge Computing Platform"
	BuildRunningCmd       = "Running build step command:\n\n"
	BuildStart            = "Building your Edge Application. This process may take a few minutes\n"
	BuildSuccessful       = "Your Edge Application was built successfully\n"
	BuildFlagHelp         = "Displays more information about the build command"
	BuildSimple           = "Skipping build step. Build isn't applied to the type 'simple'\n"
	BuildStatic           = "Skipping build step. Build isn't applied to the type 'static'\n"
	BuildNotNecessary     = "Skipping build step. There were no changes detected in your project"
	FlagTemplate          = "The Edge Application's preset; Inform this flag if you wish to change the project's preset during build"
	FlagMode              = "The Edge Application's mode; Inform this flag if you wish to change the project's mode during build"
)
