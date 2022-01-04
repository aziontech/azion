package list

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/httpmock"
	"github.com/aziontech/azion-cli/pkg/iostreams"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestList(t *testing.T) {
	t.Run("more than one service", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("GET", "edge_services/"),
			httpmock.JSONFromFile("./fixtures/services.json"),
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
			Config: viper.New(),
		}

		cmd := NewCmd(f)

		cmd.SetArgs([]string{})
		cmd.SetIn(&bytes.Buffer{})
		cmd.SetOut(ioutil.Discard)
		cmd.SetErr(ioutil.Discard)

		_, err := cmd.ExecuteC()
		require.NoError(t, err)

		assert.Equal(t,
			`ID      NAME
1718    batata
1209    ApeService
1752    ApeService
1751    Testando CLI
1750    testing new code cli
26      Service Henrique Teste
1746    jagaimo
1717    potato
1716    tst-flag
1715    tst-flag
`,
			stdout.String(),
		)
	})

	t.Run("no services", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("GET", "edge_services/"),
			httpmock.JSONFromFile("./fixtures/noservices.json"),
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
			Config: viper.New(),
		}

		cmd := NewCmd(f)

		cmd.SetArgs([]string{})
		cmd.SetIn(&bytes.Buffer{})
		cmd.SetOut(ioutil.Discard)
		cmd.SetErr(ioutil.Discard)

		_, err := cmd.ExecuteC()
		require.NoError(t, err)

		assert.Equal(t, ``, stdout.String())
	})
}
