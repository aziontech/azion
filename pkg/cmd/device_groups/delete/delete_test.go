package delete

import (
	"testing"

	"github.com/aziontech/azion-cli/pkg/httpmock"
	"github.com/aziontech/azion-cli/pkg/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreate(t *testing.T) {

	t.Run("delete instance by id", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("DELETE", "edge_applications/1234/device_groups/4321"),
			httpmock.StatusStringResponse(204, ""),
		)

		f, stdout, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)
		cmd.SetArgs([]string{"--application-id", "1234", "--group-id", "4321"})

		_, err := cmd.ExecuteC()
		require.NoError(t, err)

		assert.Equal(t, "Device Group 4321 was successfully deleted\n", stdout.String())
	})

	t.Run("try delete instance that is not found", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("DELETE", "edge_applications/1234/device_groups/4321"),
			httpmock.StatusStringResponse(404, "Not Found"),
		)

		f, _, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)

		cmd.SetArgs([]string{"-a", "1234", "-g", "4321"})

		_, err := cmd.ExecuteC()
		require.Error(t, err)
	})
}
