package deploy

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"testing"

	"github.com/aziontech/azion-cli/pkg/logger"
	"go.uber.org/zap/zapcore"

	apiapp "github.com/aziontech/azion-cli/pkg/api/edge_applications"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/httpmock"
	"github.com/aziontech/azion-cli/pkg/testutils"
	"github.com/stretchr/testify/require"
)

var successResponseInstance string = `
{
	"results": {
	  "edge_function_id": 123321,
	  "name": "varejo",
	  "args": {},
	  "id": 101001
	},
	"schema_version": 3
  }
  `

var successResponseApp string = `
{
	"results": {
	  "id": 666,
	  "name": "New Edge Application",
	  "delivery_protocol": "http,https",
	  "http_port": [
		80
	  ],
	  "https_port": [
		443
	  ],
	  "minimum_tls_version": "",
	  "active": true,
	  "debug_rules": false,
	  "http3": false,
	  "websocket": null,
	  "supported_ciphers": "all",
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

func TestDeployCmd(t *testing.T) {
	logger.New(zapcore.DebugLevel)
	t.Run("without azion.json", func(t *testing.T) {
		f, _, _ := testutils.NewFactory(nil)

		deployCmd := NewDeployCmd(f)

		deployCmd.FileReader = func(path string) ([]byte, error) {
			return nil, os.ErrNotExist
		}

		cmd := NewCobraCmd(deployCmd)

		cmd.SetArgs([]string{""})

		err := cmd.Execute()

		require.EqualError(t, err, "Failed to build your resource. Azion configuration not found. Make sure you are in the root directory of your local repository and have already initialized or linked your resource with the commands 'azion init' or 'azion link'")
	})

	t.Run("failed to create application", func(t *testing.T) {

		mock := &httpmock.Registry{}
		options := &contracts.AzionApplicationOptions{
			Name: "NotAVeryGoodName",
		}

		dat, _ := os.ReadFile("./fixtures/create_app.json")
		_ = json.Unmarshal(dat, options)

		mock.Register(
			httpmock.REST("POST", "edge_applications"),
			httpmock.StatusStringResponse(http.StatusBadRequest, "Invalid"),
		)

		f, _, _ := testutils.NewFactory(mock)
		ctx := context.Background()

		cliapp := apiapp.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))

		cmd := NewDeployCmd(f)

		_, _, err := cmd.createApplication(cliapp, ctx, options)
		require.ErrorContains(t, err, "Failed to create the Edge Application")
	})

	t.Run("create application success", func(t *testing.T) {

		mock := &httpmock.Registry{}
		options := &contracts.AzionApplicationOptions{
			Name: "LovelyName",
		}

		dat, _ := os.ReadFile("./fixtures/create_app.json")
		_ = json.Unmarshal(dat, options)

		mock.Register(
			httpmock.REST("POST", "edge_applications"),
			httpmock.JSONFromString(successResponseApp),
		)

		mock.Register(
			httpmock.REST("PATCH", "edge_applications/666"),
			httpmock.JSONFromString(successResponseApp),
		)

		mock.Register(
			httpmock.REST("POST", "edge_applications/666/functions_instances"),
			httpmock.JSONFromString(successResponseInstance),
		)

		f, _, _ := testutils.NewFactory(mock)
		ctx := context.Background()

		cliapp := apiapp.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))

		cmd := NewDeployCmd(f)

		_, _, err := cmd.createApplication(cliapp, ctx, options)
		require.NoError(t, err)
	})
}
