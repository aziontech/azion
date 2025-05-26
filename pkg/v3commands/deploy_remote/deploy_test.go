package deploy

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"testing"

	"github.com/aziontech/azion-cli/pkg/logger"
	"go.uber.org/zap/zapcore"

	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/httpmock"
	"github.com/aziontech/azion-cli/pkg/testutils"
	apiapp "github.com/aziontech/azion-cli/pkg/v3api/edge_applications"
	"github.com/stretchr/testify/require"
)

var successResponseApp string = `
{
	"results":{
		"id":1697666970,
		"name":"New Edge Applicahvjgjhgjhhgtion",
		"delivery_protocol":"http",
		"http_port":80,
		"https_port":443,
		"minimum_tls_version":"",
		"active":true,
		"application_acceleration":false,
		"caching":true,
   		"debug_rules": true,
   		"http3": false,
		"supported_ciphers": "asdf",
		"device_detection":false,
		"edge_firewall":false,
		"edge_functions":false,
		"image_optimization":false,
		"load_balancer":false,
		"raw_logs":false,
		"web_application_firewall":false,
		"l2_caching": false
	},
	"schema_version":3
}
`

func TestDeployCmd(t *testing.T) {
	logger.New(zapcore.DebugLevel)

	msgs := []string{}

	t.Run("without azion.json", func(t *testing.T) {
		f, _, _ := testutils.NewFactory(nil)

		deployCmd := NewDeployCmd(f)

		deployCmd.FileReader = func(path string) ([]byte, error) {
			return nil, os.ErrNotExist
		}

		cmd := NewCobraCmd(deployCmd)

		cmd.SetArgs([]string{""})

		err := cmd.Execute()

		require.Error(t, err)
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

		_, err := cmd.createApplication(cliapp, ctx, options, &msgs)
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
			httpmock.REST("PATCH", "edge_applications/1697666970"),
			httpmock.JSONFromString(successResponseApp),
		)

		mock.Register(
			httpmock.REST("POST", "edge_applications/1697666970/functions_instances"),
			httpmock.JSONFromString(successResponseApp),
		)

		f, _, _ := testutils.NewFactory(mock)
		ctx := context.Background()

		cliapp := apiapp.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))

		cmd := NewDeployCmd(f)

		_, err := cmd.createApplication(cliapp, ctx, options, &msgs)
		require.NoError(t, err)
	})
}
