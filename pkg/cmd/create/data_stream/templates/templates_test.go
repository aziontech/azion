package templates

import (
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/aziontech/azion-cli/pkg/logger"
	"go.uber.org/zap/zapcore"

	msg "github.com/aziontech/azion-cli/messages/templates"
	"github.com/aziontech/azion-cli/pkg/httpmock"
	"github.com/aziontech/azion-cli/pkg/testutils"
	"github.com/stretchr/testify/require"
)

var errorResponse string = `
{
  "errors": [
    {
      "status": "500",
      "code": "500",
      "title": "string",
      "detail": "string"
    }
  ]
}
`

var successResponse string = `
{
  "data": {
    "id": 1337,
    "name": "My Template",
    "last_editor": "user@example.com",
    "created_at": "2019-08-24T14:15:22Z",
    "last_modified": "2019-08-24T14:15:22Z",
    "custom": true,
    "active": true,
    "data_set": "{\"time\":\"$time\"}"
  }
}
`

func TestCreate(t *testing.T) {
	logger.New(zapcore.DebugLevel)

	t.Run("create with file", func(t *testing.T) {
		t.Parallel()

		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("POST", "workspace/stream/templates"),
			httpmock.JSONFromString(successResponse),
		)

		f, stdout, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)
		cmd.SetArgs([]string{"--file", "./fixtures/create.json"})

		err := cmd.Execute()

		require.NoError(t, err)
		require.Equal(t, fmt.Sprintf(msg.CreateOutputSuccess, 1337), stdout.String())
	})

	t.Run("create with flags", func(t *testing.T) {
		t.Parallel()

		var payload map[string]interface{}

		mock := &httpmock.Registry{}
		mock.Register(
			httpmock.REST("POST", "workspace/stream/templates"),
			httpmock.WithHeader(
				httpmock.RESTPayload(http.StatusOK, successResponse, func(p map[string]interface{}) {
					payload = p
				}),
				"Content-Type", "application/json",
			),
		)

		f, stdout, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)
		cmd.SetArgs([]string{"--name", "My Template", "--data-set", "./fixtures/data_set.json"})

		err := cmd.Execute()

		require.NoError(t, err)
		require.Equal(t, fmt.Sprintf(msg.CreateOutputSuccess, 1337), stdout.String())
		require.Equal(t, "My Template", payload["name"])

		// data_set is sent as the raw string content of the JSON file
		wantDataSet, readErr := os.ReadFile("./fixtures/data_set.json")
		require.NoError(t, readErr)
		require.Equal(t, string(wantDataSet), payload["data_set"])

		// active is only sent when the flag is explicitly set
		_, hasActive := payload["active"]
		require.False(t, hasActive)
	})

	t.Run("create with flags and active false", func(t *testing.T) {
		t.Parallel()

		var payload map[string]interface{}

		mock := &httpmock.Registry{}
		mock.Register(
			httpmock.REST("POST", "workspace/stream/templates"),
			httpmock.WithHeader(
				httpmock.RESTPayload(http.StatusOK, successResponse, func(p map[string]interface{}) {
					payload = p
				}),
				"Content-Type", "application/json",
			),
		)

		f, stdout, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)
		cmd.SetArgs([]string{"--name", "My Template", "--data-set", "./fixtures/data_set.json", "--active", "false"})

		err := cmd.Execute()

		require.NoError(t, err)
		require.Equal(t, fmt.Sprintf(msg.CreateOutputSuccess, 1337), stdout.String())
		require.Equal(t, false, payload["active"])
	})

	t.Run("data set file not found", func(t *testing.T) {
		t.Parallel()

		mock := &httpmock.Registry{}

		f, _, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)
		cmd.SetArgs([]string{"--name", "My Template", "--data-set", "./fixtures/does-not-exist.json"})

		err := cmd.Execute()

		require.ErrorIs(t, err, msg.ErrorDataSetFlag)
	})

	t.Run("invalid data set json", func(t *testing.T) {
		t.Parallel()

		mock := &httpmock.Registry{}

		f, _, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)
		cmd.SetArgs([]string{"--name", "My Template", "--data-set", "./fixtures/invalid_data_set.json"})

		err := cmd.Execute()

		require.ErrorIs(t, err, msg.ErrorParseDataSet)
	})

	t.Run("invalid active flag", func(t *testing.T) {
		t.Parallel()

		mock := &httpmock.Registry{}

		f, _, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)
		cmd.SetArgs([]string{"--name", "My Template", "--data-set", "./fixtures/data_set.json", "--active", "notabool"})

		err := cmd.Execute()

		require.Error(t, err)
	})

	t.Run("file not found", func(t *testing.T) {
		t.Parallel()

		mock := &httpmock.Registry{}

		f, _, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)
		cmd.SetArgs([]string{"--file", "./fixtures/does-not-exist.json"})

		err := cmd.Execute()

		require.Error(t, err)
	})

	t.Run("api error", func(t *testing.T) {
		t.Parallel()

		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("POST", "workspace/stream/templates"),
			httpmock.StatusStringResponse(http.StatusInternalServerError, errorResponse),
		)

		f, _, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)
		cmd.SetArgs([]string{"--file", "./fixtures/create.json"})

		err := cmd.Execute()

		require.Error(t, err)
	})
}
