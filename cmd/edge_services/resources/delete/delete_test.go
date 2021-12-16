package delete

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/httpmock"
	"github.com/aziontech/azion-cli/pkg/iostreams"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreate(t *testing.T) {
	t.Run("delete resource by id", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("DELETE", "edge_services/1234/resources/456"),
			httpmock.StatusStringResponse(204, ""),
		)

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

		cmd := NewCmd(f)

		cmd.SetArgs([]string{"1234", "456"})
		cmd.SetIn(&bytes.Buffer{})
		cmd.SetOut(ioutil.Discard)
		cmd.SetErr(ioutil.Discard)

		_, err := cmd.ExecuteC()
		require.NoError(t, err)

		assert.Equal(t, "", stdout.String())
	})

	t.Run("delete missing resource", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("DELETE", "edge_services/1234/resources/456"),
			httpmock.StatusStringResponse(404, "Not Found"),
		)

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

		cmd := NewCmd(f)

		cmd.SetArgs([]string{"1234", "456"})
		cmd.SetIn(&bytes.Buffer{})
		cmd.SetOut(ioutil.Discard)
		cmd.SetErr(ioutil.Discard)

		_, err := cmd.ExecuteC()
		require.Error(t, err)
	})
}
