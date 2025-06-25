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

const successResponseApp = `
{
  "state": "pending",
  "data": {
    "id": 1697666970,
    "name": "LovelyName",
    "last_editor": "tester",
    "last_modified": "2025-06-20T16:55:19Z",
    "modules": {
      "edge_cache_enabled": true,
      "edge_functions_enabled": false,
      "application_accelerator_enabled": false,
      "image_processor_enabled": false,
      "tiered_cache_enabled": false
    },
    "active": true,
    "debug": true,
    "product_version": "1.0.0"
  }
}`

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
			httpmock.REST("POST", "edge_application/applications"),
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
			httpmock.REST("POST", "edge_application/applications"),
			httpmock.JSONFromString(successResponseApp),
		)

		mock.Register(
			httpmock.REST("PATCH", "edge_application/applications/1697666970"),
			httpmock.JSONFromString(successResponseApp),
		)

		mock.Register(
			httpmock.REST("POST", "edge_application/applications/1697666970/functions"),
			httpmock.JSONFromString(successResponseApp),
		)

		f, _, _ := testutils.NewFactory(mock)
		ctx := context.Background()

		cliapp := apiapp.NewClient(f.HttpClient, f.Config.GetString("api_v4_url"), f.Config.GetString("token"))

		cmd := NewDeployCmd(f)

		_, err := cmd.createApplication(cliapp, ctx, options, &msgs)
		require.NoError(t, err)
	})
}
