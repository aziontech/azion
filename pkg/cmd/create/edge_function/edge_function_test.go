package edgefunction

import (
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/utils"
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

func TestCreate(t *testing.T) {
	logger.New(zapcore.DebugLevel)
	t.Run("create new Edge Function", func(t *testing.T) {
		t.Parallel()

		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("POST", "edge_functions/functions"),
			httpmock.JSONFromString(successResponse),
		)

		f, stdout, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)

		code, _ := os.CreateTemp(t.TempDir(), "func*.js")
		_, _ = code.WriteString("function which() { return 'gambit';}")

		args, _ := os.CreateTemp(t.TempDir(), "args*.json")
		_, _ = args.WriteString(`{"best_sweet": "yakitori"}`)

		cmd.SetArgs([]string{"--name", "SUPAN_FUNCTION", "--active", "true", "--args", args.Name(), "--code", code.Name()})

		err := cmd.Execute()

		require.NoError(t, err)
		require.Equal(t, fmt.Sprintf(msg.CreateOutputSuccess, 1337), stdout.String())
	})

	t.Run("bad request", func(t *testing.T) {
		t.Parallel()

		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("POST", "edge_functions/functions"),
			httpmock.StatusStringResponse(http.StatusBadRequest, "Invalid"),
		)

		f, _, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)

		file, _ := os.CreateTemp(t.TempDir(), "func*.js")
		t.TempDir()
		cmd.SetArgs([]string{"--name", "BIRD", "--active", "true", "--initiator-type", "edge_bird", "--code", file.Name()})

		err := cmd.Execute()

		require.Error(t, err)
	})

	t.Run("create with file", func(t *testing.T) {
		t.Parallel()

		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("POST", "edge_functions/functions"),
			httpmock.JSONFromString(successResponse),
		)

		f, stdout, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)

		file, _ := os.CreateTemp(t.TempDir(), "func*.js")
		t.TempDir()
		cmd.SetArgs([]string{"--name", "BIRD", "--active", "true", "--code", file.Name(), "--file", "./fixtures/create.json"})

		err := cmd.Execute()

		require.NoError(t, err)
		require.Equal(t, fmt.Sprintf(msg.CreateOutputSuccess, 1337), stdout.String())
	})

	t.Run("error file not exist", func(t *testing.T) {
		t.Parallel()

		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("POST", "edge_functions/functions"),
			httpmock.JSONFromString("{}"),
		)

		f, _, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)

		cmd.SetArgs([]string{"--file", "./fixtures/not_exist.json"})

		err := cmd.Execute()
		require.ErrorIs(t, err, utils.ErrorUnmarshalReader)
	})

	t.Run("Error Field active not is boolean", func(t *testing.T) {
		t.Parallel()

		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("POST", "edge_functions/functions"),
			httpmock.JSONFromString(successResponse),
		)

		f, _, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)

		code, _ := os.CreateTemp(t.TempDir(), "func*.js")
		_, _ = code.WriteString("function which() { return 'gambit';}")

		args, _ := os.CreateTemp(t.TempDir(), "args*.json")
		_, _ = args.WriteString(`{"best_sweet": "yakitori"}`)

		cmd.SetArgs([]string{"--name", "SUPAN_FUNCTION", "--active", "12321", "--args", args.Name(), "--code", code.Name()})

		err := cmd.Execute()
		stringErr := fmt.Sprintf("%s: %s", msg.ErrorActiveFlag, "12321")
		if stringErr == err.Error() {
			return
		}
		t.Fatalf("Error: %q", err)
	})

	t.Run("Error create function request api", func(t *testing.T) {
		t.Parallel()

		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("POST", "edge_functions/functions"),
			httpmock.StatusStringResponse(http.StatusInternalServerError, "Invalid"),
		)

		f, _, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)

		code, _ := os.CreateTemp(t.TempDir(), "func*.js")
		_, _ = code.WriteString("function which() { return 'gambit';}")

		args, _ := os.CreateTemp(t.TempDir(), "args*.json")
		_, _ = args.WriteString(`{"best_sweet": "yakitori"}`)

		cmd.SetArgs([]string{"--name", "SUPAN_FUNCTION", "--active", "true", "--args", args.Name(), "--code", code.Name()})

		err := cmd.Execute()
		stringErr := "Failed to create Edge Function: The server could not process the request because an internal and unexpected problem occurred. Wait a few seconds and try again. For more information run the command again using the '--debug' flag. If the problem persists, contact Azion’s support. Check your settings and try again. If the error persists, contact Azion support"
		if stringErr == err.Error() {
			return
		}
		t.Fatalf("Error: %q", err)
	})
}
