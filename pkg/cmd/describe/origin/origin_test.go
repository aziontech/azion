package origin

import (
	"net/http"
	"testing"

	"github.com/aziontech/azion-cli/pkg/logger"
	"go.uber.org/zap/zapcore"

	"github.com/aziontech/azion-cli/pkg/httpmock"
	"github.com/aziontech/azion-cli/pkg/testutils"
	"github.com/stretchr/testify/require"
)

func TestDescribe(t *testing.T) {
	logger.New(zapcore.DebugLevel)
	t.Run("describe an origin", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("GET", "edge_applications/123423424/origins/0000000-00000000-00a0a00s0as0-000000"),
			httpmock.JSONFromFile("./fixtures/origins.json"),
		)

		f, _, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)
		cmd.SetArgs([]string{"--application-id", "123423424", "--origin-key", "0000000-00000000-00a0a00s0as0-000000"})

		err := cmd.Execute()
		require.NoError(t, err)
	})
	t.Run("not found", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("GET", "edge_applications/123423424/origin/0000000-00000000-00a0a00s0as0-000000"),
			httpmock.StatusStringResponse(http.StatusNotFound, "Not Found"),
		)

		f, _, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)
		cmd.SetArgs([]string{"--application-id", "123423424", "--origin-key", "0000000-00000000-00a0a00s0as0-000000"})

		err := cmd.Execute()
		require.Error(t, err)
	})
	t.Run("no id sent", func(t *testing.T) {
		mock := &httpmock.Registry{}
		mock.Register(
			httpmock.REST("GET", "edge_applications/123423424/origins/0000000-00000000-00a0a00s0as0-000000"),
			httpmock.StatusStringResponse(http.StatusNotFound, "Not Found"),
		)

		f, _, _ := testutils.NewFactory(mock)
		cmd := NewCmd(f)
		cmd.SetArgs([]string{"--application-id", "123423424", "--origin-key", "0000000-00000000-00a0a00s0as0-000000"})

		err := cmd.Execute()
		require.Error(t, err)
	})
}
