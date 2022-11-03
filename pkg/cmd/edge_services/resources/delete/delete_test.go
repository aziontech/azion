package delete

import (
	"bytes"
	"io"
	"testing"

	"github.com/aziontech/azion-cli/pkg/httpmock"
	"github.com/aziontech/azion-cli/pkg/testutils"
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

		f, stdout, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)
		cmd.SetArgs([]string{"--service-id", "1234", "--resource-id", "456"})
		cmd.SetIn(&bytes.Buffer{})
		cmd.SetOut(io.Discard)
		cmd.SetErr(io.Discard)

		_, err := cmd.ExecuteC()
		require.NoError(t, err)

		assert.Equal(t, `Resource 456 was successfully deleted
`, stdout.String())
	})

	t.Run("delete missing resource", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("DELETE", "edge_services/1234/resources/456"),
			httpmock.StatusStringResponse(404, "Not Found"),
		)

		f, _, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)

		cmd.SetArgs([]string{"-s", "1234", "-r", "456"})
		cmd.SetIn(&bytes.Buffer{})
		cmd.SetOut(io.Discard)
		cmd.SetErr(io.Discard)

		_, err := cmd.ExecuteC()
		require.Error(t, err)
	})
}
