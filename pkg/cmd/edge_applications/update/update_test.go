package update

import (
	"net/http"
	"testing"

	msg "github.com/aziontech/azion-cli/messages/edge_applications"
	"github.com/aziontech/azion-cli/pkg/httpmock"
	"github.com/aziontech/azion-cli/pkg/testutils"
	"github.com/stretchr/testify/require"
)

var successResponse string = `
{
	"results": {
	  "id": 1337,
	  "name": "Update Edge Application",
	  "delivery_protocol": "http,https",
	  "http_port": 80,
	  "https_port": 443,
	  "minimum_tls_version": "",
	  "debug_rules": false,
	  "application_acceleration": false,
	  "caching": true,
	  "device_detection": false,
	  "edge_firewall": false,
	  "edge_functions": false,
	  "image_optimization": true,
	  "load_balancer": false,
	  "raw_logs": false,
	  "web_application_firewall": false
	},
	"schema_version": 3
  }
`

func TestUpdate(t *testing.T) {
	t.Run("update edge application", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("PATCH", "edge_applications/1337"),
			httpmock.JSONFromString(successResponse),
		)

		f, stdout, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)

		cmd.SetArgs([]string{"-a", "1337", "--name", "ATUALIZANDO"})

		err := cmd.Execute()

		require.NoError(t, err)
		require.Equal(t, "Updated Edge Application with ID 1337\n", stdout.String())
	})

	t.Run("bad request", func(t *testing.T) {
		mock := &httpmock.Registry{}
		mock.Register(
			httpmock.REST("PATCH", "edge_applications/1337"),
			httpmock.StatusStringResponse(http.StatusBadRequest, `{"details": "invalid field active"}`),
		)

		f, _, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)

		cmd.SetArgs([]string{"-a", "1234", "--active", "unactive"})

		err := cmd.Execute()

		require.Error(t, err)
	})

	t.Run("update with file", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("PATCH", "edge_applications/1337"),
			httpmock.JSONFromString(successResponse),
		)

		f, stdout, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)

		cmd.SetArgs([]string{"--in", "./fixtures/update.json"})

		err := cmd.Execute()

		require.NoError(t, err)
		require.Equal(t, "Updated Edge Application with ID 1337\n", stdout.String())
	})

	t.Run("return some fields", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("PATCH", "edge_applications/1337"),
			httpmock.JSONFromString(successResponse),
		)

		f, _, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)
		cmd.SetArgs([]string{"-a", "1337"})
		err := cmd.Execute()
		require.ErrorContains(t, err, msg.ErrorNoFieldInformed.Error(), nil)
	})
}
