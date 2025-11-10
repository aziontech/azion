package contracts

import (
	"os"

	sdk "github.com/aziontech/azionapi-go-sdk/edgeapplications"
	edgesdk "github.com/aziontech/azionapi-v4-go-sdk-dev/edge-api"
)

type FileOps struct {
	Path        string
	MimeType    string
	FileContent *os.File
	VersionID   string
}

type BuildInfo struct {
	Preset        string
	Entry         string
	NodePolyfills string
	OwnWorker     string
	ProjectPath   string
	SkipFramework bool
}

type BuildInfoV3 struct {
	Preset        string
	Entry         string
	NodePolyfills string
	OwnWorker     string
	ProjectPath   string
	IsFirewall    bool
}

type ListOptions struct {
	Details           bool
	OrderBy           string
	Sort              string
	Page              int64
	PageSize          int64
	NextPage          bool
	Filter            string
	ContinuationToken string
}

type DescribeOptions struct {
	OutPath string
	Format  string
}

type AzionApplicationOptions struct {
	Test          func(path string) error      `json:"-"`
	Name          string                       `json:"name"`
	Bucket        string                       `json:"bucket"`
	Preset        string                       `json:"preset"` // framework: react, next, vue, angular and etc
	Env           string                       `json:"env"`
	Prefix        string                       `json:"prefix"`
	SkipDeletion  *bool                        `json:"skip-deletion,omitempty"`
	NotFirstRun   bool                         `json:"not-first-run"`
	Function      []AzionJsonDataFunction      `json:"function"`
	Application   AzionJsonDataApplication     `json:"application"`
	Domain        AzionJsonDataDomain          `json:"domain"`
	RtPurge       AzionJsonDataPurge           `json:"rt-purge"`
	Origin        []AzionJsonDataOrigin        `json:"origin"`
	RulesEngine   AzionJsonDataRulesEngine     `json:"rules-engine"`
	CacheSettings []AzionJsonDataCacheSettings `json:"cache-settings"`
	Workloads     AzionJsonDataWorkload        `json:"workloads"`
	Connectors    []AzionJsonDataConnectors    `json:"connectors"`
}

type AzionApplicationOptionsV3 struct {
	Test          func(path string) error      `json:"-"`
	Name          string                       `json:"name"`
	Bucket        string                       `json:"bucket"`
	Preset        string                       `json:"preset"` // framework: react, next, vue, angular and etc
	Env           string                       `json:"env"`
	Prefix        string                       `json:"prefix"`
	SkipDeletion  *bool                        `json:"skip-deletion,omitempty"`
	NotFirstRun   bool                         `json:"not-first-run"`
	Function      AzionJsonDataFunction        `json:"function"`
	Application   AzionJsonDataApplication     `json:"application"`
	Domain        AzionJsonDataDomain          `json:"domain"`
	RtPurge       AzionJsonDataPurge           `json:"rt-purge"`
	Origin        []AzionJsonDataOrigin        `json:"origin"`
	RulesEngine   AzionJsonDataRulesEngine     `json:"rules-engine"`
	CacheSettings []AzionJsonDataCacheSettings `json:"cache-settings"`
}

type AzionJsonDataConnectors struct {
	Id      int64             `json:"id"`
	Name    string            `json:"name"`
	Address []edgesdk.Address `json:"address,omitempty"`
}

type Results struct {
	Result Result `json:"result"`
}

type ResultsV4 struct {
	Result ResultV4 `json:"result"`
}

type Result struct {
	Azion  *AzionApplicationOptionsV3 `json:"azion,omitempty"`
	Extras []interface{}              `json:"extras"`
	Errors *ErrorDetails              `json:"errors,omitempty"`
}

type ResultV4 struct {
	Azion       *AzionApplicationOptions `json:"azion,omitempty"`
	Extras      []interface{}            `json:"extras"`
	Errors      *ErrorDetails            `json:"errors,omitempty"`
	AzionConfig *AzionConfig             `json:"azion_config,omitempty"`
}

