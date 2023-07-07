package list

import (
	"github.com/aziontech/azion-cli/pkg/logger"
	"go.uber.org/zap/zapcore"
	"testing"

	"github.com/aziontech/azion-cli/pkg/httpmock"
	"github.com/aziontech/azion-cli/pkg/testutils"
	"github.com/stretchr/testify/require"
)

func TestNewCmd(t *testing.T) {
	logger.New(zapcore.DebugLevel)
	t.Run("list url not found", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("GET", "/edge_applications/12321/functions_instance"),
			httpmock.StatusStringResponse(404, "Not Found"),
		)

		f, _, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)

		cmd.SetArgs([]string{"-a", "12321"})

		_, err := cmd.ExecuteC()
		require.Error(t, err)
	})

	t.Run("empty", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("GET", "/edge_applications/12321/functions_instance"),
			httpmock.StatusStringResponse(204, ""),
		)

		f, _, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)

		cmd.SetArgs([]string{"-a", "12321"})

		_, err := cmd.ExecuteC()
		require.Error(t, err)
	})
	t.Run("list empty", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("GET", "/edge_applications/12321/functions_instance"),
			httpmock.JSONFromFile(".fixtures/no_resp.json"),
		)

		f, _, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)

		cmd.SetArgs([]string{"-a", "12321"})

		_, err := cmd.ExecuteC()
		require.Error(t, err)
	})
	t.Run("list success", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("GET", "/edge_applications/12321/functions_instance"),
			httpmock.JSONFromFile(".fixtures/resp.json"),
		)

		f, _, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)

		cmd.SetArgs([]string{"-a", "12321"})

		_, err := cmd.ExecuteC()
		require.Error(t, err)
	})
}
