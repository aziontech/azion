package apiversion

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSelect(t *testing.T) {
	tests := []struct {
		name string
		info AccountInfo
		want Version
	}{
		{"no client flags means the current generation", AccountInfo{}, V4},
		{"unrelated flags do not change the generation", AccountInfo{ClientFlags: []string{"other"}}, V4},
		{"block_apiv4_incompatible_endpoints blocks v4, so v3", AccountInfo{ClientFlags: []string{BlockAPIV4IncompatibleEndpoints}}, V3},
		{"block_apiv3_access blocks v3, so v4", AccountInfo{ClientFlags: []string{BlockAPIV3Access}}, V4},
		{"the v4 block wins alongside unrelated flags", AccountInfo{ClientFlags: []string{"waf_mode", BlockAPIV4IncompatibleEndpoints}}, V3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, Select(tt.info))
		})
	}
}

// The wire form is persisted in profile settings and reported in metrics.
func TestVersionWireForm(t *testing.T) {
	assert.Equal(t, "v3", V3.String())
	assert.Equal(t, "v4", V4.String())
}

func TestResolve(t *testing.T) {
	t.Run("sends the credential and resolves the account's generation", func(t *testing.T) {
		var gotAuth, gotPath string
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			gotAuth, gotPath = r.Header.Get("Authorization"), r.URL.Path
			_, _ = w.Write([]byte(`{"client_flags":["` + BlockAPIV3Access + `"]}`))
		}))
		defer srv.Close()

		got, err := Resolve(srv.Client(), srv.URL, "a-token")

		require.NoError(t, err)
		assert.Equal(t, V4, got, "block_apiv3_access leaves the account on v4")
		assert.Equal(t, "Token a-token", gotAuth)
		assert.Equal(t, "/account/info", gotPath)
	})

	t.Run("a non-200 is an error, not a version", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer srv.Close()

		got, err := Resolve(srv.Client(), srv.URL, "a-token")

		require.Error(t, err)
		assert.Empty(t, got)
	})

	t.Run("an undecodable body is an error, not a version", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte("not json"))
		}))
		defer srv.Close()

		got, err := Resolve(srv.Client(), srv.URL, "a-token")

		require.Error(t, err)
		assert.Empty(t, got)
	})

	t.Run("an unreachable service is an error, not a version", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))
		url := srv.URL
		srv.Close()

		got, err := Resolve(http.DefaultClient, url, "a-token")

		require.Error(t, err)
		assert.Empty(t, got)
	})
}

// A rejected credential must be distinguishable from a lookup that could not be
// completed: the caller stays quiet about the version in the first case and
// warns in the second.
func TestResolveReportsUnauthorizedDistinctly(t *testing.T) {
	for _, status := range []int{http.StatusUnauthorized, http.StatusForbidden} {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(status)
		}))

		_, err := Resolve(srv.Client(), srv.URL, "a-token")
		srv.Close()

		require.Error(t, err)
		assert.ErrorIs(t, err, ErrUnauthorized, "status %d must report a rejected credential", status)
	}
}

func TestResolveDoesNotReportOtherFailuresAsUnauthorized(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	_, err := Resolve(srv.Client(), srv.URL, "a-token")

	require.Error(t, err)
	assert.NotErrorIs(t, err, ErrUnauthorized, "a server fault is not a credential problem")
}

// Each client flag blocks one generation, so each points at the other. These
// cases are taken from live accounts: profiles carrying
// block_apiv4_incompatible_endpoints are v3, profiles carrying
// block_apiv3_access are v4.
//
// PR #1553 inverted this by swapping the caller's branches, which served the v4
// tree to blocked accounts and the v3 tree to everyone else.
func TestFlagsBlockTheirOwnGeneration(t *testing.T) {
	cases := []struct {
		name  string
		flags []string
		want  Version
	}{
		// observed payloads, profile name in the comment
		{"blocked from v4", []string{BlockAPIV4IncompatibleEndpoints, "waf_mode"}, V3},                                 // PabloV3Prod, ProdV3P
		{"blocked from v3", []string{BlockAPIV3Access, "waf_mode"}, V4},                                                // ConsoleTemplateDemo
		{"blocked from v3, console redirect", []string{"force_redirect_to_console", BlockAPIV3Access, "waf_mode"}, V4}, // ProdV4
		{"unrelated flag only", []string{"waf_mode"}, V4},
		{"no flags at all", nil, V4},
		// contradictory configuration: cannot use v3 wins, the legacy tree would
		// be useless to such an account
		{"both block flags", []string{BlockAPIV4IncompatibleEndpoints, BlockAPIV3Access}, V4},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			assert.Equal(t, c.want, Select(AccountInfo{ClientFlags: c.flags}))
		})
	}
}
