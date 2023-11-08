package deploy

var (
	// deploy cmd
	DeployUsage                       = "deploy"
	DeployShortDescription            = "Deploys an edge application onto the Azion platform"
	DeployLongDescription             = "Deploys an edge application onto the Azion platform"
	DeploySuccessful                  = "Your edge application was deployed successfully\n"
	SimpleDeploySuccessful            = "Your simple edge application was deployed successfully\n"
	DeployOutputDomainSuccess         = "\nTo visualize your application access the domain: %v\n"
	EdgeApplicationDeployDomainHint   = "You may now edit your domain and add your own cnames. To do this you may run 'azion domain update' command and also configure your DNS\n"
	DeployOutputCachePurge            = "Domain cache was purged\n"
	DeployOutputEdgeFunctionCreate    = "Created edge function %v with ID %v\n"
	DeployOutputEdgeFunctionUpdate    = "Updated edge function %v with ID %v\n"
	DeployOutputEdgeApplicationCreate = "Created edge application %v with ID %v\n"
	DeployOutputEdgeApplicationUpdate = "Updated edge application %v with ID %v\n"
	DeployOutputDomainCreate          = "Created domain %v with ID %v\n"
	DeployOutputDomainUpdate          = "Updated domain %v with ID %v\n"
	EdgeApplicationDeployPathFlag     = "Path to where your static files are stored"
	CacheSettingsSuccessful           = "Created Cache Settings for edge application\n"
	DeployInputAddress                = "Please inform an address to be used in the origin of this application: "
	RulesEngineSuccessful             = "Created rules engine for edge application\n"
	DeployFlagHelp                    = "Displays more information about the deploy command"
	DeployPropagation                 = "Your application is being deployed to all Azion Edge Locations and it might take a few minutes.\n"
	UploadStart                       = "Uploading static files\n"
	UploadSuccessful                  = "\nUpload completed successfully!\n"
)
