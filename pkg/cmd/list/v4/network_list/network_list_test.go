package networklist

import (
	"testing"

	"github.com/aziontech/azion-cli/pkg/logger"
	"go.uber.org/zap/zapcore"

	"github.com/aziontech/azion-cli/pkg/httpmock"
	"github.com/aziontech/azion-cli/pkg/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	tblWithNetworkList string = "ID  NAME    ACTIVE  \n0   string  true    \n"
	tblNoNetworkList   string = "ID  NAME    ACTIVE  \n"
)

func TestList(t *testing.T) {
	logger.New(zapcore.DebugLevel)
	t.Run("more than one network list", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("GET", "workspace/network_lists"),
			httpmock.JSONFromFile("./fixtures/networklists.json"),
		)

		f, stdout, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)

		cmd.SetArgs([]string{})

		_, err := cmd.ExecuteC()
		require.NoError(t, err)
		assert.Equal(t, tblWithNetworkList, stdout.String())
	})

	t.Run("list - page 0 should generate an error", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("GET", "workspace/network_lists"),
			httpmock.StatusStringResponse(404, "Not Found"),
		)

		f, _, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)

		cmd.SetArgs([]string{"--page", "0"})

		_, err := cmd.ExecuteC()
		require.Error(t, err)
	})

	t.Run("no network lists", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("GET", "workspace/network_lists"),
			httpmock.JSONFromFile("./fixtures/nonetworklist.json"),
		)

		f, stdout, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)

		cmd.SetArgs([]string{})

		_, err := cmd.ExecuteC()
		require.NoError(t, err)

		assert.Equal(t, tblNoNetworkList, stdout.String())
	})
}
