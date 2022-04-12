package build

import (
	"bytes"
	"errors"
	"os"
	"testing"

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

		build := &buildCmd{
			io: f.IOStreams,
			fileReader: func(path string) ([]byte, error) {
				return jsonContent.Bytes(), nil
			},
			commandRunner: func(cmd string, envs []string) (string, int, error) {
				return "Build completed", 0, nil
			},
			configRelativePath: "/azion/config.json",
			getWorkDir: func() (string, error) {
				return "/", nil
			},
			envLoader: func(path string) ([]string, error) {
				return envVars, nil
			},
		}

		err := build.runCmd()
		require.NoError(t, err)

		require.Equal(t, `Running build command

> npm run build
Build completed

Command exited with status code 0
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

		build := &buildCmd{
			io: f.IOStreams,
			fileReader: func(path string) ([]byte, error) {
				return jsonContent.Bytes(), nil
			},
			commandRunner: func(cmd string, envs []string) (string, int, error) {
				return "Build failed", 1, expectedErr
			},
			configRelativePath: "/azion/config.json",
			getWorkDir: func() (string, error) {
				return "/", nil
			},
			envLoader: func(path string) ([]string, error) {
				return envVars, nil
			},
		}

		err := build.runCmd()
		require.Error(t, err, expectedErr)

		require.Equal(t, `Running build command

> npm run build
Build failed

Command exited with status code 1
`, stdout.String())
	})

	t.Run("missing config file", func(t *testing.T) {
		f, _, _ := testutils.NewFactory(nil)

		build := &buildCmd{
			io: f.IOStreams,
			fileReader: func(path string) ([]byte, error) {
				return nil, os.ErrNotExist
			},
			commandRunner: func(cmd string, envs []string) (string, int, error) {
				return "", 0, nil
			},
			configRelativePath: "/azion/config.json",
			getWorkDir: func() (string, error) {
				return "/", nil
			},
			envLoader: func(path string) ([]string, error) {
				return nil, nil
			},
		}

		err := build.runCmd()
		require.ErrorIs(t, err, ErrOpeningConfigFile)
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

		build := &buildCmd{
			io: f.IOStreams,
			fileReader: func(path string) ([]byte, error) {
				return jsonContent.Bytes(), nil
			},
			commandRunner: func(cmd string, envs []string) (string, int, error) {
				return "", 0, nil
			},
			configRelativePath: "/azion/config.json",
			getWorkDir: func() (string, error) {
				return "/", nil
			},
			envLoader: func(path string) ([]string, error) {
				return nil, nil
			},
		}

		err := build.runCmd()
		require.ErrorIs(t, err, ErrUnmarshalConfigFile)
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

		build := &buildCmd{
			io: f.IOStreams,
			fileReader: func(path string) ([]byte, error) {
				return jsonContent.Bytes(), nil
			},
			commandRunner: func(cmd string, envs []string) (string, int, error) {
				return "", 0, nil
			},
			configRelativePath: "/azion/config.json",
			getWorkDir: func() (string, error) {
				return "/", nil
			},
			envLoader: func(path string) ([]string, error) {
				return nil, utils.ErrorLoadingEnvVars
			},
		}

		err := build.runCmd()
		require.ErrorIs(t, err, ErrReadEnvFile)
	})
}
