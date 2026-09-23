// Package registry constructs the API clients that commands use.
//
// It owns one thing: which service serves a given resource. Before this package
// existed, every call site spelled out a configuration key --
// f.Config.GetString("api_v4_url") and friends -- roughly two hundred times, so
// moving a resource between API generations meant editing all of them. That is
// the cost a third generation would otherwise multiply.
//
// A command asks for the client it needs and never names a base URL. When v6
// arrives, resolution changes here, in one place, and may start switching on
// f.APIVersion; the call sites do not change at all.
//
// Note that the mapping is per resource, not per account generation: variables,
// personal tokens, origins and domains are served by the legacy API even for
// accounts on v4, which is why a single "current base URL" would be wrong.
package registry

import (
	"net/http"

	apiApplications "github.com/aziontech/azion-cli/pkg/api/applications"
	apiCacheSetting "github.com/aziontech/azion-cli/pkg/api/cache_setting"
	apiConnector "github.com/aziontech/azion-cli/pkg/api/connector"
	apiCrl "github.com/aziontech/azion-cli/pkg/api/crl"
	apiCsr "github.com/aziontech/azion-cli/pkg/api/csr"
	apiCustomPages "github.com/aziontech/azion-cli/pkg/api/custom_pages"
	apiDataStream "github.com/aziontech/azion-cli/pkg/api/data_stream"
	apiDeviceGroups "github.com/aziontech/azion-cli/pkg/api/device_groups"
	apiDigitalCertificate "github.com/aziontech/azion-cli/pkg/api/digital_certificate"
	apiDnsRecord "github.com/aziontech/azion-cli/pkg/api/dns_record"
	apiDnsZone "github.com/aziontech/azion-cli/pkg/api/dns_zone"
	apiDnssec "github.com/aziontech/azion-cli/pkg/api/dnssec"
	apiDomain "github.com/aziontech/azion-cli/pkg/api/domain"
	apiFirewall "github.com/aziontech/azion-cli/pkg/api/firewall"
	apiFirewallInstance "github.com/aziontech/azion-cli/pkg/api/firewall_instance"
	apiFirewallRules "github.com/aziontech/azion-cli/pkg/api/firewall_rules"
	apiFunction "github.com/aziontech/azion-cli/pkg/api/function"
	apiFunctionInstance "github.com/aziontech/azion-cli/pkg/api/function_instance"
	apiKv "github.com/aziontech/azion-cli/pkg/api/kv"
	apiNetworkList "github.com/aziontech/azion-cli/pkg/api/network_list"
	apiOrigin "github.com/aziontech/azion-cli/pkg/api/origin"
	apiPersonalToken "github.com/aziontech/azion-cli/pkg/api/personal_token"
	apiRealtimePurge "github.com/aziontech/azion-cli/pkg/api/realtime_purge"
	apiRulesEngine "github.com/aziontech/azion-cli/pkg/api/rules_engine"
	apiStorage "github.com/aziontech/azion-cli/pkg/api/storage"
	apiVariables "github.com/aziontech/azion-cli/pkg/api/variables"
	apiWaf "github.com/aziontech/azion-cli/pkg/api/waf"
	apiWafExceptions "github.com/aziontech/azion-cli/pkg/api/waf_exceptions"
	apiWorkloads "github.com/aziontech/azion-cli/pkg/api/workloads"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
)

// The configured base URLs. These keys are named here and nowhere else.
func legacyAPIURL(f *cmdutil.Factory) string { return f.Config.GetString("api_url") }
func apiV4URL(f *cmdutil.Factory) string     { return f.Config.GetString("api_v4_url") }
func storageURL(f *cmdutil.Factory) string   { return f.Config.GetString("storage_url") }

// credential is the token the invocation authenticates with.
func credential(f *cmdutil.Factory) string { return f.Config.GetString("token") }

func httpClient(f *cmdutil.Factory) *http.Client { return f.HttpClient }

func Applications(f *cmdutil.Factory) *apiApplications.Client {
	return apiApplications.NewClient(httpClient(f), apiV4URL(f), credential(f))
}

func CacheSettings(f *cmdutil.Factory) *apiCacheSetting.ClientV4 {
	return apiCacheSetting.NewClientV4(httpClient(f), apiV4URL(f), credential(f))
}

func Connector(f *cmdutil.Factory) *apiConnector.Client {
	return apiConnector.NewClient(httpClient(f), apiV4URL(f), credential(f))
}

func CRL(f *cmdutil.Factory) *apiCrl.Client {
	return apiCrl.NewClient(httpClient(f), apiV4URL(f), credential(f))
}

func CSR(f *cmdutil.Factory) *apiCsr.Client {
	return apiCsr.NewClient(httpClient(f), apiV4URL(f), credential(f))
}

