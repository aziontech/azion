package rulesengine

import (
	"fmt"
	"testing"

	"github.com/aziontech/azion-cli/pkg/logger"
	"go.uber.org/zap/zapcore"

	msg "github.com/aziontech/azion-cli/messages/delete/rules_engine"
	"github.com/aziontech/azion-cli/pkg/httpmock"
	"github.com/aziontech/azion-cli/pkg/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreate(t *testing.T) {
	logger.New(zapcore.DebugLevel)
	t.Run("delete rules engine with success", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("DELETE", "edge_applications/4321/rules_engine/request/rules/1234"),
			httpmock.StatusStringResponse(204, ""),
		)

		f, stdout, _ := testutils.NewFactory(mock)
		cmd := NewCmd(f)
		cmd.SetArgs([]string{"--application-id", "4321", "--phase", "request", "--rule-id", "1234"})

		_, err := cmd.ExecuteC()
		require.NoError(t, err)
		assert.Equal(t, fmt.Sprintf(msg.DeleteOutputSuccess, 1234), stdout.String())
	})

	t.Run("delete rules engine that is not found", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("DELETE", "edge_applications/4321/rules_engine/request/rules/1234"),
			httpmock.StatusStringResponse(404, "Not Found"),
		)

		f, _, _ := testutils.NewFactory(mock)
		cmd := NewCmd(f)
		cmd.SetArgs([]string{"--rule-id", "1234", "--application-id", "4321", "--phase", "response"})

		_, err := cmd.ExecuteC()
		require.Error(t, err)
	})
}
