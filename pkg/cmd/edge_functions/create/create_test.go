package create

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
        "json_args":{"a":1,"b":2},"function_to_run":"",
        "initiator_type":"edge_application",
        "active":true,
        "last_editor":"fsmiamoto@gmail.com",
        "modified":"2022-01-26T12:31:09.865515Z",
        "reference_count":0
    },
    "schema_version":3
}
`

func TestCreate(t *testing.T) {
	t.Run("create new edge function", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("POST", "edge_functions"),
			httpmock.JSONFromString(successResponse),
		)

		f, stdout, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)

		file, _ := os.CreateTemp(t.TempDir(), "myfunc*.js")
		cmd.SetArgs([]string{"--name", "SUUPA_FUNCTION", "--active", "true", "--initiator-type", "edge_application", "--code", file.Name(), "--language", "javascript"})

		err := cmd.Execute()

		require.NoError(t, err)
		require.Equal(t, "Created Edge Function with ID 1337\n", stdout.String())
	})

	t.Run("bad request", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("POST", "edge_functions"),
			httpmock.StatusStringResponse(http.StatusBadRequest, "Invalid"),
		)

		f, _, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)

		file, _ := os.CreateTemp(t.TempDir(), "myfunc*.js")
		t.TempDir()
		cmd.SetArgs([]string{"--name", "BIRL", "--active", "true", "--initiator-type", "edge_birl", "--code", file.Name(), "--language", "javascript"})

		err := cmd.Execute()

		require.Error(t, err)
	})
}
