package describe

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/aziontech/azion-cli/pkg/logger"
	"go.uber.org/zap/zapcore"

	msg "github.com/aziontech/azion-cli/messages/cache_settings"
	"github.com/aziontech/azion-cli/pkg/httpmock"
	"github.com/aziontech/azion-cli/pkg/testutils"
	"github.com/stretchr/testify/require"
)

func TestDescribe(t *testing.T) {
	logger.New(zapcore.DebugLevel)
	t.Run("describe an cache settings", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("GET", "edge_applications/1673635839/cache_settings/107313"),
			httpmock.JSONFromFile("./fixtures/cache_settings.json"),
		)

		f, _, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)
		cmd.SetArgs([]string{"-a", "1673635839", "-c", "107313"})

		err := cmd.Execute()
		require.NoError(t, err)
	})
	t.Run("not found", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("GET", "edge_applications/1673635839/cache_settings/107313"),
			httpmock.StatusStringResponse(http.StatusNotFound, "Not Found"),
		)

		f, _, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)

		err := cmd.Execute()
		require.Error(t, err)
	})

	t.Run("no id sent", func(t *testing.T) {
		mock := &httpmock.Registry{}
		mock.Register(
			httpmock.REST("GET", "edge_applications/123423424/cache_settings/122149"),
			httpmock.StatusStringResponse(http.StatusNotFound, "Not Found"),
		)

		f, _, _ := testutils.NewFactory(mock)
		cmd := NewCmd(f)
		err := cmd.Execute()
		require.ErrorIs(t, err, msg.ErrorMissingArguments)
	})

	t.Run("export to a file", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("GET", "edge_applications/1673635839/cache_settings/107313"),
			httpmock.JSONFromFile("./fixtures/cache_settings.json"),
		)

		f, stdout, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)
		path := "out.json"
		cmd.SetArgs([]string{"-a", "1673635839", "-c", "107313", "--out", path})

		err := cmd.Execute()
		if err != nil {
			log.Println("error executing cmd err: ", err.Error())
		}

		_, err = os.ReadFile(path)
		if err != nil {
			t.Fatalf("error reading `out.json`: %v", err)
		}
		defer func() {
			_ = os.Remove(path)
		}()

		require.NoError(t, err)

		require.Equal(t, fmt.Sprintf(msg.CacheSettingsFileWritten, path), stdout.String())
	})
}
