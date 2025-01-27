package domain

import (
	"testing"

	"github.com/aziontech/azion-cli/pkg/httpmock"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/testutils"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
	"gotest.tools/golden"
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

		f, stdout, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)

		cmd.SetArgs([]string{})

		_, err := cmd.ExecuteC()
		require.NoError(t, err)
		golden.AssertBytes(t, stdout.Bytes(), "out_list_all.golden")
	})

	t.Run("empty list", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("GET", "domains"),
			httpmock.JSONFromFile("fixtures/nodomains.json"),
		)

		f, stdout, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)

		cmd.SetArgs([]string{})

		_, err := cmd.ExecuteC()
		require.NoError(t, err)
		golden.AssertBytes(t, stdout.Bytes(), "out_empty_list.golden")
	})
}
