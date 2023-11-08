package cachesetting

import (
	"fmt"
	"testing"

	"github.com/aziontech/azion-cli/pkg/logger"
	"go.uber.org/zap/zapcore"

	msg "github.com/aziontech/azion-cli/messages/cache_setting"
	"github.com/aziontech/azion-cli/pkg/httpmock"
	"github.com/aziontech/azion-cli/pkg/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreate(t *testing.T) {
	logger.New(zapcore.DebugLevel)
	t.Run("delete by id", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("DELETE", "edge_applications/1673635839/cache_settings/107313"),
			httpmock.StatusStringResponse(204, ""),
		)

		f, stdout, _ := testutils.NewFactory(mock)
		cmd := NewCmd(f)
		cmd.SetArgs([]string{"--application-id", "1673635839", "--cache-settings-id", "107313"})

		_, err := cmd.ExecuteC()
		require.NoError(t, err)
		assert.Equal(t, fmt.Sprintf(msg.DeleteOutputSuccess, 107313), stdout.String())
	})

	t.Run("delete that is not found", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("DELETE", "edge_applications/1234/cache_settings/107313"),
			httpmock.StatusStringResponse(404, "Not Found"),
		)

		f, _, _ := testutils.NewFactory(mock)
		cmd := NewCmd(f)
		cmd.SetArgs([]string{"--application-id", "1234", "--cache-settings-id", "107313"})

		_, err := cmd.ExecuteC()
		require.Error(t, err)
	})
}
