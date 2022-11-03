package build

import (
	"bytes"
	"errors"
	msg "github.com/aziontech/azion-cli/messages/webapp"
	"github.com/aziontech/azion-cli/utils"
	"os"
	"testing"

	"github.com/aziontech/azion-cli/pkg/testutils"
	"github.com/stretchr/testify/require"
)

func TestBuild(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		f, stdout, _ := testutils.NewFactory(nil)

		jsonContent := bytes.NewBufferString(`
        {
            "build": {
                "cmd": "npm run build"
            },
			"type": "javascript"
        }
        `)

		envVars := []string{"VAR1=PAODEBATATA", "VAR2=PAODEQUEIJO"}

		command := newBuildCmd(f)

		command.FileReader = func(path string) ([]byte, error) {
			return jsonContent.Bytes(), nil
		}

		command.CommandRunner = func(cmd string, envs []string) (string, int, error) {
			return "Build completed", 0, nil
		}
		command.EnvLoader = func(path string) ([]string, error) {
			return envVars, nil
		}

		err := command.run()
		require.NoError(t, err)

		require.Equal(t, `Running build step command:

$ npm run build
Build completed

Command exited with code 0
`, stdout.String())
	})

	t.Run("cmd failed", func(t *testing.T) {
		f, stdout, _ := testutils.NewFactory(nil)

		jsonContent := bytes.NewBufferString(`
        {
            "build": {
                "cmd": "npm run build"
            },
			"type": "javascript"
        }
        `)

		envVars := []string{"VAR1=PAODEBATATA", "VAR2=PAODEQUEIJO"}
		expectedErr := errors.New("invalid file")

		command := newBuildCmd(f)

		command.FileReader = func(path string) ([]byte, error) {
			return jsonContent.Bytes(), nil
		}
		command.CommandRunner = func(cmd string, envs []string) (string, int, error) {
			return "Command output goes here", 42, expectedErr
		}
		command.EnvLoader = func(path string) ([]string, error) {
			return envVars, nil
		}

		err := command.run()
		require.Error(t, err, expectedErr)

		require.Equal(t, `Running build step command:

$ npm run build
`, stdout.String())
	})

	t.Run("in build.cmd to run, type not informed", func(t *testing.T) {
		f, _, _ := testutils.NewFactory(nil)

		command := newBuildCmd(f)

		command.FileReader = func(path string) ([]byte, error) {
			return []byte(`{"build": {}}`), nil
		}

		err := command.run()
		require.ErrorIs(t, err, utils.ErrorUnsupportedType)
	})

	t.Run("missing config file", func(t *testing.T) {
		f, _, _ := testutils.NewFactory(nil)

		command := newBuildCmd(f)

		command.FileReader = func(path string) ([]byte, error) {
			return nil, os.ErrNotExist
		}

		err := command.run()
		require.ErrorIs(t, err, msg.ErrorOpeningAzionFile)
	})

	t.Run("invalid json", func(t *testing.T) {
		f, _, _ := testutils.NewFactory(nil)

		jsonContent := bytes.NewBufferString(`
	   {
	       "build": {
	           "cmd": rm -rm *
	       }
	   }
	   `)

		command := newBuildCmd(f)

		command.FileReader = func(path string) ([]byte, error) {
			return jsonContent.Bytes(), nil
		}

		err := command.run()
		require.ErrorIs(t, err, utils.ErrorUnsupportedType)
	})

	t.Run("invalid env", func(t *testing.T) {
		f, _, _ := testutils.NewFactory(nil)

		jsonContent := bytes.NewBufferString(`
	   {
			"build": {
		   		"cmd": "npm run build"
	   		},
			"type": "javascript"
	   }
	   `)

		command := newBuildCmd(f)

		command.FileReader = func(path string) ([]byte, error) {
			return jsonContent.Bytes(), nil
		}
		command.EnvLoader = func(path string) ([]string, error) {
			return nil, utils.ErrorLoadingEnvVars
		}

		err := command.run()
		require.ErrorIs(t, err, msg.ErrReadEnvFile)
	})
}
