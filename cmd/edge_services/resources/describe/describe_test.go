package describe

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/httpmock"
	"github.com/aziontech/azion-cli/pkg/iostreams"
	"github.com/aziontech/azion-cli/utils"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newFactory(mock *httpmock.Registry) (factory *cmdutil.Factory, out *bytes.Buffer, err *bytes.Buffer) {
	stdout, stderr := &bytes.Buffer{}, &bytes.Buffer{}
	f := &cmdutil.Factory{
		HttpClient: func() (*http.Client, error) {
			return &http.Client{Transport: mock}, nil
		},
		IOStreams: &iostreams.IOStreams{
			Out: stdout,
			Err: stderr,
		},
		Config: viper.New(),
	}
	return f, stdout, stderr
}

func TestDescribe(t *testing.T) {

	t.Run("resource not found", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("GET", "edge_services/1234/resources/666"),
			httpmock.StatusStringResponse(http.StatusNotFound, "{}"),
		)

		f, _, _ := newFactory(mock)

		cmd := NewCmd(f)

		cmd.SetArgs([]string{"1234", "666"})
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
			httpmock.StringResponse("Error: You must provide a service_id and a resource_id as arguments. Use -h or --help for more information"),
		)

		f, _, _ := newFactory(mock)

		cmd := NewCmd(f)

		cmd.SetArgs([]string{"1234"})
		cmd.SetIn(&bytes.Buffer{})
		cmd.SetOut(ioutil.Discard)
		cmd.SetErr(ioutil.Discard)

		_, err := cmd.ExecuteC()
		require.ErrorIs(t, err, utils.ErrorMissingResourceIdArgument)
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
					"type": "Install",
					"content": "echo 1"
                }`,
			),
		)

		f, stdout, _ := newFactory(mock)

		cmd := NewCmd(f)

		cmd.SetArgs([]string{"1234", "69420"})
		cmd.SetIn(&bytes.Buffer{})
		cmd.SetOut(ioutil.Discard)
		cmd.SetErr(ioutil.Discard)

		_, err := cmd.ExecuteC()
		require.NoError(t, err)

		assert.Equal(t,
			`ID: 69420
Name: Le Teste
Type: Install
Content type: Shell Script
Content: 
echo 1`,
			stdout.String(),
		)
	})

}
