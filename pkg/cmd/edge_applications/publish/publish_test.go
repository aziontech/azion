package publish

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"reflect"
	"testing"

	msg "github.com/aziontech/azion-cli/messages/edge_applications"
	apiapp "github.com/aziontech/azion-cli/pkg/api/edge_applications"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/httpmock"
	"github.com/aziontech/azion-cli/pkg/testutils"
	"github.com/stretchr/testify/require"
)

var applicationName string = "Brazilian forest Traitor"

var successResponseApp string = `
{
	"results":{
	   "id":666,
	   "name":"New Edge Applicahvjgjhgjhhgtion",
	   "delivery_protocol":"http",
	   "http_port":80,
	   "https_port":443,
	   "minimum_tls_version":"",
	   "active":true,
	   "application_acceleration":false,
	   "caching":true,
	   "device_detection":false,
	   "edge_firewall":false,
	   "edge_functions":false,
	   "image_optimization":false,
	   "load_balancer":false,
	   "raw_logs":false,
	   "web_application_firewall":false
	},
	"schema_version":3
}
`

func TestPublishCmd(t *testing.T) {
	t.Run("without azion.json", func(t *testing.T) {
		f, _, _ := testutils.NewFactory(nil)

		publishCmd := NewPublishCmd(f)

		publishCmd.FileReader = func(path string) ([]byte, error) {
			return nil, os.ErrNotExist
		}

		cmd := NewCobraCmd(publishCmd)

		cmd.SetArgs([]string{""})

		err := cmd.Execute()

		require.EqualError(t, err, "Failed to open the azion.json file. The file doesn't exist, is corrupted, or has an invalid JSON format. Verify if the file format is JSON or fix its content according to the JSON format specification at https://www.json.org/json-en.html")
	})

	t.Run("without config.json", func(t *testing.T) {
		f, _, _ := testutils.NewFactory(nil)

		cmd := NewPublishCmd(f)
		cmd.FileReader = func(path string) ([]byte, error) {
			return nil, os.ErrNotExist
		}

		err := cmd.runPublishPreCmdLine()
		require.EqualError(t, err, "Failed to open the config.json file. The file doesn't exist, is corrupted, or has an invalid JSON format. Verify if the file format is JSON or fix its content according to the JSON format specification at https://www.json.org/json-en.html")
	})

	t.Run("publish.env not exist", func(t *testing.T) {
		f, _, _ := testutils.NewFactory(nil)

		cmd := NewPublishCmd(f)

		// Specified publish.env file but it cannot be read correctly
		cmd.FileReader = func(path string) ([]byte, error) {
			return []byte(`{"publish": {"pre_cmd": "ls", "env": "./azion/publish.env"}}`), nil
		}
		cmd.EnvLoader = func(path string) ([]string, error) {
			return nil, os.ErrNotExist
		}

		err := cmd.runPublishPreCmdLine()
		require.ErrorIs(t, err, msg.ErrReadEnvFile)
	})

	t.Run("publish.env is ok", func(t *testing.T) {
		f, _, _ := testutils.NewFactory(nil)

		cmd := NewPublishCmd(f)

		// Specified publish.env file but it cannot be read correctly
		cmd.FileReader = func(path string) ([]byte, error) {
			return []byte(`{"publish": {"pre_cmd": "ls", "env": "./azion/publish.env", "output-ctrl": "on-error", "default": "ls -lia"}}`), nil
		}
		cmd.EnvLoader = func(path string) ([]string, error) {
			return []string{"UEBA=OBA", "FAZER=UM_PENSO"}, nil
		}

		err := cmd.runPublishPreCmdLine()
		require.NoError(t, err)
	})

	t.Run("without specifying publish.env", func(t *testing.T) {
		f, _, _ := testutils.NewFactory(nil)

		cmd := NewPublishCmd(f)
		cmd.FileReader = func(path string) ([]byte, error) {
			return []byte(`{"publish": {"default": "ls"}}`), nil
		}
		cmd.EnvLoader = func(path string) ([]string, error) {
			return nil, nil
		}
		err := cmd.runPublishPreCmdLine()

		require.ErrorIs(t, err, msg.EdgeApplicationsOutputErr)
	})

	t.Run("no pre_cmd.cmd", func(t *testing.T) {
		f, stdout, _ := testutils.NewFactory(nil)

		cmd := NewPublishCmd(f)
		cmd.FileReader = func(path string) ([]byte, error) {
			return []byte(`{"publish": {}}`), nil
		}

		err := cmd.runPublishPreCmdLine()
		require.NoError(t, err)
		require.NotContains(t, stdout.String(), "Running publish command")
	})

	t.Run("full", func(t *testing.T) {
		f, _, _ := testutils.NewFactory(nil)

		envs := []string{"UEBA=OBA", "FAZER=UM_PENSO"}
		cmd := NewPublishCmd(f)
		cmd.FileReader = func(path string) ([]byte, error) {
			return []byte(`{"publish": {"cmd": "ls", "env": "./azion/publish.env"}}`), nil
		}
		cmd.EnvLoader = func(path string) ([]string, error) {
			return envs, nil
		}
		cmd.CommandRunner = func(cmd string, env []string) (string, int, error) {
			if !reflect.DeepEqual(envs, env) {
				return "", -1, errors.New("unexpected env")
			}
			return "Publish pre command run", 0, os.ErrExist
		}

		err := cmd.runPublishPreCmdLine()

		require.NoError(t, err)
	})
	t.Run("failed to create application", func(t *testing.T) {

		mock := &httpmock.Registry{}
		options := &contracts.AzionApplicationOptions{}

		dat, _ := os.ReadFile("./fixtures/create_app.json")
		_ = json.Unmarshal(dat, options)

		mock.Register(
			httpmock.REST("POST", "edge_applications"),
			httpmock.StatusStringResponse(http.StatusBadRequest, "Invalid"),
		)

		f, _, _ := testutils.NewFactory(mock)
		ctx := context.Background()

		cliapp := apiapp.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))

		cmd := NewPublishCmd(f)

		_, err := cmd.createApplication(cliapp, ctx, options, applicationName)
		require.EqualError(t, err, "Failed to create the Edge Application: Invalid. Check your settings and try again. If the error persists, contact Azion support")
	})

	t.Run("create application success", func(t *testing.T) {

		mock := &httpmock.Registry{}
		options := &contracts.AzionApplicationOptions{}

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
			httpmock.JSONFromString(successResponseApp),
		)

		f, _, _ := testutils.NewFactory(mock)
		ctx := context.Background()

		cliapp := apiapp.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))

		cmd := NewPublishCmd(f)

		_, err := cmd.createApplication(cliapp, ctx, options, applicationName)
		require.NoError(t, err)
	})
}
