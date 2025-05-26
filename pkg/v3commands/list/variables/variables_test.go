package variables

import (
	"net/http"
	"os"
	"testing"

	"github.com/aziontech/azion-cli/pkg/logger"
	"go.uber.org/zap/zapcore"

	"github.com/aziontech/azion-cli/pkg/httpmock"
	"github.com/aziontech/azion-cli/pkg/testutils"
	"github.com/stretchr/testify/require"
)

func TestList(t *testing.T) {
	logger.New(zapcore.DebugLevel)

	t.Run("list page 1", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST(http.MethodGet, "variables"),
			httpmock.JSONFromFile(".fixtures/variables.json"),
		)

		f, _, _ := testutils.NewFactory(mock)
		cmd := NewCmd(f)

		_, err := cmd.ExecuteC()
		require.NoError(t, err)
	})

	t.Run("no items", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("GET", "variables"),
			httpmock.JSONFromFile(".fixtures/nocontent.json"),
		)

		f, _, _ := testutils.NewFactory(mock)
		cmd := NewCmd(f)

		_, err := cmd.ExecuteC()
		require.NoError(t, err)
	})

	t.Run("list with dump", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST(http.MethodGet, "variables"),
			httpmock.JSONFromFile(".fixtures/variables.json"),
		)

		f, _, _ := testutils.NewFactory(mock)
		cmd := NewCmd(f)

		// Set the arguments for the command to include the --dump flag
		cmd.SetArgs([]string{"--dump"})
		_, err := cmd.ExecuteC()
		require.NoError(t, err)

		// Defer the removal of the .env file after the test is finished
		defer func() {
			err := os.Remove(".env")
			require.NoError(t, err)
		}()
	})

	t.Run("invalid json response", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST(http.MethodGet, "variables"),
			httpmock.StringResponse("invalid json"),
		)

		f, _, _ := testutils.NewFactory(mock)
		cmd := NewCmd(f)

		_, err := cmd.ExecuteC()
		require.Error(t, err)
	})

	t.Run("ask for input application-id", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST(http.MethodGet, "variables"),
			httpmock.JSONFromFile(".fixtures/variables.json"),
		)

		f, _, _ := testutils.NewFactory(mock)
		listcmd := NewListCmd(f)
		cmd := NewCobraCmd(listcmd, f)

		cmd.SetArgs([]string{})
		_, err := cmd.ExecuteC()
		require.NoError(t, err)
	})
}
