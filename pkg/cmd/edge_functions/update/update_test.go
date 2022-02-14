package update

import (
	"net/http"
	"os"
	"testing"

	"github.com/aziontech/azion-cli/pkg/httpmock"
	"github.com/aziontech/azion-cli/pkg/testutils"
	"github.com/stretchr/testify/require"
)

var successResponse string = `
{
    "results":{
        "id":1337,
        "name":"SUUPA_FUNCTION",
        "language":"javascript",
        "code":"async function handleRequest(request) {return new Response(\"Hello World!\",{status:200})}",
        "json_args":{"a":1,"b":2},
        "function_to_run":"",
        "initiator_type":"edge_application",
        "active":true,
        "last_editor":"testando@azion.com",
        "modified":"2022-01-26T12:31:09.865515Z",
        "reference_count":0
    },
    "schema_version":3
}
`

func TestUpdate(t *testing.T) {
	t.Run("update edge function", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("PATCH", "edge_functions/1337"),
			httpmock.JSONFromString(successResponse),
		)

		f, stdout, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)

		cmd.SetArgs([]string{"1337", "--name", "ATUALIZANDO", "--active", "false"})

		err := cmd.Execute()

		require.NoError(t, err)
		require.Equal(t, "Updated Edge Function with ID 1337\n", stdout.String())
	})

	t.Run("update code and args", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("PATCH", "edge_functions/1337"),
			httpmock.JSONFromString(successResponse),
		)

		f, stdout, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)

		code, _ := os.CreateTemp(t.TempDir(), "myfunc*.js")
		_, _ = code.WriteString("function elevator() { return 'aufzug';}")

		args, _ := os.CreateTemp(t.TempDir(), "args*.json")
		_, _ = args.WriteString(`{"best_editor": "vim"}`)

		cmd.SetArgs([]string{"1337", "--code", code.Name(), "--args", args.Name()})

		err := cmd.Execute()

		require.NoError(t, err)
		require.Equal(t, "Updated Edge Function with ID 1337\n", stdout.String())
	})

	t.Run("bad request", func(t *testing.T) {
		mock := &httpmock.Registry{}
		mock.Register(
			httpmock.REST("PATCH", "edge_functions/1234"),
			httpmock.StatusStringResponse(http.StatusBadRequest, `{"details": "invalid field active"}`),
		)

		f, _, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)

		cmd.SetArgs([]string{"1234", "--active", "unactive"})

		err := cmd.Execute()

		require.Error(t, err)
	})

	t.Run("update with file", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("PATCH", "edge_functions/1337"),
			httpmock.JSONFromString(successResponse),
		)

		f, stdout, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)

		cmd.SetArgs([]string{"--in", "./fixtures/update.json"})

		err := cmd.Execute()

		require.NoError(t, err)
		require.Equal(t, "Updated Edge Function with ID 1337\n", stdout.String())
	})
}
