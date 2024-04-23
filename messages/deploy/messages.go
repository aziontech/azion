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
	OriginsSuccessful                 = "Created Origin for edge application\n"
	OriginsUpdateSuccessful           = "Updated Origin for edge application %v with ID %v \n"
	CacheSettingsSuccessful           = "Created Cache Settings for Edge Application\n"
	RulesEngineSuccessful             = "Created Rules Engine for Edge Application\n"
	DeployFlagHelp                    = "Displays more information about the deploy command"
	DeployFlagAuto                    = "If sent, the entire flow of the command will be run without interruptions"
	DeployFlagNoPrompt                = "If sent, whenever the CLI would display an interactive prompt due to an error, it instead just returns the error"
	DeployPropagation                 = "Your application is being deployed to all Azion Edge Locations and it might take a few minutes.\n"
	UploadStart                       = "Uploading static files\n"
	UploadSuccessful                  = "\nUpload completed successfully!\n"
	BucketInUse                       = "This bucket's name is already in use, please try another one\n"
	AppInUse                          = "This edge application's name is already in use, please try another one\n"
	DomainInUse                       = "This domain's name is already in use, please try another one\n"
	FuncInUse                         = "This edge function's name is already in use, please try another one\n"
	FuncInstInUse                     = "This function instance's name is already in use, please try another one\n"
	AskInputName                      = "Type the new name:"
	ProjectNameMessage                = "Using the same name as your project to create the bucket\n"
	AskCreateCacheSettings            = `Azion CLI offers to create the following Cache Settings specifications:
  - Browser Cache Settings: Override Cache Settings
  - Maximum TTL for Browser Cache Settings (in seconds): 7200
  - CDN Cache Settings: Override Cache Settings
  - Maximum TTL for CDN Cache Settings (in seconds): 7200

Do you wish to create a Cache Settings configuration with the above specifications? (y/N)`
	SkipUpload = "Your project does not contain a '.edge/storage' folder. Skipping upload of static files"
)
