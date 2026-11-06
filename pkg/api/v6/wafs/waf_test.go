package wafs

import (
	"context"
	"os"
	"testing"

	"github.com/aziontech/azion-cli/pkg/api/v6/v6test"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
)

func TestMain(m *testing.M) {
	logger.New(zapcore.DebugLevel)
	os.Exit(m.Run())
}

// The WAF SDK takes the version id before the WAF id; the client's methods
// take them in the same order as every other client.
func TestIDsLandInTheRightPlace(t *testing.T) {
	api := &v6test.API{}
	api.On("POST", "/workspace/wafs/9/versions/V2/build", v6test.Status(202))
	api.On("GET", "/workspace/wafs/9/versions/V2", v6test.JSON(v6test.Version(9, "w", "V2", "ready")))
	api.On("GET", "/workspace/wafs/9/versions/V2/exceptions", v6test.JSON(v6test.Page()))
	c := NewClient(api.Client(), v6test.URL, "t")
	ctx := context.Background()

	require.NoError(t, c.Build(ctx, 9, "V2", "x"))
	state, err := c.VersionState(ctx, 9, "V2")
	require.NoError(t, err)
	assert.EqualValues(t, "ready", state)
	_, err = c.ListExceptions(ctx, 9, "V2", &contracts.ListOptions{Page: 1, PageSize: 100})
	require.NoError(t, err)
}
