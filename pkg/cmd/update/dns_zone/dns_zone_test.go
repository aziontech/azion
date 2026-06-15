package dnszone

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/aziontech/azion-cli/pkg/logger"
	"go.uber.org/zap/zapcore"

	msg "github.com/aziontech/azion-cli/messages/dns_zone"
	"github.com/aziontech/azion-cli/pkg/httpmock"
	"github.com/aziontech/azion-cli/pkg/testutils"
	"github.com/stretchr/testify/require"
)

func TestUpdate(t *testing.T) {
	logger.New(zapcore.DebugLevel)

	t.Run("update DNS zone with name and active", func(t *testing.T) {
		mock := &httpmock.Registry{}
		mock.Register(
			httpmock.REST("PUT", "workspace/dns/zones/1337"),
			httpmock.JSONFromFile("./fixtures/response.json"),
		)

		f, stdout, _ := testutils.NewFactory(mock)
		cmd := NewCmd(f)
		cmd.SetArgs([]string{"--zone-id", "1337", "--name", "renamed zone", "--active", "false"})

		err := cmd.Execute()

		require.NoError(t, err)
		require.Equal(t, fmt.Sprintf(msg.DNSZoneUpdateOutputSuccess, 1337), stdout.String())
	})

	t.Run("update preserves current values via get", func(t *testing.T) {
		mock := &httpmock.Registry{}
		// Only --name is provided, so the command fetches current values first.
		mock.Register(
			httpmock.REST("GET", "workspace/dns/zones/1337"),
			httpmock.JSONFromFile("./fixtures/response.json"),
		)
		mock.Register(
			httpmock.REST("PUT", "workspace/dns/zones/1337"),
			httpmock.JSONFromFile("./fixtures/response.json"),
		)

		f, stdout, _ := testutils.NewFactory(mock)
		cmd := NewCmd(f)
		cmd.SetArgs([]string{"--zone-id", "1337", "--name", "renamed zone"})

		err := cmd.Execute()

		require.NoError(t, err)
		require.Equal(t, fmt.Sprintf(msg.DNSZoneUpdateOutputSuccess, 1337), stdout.String())
	})

	t.Run("update with file", func(t *testing.T) {
		mock := &httpmock.Registry{}
		mock.Register(
			httpmock.REST("PUT", "workspace/dns/zones/1337"),
			httpmock.JSONFromFile("./fixtures/response.json"),
		)

		f, stdout, _ := testutils.NewFactory(mock)
		cmd := NewCmd(f)
		cmd.SetArgs([]string{"--zone-id", "1337", "--file", "./fixtures/update.json"})

		err := cmd.Execute()

		require.NoError(t, err)
		require.Equal(t, fmt.Sprintf(msg.DNSZoneUpdateOutputSuccess, 1337), stdout.String())
	})

	t.Run("not found", func(t *testing.T) {
		mock := &httpmock.Registry{}
		mock.Register(
			httpmock.REST("PUT", "workspace/dns/zones/1234"),
			httpmock.StatusStringResponse(http.StatusNotFound, "Not Found"),
		)

		f, _, _ := testutils.NewFactory(mock)
		cmd := NewCmd(f)
		cmd.SetArgs([]string{"--zone-id", "1234", "--name", "renamed zone", "--active", "false"})

		err := cmd.Execute()

		require.Error(t, err)
	})
}
