package build

var (
	EdgeApplicationsBuildUsage            = "build [flags]"
	EdgeApplicationsBuildShortDescription = "Builds an Edge Application"
	EdgeApplicationsBuildLongDescription  = "Builds your Edge Application to run on Azion’s Edge Computing Platform"
	EdgeApplicationsBuildRunningCmd       = "Running build step command:\n\n"
	EdgeApplicationsBuildStart            = "Building your Edge Application. This process may take a few minutes\n"
	EdgeApplicationsBuildSuccessful       = "Your Edge Application was built successfully\n"
	EdgeApplicationsBuildFlagHelp         = "Displays more information about the build subcommand"
	EdgeApplicationsBuildSimple           = "Skipping build step. Build isn't applied to the type 'simple'\n"
	EdgeApplicationsBuildStatic           = "Skipping build step. Build isn't applied to the type 'static'\n"
	EdgeApplicationsBuildNotNecessary     = "Skipping build step. There were no changes detected in your project"
)