type AzionConfig struct {
	FileName    string `json:"file_name"`
	FileContent string `json:"file_content"`
}

type ErrorDetails struct {
	Error   int    `json:"error"`
	Message string `json:"message"`
	Stack   string `json:"stack"`
}

// LogEntry represents each log entry with content and timestamp.
type LogEntry struct {
	Content   string `json:"content,omitempty"`
	Timestamp string `json:"timestamp,omitempty"`
}

// StatusResponse represents the overall status and logs of the deployment process.
type Logs struct {
	Status string     `json:"status"`
	Logs   []LogEntry `json:"logs"`
}

type AzionApplicationSimple struct {
	Name        string                   `json:"name"`
	Type        string                   `json:"type"`
	Domain      AzionJsonDataDomain      `json:"domain"`
	Application AzionJsonDataApplication `json:"application"`
}

type AzionJsonDataFunction struct {
	ID           int64  `json:"id"`
	Name         string `json:"name"`
	File         string `json:"file"`
	Args         string `json:"args"`
	InstanceID   int64  `json:"instance-id"`
	InstanceName string `json:"instance-name"`
	CacheId      int64  `json:"cache-id"`
}

type AzionJsonDataApplication struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type AzionJsonDataOrigin struct {
	OriginId  int64    `json:"origin-id"`
	OriginKey string   `json:"origin-key"`
	Name      string   `json:"name"`
	Address   []string `json:"address,omitempty"`
}

type AzionJsonDataDomain struct {
	Id         int64  `json:"id"`
	Name       string `json:"name"`
	DomainName string `json:"domain_name"`
	Url        string `json:"url"`
}

type AzionJsonDataWorkload struct {
	Id          int64         `json:"id"`
	Name        string        `json:"name"`
	Domains     []string      `json:"domains"`
	Url         string        `json:"url"`
	Deployments []Deployments `json:"deployment_id"`
}

type Deployments struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

type AzionJsonDataPurge struct {
	PurgeOnPublish bool `json:"purge_on_publish"`
}

type AzionJsonDataRulesEngine struct {
	Created bool                 `json:"created"`
	Rules   []AzionJsonDataRules `json:"rules"`
}

type AzionJsonDataRules struct {
	Id    int64  `json:"id"`
	Name  string `json:"name"`
	Phase string `json:"phase"`
}

type AzionJsonDataCacheSettings struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

type Manifest struct {
	CacheSettings []CacheSetting `json:"cache"`
	Origins       []Origin       `json:"origin"`
	Rules         []RuleEngine   `json:"rules"`
	Domain        *Domains       `json:"domain,omitempty"`
	Purge         []Purges       `json:"purge"`
}

// BuildManifest represents the build configuration in the manifest.json file
type BuildManifest struct {
	Build Build `json:"build,omitempty"`
}

type Build struct {
	Preset    string   `json:"preset,omitempty"`    // JavaScript, etc.
	Entry     []string `json:"entry,omitempty"`     // Entry files like main.js
	Polyfills bool     `json:"polyfills,omitempty"` // Whether to include polyfills
}

type ManifestV4 struct {
	Build               BuildManifest                         `json:"build"`
	Storage             []StorageManifest                     `json:"storage"`
	Functions           []Function                            `json:"functions"`
	Applications        []Applications                        `json:"applications"`
	Connectors          []edgesdk.ConnectorPolymorphicRequest `json:"connectors"`
	Workloads           []WorkloadManifest                    `json:"workloads"`
	WorkloadDeployments []WorkloadDeployment                  `json:"workload_deployments,omitempty"`
	Purge               []PurgeManifest                       `json:"purge"`
}

type PurgeManifest struct {
	Items []string `json:"items"`
	Layer *string  `json:"layer,omitempty"`
	Type  string   `json:"type"`
}

