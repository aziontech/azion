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

func TestCreate(t *testing.T) {
	logger.New(zapcore.DebugLevel)

	t.Run("create new DNS record", func(t *testing.T) {
		t.Parallel()

		mock := &httpmock.Registry{}
		mock.Register(
			httpmock.REST("POST", "workspace/dns/zones/1337/records"),
			httpmock.JSONFromFile("./fixtures/response.json"),
		)

		f, stdout, _ := testutils.NewFactory(mock)
		cmd := NewCmd(f)
		cmd.SetArgs([]string{"--zone-id", "1337", "--name", "www", "--type", "A", "--rdata", "192.0.2.1", "--ttl", "3600"})

		err := cmd.Execute()

		require.NoError(t, err)
		require.Equal(t, fmt.Sprintf(msg.DNSRecordCreateOutputSuccess, 1111), stdout.String())
	})

	t.Run("create with multiple rdata values", func(t *testing.T) {
		t.Parallel()

		mock := &httpmock.Registry{}
		mock.Register(
			httpmock.REST("POST", "workspace/dns/zones/1337/records"),
			httpmock.JSONFromFile("./fixtures/response.json"),
		)

		f, stdout, _ := testutils.NewFactory(mock)
		cmd := NewCmd(f)
		cmd.SetArgs([]string{"--zone-id", "1337", "--name", "www", "--type", "A", "--rdata", "192.0.2.1,192.0.2.2"})

		err := cmd.Execute()

		require.NoError(t, err)
		require.Equal(t, fmt.Sprintf(msg.DNSRecordCreateOutputSuccess, 1111), stdout.String())
	})

	t.Run("create with file", func(t *testing.T) {
		t.Parallel()

		mock := &httpmock.Registry{}
		mock.Register(
			httpmock.REST("POST", "workspace/dns/zones/1337/records"),
			httpmock.JSONFromFile("./fixtures/response.json"),
		)

		f, stdout, _ := testutils.NewFactory(mock)
		cmd := NewCmd(f)
		cmd.SetArgs([]string{"--zone-id", "1337", "--file", "./fixtures/create.json"})

		err := cmd.Execute()

		require.NoError(t, err)
		require.Equal(t, fmt.Sprintf(msg.DNSRecordCreateOutputSuccess, 1111), stdout.String())
	})

	t.Run("bad request", func(t *testing.T) {
		t.Parallel()

		mock := &httpmock.Registry{}
		mock.Register(
			httpmock.REST("POST", "workspace/dns/zones/1337/records"),
			httpmock.StatusStringResponse(http.StatusBadRequest, "Invalid"),
		)

		f, _, _ := testutils.NewFactory(mock)
		cmd := NewCmd(f)
		cmd.SetArgs([]string{"--zone-id", "1337", "--name", "www", "--type", "A", "--rdata", "192.0.2.1"})

		err := cmd.Execute()

		require.Error(t, err)
	})

	t.Run("internal server error", func(t *testing.T) {
		t.Parallel()

		mock := &httpmock.Registry{}
		mock.Register(
			httpmock.REST("POST", "workspace/dns/zones/1337/records"),
			httpmock.StatusStringResponse(http.StatusInternalServerError, "Internal Server Error"),
		)

		f, _, _ := testutils.NewFactory(mock)
		cmd := NewCmd(f)
		cmd.SetArgs([]string{"--zone-id", "1337", "--name", "www", "--type", "A", "--rdata", "192.0.2.1"})

		err := cmd.Execute()

		require.Error(t, err)
	})
}
