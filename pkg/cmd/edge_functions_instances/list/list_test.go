package list

import (
	"testing"

	"github.com/aziontech/azion-cli/pkg/httpmock"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/testutils"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
)

func TestList(t *testing.T) {
	logger.New(zapcore.DebugLevel)
	t.Run("command list with successes", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("GET", "edge_applications/1674040168/functions_instances"),
			httpmock.JSONFromFile(".fixtures/resp.json"),
		)

		f, _, _ := testutils.NewFactory(mock)
		cmd := NewCmd(f)

		cmd.SetArgs([]string{"-a", "1674040168"})

		_, err := cmd.ExecuteC()
		require.NoError(t, err)
	})

	t.Run("command list response without items", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("GET", "edge_applications/1674040168/functions_instances"),
			httpmock.JSONFromFile(".fixtures/no_resp.json"),
		)

		f, _, _ := testutils.NewFactory(mock)
		cmd := NewCmd(f)

		cmd.SetArgs([]string{"-a", "1674040168"})

		_, err := cmd.ExecuteC()
		require.NoError(t, err)
	})
}
