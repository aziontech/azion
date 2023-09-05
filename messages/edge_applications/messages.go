package edge_applications

var (

	//used by more than one cmd
	FlagId      = "Unique identifier of the Edge Application"
	FileWritten = "File successfully written to: %s\n"

	//edge_applications cmd
	Usage            = "edge_applications"
	ShortDescription = "Creates Edge Applications on Azion's platform"
	LongDescription  = "Build your Edge applications in minutes without the need to manage infrastructure or security"
	FlagHelp         = "Displays more information about the edge_application command"
	AutoDetectec     = "Auto-detected Project Settings (%s)\n"

	//build cmd
	BuildUsage            = "build [flags]"
	BuildShortDescription = "Builds an Edge Application"
	BuildLongDescription  = "Builds your Edge Application to run on Azion’s Edge Computing Platform"
	BuildRunningCmd       = "Running build step command:\n\n"
	BuildStart            = "Building your Edge Application. This process may take a few minutes\n"
	BuildSuccessful       = "Your Edge Application was built successfully\n"
	BuildFlagHelp         = "Displays more information about the build subcommand"
	BuildSimple           = "Skipping build step. Build isn't applied to the type 'Simple'\n"
	BuildNotNecessary     = "Skipping build step. There were no changes detected in your project"

	UploadStart      = "Uploading static files\n"
	UploadSuccessful = "Upload completed successfully!\n"

	//init cmd
	InitUsage = `init [flags]
	--type string       The type of Edge application
	simple              Create an edge application to cache and deliver your content.
	static          	Create a static page in application on edge.
	nextjs              Create a serverless NextJS edge-runtime application on edge.`
	InitShortDescription        = "Initializes an Edge Application"
	InitLongDescription         = "Defines primary parameters based on a given name and application type to start an Edge Application on the Azion’s platform"
	InitRunningCmd              = "Running init step command:\n\n"
	InitFlagName                = "The Edge application's name"
	InitFlagTemplate            = "The type of Edge Application"
	InitFlagMode                = "The mode of Edge Application"
	InitFlagYes                 = "Forces the automatic response 'yes' to all user input"
	InitFlagNo                  = "Forces the automatic response 'no' to all user input"
	WebAppInitContentOverridden = "This project was already configured. Do you want to override the previous configuration? <yes | no> (default: no) "
	WebAppInitCmdSuccess        = "Template successfully fetched and configured\n\n"
	InitFlagHelp                = "Displays more information about the init subcommand"
	InitSuccessful              = "Your project %s was initialized successfully\n"
	InitNameNotSent             = "The Project Name was not sent through the --name flag; By default when --name is not informed the one found in your package.json file or working directory is used\n\n"
	InitNameNotSentSimple       = "The project name was not sent by the --name flag; By default, when --name is not given, the working directory is used\n"
	InitNameNotSentStatic       = "The project name was not sent by the --name flag; By default, when --name is not given, the working directory is used\n"
	UpdateNamePackageJson       = "Updating your package.json name field with the value informed through the --name flag\n"
	InitTypeNotSent             = "The Project Type was not sent through the --type flag; By default when --type is not informed it is auto-detected based on the framework used by the user\n\n"

	//publish cmd
	PublishUsage                    = "publish"
	PublishShortDescription         = "Publishes an Edge Application on the Azion platform"
	PublishLongDescription          = "Publishes an Edge Application based on the Azion’s Platform"
	PublishRunningCmd               = "Running pre-publish command:\n\n"
	PublishSuccessful               = "Your Edge Application was published successfully\n"
	SimplePublishSuccessful         = "Your Simple Edge Application was published successfully\n"
	PublishOutputDomainSuccess      = "\nTo visualize your application access the domain: %v\n"
	PublishDomainHint               = "You may now edit your domain and add your own cnames. To do this you may run 'azioncli domain update' command and also configure your DNS\n"
	PublishOutputCachePurge         = "Domain cache was purged\n"
	PublishOutputEdgeFunctionCreate = "Created Edge Function %v with ID %v\n"
	PublishOutputEdgeFunctionUpdate = "Updated Edge Function %v with ID %v\n"
	PublishOutputCreate             = "Created Edge Application %v with ID %v\n"
	PublishOutputUpdate             = "Updated Edge Application %v with ID %v\n"
	PublishOutputDomainCreate       = "Created Domain %v with ID %v\n"
	PublishOutputDomainUpdate       = "Updated Domain %v with ID %v\n"
	PublishPathFlag                 = "Path to where your static files are stored"
	CacheSettingsSuccessful         = "Created Cache Settings for web application\n"
	PublishInputAddress             = "Please inform an address to be used in the origin of this application: "
	RulesEngineSuccessful           = "Created Rules Engine for web application\n"
	PublishFlagHelp                 = "Displays more information about the publish subcommand"
	PublishPropagation              = "Content is being propagated to all Azion POPs and it might take a few minutes for all edges to be up to date\n"

	//CRUD
	//list cmd
	ListUsage            = "list [flags]"
	ListShortDescription = "Displays your account's Edge Applications"
	ListLongDescription  = "Displays all Applications in the user account’s Edge Applications library"
	ListHelpFlag         = "Displays more information about the list subcommand"

	//describe cmd
	DescribeUsage            = "describe --application-id <application_id> [flags]"
	DescribeShortDescription = "Returns the Edge Application data"
	DescribeLongDescription  = "Displays information about the Edge Application via a given ID to show the application’s attributes in detail"
	DescribeFlagOut          = "Exports the output to the given <file_path/file_name.ext>"
	DescribeFlagFormat       = "Changes the output format passing the json value to the flag"
	DescribeHelpFlag         = "Displays more information about the describe command"

	//delete cmd
	DeleteUsage            = "delete --application-id <application_id> [flags]"
	DeleteShortDescription = "Removes an Edge Application"
	DeleteLongDescription  = "Removes an Edge Application from the Edge Applications library based on its given ID"
	DeleteOutputSuccess    = "Edge Application %d was successfully deleted\n"
	DeleteHelpFlag         = "Displays more information about the delete subcommand"
	DeleteCascadeFlag      = "Deletes all resources created through the command azioncli edge_applications publish"
	DeleteMissingFunction  = "Missing Edge Function ID in azion.json file. Skipping deletion"
	DeleteCascadeSuccess   = "Cascade delete carried out successfully"

	//update cmd
	UpdateUsage                       = "update --application-id <application_id> [flags]"
	UpdateShortDescription            = "Modifies an Edge Application"
	UpdateLongDescription             = "Modifies an Edge Application based on its ID to update its name, activity status, and other attributes"
	UpdateFlagName                    = "The Edge Application's name"
	UpdateFlagDeliveryProtocol        = "The Edge Application's Delivery Protocol"
	UpdateFlagHttpPort                = "The Edge Application's Http Port"
	UpdateFlagHttpsPort               = "The Edge Application's Https Port"
	UpdateFlagMinimumTlsVersion       = "The Edge Application's Minimum Tls Version"
	UpdateFlagApplicationAcceleration = "Whether the Edge Application has Application Acceleration active or not"
	UpdateFlagCaching                 = "Whether the Edge Application has Caching active or not"
	UpdateFlagDeviceDetection         = "Whether the Edge Application has Device Detection active or not"
	UpdateFlagEdgeFirewall            = "Whether the Edge Application has Edge Firewall active or not"
	UpdateFlagEdgeFunctions           = "Whether the Edge Application has Edge Functions active or not"
	UpdateFlagImageOptimization       = "Whether the Edge Application has Image Optimization active or not"
	UpdateFlagL2Caching               = "Whether the Edge Application has L2 Caching active or not"
	UpdateFlagLoadBalancer            = "Whether the Edge Application has Load Balancer active or not"
	UpdateRawLogs                     = "Whether the Edge Application has Raw Logs active or not"
	UpdateWebApplicationFirewall      = "Whether the Edge Application has Web Application Firewall active or not"
	UpdateFlagIn                      = "Given path and JSON file to automatically update the Edge Application attributes; you can use - for reading from stdin"
	UpdateOutputSuccess               = "Updated Edge Application with ID %d\n"
	UpdateHelpFlag                    = "Displays more information about the update subcommand"

	LsUsage            = "ls"
	LsShortDescription = "Displays presets accepted by Vulcan"
	LsLongDescription  = "Displays presets accepted by Vulcan"
	InstallingVulcan   = "Vulcan was not found in your machine. Please wait while vulcan is being installed"
)
