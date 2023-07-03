package describe

import (
	"github.com/aziontech/azion-cli/pkg/logger"
	"go.uber.org/zap/zapcore"
	"net/http"
	"os"
	"testing"

	"github.com/aziontech/azion-cli/pkg/httpmock"
	"github.com/aziontech/azion-cli/pkg/testutils"
	"github.com/stretchr/testify/require"
)

var successResponse string = `
{
	"results": {
	  "id": 666,
	  "name": "ssass",
	  "delivery_protocol": "http",
	  "http_port": [80],
	  "https_port": [443],
	  "minimum_tls_version": "",
	  "active": true,
	  "debug_rules": false,
	  "application_acceleration": false,
	  "caching": true,
	  "device_detection": false,
	  "edge_firewall": false,
	  "edge_functions": false,
	  "image_optimization": false,
	  "l2_caching": false,
	  "load_balancer": false,
	  "raw_logs": false,
	  "web_application_firewall": false
	},
	"schema_version": 3
  }
`

func TestDescribe(t *testing.T) {
	logger.New(zapcore.DebugLevel)
	t.Run("describe an edge application", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("GET", "edge_applications/666"),
			httpmock.JSONFromString(successResponse),
		)

		f, _, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)

		cmd.SetArgs([]string{"-a", "666"})

		err := cmd.Execute()
		require.NoError(t, err)

		require.NoError(t, err)

	})
	t.Run("not found", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("GET", "edge_applications/1234"),
			httpmock.StatusStringResponse(http.StatusNotFound, "Not Found"),
		)

		f, _, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)

		cmd.SetArgs([]string{"-a", "1234"})

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

	t.Run("export to a file", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("GET", "edge_applications/123"),
			httpmock.JSONFromString(successResponse),
		)

		f, stdout, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)

		path := "./out.json"

		cmd.SetArgs([]string{"--application-id", "123", "--out", path})

		err := cmd.Execute()
		if err != nil {
			t.Fatalf("error executing cmd")
		}

		_, err = os.ReadFile(path)
		if err != nil {
			t.Fatalf("error reading `out.json`: %v", err)
		}
		defer func() {
			_ = os.Remove(path)
		}()

		require.NoError(t, err)

		require.Equal(t, `File successfully written to: out.json
`, stdout.String())

	})

}
