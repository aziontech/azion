package ruleengine

import (
	"testing"

	"github.com/aziontech/azion-cli/pkg/logger"
	"go.uber.org/zap/zapcore"

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

		cmd.SetArgs([]string{"--application-id", "1678743802", "--phase", "request"})

		_, err := cmd.ExecuteC()
		require.NoError(t, err)
	})

	t.Run("list all rules engines - ask for input", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("GET", "edge_applications/1678743802/rules_engine/request/rules"),
			httpmock.JSONFromFile("./fixtures/rules.json"),
		)

		f, _, _ := testutils.NewFactory(mock)
		listCmd := NewListCmd(f)
		listCmd.AskInput = func(s string) (string, error) {
			return "1678743802", nil
		}
		cmd := NewCobraCmd(listCmd, f)

		cmd.SetArgs([]string{"--phase", "request"})

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

		cmd.SetArgs([]string{"--application-id", "1678743802", "--phase", "request"})

		_, err := cmd.ExecuteC()
		require.NoError(t, err)
	})
}
