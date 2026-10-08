package dryrun

// Used by both the v3 and the v4 command trees.
var (
	SkipManifest         = "This project has not been built yet. Skipping the simulation for resources found in your azion.config file\n"
	CreateEdgeApp        = "Creating Application named '%s'\n"
	UpdateEdgeApp        = "Updating Application with ID '%d', named '%s'\n"
	CreateBucket         = "Creating Bucket named '%s'\n"
	CreateDomain         = "Creating Domain named '%s'\n"
	UpdateDomain         = "Updating Domain with ID '%d', named '%s'\n"
	CreateOriginSingle   = "Creating single Origin named '%s'\n"
	UpdateDefaultRule    = "Updating default Rule Engine - Set Origin '%s'\n"
	DeletingRuleEngine   = "Deleting Rule Engine with ID '%d', named '%s'"
	DeletingOrigin       = "Deleting Origin with ID '%d' and Key '%s', named '%s'"
	DeletingCacheSetting = "Deleting Cache Setting with ID '%d', named '%s'"
	CreateRulesCache     = "Presenting the option to create Cache Setting (details below) and Rule Engine setting said Cache Setting\n"
	// Steps of the deploy pipelines that had no message before, because the
	// simulation they replace did not describe them.
	Header            = "The following would be done, without changing anything:\n"
	Build             = "Building your application\n"
	BundlerInit       = "Initializing the bundler and building your application\n"
	UploadStaticFiles = "Uploading static files to Bucket '%s'\n"
	CreateFunction    = "Creating Function named '%s'\n"
	UpdateFunction    = "Updating Function with ID '%d', named '%s'\n"
	CreateWorkload    = "Creating Workload named '%s'\n"
	// v4 only: its rules deploy creates the preset's default rules and does not
	// offer a cache setting. The prompt in CreateRulesCache is a v3 behaviour.
	CreateRulesEngine = "Creating the default Rules Engine for the '%s' preset\n"
	UpdateWorkload    = "Updating Workload with ID '%d', named '%s'\n"

	// Steps of the manifest pipeline, which the simulation never reached.
	ManifestHeader = "Applying the resources declared in your manifest:\n"
	ManifestStep   = "  - %s\n"
	// %s is the resource, %d how many the manifest declares.
	ManifestStepCount = "  - %s (%d)\n"
	// A step whose section the manifest does not declare still runs on v3, and
	// empties that section of azion.json.
	ManifestStepCleared = "  - %s: none declared, this section will be cleared\n"

	AskCreateCacheSettings = `Cache Settings specifications:
  - Browser Cache Settings: Override Cache Settings
  - Maximum TTL for Browser Cache Settings (in seconds): 7200
  - Application Cache Settings: Override Cache Settings
  - Maximum TTL for Application Cache Settings (in seconds): 7200`
)
