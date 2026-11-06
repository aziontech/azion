package contracts

// ============================================================================
// v6 — owned by API v6. The v4 types above are not reused: in v6 an id alone
// does not say what is live, so every resource also records its version.
// ============================================================================

import (
	"encoding/json"

	sdkv6 "github.com/aziontech/azionapi-v6-go-sdk/azion-api"
)

// ManifestVersionV6 is the "version" of a manifest the bundler writes with
// --config-version 6.
const ManifestVersionV6 = 6

// AzionApplicationOptionsV6 is azion.json for a v6 project. It is written after
// every step, so a run that stops halfway leaves enough to continue from.
type AzionApplicationOptionsV6 struct {
	Test         func(path string) error `json:"-"`
	Name         string                  `json:"name"`
	Bucket       string                  `json:"bucket"`
	Preset       string                  `json:"preset"`
	Env          string                  `json:"env"`
	Prefix       string                  `json:"prefix"`
	RotatePrefix *bool                   `json:"rotate-prefix,omitempty"`
	NotFirstRun  bool                    `json:"not-first-run,omitempty"`

	// Environments and Deployments are the ones the manifest declares, by
	// name, with the ids this project resolved them to.
	Environments []AzionJsonV6Environment `json:"environments,omitempty"`
	Deployments  []AzionJsonV6Deployment  `json:"deployments,omitempty"`
	Workloads    []AzionJsonV6Workload    `json:"workloads,omitempty"`

	Applications []AzionJsonV6Resource `json:"applications,omitempty"`
	Functions    []AzionJsonV6Resource `json:"functions,omitempty"`
	Connectors   []AzionJsonV6Resource `json:"connectors,omitempty"`
	Firewalls    []AzionJsonV6Resource `json:"firewalls,omitempty"`
	CustomPages  []AzionJsonV6Resource `json:"custom-pages,omitempty"`
	NetworkLists []AzionJsonV6Resource `json:"network-lists,omitempty"`
	WAFs         []AzionJsonV6Resource `json:"wafs,omitempty"`
}

// AzionJsonV6Resource is one versioned resource the project owns: which one,
// and the version this project last wrote or built.
type AzionJsonV6Resource struct {
	ID           int64  `json:"id"`
	Name         string `json:"name"`
	VersionID    string `json:"version-id,omitempty"`
	VersionState string `json:"version-state,omitempty"`
}

type AzionJsonV6Environment struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Policy string `json:"policy,omitempty"`
}

// AzionJsonV6Deployment is one Deployment of the project and the last release
// the CLI made on it.
type AzionJsonV6Deployment struct {
	ID      string             `json:"id"`
	Name    string             `json:"name"`
	Policy  string             `json:"policy,omitempty"`
	Release AzionJsonV6Release `json:"release"`
}

type AzionJsonV6Release struct {
	ID        string `json:"id,omitempty"`
	State     string `json:"state,omitempty"`
	Activated bool   `json:"activated,omitempty"`
}

type AzionJsonV6Workload struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Url  string `json:"url,omitempty"`
}

// ManifestV6 is .edge/manifest.json for a v6 project, as the bundler writes it
// with --config-version 6 (@aziontech/config, branch feat/config-api-v6).
//
// Resource bodies follow the v6 API requests, except that wherever the API
// takes the id of another resource the project declares, the manifest gives its
// name and the CLI resolves it once the target exists:
//
//   - functions_instances[].function: a Function;
//   - rule behaviors: run_function (a Function, run through the instance of it
//     in the same Application or Firewall), set_cache_policy (a cache setting),
//     set_connector (a Connector), set_waf.waf_id (a WAF);
//   - firewall criteria is_in_list / is_not_in_list: a Network List;
//   - custom_pages[].pages[].page.attributes.connector: a Connector;
//   - workloads[].bindings[]: environment_name and deployment_name.
//
// Rules and custom pages are kept as raw JSON for that reason: they are only
// decoded into SDK types after the names are replaced.
//
// The link between a Deployment and the Application, Firewall and Custom Page
// it serves is not in the manifest: the CLI makes it with a release. A release
// serves at most one of each, so the manifest may declare at most one of each.
type ManifestV6 struct {
	Version int               `json:"version"`
	Build   Build             `json:"build"`
	Storage []StorageManifest `json:"storage,omitempty"`

	Functions    []ManifestV6Function             `json:"functions,omitempty"`
	Connectors   []sdkv6.ConnectorRequest         `json:"connectors,omitempty"`
	NetworkLists []sdkv6.NetworkListRequest       `json:"network_list,omitempty"`
	WAFs         []ManifestV6WAF                  `json:"waf,omitempty"`
	Applications []ManifestV6Application          `json:"applications,omitempty"`
	Firewalls    []ManifestV6Firewall             `json:"firewall,omitempty"`
	CustomPages  []json.RawMessage                `json:"custom_pages,omitempty"`
	Environments []sdkv6.EnvironmentCreateRequest `json:"environments,omitempty"`
	Deployments  []sdkv6.DeploymentCreateRequest  `json:"deployments,omitempty"`
	Workloads    []ManifestV6Workload             `json:"workloads,omitempty"`
	Purge        []PurgeManifest                  `json:"purge,omitempty"`
}