// WorkloadManifest represents a workload in the manifest.json file
type WorkloadManifest struct {
	Name                      string                      `json:"name"`
	Active                    *bool                       `json:"active,omitempty"`
	Infrastructure            int64                       `json:"infrastructure,omitempty"`
	WorkloadDomainAllowAccess *bool                       `json:"workload_domain_allow_access,omitempty"`
	Domains                   []string                    `json:"domains,omitempty"`
	Tls                       *edgesdk.TLSWorkloadRequest `json:"tls,omitempty"`
	Protocols                 *edgesdk.ProtocolsRequest   `json:"protocols,omitempty"`
	Mtls                      *edgesdk.MTLSRequest        `json:"mtls,omitempty"`
	NetworkMap                *string                     `json:"network_map,omitempty"`
}

type WorkloadDeployment struct {
	Name     string           `json:"name"`
	Current  bool             `json:"current"`
	Active   bool             `json:"active"`
	Strategy WorkloadStrategy `json:"strategy"`
}

type WorkloadStrategy struct {
	Type       string                `json:"type"`
	Attributes WorkloadStrategyAttrs `json:"attributes"`
}

type WorkloadStrategyAttrs struct {
	Application *string `json:"application"`
	Firewall    *string `json:"firewall"`
	CustomPage  *string `json:"custom_page"`
}

type Purges struct {
	Type   string   `json:"type"`
	Urls   []string `json:"urls"`
	Method string   `json:"method"`
}

type Domains struct {
	Name                       string   `json:"name,omitempty"`
	Cnames                     []string `json:"cnames,omitempty"`
	CnameAccessOnly            *bool    `json:"cname_access_only,omitempty"`
	IsActive                   *bool    `json:"is_active,omitempty"`
	EdgeApplicationId          int64    `json:"edge_application_id,omitempty"`
	DigitalCertificateId       *string  `json:"digital_certificate_id,omitempty"`
	Environment                *string  `json:"environment,omitempty"`
	IsMtlsEnabled              *bool    `json:"is_mtls_enabled,omitempty"`
	MtlsTrustedCaCertificateId int64    `json:"mtls_trusted_ca_certificate_id,omitempty"`
	EdgeFirewallId             int64    `json:"edge_firewall_id,omitempty"`
	MtlsVerification           *string  `json:"mtls_verification,omitempty"`
	CrlList                    []int64  `json:"crl_list,omitempty"`
}

type CacheSetting struct {
	Name                           *string  `json:"name,omitempty"`
	BrowserCacheSettings           *string  `json:"browser_cache_settings,omitempty"`
	BrowserCacheSettingsMaximumTtl *int64   `json:"browser_cache_settings_maximum_ttl,omitempty"`
	CdnCacheSettings               *string  `json:"cdn_cache_settings,omitempty"`
	CdnCacheSettingsMaximumTtl     *int64   `json:"cdn_cache_settings_maximum_ttl,omitempty"`
	CacheByQueryString             *string  `json:"cache_by_query_string,omitempty"`
	QueryStringFields              []string `json:"query_string_fields,omitempty"`
	EnableQueryStringSort          *bool    `json:"enable_query_string_sort,omitempty"`
	CacheByCookies                 *string  `json:"cache_by_cookies,omitempty"`
	CookieNames                    []string `json:"cookie_names,omitempty"`
	AdaptiveDeliveryAction         *string  `json:"adaptive_delivery_action,omitempty"`
	DeviceGroup                    []int32  `json:"device_group,omitempty"`
	EnableCachingForPost           *bool    `json:"enable_caching_for_post,omitempty"`
	L2CachingEnabled               *bool    `json:"l2_caching_enabled,omitempty"`
	IsSliceConfigurationEnabled    *bool    `json:"is_slice_configuration_enabled,omitempty"`
	IsSliceEdgeCachingEnabled      *bool    `json:"is_slice_edge_caching_enabled,omitempty"`
	IsSliceL2CachingEnabled        *bool    `json:"is_slice_l2_caching_enabled,omitempty"`
	SliceConfigurationRange        *int64   `json:"slice_configuration_range,omitempty"`
	EnableCachingForOptions        *bool    `json:"enable_caching_for_options,omitempty"`
	EnableStaleCache               *bool    `json:"enable_stale_cache,omitempty"`
	L2Region                       *string  `json:"l2_region,omitempty"`
}

