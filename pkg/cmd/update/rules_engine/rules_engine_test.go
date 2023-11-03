package rules_engine

import (
	msg "github.com/aziontech/azion-cli/pkg/messages/update/rules_engine"
	"net/http"
	"testing"

	"github.com/aziontech/azion-cli/pkg/logger"
	"go.uber.org/zap/zapcore"

	"github.com/aziontech/azion-cli/pkg/httpmock"
	"github.com/aziontech/azion-cli/pkg/testutils"
	"github.com/stretchr/testify/require"
)

func TestUpdate(t *testing.T) {
	logger.New(zapcore.DebugLevel)
	t.Run("update ", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST(http.MethodPatch, "edge_applications/1673635839/rules_engine/request/rules/1234"),
			httpmock.JSONFromFile("./fixtures/response.json"),
		)

		f, _, _ := testutils.NewFactory(mock)
		cmd := NewCmd(f)

		cmd.SetArgs([]string{
			"--application-id", "1673635839",
			"--phase", "request",
			"--rule-id", "1234",
			"--in", "./fixtures/update.json",
		})

		err := cmd.Execute()
		require.NoError(t, err)
	})

	t.Run("missing fields", func(t *testing.T) {
		mock := &httpmock.Registry{}
		mock.Register(
			httpmock.REST(http.MethodPatch, "edge_applications/1673635839/rules_engine/request/rules/1234"),
			httpmock.StatusStringResponse(http.StatusBadRequest, `{}`),
		)

		f, _, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)

		cmd.SetArgs([]string{
			"--application-id", "1673635839",
			"--rule-id", "1234",
			"--phase", "request",
			"--in", "./fixtures/missing.json",
		})

		err := cmd.Execute()
		require.ErrorIs(t, err, msg.ErrorVariableEmpty)
	})
}
