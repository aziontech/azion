package delete

import (
	"fmt"
	"github.com/aziontech/azion-cli/pkg/logger"
	"go.uber.org/zap/zapcore"
	"testing"

	msg "github.com/aziontech/azion-cli/messages/device_groups"
	"github.com/aziontech/azion-cli/pkg/httpmock"
	"github.com/aziontech/azion-cli/pkg/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreate(t *testing.T) {
	logger.New(zapcore.DebugLevel)
	t.Run("delete device group by id", func(t *testing.T) {
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

		assert.Equal(t, fmt.Sprintf(msg.DeviceGroupsDeleteOutputSuccess, 4321), stdout.String())
	})

	t.Run("try delete device group that is not found", func(t *testing.T) {
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
