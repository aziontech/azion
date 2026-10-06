package realtime_purge

import (
	"context"
	"net/http"
	"testing"

	"github.com/aziontech/azion-cli/pkg/httpmock"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/utils"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
)

// TestPurgeWithoutResponseReturnsAnErrorInsteadOfPanicking covers the case
// where the request never reaches a server — no network, DNS failure, refused
// connection, timeout.
//
// The generated SDK returns a nil *http.Response alongside the error in that
// case. Before the nil guard, the error path called utils.LogAndRewindBody
// with that nil response, which dereferences it on its first line, so the CLI
// died with a Go stack trace instead of reporting the failure. These calls now
// skip straight to utils.ErrorPerStatusCode, which has handled a nil response
// all along.
func TestPurgeWithoutResponseReturnsAnErrorInsteadOfPanicking(t *testing.T) {
	logger.New(zapcore.DebugLevel)

	// An empty registry matches nothing, so every request fails at transport
	// level and yields a nil response.
	mock := &httpmock.Registry{}
	client := NewClient(&http.Client{Transport: mock}, "http://api.example.invalid", "token")

	calls := map[string]func() error{
		"PurgeUrls":     func() error { return client.PurgeUrls(context.Background(), []string{"http://x/a"}) },
		"PurgeWildcard": func() error { return client.PurgeWildcard(context.Background(), []string{"http://x/*"}) },
		"PurgeCacheKey": func() error {
			return client.PurgeCacheKey(context.Background(), []string{"http://x/a"}, "edge_caching")
		},
	}

	for name, call := range calls {
		t.Run(name, func(t *testing.T) {
			var err error
			require.NotPanics(t, func() { err = call() }, "a response-less failure must not panic")
			require.Error(t, err, "the failure must be reported as an error")
			// utils.ErrorPerStatusCode maps a nil response through
			// checkStatusCode500Error, so the user gets the timeout message or
			// the internal-error one, never a crash.
			require.True(t,
				err == utils.ErrorInternalServerError || err == utils.ErrorTimeoutAPICall,
				"got %q, want the internal-error or timeout message", err)
		})
	}
}