type Origin struct {
	Name                 string    `json:"name"`
	OriginType           string    `json:"origin_type,omitempty"`
	Bucket               string    `json:"bucket,omitempty"`
	Prefix               string    `json:"prefix,omitempty"`
	Addresses            []Address `json:"addresses,omitempty"`
	HostHeader           string    `json:"host_header,omitempty"`
	OriginProtocolPolicy *string   `json:"origin_protocol_policy,omitempty"`
	OriginPath           *string   `json:"origin_path,omitempty"`
	HmacAuthentication   *bool     `json:"hmac_authentication,omitempty"`
	HmacRegionName       *string   `json:"hmac_region_name,omitempty"`
	HmacAccessKey        *string   `json:"hmac_access_key,omitempty"`
	HmacSecretKey        *string   `json:"hmac_secret_key,omitempty"`
}

type Address struct {
	Address string `json:"address"`
	Weight  int    `json:"weight,omitempty"`
}

type RuleEngine struct {
	Name        string                         `json:"name"`
	Description *string                        `json:"description,omitempty"`
	Phase       string                         `json:"phase,omitempty"`
	Order       int64                          `json:"order,omitempty"`
	IsActive    bool                           `json:"is_active,omitempty"`
	Criteria    [][]sdk.RulesEngineCriteria    `json:"criteria,omitempty"`
	Behaviors   []sdk.RulesEngineBehaviorEntry `json:"behaviors,omitempty"`
}

type SyncOpts struct {
	RuleIds   map[string]RuleIdsStruct
	CacheIds  map[string]AzionJsonDataCacheSettings
	OriginIds map[string]AzionJsonDataOrigin
	Conf      *AzionApplicationOptions
}

type SyncOptsV3 struct {
	RuleIds   map[string]RuleIdsStruct
	CacheIds  map[string]AzionJsonDataCacheSettings
	OriginIds map[string]AzionJsonDataOrigin
	Conf      *AzionApplicationOptionsV3
}

type RuleIdsStruct struct {
	Id    int64
	Phase string
}

type ManifestCacheSetting struct {
	Name         string                               `json:"name"`
	BrowserCache *edgesdk.BrowserCacheModuleRequest   `json:"browser_cache,omitempty"`
	Modules      *edgesdk.CacheSettingsModulesRequest `json:"modules,omitempty"`
}

type CacheSettingManifest struct {
	Name         string                `json:"name"`
	BrowserCache *BrowserCacheSettings `json:"browser_cache,omitempty"`
	Modules      *CacheSettingModules  `json:"modules,omitempty"`
}

type BrowserCacheSettings struct {
	Behavior string `json:"behavior,omitempty"`
	MaxAge   int64  `json:"max_age,omitempty"`
}

type CacheSettingModules struct {
	EdgeCache              *EdgeCacheSettings              `json:"edge_cache,omitempty"`
	TieredCache            *TieredCacheSettings            `json:"tiered_cache,omitempty"`
	ApplicationAccelerator *ApplicationAcceleratorSettings `json:"application_accelerator,omitempty"`
}

type EdgeCacheSettings struct {
	Behavior       string                  `json:"behavior,omitempty"`
	MaxAge         int64                   `json:"max_age,omitempty"`
	StaleCache     *StaleCacheSettings     `json:"stale_cache,omitempty"`
	LargeFileCache *LargeFileCacheSettings `json:"large_file_cache,omitempty"`
}

type StaleCacheSettings struct {
	Enabled bool `json:"enabled"`
}

type LargeFileCacheSettings struct {
	Enabled bool  `json:"enabled"`
	Offset  int64 `json:"offset,omitempty"`
}

type TieredCacheSettings struct {
	Topology string `json:"topology,omitempty"`
}

type ApplicationAcceleratorSettings struct {
	CacheVaryByMethod      []string                  `json:"cache_vary_by_method,omitempty"`
	CacheVaryByQuerystring *QuerystringCacheSettings `json:"cache_vary_by_querystring,omitempty"`
	CacheVaryByCookies     *CookiesCacheSettings     `json:"cache_vary_by_cookies,omitempty"`
	CacheVaryByDevices     *DevicesCacheSettings     `json:"cache_vary_by_devices,omitempty"`
}

