package cachesetting

import (
	"strconv"
	"testing"

	"github.com/aziontech/azion-cli/pkg/logger"
	"go.uber.org/zap/zapcore"

	msg "github.com/aziontech/azion-cli/messages/cache_setting"
	"github.com/aziontech/azion-cli/pkg/httpmock"
	"github.com/aziontech/azion-cli/pkg/testutils"
	"github.com/stretchr/testify/require"
)

func TestList(t *testing.T) {
	logger.New(zapcore.DebugLevel)

	t.Run("list default", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("GET", "edge_applications/1673635839/cache_settings"),
			httpmock.JSONFromFile("./fixtures/caches.json"),
		)

		f, _, _ := testutils.NewFactory(mock)
		cmd := NewCmd(f)

		cmd.SetArgs([]string{"--application-id", "1673635839"})

		_, err := cmd.ExecuteC()
		require.NoError(t, err)
	})

	t.Run("list page 1 with item 1", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("GET", "edge_applications/1673635839/cache_settings"),
			httpmock.JSONFromFile("./fixtures/list_3_itens.json"),
		)

		f, _, _ := testutils.NewFactory(mock)
		cmd := NewCmd(f)

		cmd.SetArgs([]string{"--application-id", "1673635839", "--page", "1"})

		_, err := cmd.ExecuteC()
		require.NoError(t, err)
	})

	t.Run("list page 3 with item 1", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("GET", "edge_applications/1673635839/cache_settings"),
			httpmock.JSONFromFile("./fixtures/list_3_itens.json"),
		)

		f, _, _ := testutils.NewFactory(mock)
		cmd := NewCmd(f)

		cmd.SetArgs([]string{"--application-id", "1673635839", "--page", "3"})

		_, err := cmd.ExecuteC()
		require.NoError(t, err)
	})

	t.Run("no items", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("GET", "edge_applications/1673635839/cache_settings"),
			httpmock.JSONFromFile("./fixtures/nocaches.json"),
		)

		f, _, _ := testutils.NewFactory(mock)
		cmd := NewCmd(f)

		cmd.SetArgs([]string{"--application-id", "1673635839"})

		_, err := cmd.ExecuteC()
		require.NoError(t, err)
	})

	t.Run("list with AskInput function", func(t *testing.T) {
		mock := &httpmock.Registry{}
		expectedAppID := int64(1673635839)

		// Create a mock for AskInput
		mockAskInput := func(prompt string) (string, error) {
			return strconv.FormatInt(expectedAppID, 10), nil
		}

		mock.Register(
			httpmock.REST("GET", "edge_applications/1673635839/cache_settings"),
			httpmock.JSONFromFile("./fixtures/caches.json"),
		)

		f, _, _ := testutils.NewFactory(mock)
		listCmd := NewListCmd(f)
		listCmd.AskInput = mockAskInput // Use the mocked AskInput function

		cmd := NewCobraCmd(listCmd, f)

		// Simulate the --application-id flag not being set
		cmd.SetArgs([]string{})

		_, err := cmd.ExecuteC()
		require.NoError(t, err, msg.ErrorGetCache)
	})
}
