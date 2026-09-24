package dev

// Used by both the v3 and the v4 command trees.
var (
	DevFlagHelp         = "Displays more information about the dev command"
	DevUsage            = "dev [flags]"
	DevShortDescription = "Starts a local development server for the current application"
	DevLongDescription  = "Starts a local development server for the current application, so it's possible to preview and test it locally before the deployment"
	PortFlag            = "Indicates which port to use when starting localhost environment"
	SkipFrameworkBuild  = "Indicates whether to bypass the framework build phase"
)

// Used only by the v3 command tree (bug-fix-only; see doc/plan.md).
var (
	IsFirewall = "Indicates whether the function to be run is intended for the Firewall"
)
