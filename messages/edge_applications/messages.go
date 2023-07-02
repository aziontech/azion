package edge_applications

var (

	//used by more than one cmd
	EdgeApplicationFlagId      = "Unique identifier of the Edge Application"
	EdgeApplicationFileWritten = "File successfully written to: %s\n"

	//edge_applications cmd
	EdgeApplicationsUsage            = "edge_applications"
	EdgeApplicationsShortDescription = "Creates Edge Applications on Azion's platform"
	EdgeApplicationsLongDescription  = "Build your Edge applications in minutes without the need to manage infrastructure or security"
	EdgeApplicationsFlagHelp         = "Displays more information about the edge_application command"
	EdgeApplicationsAutoDetectec     = "Auto-detected Project Settings (%s)\n"

	//build cmd
	EdgeApplicationsBuildUsage            = "build [flags]"
	EdgeApplicationsBuildShortDescription = "Builds an Edge Application"
	EdgeApplicationsBuildLongDescription  = "Builds your Edge Application to run on Azion’s Edge Computing Platform"
	EdgeApplicationsBuildRunningCmd       = "Running build step command:\n\n"
	EdgeApplicationsBuildStart            = "Building your Edge Application. This process may take a few minutes\n"
	EdgeApplicationsBuildSuccessful       = "Your Edge Application was built successfully\n"
	EdgeApplicationsBuildFlagHelp         = "Displays more information about the build subcommand"
	EdgeApplicationsBuildCdn              = "Skipping build step. Build isn't applied to the type 'CDN'"
	EdgeApplicationsBuildNotNecessary     = "Skipping build step. There were no changes detected in your project"

	UploadStart      = "Uploading static files\n"
	UploadSuccessful = "Upload completed successfully!\n"

	//init cmd
	EdgeApplicationsInitUsage = `init [flags]
	--type string       The type of Edge application
	cdn                 Create an edge application to cache and deliver your content.
	static          	Create a static page in application on edge.
	nextjs              Create a serverless NextJS edge-runtime application on edge.`
	EdgeApplicationsInitShortDescription  = "Initializes an Edge Application"
	EdgeApplicationsInitLongDescription   = "Defines primary parameters based on a given name and application type to start a Edge Application on Azion’s platform"
	EdgeApplicationsInitRunningCmd        = "Running init step command:\n\n"
	EdgeApplicationsInitFlagName          = "The Edge application's name"
	EdgeApplicationsInitFlagType          = "The type of Edge application <cdn | static | nextjs>"
	EdgeApplicationsInitFlagYes           = "Forces the automatic response 'yes' to all user input"
	EdgeApplicationsInitFlagNo            = "Forces the automatic response 'no' to all user input"
	WebAppInitContentOverridden           = "This project was already configured. Do you want to override the previous configuration? <yes | no> (default: no) "
	WebAppInitCmdSuccess                  = "Template successfully fetched and configured\n"
	EdgeApplicationsInitFlagHelp          = "Displays more information about the init subcommand"
	EdgeApplicationsInitSuccessful        = "Your project %s was initialized successfully"
	EdgeApplicationsInitNameNotSent       = "The Project Name was not sent through the --name flag; By default when --name is not informed the one found in your package.json file or working directory is used\n"
	EdgeApplicationsInitNameNotSentCdn    = "The project name was not sent by the --name flag; By default, when --name is not given, the working directory is used\n"
	EdgeApplicationsInitNameNotSentStatic = "The project name was not sent by the --name flag; By default, when --name is not given, the working directory is used\n"
	EdgeApplicationsUpdateNamePackageJson = "Updating your package.json name field with the value informed through the --name flag"
	EdgeApplicationsInitTypeNotSent       = "The Project Type was not sent through the --type flag; By default when --type is not informed it is auto-detected based on the framework used by the user\n"

	//publish cmd
	EdgeApplicationsPublishUsage                       = "publish"
	EdgeApplicationsPublishShortDescription            = "Publishes an Edge Application on the Azion platform"
	EdgeApplicationsPublishLongDescription             = "Publishes an Edge Application based on the Azion’s Platform"
	EdgeApplicationsPublishRunningCmd                  = "Running pre-publish command:\n\n"
	EdgeApplicationsPublishSuccessful                  = "Your Edge Application was published successfully\n"
	EdgeApplicationsCdnPublishSuccessful               = "Your CDN Edge Application was published successfully\n"
	EdgeApplicationsPublishOutputDomainSuccess         = "\nTo visualize your application access the domain: %s\n"
	EdgeApplicationPublishDomainHint                   = "You may now edit your domain and add your own cnames. To do this you may run 'azioncli domain update' command and also configure your DNS\n"
	EdgeApplicationsPublishOutputCachePurge            = "Domain cache was purged\n"
	EdgeApplicationsPublishOutputEdgeFunctionCreate    = "Created Edge Function %s with ID %d\n"
	EdgeApplicationsPublishOutputEdgeFunctionUpdate    = "Updated Edge Function %s with ID %d\n"
	EdgeApplicationsPublishOutputEdgeApplicationCreate = "Created Edge Application %s with ID %d\n"
	EdgeApplicationsPublishOutputEdgeApplicationUpdate = "Updated Edge Application %s with ID %d\n"
	EdgeApplicationsPublishOutputDomainCreate          = "Created Domain %s with ID %d\n"
	EdgeApplicationsPublishOutputDomainUpdate          = "Updated Domain %s with ID %d\n"
	EdgeApplicationPublishPathFlag                     = "Path to where your static files are stored"
	EdgeApplicationsCacheSettingsSuccessful            = "Created Cache Settings for web application"
	EdgeApplicationsPublishInputAddress                = "Please inform an address to be used in the origin of this application: "
	EdgeApplicationsRulesEngineSuccessful              = "Created Rules Engine for web application"
	EdgeApplicationsPublishFlagHelp                    = "Displays more information about the publish subcommand"
	EdgeApplicationsPublishPropagation                 = "Content is being propagated to all Azion POPs and it might take a few minutes for all edges to be up to date\n"
	EdgeApplicationPublishIgnoreFlag                   = "Files and directories to ignore when publishing the application; follows the gitignore pattern"

	//CRUD
	//list cmd
	EdgeApplicationsListUsage            = "list [flags]"
	EdgeApplicationsListShortDescription = "Displays your account's Edge Applications"
	EdgeApplicationsListLongDescription  = "Displays all Applications in the user account’s Edge Applications library"
	EdgeApplicationsListHelpFlag         = "Displays more information about the list subcommand"

	//describe cmd
	EdgeApplicationDescribeUsage            = "describe --application-id <application_id> [flags]"
	EdgeApplicationDescribeShortDescription = "Returns the Edge Application data"
	EdgeApplicationDescribeLongDescription  = "Displays information about the Edge Application via a given ID to show the application’s attributes in detail"
	EdgeApplicationDescribeFlagOut          = "Exports the output to the given <file_path/file_name.ext>"
	EdgeApplicationDescribeFlagFormat       = "Changes the output format passing the json value to the flag"
	EdgeApplicationDescribeHelpFlag         = "Displays more information about the describe command"

	//delete cmd
	EdgeApplicationDeleteUsage            = "delete --application-id <application_id> [flags]"
	EdgeApplicationDeleteShortDescription = "Removes an Edge Application"
	EdgeApplicationDeleteLongDescription  = "Removes an Edge Application from the Edge Applications library based on its given ID"
	EdgeApplicationDeleteOutputSuccess    = "Edge Application %d was successfully deleted\n"
	EdgeApplicationDeleteHelpFlag         = "Displays more information about the delete subcommand"
	EdgeApplicationDeleteCascadeFlag      = "Deletes all resources created through the command azioncli edge_applications publish"
	EdgeApplicationDeleteMissingFunction  = "Missing Edge Function ID in azion.json file. Skipping deletion"
	EdgeApplicationDeleteCascadeSuccess   = "Cascade delete carried out successfully"

	//update cmd
	EdgeApplicationUpdateUsage                       = "update --application-id <application_id> [flags]"
	EdgeApplicationUpdateShortDescription            = "Modifies an Edge Application"
	EdgeApplicationUpdateLongDescription             = "Modifies an Edge Application based on its ID to update its name, activity status, and other attributes"
	EdgeApplicationUpdateFlagName                    = "The Edge Application's name"
	EdgeApplicationUpdateFlagDeliveryProtocol        = "The Edge Application's Delivery Protocol"
	EdgeApplicationUpdateFlagHttpPort                = "The Edge Application's Http Port"
	EdgeApplicationUpdateFlagHttpsPort               = "The Edge Application's Https Port"
	EdgeApplicationUpdateFlagMinimumTlsVersion       = "The Edge Application's Minimum Tls Version"
	EdgeApplicationUpdateFlagApplicationAcceleration = "Whether the Edge Application has Application Acceleration active or not"
	EdgeApplicationUpdateFlagCaching                 = "Whether the Edge Application has Caching active or not"
	EdgeApplicationUpdateFlagDeviceDetection         = "Whether the Edge Application has Device Detection active or not"
	EdgeApplicationUpdateFlagEdgeFirewall            = "Whether the Edge Application has Edge Firewall active or not"
	EdgeApplicationUpdateFlagEdgeFunctions           = "Whether the Edge Application has Edge Functions active or not"
	EdgeApplicationUpdateFlagImageOptimization       = "Whether the Edge Application has Image Optimization active or not"
	EdgeApplicationUpdateFlagL2Caching               = "Whether the Edge Application has L2 Caching active or not"
	EdgeApplicationUpdateFlagLoadBalancer            = "Whether the Edge Application has Load Balancer active or not"
	EdgeApplicationUpdateRawLogs                     = "Whether the Edge Application has Raw Logs active or not"
	EdgeApplicationUpdateWebApplicationFirewall      = "Whether the Edge Application has Web Application Firewall active or not"
	EdgeApplicationUpdateFlagIn                      = "Given path and JSON file to automatically update the Edge Application attributes; you can use - for reading from stdin"
	EdgeApplicationUpdateOutputSuccess               = "Updated Edge Application with ID %d\n"
	EdgeApplicationUpdateHelpFlag                    = "Displays more information about the update subcommand"
)
