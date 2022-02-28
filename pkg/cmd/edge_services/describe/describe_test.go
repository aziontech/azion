package describe

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"

	errmsg "github.com/aziontech/azion-cli/pkg/cmd/edge_services/error_messages"
	"github.com/aziontech/azion-cli/pkg/httpmock"
	"github.com/aziontech/azion-cli/pkg/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDescribe(t *testing.T) {

	t.Run("service_id not sent", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("GET", "edge_services/1234/resources/666"),
			httpmock.StringResponse("Error: You must provide a service_id and a resource_id as arguments. Use -h or --help for more information"),
		)

		f, _, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)

		cmd.SetIn(&bytes.Buffer{})
		cmd.SetOut(ioutil.Discard)
		cmd.SetErr(ioutil.Discard)

		_, err := cmd.ExecuteC()
		require.ErrorIs(t, err, errmsg.ErrorMissingServiceIdArgument)
	})

	t.Run("service not found", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("GET", "edge_services/1234"),
			httpmock.StatusStringResponse(http.StatusNotFound, "{}"),
		)

		f, _, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)

		cmd.SetArgs([]string{"1234"})
		cmd.SetIn(&bytes.Buffer{})
		cmd.SetOut(ioutil.Discard)
		cmd.SetErr(ioutil.Discard)

		_, err := cmd.ExecuteC()
		require.Error(t, err)
	})

	t.Run("valid service", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("GET", "edge_services/1234"),
			httpmock.JSONFromString(
				`{
                    "id":1209,
                    "name":"ApeService",
                    "updated_at":"2021-12-15T21:03:54Z",
                    "last_editor":"azion-alfreds",
                    "active":true,
                    "bound_nodes":4,
                    "permissions":["read","write"]
                }`,
			),
		)

		f, stdout, _ := testutils.NewFactory(mock)
		cmd := NewCmd(f)

		cmd.SetArgs([]string{"1234"})
		cmd.SetIn(&bytes.Buffer{})
		cmd.SetOut(ioutil.Discard)
		cmd.SetErr(ioutil.Discard)

		_, err := cmd.ExecuteC()
		require.NoError(t, err)

		assert.Equal(t,
			`ID: 1209
Name: ApeService
Updated at: 2021-12-15T21:03:54Z
Last Editor: azion-alfreds
Active: true
Bound Nodes: 4
Permissions: [read write]
`,
			stdout.String(),
		)
	})

	t.Run("valid service with vars", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("GET", "edge_services/1234"),
			httpmock.JSONFromString(
				`{
                    "id":1209,
                    "name":"ApeService",
                    "updated_at":"2021-12-15T21:03:54Z",
                    "last_editor":"azion-alfreds",
                    "active":true,
                    "bound_nodes":4,
                    "permissions":["read","write"],
					"variables": [
						{
						"name": "teste",
						"value": "oteste"
					    }
					]
                }`,
			),
		)

		f, stdout, _ := testutils.NewFactory(mock)
		cmd := NewCmd(f)

		cmd.SetArgs([]string{"1234", "--with-variables", "True"})
		cmd.SetIn(&bytes.Buffer{})
		cmd.SetOut(ioutil.Discard)
		cmd.SetErr(ioutil.Discard)

		_, err := cmd.ExecuteC()
		require.NoError(t, err)

		assert.Equal(t,
			`ID: 1209
Name: ApeService
Updated at: 2021-12-15T21:03:54Z
Last Editor: azion-alfreds
Active: true
Bound Nodes: 4
Permissions: [read write]
Variables:
 Name: teste	Value: oteste
`,
			stdout.String(),
		)
	})

}
