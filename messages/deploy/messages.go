package deploy

var (
	// publish cmd
	EdgeApplicationsPublishUsage                       = "deploy"
	EdgeApplicationsPublishShortDescription            = "Publishes an Edge Application on the Azion platform"
	EdgeApplicationsPublishLongDescription             = "Publishes an Edge Application based on the Azion’s Platform"
	EdgeApplicationsPublishRunningCmd                  = "Running pre-publish command:\n\n"
	EdgeApplicationsPublishSuccessful                  = "Your Edge Application was published successfully\n"
	EdgeApplicationsSimplePublishSuccessful            = "Your Simple Edge Application was published successfully\n"
	EdgeApplicationsPublishOutputDomainSuccess         = "\nTo visualize your application access the domain: %v\n"
	EdgeApplicationPublishDomainHint                   = "You may now edit your domain and add your own cnames. To do this you may run 'azioncli domain update' command and also configure your DNS\n"
	EdgeApplicationsPublishOutputCachePurge            = "Domain cache was purged\n"
	EdgeApplicationsPublishOutputEdgeFunctionCreate    = "Created Edge Function %v with ID %v\n"
	EdgeApplicationsPublishOutputEdgeFunctionUpdate    = "Updated Edge Function %v with ID %v\n"
	EdgeApplicationsPublishOutputEdgeApplicationCreate = "Created Edge Application %v with ID %v\n"
	EdgeApplicationsPublishOutputEdgeApplicationUpdate = "Updated Edge Application %v with ID %v\n"
	EdgeApplicationsPublishOutputDomainCreate          = "Created Domain %v with ID %v\n"
	EdgeApplicationsPublishOutputDomainUpdate          = "Updated Domain %v with ID %v\n"
	EdgeApplicationPublishPathFlag                     = "Path to where your static files are stored"
	EdgeApplicationsCacheSettingsSuccessful            = "Created Cache Settings for web application\n"
	EdgeApplicationsPublishInputAddress                = "Please inform an address to be used in the origin of this application: "
	EdgeApplicationsRulesEngineSuccessful              = "Created Rules Engine for web application\n"
	EdgeApplicationsPublishFlagHelp                    = "Displays more information about the publish subcommand"
	EdgeApplicationsPublishPropagation                 = "Content is being propagated to all Azion POPs and it might take a few minutes for all edges to be up to date\n"
	UploadStart                                        = "Uploading static files\n"
	UploadSuccessful                                   = "\nUpload completed successfully!\n"
)
