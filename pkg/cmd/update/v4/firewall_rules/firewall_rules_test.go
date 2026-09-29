package firewallrules

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/aziontech/azion-cli/pkg/logger"
	"go.uber.org/zap/zapcore"

	msg "github.com/aziontech/azion-cli/messages/update/firewall_rules"
	"github.com/aziontech/azion-cli/pkg/httpmock"
	"github.com/aziontech/azion-cli/pkg/testutils"
	"github.com/stretchr/testify/require"
)

func TestUpdate(t *testing.T) {
	logger.New(zapcore.DebugLevel)

	t.Run("update Firewall Rule", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("PATCH", "workspace/firewalls/1234/request_rules/111"),
			httpmock.JSONFromFile("./fixtures/response.json"),
		)

		f, stdout, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)

		cmd.SetArgs([]string{"--firewall-id", "1234", "--rule-id", "111", "--file", "./fixtures/update.json"})

		err := cmd.Execute()

		require.NoError(t, err)
		require.Equal(t, fmt.Sprintf(msg.OutputSuccess, 111), stdout.String())
	})

	t.Run("bad request", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("PATCH", "workspace/firewalls/1234/request_rules/111"),
			httpmock.StatusStringResponse(http.StatusBadRequest, "Invalid"),
		)

		f, _, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)

		cmd.SetArgs([]string{"--firewall-id", "1234", "--rule-id", "111", "--file", "./fixtures/update.json"})

		err := cmd.Execute()

		require.Error(t, err)
	})

	t.Run("internal server error", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("PATCH", "workspace/firewalls/1234/request_rules/111"),
			httpmock.StatusStringResponse(http.StatusInternalServerError, "Internal Server Error"),
		)

		f, _, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)

		cmd.SetArgs([]string{"--firewall-id", "1234", "--rule-id", "111", "--file", "./fixtures/update.json"})

		err := cmd.Execute()
		require.Error(t, err)
	})
}
