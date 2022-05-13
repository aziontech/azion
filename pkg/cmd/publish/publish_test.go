package publish

import (
	"errors"
	"os"
	"reflect"
	"testing"

	"github.com/aziontech/azion-cli/pkg/testutils"
	"github.com/stretchr/testify/require"
)

func TestPublishCmd(t *testing.T) {
	t.Run("without package.json", func(t *testing.T) {
		f, _, _ := testutils.NewFactory(nil)

		publishCmd := newPublishCmd(f)

		publishCmd.fileReader = func(path string) ([]byte, error) {
			return nil, os.ErrNotExist
		}

		cmd := newCobraCmd(publishCmd)

		cmd.SetArgs([]string{""})

		err := cmd.Execute()

		require.EqualError(t, err, "Failed to open config.json file")
	})

	t.Run("without config.json", func(t *testing.T) {
		f, _, _ := testutils.NewFactory(nil)

		cmd := newPublishCmd(f)
		cmd.fileReader = func(path string) ([]byte, error) {
			return nil, os.ErrNotExist
		}

		err := cmd.runPublishPreCmdLine()
		require.EqualError(t, err, "Failed to open config.json file")
	})

	t.Run("publish.env not exist", func(t *testing.T) {
		f, _, _ := testutils.NewFactory(nil)

		cmd := newPublishCmd(f)

		// Specified publish.env file but it cannot be read correctly
		cmd.fileReader = func(path string) ([]byte, error) {
			return []byte(`{"publish": {"pre_cmd": "ls", "env": "./azion/publish.env"}}`), nil
		}
		cmd.envLoader = func(path string) ([]string, error) {
			return nil, os.ErrNotExist
		}

		err := cmd.runPublishPreCmdLine()
		require.ErrorIs(t, err, os.ErrNotExist)
	})

	t.Run("publish.env is ok", func(t *testing.T) {
		f, stdout, _ := testutils.NewFactory(nil)

		cmd := newPublishCmd(f)

		// Specified publish.env file but it cannot be read correctly
		cmd.fileReader = func(path string) ([]byte, error) {
			return []byte(`{"publish": {"pre_cmd": "ls", "env": "./azion/publish.env"}}`), nil
		}
		cmd.envLoader = func(path string) ([]string, error) {
			return []string{"UEBA=OBA", "FAZER=UM_PENSO"}, nil
		}

		err := cmd.runPublishPreCmdLine()
		require.NoError(t, err)
		require.Contains(t, stdout.String(), "Command exited with code 0")
	})

	t.Run("without specifing publish.env", func(t *testing.T) {
		f, _, _ := testutils.NewFactory(nil)

		cmd := newPublishCmd(f)
		cmd.fileReader = func(path string) ([]byte, error) {
			return []byte(`{"publish": {"cmd": "ls"}}`), nil
		}
		cmd.envLoader = func(path string) ([]string, error) {
			return nil, nil
		}
		cmd.commandRunner = func(cmd string, env []string) (string, int, error) {
			if env != nil {
				return "", -1, errors.New("unexpected env")
			}
			return "", 0, nil
		}

		err := cmd.runPublishPreCmdLine()

		require.NoError(t, err)
	})

	t.Run("no pre_cmd.cmd", func(t *testing.T) {
		f, stdout, _ := testutils.NewFactory(nil)

		cmd := newPublishCmd(f)
		cmd.fileReader = func(path string) ([]byte, error) {
			return []byte(`{"publish": {}}`), nil
		}

		err := cmd.runPublishPreCmdLine()
		require.NoError(t, err)
		require.NotContains(t, stdout.String(), "Running publish command")
	})

	t.Run("full", func(t *testing.T) {
		f, _, _ := testutils.NewFactory(nil)

		envs := []string{"UEBA=OBA", "FAZER=UM_PENSO"}
		cmd := newPublishCmd(f)
		cmd.fileReader = func(path string) ([]byte, error) {
			return []byte(`{"publish": {"cmd": "ls", "env": "./azion/publish.env"}}`), nil
		}
		cmd.envLoader = func(path string) ([]string, error) {
			return envs, nil
		}
		cmd.commandRunner = func(cmd string, env []string) (string, int, error) {
			if !reflect.DeepEqual(envs, env) {
				return "", -1, errors.New("unexpected env")
			}
			return "Publish pre command run", 0, os.ErrExist
		}

		err := cmd.runPublishPreCmdLine()

		require.NoError(t, err)
	})
}
