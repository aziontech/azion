package cachesetting

import (
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
	t.Run("list defatul", func(t *testing.T) {
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

	t.Run("list page 1 with iten 1", func(t *testing.T) {
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

	t.Run("list page 3 with iten 1", func(t *testing.T) {
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

	t.Run("no itens", func(t *testing.T) {
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

	t.Run("list page 1 with iten 1", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("GET", "edge_applications/1673635839/cache_settings"),
			httpmock.JSONFromFile("./fixtures/caches.json"),
		)

		f, _, _ := testutils.NewFactory(mock)
		cmd := NewCmd(f)

		cmd.SetArgs([]string{"--application-id", "1673635839"})

		_, err := cmd.ExecuteC()
		require.NoError(t, err, msg.ErrorGetCache)
	})
}
