package variable

import (
	"net/http"
	"testing"

	"github.com/aziontech/azion-cli/pkg/httpmock"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/testutils"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
)

func TestNewCmd(t *testing.T) {
	logger.New(zapcore.DebugLevel)
	t.Run("create new variable", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("POST", "variables"),
			httpmock.JSONFromFile("fixtures/response.json"),
		)

		f, stdout, _ := testutils.NewFactory(mock)
		cmd := NewCmd(f)
		cmd.SetArgs([]string{
			"--key", "Content-Type",
			"--value", "json",
			"--secret", "true",
		})

		err := cmd.Execute()
		require.NoError(t, err)
		require.Equal(t, "🚀 Created variable with UUID bea9d757-8b83-4b4a-a3b1-49dfd6111303\n", stdout.String())

	})

	t.Run("create with file", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("POST", "variables"),
			httpmock.JSONFromFile("fixtures/response.json"),
		)

		f, stdout, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)
		cmd.SetArgs([]string{
			"--in", "fixtures/create.json",
		})

		err := cmd.Execute()
		require.NoError(t, err)
		require.Equal(t, "🚀 Created variable with UUID bea9d757-8b83-4b4a-a3b1-49dfd6111303\n", stdout.String())
	})

	t.Run("bad request status 400", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("POST", "variables"),
			httpmock.StatusStringResponse(http.StatusBadRequest, "Invalid"),
		)

		f, _, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)
		err := cmd.Execute()
		require.Error(t, err)
	})

	t.Run("internal server error 500", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("POST", "variables"),
			httpmock.StatusStringResponse(http.StatusInternalServerError, "Invalid"),
		)

		f, _, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)
		err := cmd.Execute()
		require.Error(t, err)
	})
}
