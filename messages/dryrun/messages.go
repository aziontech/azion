package dryrun

// Used by both the v3 and the v4 command trees.
var (
	SkipManifest           = "This project has not been built yet. Skipping the simulation for resources found in your azion.config file\n"
	CreateEdgeApp          = "Creating Application named '%s'\n"
	UpdateEdgeApp          = "Updating Application with ID '%d', named '%s'\n"
	CreateBucket           = "Creating Bucket named '%s'\n"
	CreateDomain           = "Creating Domain named '%s'\n"
	UpdateDomain           = "Updating Domain with ID '%d', named '%s'\n"
	CreateOriginSingle     = "Creating single Origin named '%s'\n"
	UpdateDefaultRule      = "Updating default Rule Engine - Set Origin '%s'\n"
	DeletingRuleEngine     = "Deleting Rule Engine with ID '%d', named '%s'"
	DeletingOrigin         = "Deleting Origin with ID '%d' and Key '%s', named '%s'"
	DeletingCacheSetting   = "Deleting Cache Setting with ID '%d', named '%s'"
	CreateRulesCache       = "Presenting the option to create Cache Setting (details below) and Rule Engine setting said Cache Setting\n"
	AskCreateCacheSettings = `Cache Settings specifications:
  - Browser Cache Settings: Override Cache Settings
  - Maximum TTL for Browser Cache Settings (in seconds): 7200
  - Application Cache Settings: Override Cache Settings
  - Maximum TTL for Application Cache Settings (in seconds): 7200`
)
