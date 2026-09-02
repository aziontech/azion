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

var successResponse string = `
{
  "data": {
    "id": 1337,
    "name": "Updated Template",
    "last_editor": "user@example.com",
    "created_at": "2019-08-24T14:15:22Z",
    "last_modified": "2019-08-24T14:15:22Z",
    "custom": true,
    "active": false,
    "data_set": "{\"time\":\"$time\"}"
  }
}
`

func TestUpdate(t *testing.T) {
	logger.New(zapcore.DebugLevel)

	t.Run("update with file", func(t *testing.T) {
		t.Parallel()

		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("PATCH", "workspace/stream/templates/1337"),
			httpmock.JSONFromString(successResponse),
		)

		f, stdout, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)
		cmd.SetArgs([]string{"--template-id", "1337", "--file", "./fixtures/update.json"})

		err := cmd.Execute()

		require.NoError(t, err)
		require.Equal(t, fmt.Sprintf(msg.UpdateOutputSuccess, 1337), stdout.String())
	})

	t.Run("update with name flag", func(t *testing.T) {
		t.Parallel()

		var payload map[string]interface{}

		mock := &httpmock.Registry{}
		mock.Register(
			httpmock.REST("PATCH", "workspace/stream/templates/1337"),
			httpmock.WithHeader(
				httpmock.RESTPayload(http.StatusOK, successResponse, func(p map[string]interface{}) {
					payload = p
				}),
				"Content-Type", "application/json",
			),
		)

		f, stdout, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)
		cmd.SetArgs([]string{"--template-id", "1337", "--name", "New Name"})

		err := cmd.Execute()

		require.NoError(t, err)
		require.Equal(t, fmt.Sprintf(msg.UpdateOutputSuccess, 1337), stdout.String())
		require.Equal(t, "New Name", payload["name"])
		// only the changed field is sent in the patch
		_, hasDataSet := payload["data_set"]
		require.False(t, hasDataSet)
		_, hasActive := payload["active"]
		require.False(t, hasActive)
	})

	t.Run("update with data-set and active flags", func(t *testing.T) {
		t.Parallel()

		var payload map[string]interface{}

		mock := &httpmock.Registry{}
		mock.Register(
			httpmock.REST("PATCH", "workspace/stream/templates/1337"),
			httpmock.WithHeader(
				httpmock.RESTPayload(http.StatusOK, successResponse, func(p map[string]interface{}) {
					payload = p
				}),
				"Content-Type", "application/json",
			),
		)

		f, stdout, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)
		cmd.SetArgs([]string{"--template-id", "1337", "--data-set", "./fixtures/data_set.json", "--active", "false"})

		err := cmd.Execute()

		require.NoError(t, err)
		require.Equal(t, fmt.Sprintf(msg.UpdateOutputSuccess, 1337), stdout.String())

		wantDataSet, readErr := os.ReadFile("./fixtures/data_set.json")
		require.NoError(t, readErr)
		require.Equal(t, string(wantDataSet), payload["data_set"])
		require.Equal(t, false, payload["active"])
		// name was not provided, so it must not be part of the patch
		_, hasName := payload["name"]
		require.False(t, hasName)
	})

	t.Run("no fields provided", func(t *testing.T) {
		t.Parallel()

		mock := &httpmock.Registry{}

		f, _, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)
		cmd.SetArgs([]string{"--template-id", "1337"})

		err := cmd.Execute()

		require.ErrorIs(t, err, msg.ErrorUpdateNoFields)
	})

	t.Run("invalid active flag", func(t *testing.T) {
		t.Parallel()

		mock := &httpmock.Registry{}

		f, _, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)
		cmd.SetArgs([]string{"--template-id", "1337", "--active", "notabool"})

		err := cmd.Execute()

		require.Error(t, err)
	})

	t.Run("data set file not found", func(t *testing.T) {
		t.Parallel()

		mock := &httpmock.Registry{}

		f, _, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)
		cmd.SetArgs([]string{"--template-id", "1337", "--data-set", "./fixtures/does-not-exist.json"})

		err := cmd.Execute()

		require.ErrorIs(t, err, msg.ErrorDataSetFlag)
	})

	t.Run("file not found", func(t *testing.T) {
		t.Parallel()

		mock := &httpmock.Registry{}

		f, _, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)
		cmd.SetArgs([]string{"--template-id", "1337", "--file", "./fixtures/does-not-exist.json"})

		err := cmd.Execute()

		require.Error(t, err)
	})

	t.Run("not found", func(t *testing.T) {
		t.Parallel()

		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("PATCH", "workspace/stream/templates/9999"),
			httpmock.StatusStringResponse(http.StatusNotFound, `{"details": "Template not found"}`),
		)

		f, _, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)
		cmd.SetArgs([]string{"--template-id", "9999", "--file", "./fixtures/update.json"})

		err := cmd.Execute()

		require.Error(t, err)
	})
}
