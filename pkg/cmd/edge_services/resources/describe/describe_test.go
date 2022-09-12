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

	t.Run("resource not found", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("GET", "edge_services/1234/resources/666"),
			httpmock.StatusStringResponse(http.StatusNotFound, "{}"),
		)

		f, _, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)

		cmd.SetArgs([]string{"--service-id", "1234", "--resource-id", "666"})
		cmd.SetIn(&bytes.Buffer{})
		cmd.SetOut(ioutil.Discard)
		cmd.SetErr(ioutil.Discard)

		_, err := cmd.ExecuteC()
		require.Error(t, err)
	})

	t.Run("not all arguments were sent", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("GET", "edge_services/1234/resources/666"),
			httpmock.StringResponse("Error: You must provide a service_id and a resource_id. Use -h or --help for more information"),
		)

		f, _, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)

		cmd.SetArgs([]string{"--service-id", "1234"})
		cmd.SetIn(&bytes.Buffer{})
		cmd.SetOut(ioutil.Discard)
		cmd.SetErr(ioutil.Discard)

		_, err := cmd.ExecuteC()
		require.ErrorIs(t, err, errmsg.ErrorMissingResourceIdArgument)
	})

	t.Run("valid resource", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("GET", "edge_services/1234/resources/69420"),
			httpmock.JSONFromString(
				`{
                    "id": 69420,
                    "name": "Le Teste",
					"content_type": "Shell Script",
					"trigger": "Install",
					"content": "echo 1"
                }`,
			),
		)

		f, stdout, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)

		cmd.SetArgs([]string{"-s", "1234", "-r", "69420"})
		cmd.SetIn(&bytes.Buffer{})
		cmd.SetOut(ioutil.Discard)
		cmd.SetErr(ioutil.Discard)

		_, err := cmd.ExecuteC()
		require.NoError(t, err)

		assert.Equal(t,
			`ID: 69420
Name: Le Teste
Trigger: Install
Content type: Shell Script
Content: 
echo 1`,
			stdout.String(),
		)
	})

}
