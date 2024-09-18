package dryrun

var (
	SkipManifest           = "This project has not been built yet. Skipping the simulation for resources found in your azion.config file\n"
	CreateEdgeApp          = "Creating Edge Application named '%s'\n"
	UpdateEdgeApp          = "Updating Edge Application with ID '%d', named '%s'\n"
	CreateOrigin           = "Creating Origin named '%s'\n"
	UpdateOrigin           = "Updating Origin with ID '%d' and Key '%s', named '%s'\n"
	CreateCacheSetting     = "Creating Cache Setting named '%s'\n"
	UpdateCacheSetting     = "Updating Cache Setting with ID '%d', named '%s'\n"
	CreateRule             = "Creating Rule Engine named '%s'\n"
	UpdateRule             = "Updating Rule Engine with ID '%d', named '%s'\n"
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
  - Edge Application Cache Settings: Override Cache Settings
  - Maximum TTL for Edge Application Cache Settings (in seconds): 7200`
)
