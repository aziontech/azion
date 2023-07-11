package list

import (
	"github.com/aziontech/azion-cli/pkg/logger"
	"go.uber.org/zap/zapcore"
	"testing"

	msg "github.com/aziontech/azion-cli/messages/rules_engine"
	"github.com/aziontech/azion-cli/pkg/httpmock"
	"github.com/aziontech/azion-cli/pkg/testutils"
	"github.com/stretchr/testify/require"
)

func TestNewCmd(t *testing.T) {
	logger.New(zapcore.DebugLevel)
	t.Run("list all rules engines", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("GET", "edge_applications/1678743802/rules_engine/request/rules"),
			httpmock.JSONFromFile("./fixtures/rules.json"),
		)

		f, _, _ := testutils.NewFactory(mock)
		cmd := NewCmd(f)

		cmd.SetArgs([]string{"-a", "1678743802", "-p", "request"})

		_, err := cmd.ExecuteC()
		require.NoError(t, err)
	})

	t.Run("no itens", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("GET", "edge_applications/1678743802/rules_engine/request/rules"),
			httpmock.JSONFromFile("./fixtures/norules.json"),
		)

		f, _, _ := testutils.NewFactory(mock)
		cmd := NewCmd(f)

		cmd.SetArgs([]string{"-a", "1678743802", "-p", "request"})

		_, err := cmd.ExecuteC()
		require.NoError(t, err)
	})

	t.Run("missing mandatory flags", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("GET", "edge_applications/1678743802/rules_engine/request/rules"),
			httpmock.JSONFromFile("./fixtures/norules.json"),
		)

		f, _, _ := testutils.NewFactory(mock)
		cmd := NewCmd(f)

		cmd.SetArgs([]string{"-a", "1678743802"})

		_, err := cmd.ExecuteC()
		require.ErrorIs(t, err, msg.ErrorMandatoryListFlags)
	})
}
