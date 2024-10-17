package contracts

import (
	"os"

	sdk "github.com/aziontech/azionapi-go-sdk/edgeapplications"
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
	IsFirewall    bool
}

type DevInfo struct {
	IsFirewall string
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
	NotFirstRun   bool                         `json:"not-first-run"`
	Function      AzionJsonDataFunction        `json:"function"`
	Application   AzionJsonDataApplication     `json:"application"`
	Domain        AzionJsonDataDomain          `json:"domain"`
	RtPurge       AzionJsonDataPurge           `json:"rt-purge"`
	Origin        []AzionJsonDataOrigin        `json:"origin"`
	RulesEngine   AzionJsonDataRulesEngine     `json:"rules-engine"`
	CacheSettings []AzionJsonDataCacheSettings `json:"cache-settings"`
}

type Results struct {
	Result Result `json:"result"`
}

type Result struct {
	Azion  *AzionApplicationOptions `json:"azion,omitempty"`  // Pointer and omitempty tag
	Extras []interface{}            `json:"extras"`           // Assuming Extras can contain any data
	Errors *ErrorDetails            `json:"errors,omitempty"` // Pointer and omitempty for optional errors
}

type ErrorDetails struct {
	Error   int    `json:"error"`
	Message string `json:"message"`
	Stack   string `json:"stack"`
}

// type Logs struct {
// 	Status string `json:"status"`
// }

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

type AzionApplicationConfig struct {
	InitData    InitConf    `json:"init"`
	BuildData   BuildConf   `json:"build"`
	PublishData PublishConf `json:"publish"`
}

type InitConf struct {
	Cmd        string `json:"cmd"`
	Env        string `json:"env"`
	OutputCtrl string `json:"output-ctrl"`
	Default    string `json:"default"`
}

type BuildConf struct {
	Cmd        string `json:"cmd"`
	Env        string `json:"env"`
	OutputCtrl string `json:"output-ctrl"`
	Default    string `json:"default"`
}

type PublishConf struct {
	Cmd        string `json:"pre_cmd"`
	Env        string `json:"env"`
	OutputCtrl string `json:"output-ctrl"`
	Default    string `json:"default"`
}

type CacheConf struct {
	PurgeOnPublish bool `json:"purge_on_publish"`
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
	Domain        Domains        `json:"domain"`
	Purge         []Purges       `json:"purge"`
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
	Name                 string                              `json:"name"`
	OriginType           string                              `json:"origin_type,omitempty"`
	Bucket               string                              `json:"bucket,omitempty"`
	Prefix               string                              `json:"prefix,omitempty"`
	Addresses            []sdk.CreateOriginsRequestAddresses `json:"addresses,omitempty"`
	HostHeader           string                              `json:"host_header,omitempty"`
	OriginProtocolPolicy *string                             `json:"origin_protocol_policy,omitempty"`
	OriginPath           *string                             `json:"origin_path,omitempty"`
	HmacAuthentication   *bool                               `json:"hmac_authentication,omitempty"`
	HmacRegionName       *string                             `json:"hmac_region_name,omitempty"`
	HmacAccessKey        *string                             `json:"hmac_access_key,omitempty"`
	HmacSecretKey        *string                             `json:"hmac_secret_key,omitempty"`
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

type RuleIdsStruct struct {
	Id    int64
	Phase string
}
