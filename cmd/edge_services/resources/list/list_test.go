package list

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/httpmock"
	"github.com/aziontech/azion-cli/pkg/iostreams"
	"github.com/aziontech/azion-cli/utils"
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
	}
	return f, stdout, stderr
}

func TestList(t *testing.T) {
	t.Run("more than one resource", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("GET", "edge_services/1234/resources"),
			httpmock.JSONFromFile("./fixtures/resources.json"),
		)

		f, stdout, _ := newFactory(mock)

		cmd := NewCmd(f)

		cmd.SetArgs([]string{"1234"})
		cmd.SetIn(&bytes.Buffer{})
		cmd.SetOut(ioutil.Discard)
		cmd.SetErr(ioutil.Discard)

		_, err := cmd.ExecuteC()
		require.NoError(t, err)

		assert.Equal(t,
			`ID: 82587     Name: /tmp/abacatito 
ID: 82588     Name: /tmp/abacatito 
ID: 82592     Name: /tmp/test/asasa 
ID: 82603     Name: /tmp/namechanged 
ID: 82606     Name: /tmp/abacatito 
ID: 82611     Name: /tmp/test/assssas 
`,
			stdout.String(),
		)
	})

	t.Run("no resources", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("GET", "edge_services/1234/resources"),
			httpmock.JSONFromFile("./fixtures/noresources.json"),
		)

		f, stdout, _ := newFactory(mock)

		cmd := NewCmd(f)

		cmd.SetArgs([]string{"1234"})
		cmd.SetIn(&bytes.Buffer{})
		cmd.SetOut(ioutil.Discard)
		cmd.SetErr(ioutil.Discard)

		_, err := cmd.ExecuteC()
		require.NoError(t, err)

		assert.Equal(t, ``, stdout.String())
	})

	t.Run("no resource_id sent", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("GET", "edge_services/1234/resources"),
			httpmock.StringResponse("Error: You must provide a service_id as an argument. Use -h or --help for more information"),
		)

		f, _, _ := newFactory(mock)

		cmd := NewCmd(f)

		cmd.SetIn(&bytes.Buffer{})
		cmd.SetOut(ioutil.Discard)
		cmd.SetErr(ioutil.Discard)

		_, err := cmd.ExecuteC()
		require.ErrorIs(t, err, utils.ErrorMissingServiceIdArgument)
	})

	t.Run("invalid resource_id", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("GET", "edge_services/666/resources"),
			httpmock.StatusStringResponse(http.StatusNotFound, "Error: 404 Not Found"),
		)

		f, _, _ := newFactory(mock)

		cmd := NewCmd(f)

		cmd.SetArgs([]string{"666"})
		cmd.SetIn(&bytes.Buffer{})
		cmd.SetOut(ioutil.Discard)
		cmd.SetErr(ioutil.Discard)

		_, err := cmd.ExecuteC()
		require.Error(t, err, "Error: 404 Not Found")
	})
}
