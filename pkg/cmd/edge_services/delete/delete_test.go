package delete

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/aziontech/azion-cli/pkg/httpmock"
	"github.com/aziontech/azion-cli/pkg/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreate(t *testing.T) {
	t.Run("delete service by id", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("DELETE", "edge_services/1234"),
			httpmock.StatusStringResponse(204, ""),
		)

		f, stdout, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)
		cmd.SetArgs([]string{"1234"})
		cmd.SetIn(&bytes.Buffer{})
		cmd.SetOut(ioutil.Discard)
		cmd.SetErr(ioutil.Discard)

		_, err := cmd.ExecuteC()
		require.NoError(t, err)

		assert.Equal(t, `Service 1234 was successfully deleted
`, stdout.String())
	})

	t.Run("delete service that is not found", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("DELETE", "edge_services/1234"),
			httpmock.StatusStringResponse(404, "Not Found"),
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
}