func CustomPages(f *cmdutil.Factory) *apiCustomPages.Client {
	return apiCustomPages.NewClient(httpClient(f), apiV4URL(f), credential(f))
}

func DataStream(f *cmdutil.Factory) *apiDataStream.Client {
	return apiDataStream.NewClient(httpClient(f), apiV4URL(f), credential(f))
}

func DeviceGroups(f *cmdutil.Factory) *apiDeviceGroups.Client {
	return apiDeviceGroups.NewClient(httpClient(f), apiV4URL(f), credential(f))
}

func DigitalCertificate(f *cmdutil.Factory) *apiDigitalCertificate.Client {
	return apiDigitalCertificate.NewClient(httpClient(f), apiV4URL(f), credential(f))
}

func DNSRecord(f *cmdutil.Factory) *apiDnsRecord.Client {
	return apiDnsRecord.NewClient(httpClient(f), apiV4URL(f), credential(f))
}

func DNSZone(f *cmdutil.Factory) *apiDnsZone.Client {
	return apiDnsZone.NewClient(httpClient(f), apiV4URL(f), credential(f))
}

func DNSSEC(f *cmdutil.Factory) *apiDnssec.Client {
	return apiDnssec.NewClient(httpClient(f), apiV4URL(f), credential(f))
}

func Domain(f *cmdutil.Factory) *apiDomain.Client {
	return apiDomain.NewClient(httpClient(f), legacyAPIURL(f), credential(f))
}

func Firewall(f *cmdutil.Factory) *apiFirewall.Client {
	return apiFirewall.NewClient(httpClient(f), apiV4URL(f), credential(f))
}

func FirewallInstance(f *cmdutil.Factory) *apiFirewallInstance.Client {
	return apiFirewallInstance.NewClient(httpClient(f), apiV4URL(f), credential(f))
}

func FirewallRules(f *cmdutil.Factory) *apiFirewallRules.Client {
	return apiFirewallRules.NewClient(httpClient(f), apiV4URL(f), credential(f))
}

func Function(f *cmdutil.Factory) *apiFunction.Client {
	return apiFunction.NewClient(httpClient(f), apiV4URL(f), credential(f))
}

func FunctionInstance(f *cmdutil.Factory) *apiFunctionInstance.Client {
	return apiFunctionInstance.NewClient(httpClient(f), apiV4URL(f), credential(f))
}

func KV(f *cmdutil.Factory) *apiKv.Client {
	return apiKv.NewClient(httpClient(f), apiV4URL(f), credential(f))
}

func NetworkList(f *cmdutil.Factory) *apiNetworkList.Client {
	return apiNetworkList.NewClient(httpClient(f), apiV4URL(f), credential(f))
}

func Origin(f *cmdutil.Factory) *apiOrigin.Client {
	return apiOrigin.NewClient(httpClient(f), legacyAPIURL(f), credential(f))
}

// PersonalToken builds the client for the credential this invocation uses.
func PersonalToken(f *cmdutil.Factory) *apiPersonalToken.Client {
	return PersonalTokenFor(f, credential(f))
}

// PersonalTokenFor builds the client for a specific credential. The profile commands
// need it when acting on a profile other than the active one.
func PersonalTokenFor(f *cmdutil.Factory, token string) *apiPersonalToken.Client {
	return apiPersonalToken.NewClient(httpClient(f), legacyAPIURL(f), token)
}

func RealtimePurge(f *cmdutil.Factory) *apiRealtimePurge.Client {
	return apiRealtimePurge.NewClient(httpClient(f), apiV4URL(f), credential(f))
}

func RulesEngine(f *cmdutil.Factory) *apiRulesEngine.Client {
	return apiRulesEngine.NewClient(httpClient(f), apiV4URL(f), credential(f))
}

func Storage(f *cmdutil.Factory) *apiStorage.Client {
	return apiStorage.NewClient(httpClient(f), storageURL(f), credential(f))
}

func Variables(f *cmdutil.Factory) *apiVariables.Client {
	return apiVariables.NewClient(httpClient(f), legacyAPIURL(f), credential(f))
}

func WAF(f *cmdutil.Factory) *apiWaf.Client {
	return apiWaf.NewClient(httpClient(f), apiV4URL(f), credential(f))
}

func WAFExceptions(f *cmdutil.Factory) *apiWafExceptions.Client {
	return apiWafExceptions.NewClient(httpClient(f), apiV4URL(f), credential(f))
}

func Workloads(f *cmdutil.Factory) *apiWorkloads.Client {
	return apiWorkloads.NewClient(httpClient(f), apiV4URL(f), credential(f))
}
