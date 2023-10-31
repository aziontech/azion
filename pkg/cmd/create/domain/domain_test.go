package domain

import (
	"net/http"
	"testing"

	"github.com/aziontech/azion-cli/pkg/logger"
	"go.uber.org/zap/zapcore"

	"github.com/aziontech/azion-cli/pkg/httpmock"
	"github.com/aziontech/azion-cli/pkg/testutils"
	"github.com/stretchr/testify/require"
)

func TestCreate(t *testing.T) {
	logger.New(zapcore.DebugLevel)
	t.Run("create new domains", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("POST", "domains"),
			httpmock.JSONFromFile("./fixtures/response.json"),
		)

		f, _, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)
		cmd.SetArgs([]string{
			"--name", "one piece is the best",
			"--application-id", "1673635841",
			"--cname-access-only", "false",
			"--active", "true",
		})

		err := cmd.Execute()
		require.NoError(t, err)
	})

	t.Run("create with file", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("POST", "domains"),
			httpmock.JSONFromFile("./fixtures/response.json"),
		)

		f, _, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)
		cmd.SetArgs([]string{"--in", "./fixtures/create.json"})

		err := cmd.Execute()
		require.NoError(t, err)
	})

	t.Run("bad request status 400", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("POST", "domains"),
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
			httpmock.REST("POST", "domains"),
			httpmock.StatusStringResponse(http.StatusInternalServerError, "Invalid"),
		)

		f, _, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)
		err := cmd.Execute()
		require.Error(t, err)
	})
}
