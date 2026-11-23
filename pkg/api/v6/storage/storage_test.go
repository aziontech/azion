package storage

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/aziontech/azion-cli/pkg/api/v6/v6test"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
)

func TestMain(m *testing.M) {
	logger.New(zapcore.DebugLevel)
	os.Exit(m.Run())
}

const buckets = "/workspace/storage/buckets"

func client(api *v6test.API) *Client { return NewClient(api.Client(), v6test.URL, "token") }

func bucket(name string) map[string]any {
	return map[string]any{
		"name":             name,
		"workloads_access": "read_only",
		"state":            "executed",
		"product_version":  "1.0",
		"last_editor":      "someone",
		"last_modified":    "2026-11-23T15:00:39Z",
	}
}

// The v4 SDK rejected the "state" field the API now returns, which failed every
// first deploy after the bucket had been created.
func TestCreateBucketAcceptsState(t *testing.T) {
	api := &v6test.API{}
	api.On("POST", buckets, v6test.Status(201, v6test.Data(bucket("site"))))

	err := client(api).CreateBucket(context.Background(), "site", "read_write")

	require.NoError(t, err)
	assert.Equal(t, map[string]any{"name": "site", "workloads_access": "read_write"}, api.Bodies("POST", buckets)[0])
}

func TestCreateBucketAccepted(t *testing.T) {
	api := &v6test.API{}
	api.On("POST", buckets, v6test.Status(202, map[string]any{"state": "pending", "data": bucket("site")}))

	require.NoError(t, client(api).CreateBucket(context.Background(), "site", "read_only"))
}

func TestCreateBucketNameInUse(t *testing.T) {
	api := &v6test.API{}
	api.On("POST", buckets, v6test.Status(409, map[string]any{"errors": []map[string]any{{"detail": "bucket name is already in use"}}}))

	err := client(api).CreateBucket(context.Background(), "site", "read_only")

	require.ErrorIs(t, err, utils.ErrorNameInUse)
}

func TestCreateCredentials(t *testing.T) {
	const path = "/workspace/storage/credentials"
	api := &v6test.API{}
	api.On("POST", path, v6test.Status(201, v6test.Data(map[string]any{
		"id": 7, "name": "site", "access_key": "AK", "secret_key": "SK",
		"capabilities": []string{"readFiles"}, "buckets": []string{"site"},
		"expiration_date": "2027-11-23T15:00:39Z",
	})))
	expires := time.Date(2027, 11, 23, 15, 0, 39, 0, time.UTC)

	creds, err := client(api).CreateCredentials(context.Background(), CredentialRequest{
		Name: "site", Capabilities: []string{"readFiles"}, Buckets: []string{"site"}, ExpiresAt: &expires,
	})

	require.NoError(t, err)
	assert.Equal(t, "AK", creds.GetAccessKey())
	assert.Equal(t, "SK", creds.GetSecretKey())
	body := api.Bodies("POST", path)[0]
	assert.Equal(t, "site", body["name"])
	assert.Equal(t, []any{"site"}, body["buckets"])
	assert.Equal(t, "2027-11-23T15:00:39Z", body["expiration_date"])
}
