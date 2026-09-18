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
		{"no client flags stays on the legacy generation", AccountInfo{}, V3},
		{"unrelated flags do not opt the account into v4", AccountInfo{ClientFlags: []string{"other"}}, V3},
		{"block_apiv4_incompatible_endpoints opts into v4", AccountInfo{ClientFlags: []string{BlockAPIV4IncompatibleEndpoints}}, V4},
		{"block_apiv3_access opts into v4", AccountInfo{ClientFlags: []string{BlockAPIV3Access}}, V4},
		{"either flag is enough alongside others", AccountInfo{ClientFlags: []string{"other", BlockAPIV3Access}}, V4},
		{"both flags together resolve to v4", AccountInfo{ClientFlags: []string{BlockAPIV4IncompatibleEndpoints, BlockAPIV3Access}}, V4},
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
		assert.Equal(t, V4, got)
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
