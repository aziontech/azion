package deploy

var (
	// deploy cmd
	DeployUsage                       = "deploy"
	DeployShortDescription            = "Deploys an Edge Application"
	DeployLongDescription             = "Deploys an Edge Application"
	DeploySuccessful                  = "Your Edge Application was deployed successfully\n"
	SimpleDeploySuccessful            = "Your simple Edge Application was deployed successfully\n"
	DeployOutputDomainSuccess         = "\nTo visualize your application access the Domain: %v\n"
	EdgeApplicationDeployDomainHint   = "You may now edit your Domain and add your own CNAMES. To do this you may run 'azion domain update' command and also configure your DNS\n"
	DeployOutputCachePurge            = "Domain cache was purged\n"
	DeployOutputEdgeFunctionCreate    = "Created Edge Function %v with ID %v\n"
	DeployOutputEdgeFunctionUpdate    = "Updated Edge Function %v with ID %v\n"
	DeployOutputEdgeApplicationCreate = "Created Edge Application %v with ID %v\n"
	DeployOutputEdgeApplicationUpdate = "Updated Edge Application %v with ID %v\n"
	DeployOutputDomainCreate          = "Created Domain %v with ID %v\n"
	DeployOutputDomainUpdate          = "Updated Domain %v with ID %v\n"
	EdgeApplicationDeployPathFlag     = "Path to where your static files are stored"
	CacheSettingsSuccessful           = "Created Cache Settings for Edge Application\n"
	DeployInputAddress                = "Please inform an address to be used in the Origin of this application: "
	RulesEngineSuccessful             = "Created Rules Engine for Edge Application\n"
	DeployFlagHelp                    = "Displays more information about the deploy command"
	DeployPropagation                 = "Your application is being deployed to all Azion Edge Locations and it might take a few minutes.\n"
	UploadStart                       = "Uploading static files\n"
	UploadSuccessful                  = "\nUpload completed successfully!\n"
)
