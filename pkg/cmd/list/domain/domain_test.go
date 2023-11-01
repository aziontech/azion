package domain

import (
	"testing"

	"github.com/aziontech/azion-cli/pkg/httpmock"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/testutils"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
)

func TestNewCmd(t *testing.T) {
	logger.New(zapcore.DebugLevel)
	t.Run("list url not found", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("GET", "domains"),
			httpmock.StatusStringResponse(404, "Not Found"),
		)

		f, _, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)

		cmd.SetArgs([]string{})

		_, err := cmd.ExecuteC()
		require.Error(t, err)
	})
	t.Run("list all domains", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("GET", "domains"),
			httpmock.JSONFromFile("fixtures/domains.json"),
		)

		f, _, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)

		cmd.SetArgs([]string{})

		_, err := cmd.ExecuteC()
		require.NoError(t, err)
	})
	t.Run("empty list", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("GET", "domains"),
			httpmock.JSONFromFile("fixtures/nodomains.json"),
		)

		f, _, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)

		cmd.SetArgs([]string{})

		_, err := cmd.ExecuteC()
		require.NoError(t, err)
	})
}
