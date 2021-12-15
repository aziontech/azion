package describe

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

func TestDescribe(t *testing.T) {
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

		cmd.SetArgs([]string{"1234"})
		cmd.SetIn(&bytes.Buffer{})
		cmd.SetOut(ioutil.Discard)
		cmd.SetErr(ioutil.Discard)

		_, err := cmd.ExecuteC()
		require.NoError(t, err)

		assert.Equal(t,
			stdout.String(),
			"ID: 1209\nName: ApeService\nLast Editor: azion-alfreds\nUpdated at: 2021-12-15T21:03:54Z\nActive: true\nBound Nodes: 4\n",
		)
	})

}