// ManifestV6Function is a Function whose code is a file of the build output.
type ManifestV6Function struct {
	Name                 string                 `json:"name"`
	Path                 string                 `json:"path"`
	Runtime              string                 `json:"runtime,omitempty"`
	ExecutionEnvironment string                 `json:"execution_environment,omitempty"`
	DefaultArgs          map[string]interface{} `json:"default_args,omitempty"`
	Active               *bool                  `json:"active,omitempty"`
}

type ManifestV6WAF struct {
	Name           string                               `json:"name"`
	Active         *bool                                `json:"active,omitempty"`
	EngineSettings *sdkv6.WAFEngineSettingsFieldRequest `json:"engine_settings,omitempty"`
	Exceptions     []sdkv6.WAFRuleRequest               `json:"exceptions,omitempty"`
}

type ManifestV6Application struct {
	Name              string                           `json:"name"`
	Modules           *sdkv6.ApplicationModulesRequest `json:"modules,omitempty"`
	Active            *bool                            `json:"active,omitempty"`
	Debug             *bool                            `json:"debug,omitempty"`
	CacheSettings     []sdkv6.CacheSettingRequest      `json:"cache_settings,omitempty"`
	DeviceGroups      []sdkv6.DeviceGroupRequest       `json:"device_groups,omitempty"`
	FunctionInstances []ManifestV6FunctionInstance     `json:"functions_instances,omitempty"`
	// Rules are request and response rules in one list, in order; Phase
	// says which.
	Rules []ManifestV6Rule `json:"rules,omitempty"`
}

// Rule phases of an Application.
const (
	PhaseRequest  = "request"
	PhaseResponse = "response"
)

type ManifestV6Rule struct {
	Phase string          `json:"phase"`
	Rule  json.RawMessage `json:"rule"`
}

type ManifestV6Firewall struct {
	Name              string                        `json:"name"`
	Modules           *sdkv6.FirewallModulesRequest `json:"modules,omitempty"`
	Debug             *bool                         `json:"debug,omitempty"`
	Active            *bool                         `json:"active,omitempty"`
	FunctionInstances []ManifestV6FunctionInstance  `json:"functions_instances,omitempty"`
	Rules             []json.RawMessage             `json:"rules_engine,omitempty"`
}

// ManifestV6FunctionInstance names the Function it runs; FunctionReference
// accepts either the Function's name or its id.
type ManifestV6FunctionInstance struct {
	Name     string                 `json:"name"`
	Function FunctionReference      `json:"function"`
	Args     map[string]interface{} `json:"args,omitempty"`
	Active   *bool                  `json:"active,omitempty"`
}

type ManifestV6Workload struct {
	Name                      string                    `json:"name"`
	Active                    *bool                     `json:"active,omitempty"`
	Infrastructure            *int32                    `json:"infrastructure,omitempty"`
	WorkloadDomainAllowAccess *bool                     `json:"workload_domain_allow_access,omitempty"`
	Tls                       *sdkv6.TLSWorkloadRequest `json:"tls,omitempty"`
	Protocols                 *sdkv6.ProtocolsRequest   `json:"protocols,omitempty"`
	Mtls                      *sdkv6.MTLSRequest        `json:"mtls,omitempty"`
	Bindings                  []ManifestV6Binding       `json:"bindings,omitempty"`
}

// ManifestV6Binding serves a Deployment in an Environment, on some domains.
// Both are named; the CLI resolves them to the ids it created.
type ManifestV6Binding struct {
	EnvironmentName string             `json:"environment_name"`
	DeploymentName  string             `json:"deployment_name"`
	Domains         []ManifestV6Domain `json:"domains,omitempty"`
}

// ManifestV6Domain is a domain a binding answers on, with the id of the
// certificate to serve or none.
type ManifestV6Domain struct {
	Domain      string `json:"domain"`
	Certificate *int64 `json:"certificate,omitempty"`
}
