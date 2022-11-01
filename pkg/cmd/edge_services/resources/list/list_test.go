package list

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	errmsg "github.com/aziontech/azion-cli/messages/edge_services"
	"github.com/aziontech/azion-cli/pkg/httpmock"
	"github.com/aziontech/azion-cli/pkg/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestList(t *testing.T) {
	t.Run("more than one resource", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("GET", "edge_services/1234/resources"),
			httpmock.JSONFromFile("./fixtures/resources.json"),
		)

		f, stdout, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)

		cmd.SetArgs([]string{"--service-id", "1234"})
		cmd.SetIn(&bytes.Buffer{})
		cmd.SetOut(io.Discard)
		cmd.SetErr(io.Discard)

		_, err := cmd.ExecuteC()
		require.NoError(t, err)

		assert.Equal(t,
			`ID       NAME
82587    /tmp/abacatito
82588    /tmp/abacatito
82592    /tmp/test/asasa
82603    /tmp/namechanged
82606    /tmp/abacatito
82611    /tmp/test/assssas
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

		f, stdout, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)

		cmd.SetArgs([]string{"-s", "1234"})
		cmd.SetIn(&bytes.Buffer{})
		cmd.SetOut(io.Discard)
		cmd.SetErr(io.Discard)

		_, err := cmd.ExecuteC()
		require.NoError(t, err)

		assert.Equal(t, `ID    NAME
`, stdout.String())
	})

	t.Run("no resource_id sent", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("GET", "edge_services/1234/resources"),
			httpmock.StringResponse("Error: You must provide a service_id as an argument. Use -h or --help for more information"),
		)

		f, _, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)

		cmd.SetIn(&bytes.Buffer{})
		cmd.SetOut(io.Discard)
		cmd.SetErr(io.Discard)

		_, err := cmd.ExecuteC()
		require.ErrorIs(t, err, errmsg.ErrorMissingServiceIdArgument)
	})

	t.Run("invalid resource_id", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("GET", "edge_services/666/resources"),
			httpmock.StatusStringResponse(http.StatusNotFound, "Error: 404 Not Found"),
		)

		f, _, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)

		cmd.SetArgs([]string{"--service-id", "666"})
		cmd.SetIn(&bytes.Buffer{})
		cmd.SetOut(io.Discard)
		cmd.SetErr(io.Discard)

		_, err := cmd.ExecuteC()
		require.Error(t, err, "Error: 404 Not Found")
	})
}
