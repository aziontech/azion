package function

import (
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/aziontech/azion-cli/pkg/logger"
	"go.uber.org/zap/zapcore"

	msg "github.com/aziontech/azion-cli/messages/function"
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
      "detail": "string",
      "source": {
        "pointer": "string",
        "parameter": "string",
        "header": "string"
      },
      "meta": {
        "property1": null,
        "property2": null
      }
    }
  ]
}
`

var successResponse string = `
{
	"state": "executed",
	"data": {
	  "id": 1337,
	  "name": "string",
	  "last_editor": "string",
	  "last_modified": "2019-08-24T14:15:22Z",
	  "product_version": "string",
	  "active": true,
	  "runtime": "azion_js",
	  "execution_environment": "firewall",
	  "code": "string",
	  "default_args": {
		"arg_01": "value_01"
	  },
	  "reference_count": 0,
	  "version": "string",
	  "vendor": "string"
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
			httpmock.StatusStringResponse(http.StatusInternalServerError, errorResponse),
		)

		f, _, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)

		code, _ := os.CreateTemp(t.TempDir(), "func*.js")
		_, _ = code.WriteString("function which() { return 'gambit';}")

		args, _ := os.CreateTemp(t.TempDir(), "args*.json")
		_, _ = args.WriteString(`{"best_sweet": "yakitori"}`)

		cmd.SetArgs([]string{"--name", "SUPAN_FUNCTION", "--active", "true", "--args", args.Name(), "--code", code.Name()})

		err := cmd.Execute()
		if err != nil {
			return
		}
		t.Fatalf("Error: %q", err)
	})
}
