package edgefunction

import (
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/aziontech/azion-cli/pkg/logger"
	"go.uber.org/zap/zapcore"

	msg "github.com/aziontech/azion-cli/messages/edge_function"
	"github.com/aziontech/azion-cli/pkg/httpmock"
	"github.com/aziontech/azion-cli/pkg/testutils"
	"github.com/stretchr/testify/require"
)

var successResponse string = `
{
  "state": "pending",
  "data": {
    "id": 1337,
    "name": "string",
    "language": "javascript",
    "code": "string",
    "json_args": {
      "arg_01": "value_01"
    },
    "initiator_type": "edge_application",
    "active": true,
    "reference_count": 0,
    "version": "string",
    "vendor": "string",
    "last_editor": "string",
    "last_modified": "2019-08-24T14:15:22Z",
    "product_version": "string"
  }
}
`

func TestUpdate(t *testing.T) {
	logger.New(zapcore.DebugLevel)
	t.Run("update Function", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("PATCH", "edge_functions/functions/1337"),
			httpmock.JSONFromString(successResponse),
		)

		f, stdout, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)

		cmd.SetArgs([]string{"--function-id", "1337", "--name", "ATUALIZANDO", "--active", "false"})

		err := cmd.Execute()

		require.NoError(t, err)
		require.Equal(t, fmt.Sprintf(msg.UpdateOutputSuccess, 1337), stdout.String())
	})

	t.Run("update code and args", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("PATCH", "edge_functions/functions/1337"),
			httpmock.JSONFromString(successResponse),
		)

		f, stdout, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)

		code, _ := os.CreateTemp(t.TempDir(), "myfunc*.js")
		_, _ = code.WriteString("function elevator() { return 'aufzug';}")

		args, _ := os.CreateTemp(t.TempDir(), "args*.json")
		_, _ = args.WriteString(`{"best_editor": "vim"}`)

		cmd.SetArgs([]string{"--function-id", "1337", "--code", code.Name(), "--args", args.Name()})

		err := cmd.Execute()

		require.NoError(t, err)
		require.Equal(t, fmt.Sprintf(msg.UpdateOutputSuccess, 1337), stdout.String())
	})

	t.Run("bad request", func(t *testing.T) {
		mock := &httpmock.Registry{}
		mock.Register(
			httpmock.REST("PATCH", "edge_functions/functions/1234"),
			httpmock.StatusStringResponse(http.StatusBadRequest, `{"details": "invalid field active"}`),
		)

		f, _, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)

		cmd.SetArgs([]string{"--function-id", "1234", "--active", "unactive"})

		err := cmd.Execute()

		require.Error(t, err)
	})

	t.Run("update with file", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("PATCH", "edge_functions/functions/1337"),
			httpmock.JSONFromString(successResponse),
		)

		f, stdout, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)

		cmd.SetArgs([]string{"--function-id", "1337", "--file", "./fixtures/update.json"})

		err := cmd.Execute()

		require.NoError(t, err)
		require.Equal(t, fmt.Sprintf(msg.UpdateOutputSuccess, 1337), stdout.String())
	})
}