type QuerystringCacheSettings struct {
	Behavior    string   `json:"behavior,omitempty"`
	Fields      []string `json:"fields,omitempty"`
	SortEnabled bool     `json:"sort_enabled"`
}

type CookiesCacheSettings struct {
	Behavior    string   `json:"behavior,omitempty"`
	CookieNames []string `json:"cookie_names,omitempty"`
}

type DevicesCacheSettings struct {
	Behavior    string   `json:"behavior,omitempty"`
	DeviceGroup []string `json:"device_group,omitempty"`
}

// FunctionInstance represents an edge function instance in the manifest.json file
type FunctionInstance struct {
	Name     string                 `json:"name"`
	Function string                 `json:"function"`
	Active   bool                   `json:"active"`
	Args     map[string]interface{} `json:"args,omitempty"`
}

// StorageManifest represents an edge storage entry in the manifest.json file
type StorageManifest struct {
	Name       string `json:"name"`
	EdgeAccess string `json:"edge_access"` // read_write, read_only, etc.
	Dir        string `json:"dir"`         // Directory path
	Prefix     string `json:"prefix"`
}

type Applications struct {
	Name               string                             `json:"name"`
	Modules            *edgesdk.ApplicationModulesRequest `json:"modules,omitempty"`
	Active             *bool                              `json:"active,omitempty"`
	Debug              *bool                              `json:"debug,omitempty"`
	Rules              []ManifestRulesEngine              `json:"rules"`
	CacheSettings      []ManifestCacheSetting             `json:"cache_settings"`
	FunctionsInstances []FunctionInstance                 `json:"functions_instances,omitempty"`
}

type StorageBinding struct {
	Bucket string `json:"bucket,omitempty"`
	Prefix string `json:"prefix,omitempty"`
}

type FunctionBindings struct {
	Storage StorageBinding `json:"storage,omitempty"`
}

type Function struct {
	Name                 string                 `json:"name"`
	Path                 string                 `json:"path"`
	Runtime              string                 `json:"runtime,omitempty"`               // azion_js, etc.
	DefaultArgs          map[string]interface{} `json:"default_args,omitempty"`          // Default arguments
	ExecutionEnvironment string                 `json:"execution_environment,omitempty"` // application, etc.
	Active               bool                   `json:"active,omitempty"`                // Whether the function is active
	Bindings             FunctionBindings       `json:"bindings,omitempty"`              // Function bindings
	// Keep the old fields for backward compatibility
	Argument string `json:"argument,omitempty"`
}

type Modules struct {
	EdgeCacheEnabled              *bool `json:"edge_cache_enabled,omitempty"`
	EdgeFunctionsEnabled          *bool `json:"edge_functions_enabled,omitempty"`
	ApplicationAcceleratorEnabled *bool `json:"application_accelerator_enabled,omitempty"`
	ImageProcessorEnabled         *bool `json:"image_processor_enabled,omitempty"`
	TieredCacheEnabled            *bool `json:"tiered_cache_enabled,omitempty"`
}

// ManifestRulesEngine represents a rule engine entry in the manifest.json file
type ManifestRulesEngine struct {
	Phase string       `json:"phase"`
	Rule  ManifestRule `json:"rule"`
}

// ManifestRule represents a rule in the manifest.json file
type ManifestRule struct {
	Name        string                                           `json:"name"`
	Description string                                           `json:"description,omitempty"`
	Active      bool                                             `json:"active,omitempty"`
	Criteria    [][]edgesdk.EdgeApplicationCriterionFieldRequest `json:"criteria"`
	Behaviors   []ManifestRuleBehavior                           `json:"behaviors"`
}

// ManifestRuleBehavior represents a behavior in a rule
type ManifestRuleBehavior struct {
	Type       string                 `json:"type,omitempty"`
	Attributes map[string]interface{} `json:"attributes,omitempty"`
}
