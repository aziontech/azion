package edge_applications

var (

	//used by more than one cmd
	EdgeApplicationFlagId      = "Unique identifier of the Edge Application"
	EdgeApplicationFileWritten = "File successfully written to: %s\n"

	//edge_applications cmd
	EdgeApplicationsUsage            = "edge_applications"
	EdgeApplicationsShortDescription = "Creates Web Applications on Azion's platform"
	EdgeApplicationsLongDescription  = "Build your Web applications in minutes without the need to manage infrastructure or security"
	EdgeApplicationsFlagHelp         = "Displays more information about the edge_application command"
	EdgeApplicationsAutoDetectec     = "Auto-detected Project Settings (%s)\n"

	//build cmd
	EdgeApplicationsBuildUsage            = "build [flags]"
	EdgeApplicationsBuildShortDescription = "Builds a Web Application"
	EdgeApplicationsBuildLongDescription  = "Builds your Web Application to run on Azion’s Edge Computing Platform"
	EdgeApplicationsBuildRunningCmd       = "Running build step command:\n\n"
	EdgeApplicationsBuildStart            = "Building your Web Application\n"
	EdgeApplicationsBuildSuccessful       = "Your Web Application was built successfully\n"
	EdgeApplicationsBuildFlagHelp         = "Displays more information about the build subcommand"

	//init cmd
	EdgeApplicationsInitUsage = `init [flags]
	--type string       The type of Edge application
	cdn                 Create an edge application to cache and deliver your content.
	javascript          Create a serverless Javascript application on edge.
	flareact            Create a serverless Flareact application on edge.
	nextjs              Create a serverless NextJS edge-runtime application on edge.`
	EdgeApplicationsInitShortDescription  = "Initializes a Web Application"
	EdgeApplicationsInitLongDescription   = "Defines primary parameters based on a given name and application type to start a Web Application on Azion’s platform"
	EdgeApplicationsInitRunningCmd        = "Running init step command:\n\n"
	EdgeApplicationsInitFlagName          = "The Web application's name"
	EdgeApplicationsInitFlagType          = "The type of  Web application <javascript | flareact | nextjs>"
	EdgeApplicationsInitFlagYes           = "Forces the automatic response 'yes' to all user input"
	EdgeApplicationsInitFlagNo            = "Forces the automatic response 'no' to all user input"
	WebAppInitContentOverridden           = "This project was already configured. Do you want to override the previous configuration? <yes | no> (default: no) "
	WebAppInitCmdSuccess                  = "Template successfully fetched and configured\n"
	EdgeApplicationsInitFlagHelp          = "Displays more information about the init subcommand"
	EdgeApplicationsInitSuccessful        = "Your project %s was initialized successfully"
	EdgeApplicationsInitNameNotSent       = "The Project Name was not sent through the --name flag; By default when --name is not informed the one found in your package.json file is used\n"
	EdgeApplicationsUpdateNamePackageJson = "Updating your package.json name field with the value informed through the --name flag"
	EdgeApplicationsInitTypeNotSent       = "The Project Type was not sent through the --type flag; By default when --type is not informed it is auto-detected based on the framework used by the user\n"

	//publish cmd
	EdgeApplicationsPublishUsage                       = "publish"
	EdgeApplicationsPublishShortDescription            = "Publishes a Web Application on the Azion platform"
	EdgeApplicationsPublishLongDescription             = "Publishes a Web Application based on the Azion’s Platform"
	EdgeApplicationsPublishRunningCmd                  = "Running pre publish command:\n\n"
	EdgeApplicationsPublishSuccessful                  = "Your Web Application was published successfully\n"
	EdgeApplicationsPublishOutputDomainSuccess         = "\nTo visualize your application access the domain: %s\n"
	EdgeApplicationsPublishOutputCachePurge            = "Domain cache was purged\n"
	EdgeApplicationsPublishOutputEdgeFunctionCreate    = "Created Edge Function %s with ID %d\n"
	EdgeApplicationsPublishOutputEdgeFunctionUpdate    = "Updated Edge Function %s with ID %d\n"
	EdgeApplicationsPublishOutputEdgeApplicationCreate = "Created Edge Application %s with ID %d\n"
	EdgeApplicationsPublishOutputEdgeApplicationUpdate = "Updated Edge Application %s with ID %d\n"
	EdgeApplicationsPublishOutputDomainCreate          = "Created Domain %s with ID %d\n"
	EdgeApplicationsPublishOutputDomainUpdate          = "Updated Domain %s with ID %d\n"
	EdgeApplicationsPublishFlagHelp                    = "Displays more information about the publish subcommand"
	EdgeApplicationsPublishPropagation                 = "Content is being propagated to all Azion POPs and it might take a few minutes for all edges to be up to date\n"

	EdgeApplicationsAWSMesaage = "Please inform your AWS credentials below:\n"
	EdgeApplicationsAWSSecret  = "AWS_SECRET_ACCESS_KEY: "
	EdgeApplicationsAWSAcess   = "AWS_ACCESS_KEY_ID: "

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
	EdgeApplicationDeleteShortDescription = "Removes an Edge Function"
	EdgeApplicationDeleteLongDescription  = "Removes an Edge Application from the Edge Applications library based on its given ID"
	EdgeApplicationDeleteOutputSuccess    = "Edge Application %s was successfully deleted\n"
	EdgeApplicationDeleteHelpFlag         = "Displays more information about the delete subcommand"
)
