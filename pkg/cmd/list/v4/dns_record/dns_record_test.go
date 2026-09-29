package dnsrecord

import (
	"testing"

	"github.com/aziontech/azion-cli/pkg/logger"
	"go.uber.org/zap/zapcore"

	"github.com/aziontech/azion-cli/pkg/httpmock"
	"github.com/aziontech/azion-cli/pkg/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestList(t *testing.T) {
	logger.New(zapcore.DebugLevel)

	t.Run("list DNS records", func(t *testing.T) {
		mock := &httpmock.Registry{}
		mock.Register(
			httpmock.REST("GET", "workspace/dns/zones/1337/records"),
			httpmock.JSONFromFile("./fixtures/response.json"),
		)

		f, stdout, _ := testutils.NewFactory(mock)
		cmd := NewCmd(f)
		cmd.SetArgs([]string{"--zone-id", "1337"})

		_, err := cmd.ExecuteC()
		require.NoError(t, err)
		assert.Contains(t, stdout.String(), "ID")
		assert.Contains(t, stdout.String(), "1111")
		assert.Contains(t, stdout.String(), "www")
		assert.Contains(t, stdout.String(), "192.0.2.1")
	})

	t.Run("no DNS records", func(t *testing.T) {
		mock := &httpmock.Registry{}
		mock.Register(
			httpmock.REST("GET", "workspace/dns/zones/1337/records"),
			httpmock.JSONFromFile("./fixtures/norecord.json"),
		)

		f, stdout, _ := testutils.NewFactory(mock)
		cmd := NewCmd(f)
		cmd.SetArgs([]string{"--zone-id", "1337"})

		_, err := cmd.ExecuteC()
		require.NoError(t, err)
		assert.NotContains(t, stdout.String(), "192.0.2.1")
	})

	t.Run("list with details", func(t *testing.T) {
		mock := &httpmock.Registry{}
		mock.Register(
			httpmock.REST("GET", "workspace/dns/zones/1337/records"),
			httpmock.JSONFromFile("./fixtures/response.json"),
		)

		f, _, _ := testutils.NewFactory(mock)
		cmd := NewCmd(f)
		cmd.SetArgs([]string{"--zone-id", "1337", "--details"})

		_, err := cmd.ExecuteC()
		require.NoError(t, err)
	})

	t.Run("not found", func(t *testing.T) {
		mock := &httpmock.Registry{}
		mock.Register(
			httpmock.REST("GET", "workspace/dns/zones/1337/records"),
			httpmock.StatusStringResponse(404, "Not Found"),
		)

		f, _, _ := testutils.NewFactory(mock)
		cmd := NewCmd(f)
		cmd.SetArgs([]string{"--zone-id", "1337"})

		_, err := cmd.ExecuteC()
		require.Error(t, err)
	})
}
