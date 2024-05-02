package domains

import (
	"net/http"
	"testing"

	"github.com/aziontech/azion-cli/pkg/logger"
	"go.uber.org/zap/zapcore"

	"github.com/aziontech/azion-cli/pkg/httpmock"
	"github.com/aziontech/azion-cli/pkg/testutils"
	"github.com/stretchr/testify/require"
)

func TestDescribe(t *testing.T) {
	logger.New(zapcore.DebugLevel)
	t.Run("describe an domains", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("GET", "domains/1675272891"),
			httpmock.JSONFromFile("./fixtures/domain.json"),
		)

		f, _, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)

		cmd.SetArgs([]string{"--domain-id", "1675272891"})

		err := cmd.Execute()
		require.NoError(t, err)
	})
	t.Run("not found", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("GET", "domains/878"),
			httpmock.StatusStringResponse(http.StatusNotFound, "Not Found"),
		)

		f, _, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)
		cmd.SetArgs([]string{"--domain-id", "878"})

		err := cmd.Execute()
		require.Error(t, err)
	})

	t.Run("no id sent", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("GET", "edge_applications/1234"),
			httpmock.StatusStringResponse(http.StatusNotFound, "Not Found"),
		)

		f, _, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)

		err := cmd.Execute()

		require.Error(t, err)
	})
}
