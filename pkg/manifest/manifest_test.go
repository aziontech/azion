package manifest

import (
	"testing"

	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/httpmock"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/testutils"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
)

func TestDeployCmd(t *testing.T) {
	logger.New(zapcore.DebugLevel)
	msgs := []string{}
	t.Run("get manifest path", func(t *testing.T) {
		interpreter := NewManifestInterpreter()

		_, err := interpreter.ManifestPath()
		require.NoError(t, err)
	})

	t.Run("read manifest structure", func(t *testing.T) {
		f, _, _ := testutils.NewFactory(nil)

		interpreter := NewManifestInterpreter()

		pathManifest := "fixtures/manifest.json"
		_, err := interpreter.ReadManifest(pathManifest, f, &msgs)
		require.NoError(t, err)
	})

	t.Run("create resources", func(t *testing.T) {
		mock := &httpmock.Registry{}
		options := &contracts.AzionApplicationOptions{
			Name: "NotAVeryGoodName",
			Application: contracts.AzionJsonDataApplication{
				ID: 1673635841,
			},
		}

		mock.Register(
			httpmock.REST("POST", "edge_applications/1673635841/cache_settings"),
			httpmock.JSONFromFile("./fixtures/cachesuccess.json"),
		)

		mock.Register(
			httpmock.REST("POST", "edge_applications/1673635841/rules_engine/request/rules"),
			httpmock.JSONFromFile("./fixtures/rulessuccess.json"),
		)

		f, _, _ := testutils.NewFactory(mock)

		interpreter := NewManifestInterpreter()

		interpreter.WriteAzionJsonContent = func(conf *contracts.AzionApplicationOptions, confPath string) error {
			return nil
		}

		pathManifest := "fixtures/manifest.json"
		manifest, err := interpreter.ReadManifest(pathManifest, f, &msgs)
		require.NoError(t, err)
		funcs := make(map[string]contracts.AzionJsonDataFunction)
		err = interpreter.CreateResources(options, manifest, funcs, f, "azion", &msgs)
		require.NoError(t, err)
	})

}
