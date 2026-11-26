package connectors

import (
	"context"
	"os"
	"testing"

	"github.com/aziontech/azion-cli/pkg/api/v6/v6test"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
)

func TestMain(m *testing.M) {
	logger.New(zapcore.DebugLevel)
	os.Exit(m.Run())
}

// A Connector is a oneOf; the client returns the variant it holds.
func TestConnectorVariantIsReturned(t *testing.T) {
	api := &v6test.API{}
	http := map[string]any{"id": 11, "name": "origin", "type": "http", "version_id": "C1", "version_state": "draft",
		"attributes": map[string]any{"addresses": []any{map[string]any{"address": "example.com"}}}}
	api.On("GET", "/workspace/connectors", v6test.JSON(v6test.Page(http)))
	c := NewClient(api.Client(), v6test.URL, "t")

	conn, found, err := c.FindByName(context.Background(), "origin")

	require.NoError(t, err)
	require.True(t, found)
	assert.Equal(t, int64(11), conn.GetId())
	assert.Equal(t, "C1", conn.GetVersionId())
}

// The API answered 200 with a field the SDK does not know, so the SDK failed
// to decode it. That is a failed call: the client returns an error that prints,
// never a nil error (which made the caller go on with an empty result) or a
// StatusError around nil (which panicked when printed).
func TestUndecodableSuccessIsAnError(t *testing.T) {
	body := map[string]any{
		"version_id": "A32XCE2E", "version_state": "draft", "name": "eleventao-v6-test", "type": "storage",
		"attributes":        map[string]any{"bucket": "eleventao-v6-test", "prefix": "20261105111008"},
		"source_version_id": nil, "snapshot": nil, "active": true,
	}
	// SDK 0.3.0 accepts the source_version_id that first broke this, and
	// unknown fields too; a field of the wrong type still cannot be decoded.
	body["version_id"] = 12345
	api := &v6test.API{}
	api.On("GET", "/workspace/connectors/11/versions/A32XCE2E", v6test.JSON(body))
	c := NewClient(api.Client(), v6test.URL, "t")

	_, err := c.GetVersion(context.Background(), 11, "A32XCE2E")
	require.Error(t, err)
	assert.NotPanics(t, func() { _ = err.Error() })
	assert.Contains(t, err.Error(), "oneOf(ConnectorVersion)")
}
