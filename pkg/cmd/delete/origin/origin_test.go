package origin

import (
	"fmt"
	"testing"

	"github.com/aziontech/azion-cli/pkg/logger"
	"go.uber.org/zap/zapcore"

	msg "github.com/aziontech/azion-cli/messages/origin"
	"github.com/aziontech/azion-cli/pkg/httpmock"
	"github.com/aziontech/azion-cli/pkg/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreate(t *testing.T) {
	logger.New(zapcore.DebugLevel)
	t.Run("delete origin by key", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("DELETE", "edge_applications/1673635839/origins/58755fef-e830-4ea4-b9e0-6481f1ef496d"),
			httpmock.StatusStringResponse(204, ""),
		)

		f, stdout, _ := testutils.NewFactory(mock)
		cmd := NewCmd(f)
		cmd.SetArgs([]string{"--application-id", "1673635839", "--origin-key", "58755fef-e830-4ea4-b9e0-6481f1ef496d"})

		_, err := cmd.ExecuteC()
		require.NoError(t, err)

		assert.Equal(t, fmt.Sprintf(msg.DeleteOutputSuccess, "58755fef-e830-4ea4-b9e0-6481f1ef496d"), stdout.String())
	})

	t.Run("delete domain - not found", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("DELETE", "edge_applications/1673635839/origins/58755fef-e830-4ea4-b9e0-6481f1ef496d"),
			httpmock.StatusStringResponse(404, "Not Found"),
		)

		f, _, _ := testutils.NewFactory(mock)
		cmd := NewCmd(f)
		cmd.SetArgs([]string{"--application-id", "1673635839", "--origin-key", "58755fef-e830-4ea4-b9e0-6481f1ef496d"})

		_, err := cmd.ExecuteC()
		require.Error(t, err)
	})
}
