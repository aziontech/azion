package origin

import (
	"testing"

	"github.com/aziontech/azion-cli/pkg/logger"
	"go.uber.org/zap/zapcore"

	"github.com/aziontech/azion-cli/pkg/httpmock"
	"github.com/aziontech/azion-cli/pkg/testutils"
	"github.com/stretchr/testify/require"
)

func TestList(t *testing.T) {
	logger.New(zapcore.DebugLevel)
	t.Run("list page 1 with flag", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("GET", "edge_applications/123423424/origins"),
			httpmock.JSONFromFile("./fixtures/origins.json"),
		)

		f, _, _ := testutils.NewFactory(mock)
		cmd := NewCmd(f)

		cmd.SetArgs([]string{"--application-id", "123423424"})

		_, err := cmd.ExecuteC()
		require.NoError(t, err)
	})

	t.Run("list page 1 with AskForInput", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("GET", "edge_applications/123423424/origins"),
			httpmock.JSONFromFile("./fixtures/origins.json"),
		)

		f, _, _ := testutils.NewFactory(mock)

		listCmd := NewListCmd(f)
		listCmd.AskInput = func(s string) (string, error) {
			return "123423424", nil
		}
		cmd := NewCobraCmd(listCmd, f)
		cmd.SetArgs([]string{})

		_, err := cmd.ExecuteC()
		require.NoError(t, err)
	})

	t.Run("list - page 0 should generate an error", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("GET", "edge_applications/123423424/origins"),
			httpmock.StatusStringResponse(404, "Not Found"),
		)

		f, _, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)

		cmd.SetArgs([]string{"--application-id", "1673635839", "--page", "0"})

		_, err := cmd.ExecuteC()
		require.Error(t, err)
	})

	t.Run("no items with flag", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("GET", "edge_applications/123423424/origins"),
			httpmock.JSONFromFile("./fixtures/noorigins.json"),
		)

		f, _, _ := testutils.NewFactory(mock)
		cmd := NewCmd(f)

		cmd.SetArgs([]string{"--application-id", "123423424"})

		_, err := cmd.ExecuteC()
		require.NoError(t, err)
	})

	t.Run("no items with AskForInput", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("GET", "edge_applications/123423424/origins"),
			httpmock.JSONFromFile("./fixtures/noorigins.json"),
		)

		f, _, _ := testutils.NewFactory(mock)

		listCmd := NewListCmd(f)
		listCmd.AskInput = func(s string) (string, error) {
			return "123423424", nil
		}
		cmd := NewCobraCmd(listCmd, f)
		cmd.SetArgs([]string{})

		_, err := cmd.ExecuteC()
		require.NoError(t, err)
	})
}
