package build

import (
	"bytes"
	"errors"
	"os"
	"reflect"
	"testing"

	msg "github.com/aziontech/azion-cli/messages/webapp"
	"github.com/aziontech/azion-cli/pkg/testutils"
	"github.com/aziontech/azion-cli/utils"
	"github.com/stretchr/testify/require"
)

func TestBuild(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		f, stdout, _ := testutils.NewFactory(nil)

		jsonContent := bytes.NewBufferString(`
        {
            "build": {
                "cmd": "npm run build"
            }
        }
        `)

		envVars := []string{"VAR1=PAODEBATATA", "VAR2=PAODEQUEIJO"}

		command := newBuildCmd(f)

		command.fileReader = func(path string) ([]byte, error) {
			return jsonContent.Bytes(), nil
		}
		command.commandRunner = func(cmd string, envs []string) (string, int, error) {
			if cmd != "npm run build" {
				return "", -1, errors.New("unexpected command")
			}
			if !reflect.DeepEqual(envs, envVars) {
				return "", -1, errors.New("unexpected envvars")
			}
			return "Build completed", 0, nil
		}
		command.envLoader = func(path string) ([]string, error) {
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
            }
        }
        `)

		envVars := []string{"VAR1=PAODEBATATA", "VAR2=PAODEQUEIJO"}
		expectedErr := errors.New("invalid file")

		command := newBuildCmd(f)

		command.fileReader = func(path string) ([]byte, error) {
			return jsonContent.Bytes(), nil
		}
		command.commandRunner = func(cmd string, envs []string) (string, int, error) {
			return "Command output goes here", 42, expectedErr
		}
		command.envLoader = func(path string) ([]string, error) {
			return envVars, nil
		}

		err := command.run()
		require.Error(t, err, expectedErr)

		require.Equal(t, `Running build step command:

$ npm run build
Command output goes here

Command exited with code 42
`, stdout.String())
	})

	t.Run("no build.cmd to execute", func(t *testing.T) {
		f, stdout, _ := testutils.NewFactory(nil)

		command := newBuildCmd(f)

		command.fileReader = func(path string) ([]byte, error) {
			return []byte(`{"build": {}}`), nil
		}

		err := command.run()
		require.NoError(t, err)
		require.NotContains(t, stdout.String(), "Running build step command")
	})

	t.Run("missing config file", func(t *testing.T) {
		f, _, _ := testutils.NewFactory(nil)

		command := newBuildCmd(f)

		command.fileReader = func(path string) ([]byte, error) {
			return nil, os.ErrNotExist
		}

		err := command.run()
		require.ErrorIs(t, err, msg.ErrOpeningConfigFile)
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

		command.fileReader = func(path string) ([]byte, error) {
			return jsonContent.Bytes(), nil
		}

		err := command.run()
		require.ErrorIs(t, err, msg.ErrUnmarshalConfigFile)
	})

	t.Run("invalid env", func(t *testing.T) {
		f, _, _ := testutils.NewFactory(nil)

		jsonContent := bytes.NewBufferString(`
        {
            "build": {
                "cmd": "npm run build"
            }
        }
        `)

		command := newBuildCmd(f)

		command.fileReader = func(path string) ([]byte, error) {
			return jsonContent.Bytes(), nil
		}
		command.envLoader = func(path string) ([]string, error) {
			return nil, utils.ErrorLoadingEnvVars
		}

		err := command.run()
		require.ErrorIs(t, err, msg.ErrReadEnvFile)
	})
}
