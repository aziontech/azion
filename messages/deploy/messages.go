package deploy

var (
	// deploy cmd
	DeployUsage                       = "deploy"
	DeployShortDescription            = "Deploys a project on the Azion platform"
	DeployLongDescription             = "Deploys an Edge Application based on the Azion’s Platform"
	DeploySuccessful                  = "Your Edge Application was deployed successfully\n"
	SimpleDeploySuccessful            = "Your Simple Edge Application was deployed successfully\n"
	DeployOutputDomainSuccess         = "\nTo visualize your application access the domain: %v\n"
	EdgeApplicationDeployDomainHint   = "You may now edit your domain and add your own cnames. To do this you may run 'azioncli domain update' command and also configure your DNS\n"
	DeployOutputCachePurge            = "Domain cache was purged\n"
	DeployOutputEdgeFunctionCreate    = "Created Edge Function %v with ID %v\n"
	DeployOutputEdgeFunctionUpdate    = "Updated Edge Function %v with ID %v\n"
	DeployOutputEdgeApplicationCreate = "Created Edge Application %v with ID %v\n"
	DeployOutputEdgeApplicationUpdate = "Updated Edge Application %v with ID %v\n"
	DeployOutputDomainCreate          = "Created Domain %v with ID %v\n"
	DeployOutputDomainUpdate          = "Updated Domain %v with ID %v\n"
	EdgeApplicationDeployPathFlag     = "Path to where your static files are stored"
	CacheSettingsSuccessful           = "Created Cache Settings for web application\n"
	DeployInputAddress                = "Please inform an address to be used in the origin of this application: "
	RulesEngineSuccessful             = "Created Rules Engine for web application\n"
	DeployFlagHelp                    = "Displays more information about the deploy command"
	DeployPropagation                 = "Content is being propagated to all Azion POPs and it might take a few minutes for all edges to be up to date\n"
	UploadStart                       = "Uploading static files\n"
	UploadSuccessful                  = "\nUpload completed successfully!\n"
)
