package registry

import (
	"net/http"
	"testing"

	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// recordingConfig satisfies config.Config and remembers every key looked up, so
// a test can assert which base URL a client was built from.
type recordingConfig struct{ keys []string }

func (r *recordingConfig) GetString(key string) string {
	r.keys = append(r.keys, key)
	return "https://example.test/" + key
}

func (r *recordingConfig) urlKey(t *testing.T) string {
	t.Helper()
	var urls []string
	for _, k := range r.keys {
		if k != "token" {
			urls = append(urls, k)
		}
	}
	require.Len(t, urls, 1, "a client must be built from exactly one base URL, got %v", r.keys)
	return urls[0]
}

// Every resource is pinned to the service that serves it. These are the
// mappings the call sites used to spell out by hand; changing one here
// repoints every command that asks for that client, which is the point of the
// package and also the reason it is worth asserting.
func TestClientsUseTheExpectedBaseURL(t *testing.T) {
	const (
		legacy  = "api_url"
		v4      = "api_v4_url"
		storage = "storage_url"
	)

	cases := []struct {
		name  string
		build func(*cmdutil.Factory)
		want  string
	}{
		{"Applications", func(f *cmdutil.Factory) { Applications(f) }, v4},
		{"CacheSettings", func(f *cmdutil.Factory) { CacheSettings(f) }, v4},
		{"Connector", func(f *cmdutil.Factory) { Connector(f) }, v4},
		{"CRL", func(f *cmdutil.Factory) { CRL(f) }, v4},
		{"CSR", func(f *cmdutil.Factory) { CSR(f) }, v4},
		{"CustomPages", func(f *cmdutil.Factory) { CustomPages(f) }, v4},
		{"DataStream", func(f *cmdutil.Factory) { DataStream(f) }, v4},
		{"DeviceGroups", func(f *cmdutil.Factory) { DeviceGroups(f) }, v4},
		{"DigitalCertificate", func(f *cmdutil.Factory) { DigitalCertificate(f) }, v4},
		{"DNSRecord", func(f *cmdutil.Factory) { DNSRecord(f) }, v4},
		{"DNSZone", func(f *cmdutil.Factory) { DNSZone(f) }, v4},
		{"DNSSEC", func(f *cmdutil.Factory) { DNSSEC(f) }, v4},
		{"Firewall", func(f *cmdutil.Factory) { Firewall(f) }, v4},
		{"FirewallInstance", func(f *cmdutil.Factory) { FirewallInstance(f) }, v4},
		{"FirewallRules", func(f *cmdutil.Factory) { FirewallRules(f) }, v4},
		{"Function", func(f *cmdutil.Factory) { Function(f) }, v4},
		{"FunctionInstance", func(f *cmdutil.Factory) { FunctionInstance(f) }, v4},
		{"KV", func(f *cmdutil.Factory) { KV(f) }, v4},
		{"NetworkList", func(f *cmdutil.Factory) { NetworkList(f) }, v4},
		{"RealtimePurge", func(f *cmdutil.Factory) { RealtimePurge(f) }, v4},
		{"RulesEngine", func(f *cmdutil.Factory) { RulesEngine(f) }, v4},
		{"WAF", func(f *cmdutil.Factory) { WAF(f) }, v4},
		{"WAFExceptions", func(f *cmdutil.Factory) { WAFExceptions(f) }, v4},
		{"Workloads", func(f *cmdutil.Factory) { Workloads(f) }, v4},

		// served by the legacy API even for accounts on v4, which is why a
		// single "current generation" base URL would be wrong
		{"Domain", func(f *cmdutil.Factory) { Domain(f) }, legacy},
		{"Origin", func(f *cmdutil.Factory) { Origin(f) }, legacy},
		{"PersonalToken", func(f *cmdutil.Factory) { PersonalToken(f) }, legacy},
		{"Variables", func(f *cmdutil.Factory) { Variables(f) }, legacy},

		{"Storage", func(f *cmdutil.Factory) { Storage(f) }, storage},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &recordingConfig{}
			tt.build(&cmdutil.Factory{HttpClient: &http.Client{}, Config: cfg})
			assert.Equal(t, tt.want, cfg.urlKey(t))
		})
	}
}

// Every client authenticates with the configured credential.
func TestClientsUseTheConfiguredCredential(t *testing.T) {
	cfg := &recordingConfig{}
	Applications(&cmdutil.Factory{HttpClient: &http.Client{}, Config: cfg})
	assert.Contains(t, cfg.keys, "token")
}

// The profile commands revoke the token of a profile other than the active one,
// so that credential has to be passed in rather than read from config.
func TestPersonalTokenForUsesTheGivenCredential(t *testing.T) {
	cfg := &recordingConfig{}
	PersonalTokenFor(&cmdutil.Factory{HttpClient: &http.Client{}, Config: cfg}, "another-profiles-token")

	assert.Equal(t, "api_url", cfg.urlKey(t))
	assert.NotContains(t, cfg.keys, "token", "the active credential must not be read")
}
