package dnsrecord

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/aziontech/azion-cli/pkg/logger"
	"go.uber.org/zap/zapcore"

	msg "github.com/aziontech/azion-cli/messages/dns_record"
	"github.com/aziontech/azion-cli/pkg/httpmock"
	"github.com/aziontech/azion-cli/pkg/testutils"
	"github.com/stretchr/testify/require"
)

func TestUpdate(t *testing.T) {
	logger.New(zapcore.DebugLevel)

	t.Run("update DNS record preserving current values", func(t *testing.T) {
		mock := &httpmock.Registry{}
		// Only --rdata is provided, so the command fetches current values first.
		mock.Register(
			httpmock.REST("GET", "workspace/dns/zones/1337/records/1111"),
			httpmock.JSONFromFile("./fixtures/response.json"),
		)
		mock.Register(
			httpmock.REST("PATCH", "workspace/dns/zones/1337/records/1111"),
			httpmock.JSONFromFile("./fixtures/response.json"),
		)

		f, stdout, _ := testutils.NewFactory(mock)
		cmd := NewCmd(f)
		cmd.SetArgs([]string{"--zone-id", "1337", "--record-id", "1111", "--rdata", "192.0.2.3"})

		err := cmd.Execute()

		require.NoError(t, err)
		require.Equal(t, fmt.Sprintf(msg.DNSRecordUpdateOutputSuccess, 1111), stdout.String())
	})

	t.Run("update with file", func(t *testing.T) {
		mock := &httpmock.Registry{}
		mock.Register(
			httpmock.REST("PATCH", "workspace/dns/zones/1337/records/1111"),
			httpmock.JSONFromFile("./fixtures/response.json"),
		)

		f, stdout, _ := testutils.NewFactory(mock)
		cmd := NewCmd(f)
		cmd.SetArgs([]string{"--zone-id", "1337", "--record-id", "1111", "--file", "./fixtures/update.json"})

		err := cmd.Execute()

		require.NoError(t, err)
		require.Equal(t, fmt.Sprintf(msg.DNSRecordUpdateOutputSuccess, 1111), stdout.String())
	})

	t.Run("not found", func(t *testing.T) {
		mock := &httpmock.Registry{}
		mock.Register(
			httpmock.REST("GET", "workspace/dns/zones/1337/records/1234"),
			httpmock.StatusStringResponse(http.StatusNotFound, "Not Found"),
		)

		f, _, _ := testutils.NewFactory(mock)
		cmd := NewCmd(f)
		cmd.SetArgs([]string{"--zone-id", "1337", "--record-id", "1234", "--rdata", "192.0.2.3"})

		err := cmd.Execute()

		require.Error(t, err)
	})
}
