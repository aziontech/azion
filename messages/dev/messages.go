package dev

var (
	DevFlagHelp         = "Displays more information about the dev command"
	DevUsage            = "dev [flags]"
	DevShortDescription = "Starts a local development server for the current application"
	DevLongDescription  = "Starts a local development server for the current application, so it's possible to preview and test it locally before the deployment"
	IsFirewall          = "Indicates whether the function to be run is intended for the Edge Firewall"
	PortFlag            = "Indicates which port to use when starting localhost environment"
	SkipFrameworkBuild  = "Indicates whether to bypass the framework build phase before executing 'azion build'"
)
