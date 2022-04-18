package init

import (
	"bytes"
	"errors"
	"io"
	"io/fs"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/aziontech/azion-cli/pkg/httpmock"
	"github.com/aziontech/azion-cli/pkg/testutils"
	"github.com/aziontech/azion-cli/utils"
	"github.com/stretchr/testify/require"
)

func TestCobraCmd(t *testing.T) {
	t.Run("without package.json", func(t *testing.T) {
		f, _, _ := testutils.NewFactory(nil)

		initCmd := newInitCmd(f)

		initCmd.fileReader = func(path string) ([]byte, error) {
			return nil, os.ErrNotExist
		}

		cmd := newCobraCmd(initCmd)

		cmd.SetArgs([]string{"--name", "SUUPA_DOOPA", "--type", "javascript"})

		err := cmd.Execute()

		require.ErrorIs(t, err, ErrorPackageJsonNotFound)
	})

	t.Run("with unsupported type", func(t *testing.T) {
		f, _, _ := testutils.NewFactory(nil)

		initCmd := newInitCmd(f)

		cmd := newCobraCmd(initCmd)

		cmd.SetArgs([]string{"--name", "BLEBLEBLE", "--type", "demeuamor"})

		err := cmd.Execute()

		require.ErrorIs(t, err, utils.ErrorUnsupportedType)
	})

	t.Run("with -y and -n flags", func(t *testing.T) {
		f, _, _ := testutils.NewFactory(nil)

		initCmd := newInitCmd(f)

		cmd := newCobraCmd(initCmd)

		cmd.SetArgs([]string{"--name", "BLEBLEBLE", "--type", "demeuamor", "-y", "-n"})

		err := cmd.Execute()

		require.ErrorIs(t, err, ErrorYesAndNoOptions)
	})

	t.Run("success with javascript", func(t *testing.T) {
		mock := &httpmock.Registry{}
		f, stdout, _ := testutils.NewFactory(mock)
		initCmd := newInitCmd(f)

		initCmd.commandRunner = func(cmd string, envs []string) (string, int, error) {
			if !strings.HasPrefix(cmd, GIT) && !strings.HasPrefix(cmd, "ls") {
				return "", -1, errors.New("unexpected command")
			}
			return "", 0, nil
		}
		initCmd.fileReader = func(path string) ([]byte, error) {
			return []byte(`{"init": {"cmd": "ls"}}`), nil
		}
		initCmd.writeFile = func(filename string, data []byte, perm fs.FileMode) error {
			return nil
		}
		initCmd.rename = func(oldpath string, newpath string) error {
			return nil
		}
		initCmd.stat = func(path string) (fs.FileInfo, error) {
			if !strings.HasSuffix(path, "package.json") {
				return nil, os.ErrNotExist
			}
			return nil, nil
		}

		cmd := newCobraCmd(initCmd)

		cmd.SetArgs([]string{"--name", "SUUPA_DOOPA", "--type", "javascript"})

		in := bytes.NewBuffer(nil)
		in.WriteString("yes\n")
		f.IOStreams.In = io.NopCloser(in)

		err := cmd.Execute()

		require.NoError(t, err)
		require.Contains(t, stdout.String(), `Template successfully fetched and configured
`)
	})

	t.Run("success with javascript using flag -y", func(t *testing.T) {
		f, stdout, _ := testutils.NewFactory(nil)

		initCmd := newInitCmd(f)

		initCmd.commandRunner = func(cmd string, envs []string) (string, int, error) {
			if !strings.HasPrefix(cmd, GIT) && !strings.HasPrefix(cmd, "ls") {
				return "", -1, errors.New("unexpected command")
			}
			return "", 0, nil
		}
		initCmd.fileReader = func(path string) ([]byte, error) {
			return []byte(`{"init": {"cmd": "ls"}}`), nil
		}
		initCmd.writeFile = func(filename string, data []byte, perm fs.FileMode) error {
			return nil
		}
		initCmd.rename = func(oldpath string, newpath string) error {
			return nil
		}
		initCmd.stat = func(path string) (fs.FileInfo, error) {
			if !strings.HasSuffix(path, "package.json") {
				return nil, os.ErrNotExist
			}
			return nil, nil
		}

		cmd := newCobraCmd(initCmd)

		cmd.SetArgs([]string{"--name", "SUUPA_DOOPA", "--type", "javascript", "-y"})

		err := cmd.Execute()

		require.NoError(t, err)
		require.Contains(t, stdout.String(), `Template successfully fetched and configured
`)
	})

	t.Run("does not overwrite contents", func(t *testing.T) {
		mock := &httpmock.Registry{}
		f, _, _ := testutils.NewFactory(mock)

		initCmd := newInitCmd(f)

		initCmd.commandRunner = func(cmd string, envs []string) (string, int, error) {
			if !strings.HasPrefix(cmd, GIT) && !strings.HasPrefix(cmd, "ls") {
				return "", -1, errors.New("unexpected command")
			}
			return "", 0, nil
		}
		initCmd.fileReader = func(path string) ([]byte, error) {
			return []byte(`{"init": {"cmd": "ls"}}`), nil
		}
		initCmd.writeFile = func(filename string, data []byte, perm fs.FileMode) error {
			return nil
		}
		initCmd.rename = func(oldpath string, newpath string) error {
			return errors.New("unexpected rename")
		}
		initCmd.stat = func(path string) (fs.FileInfo, error) {
			if !strings.HasSuffix(path, "package.json") {
				return nil, os.ErrNotExist
			}
			return nil, nil
		}
		initCmd.isDirEmpty = func(dirpath string) (bool, error) {
			return false, nil
		}

		cmd := newCobraCmd(initCmd)

		cmd.SetArgs([]string{"--name", "SUUPA_DOOPA", "--type", "javascript"})

		in := bytes.NewBuffer(nil)
		in.WriteString("no\n")
		f.IOStreams.In = io.NopCloser(in)

		err := cmd.Execute()

		require.NoError(t, err)
	})

	t.Run("does not overwrite contents using flag -n", func(t *testing.T) {
		f, _, _ := testutils.NewFactory(nil)

		initCmd := newInitCmd(f)

		initCmd.commandRunner = func(cmd string, envs []string) (string, int, error) {
			if !strings.HasPrefix(cmd, GIT) && !strings.HasPrefix(cmd, "ls") {
				return "", -1, errors.New("unexpected command")
			}
			return "", 0, nil
		}
		initCmd.fileReader = func(path string) ([]byte, error) {
			return []byte(`{"init": {"cmd": "ls"}}`), nil
		}
		initCmd.writeFile = func(filename string, data []byte, perm fs.FileMode) error {
			return errors.New("unexpected write")
		}
		initCmd.rename = func(oldpath string, newpath string) error {
			return errors.New("unexpected rename")
		}
		initCmd.stat = func(path string) (fs.FileInfo, error) {
			if !strings.HasSuffix(path, "package.json") {
				return nil, os.ErrNotExist
			}
			return nil, nil
		}
		initCmd.isDirEmpty = func(dirpath string) (bool, error) {
			return false, nil
		}

		cmd := newCobraCmd(initCmd)

		cmd.SetArgs([]string{"--name", "SUUPA_DOOPA", "--type", "javascript", "-n"})

		err := cmd.Execute()

		require.NoError(t, err)
	})

	t.Run("invalid option", func(t *testing.T) {
		f, _, _ := testutils.NewFactory(nil)

		initCmd := newInitCmd(f)

		cmd := newCobraCmd(initCmd)

		initCmd.stat = func(path string) (fs.FileInfo, error) {
			if !strings.HasSuffix(path, "package.json") {
				return nil, os.ErrNotExist
			}
			return nil, nil
		}
		initCmd.isDirEmpty = func(dirpath string) (bool, error) {
			return false, nil
		}

		cmd.SetArgs([]string{"--name", "SUUPA_DOOPA", "--type", "javascript"})

		in := bytes.NewBuffer(nil)
		in.WriteString("pix\n")
		f.IOStreams.In = io.NopCloser(in)

		err := cmd.Execute()

		require.ErrorIs(t, err, utils.ErrorInvalidOption)
	})

}

func TestInitCmd(t *testing.T) {
	t.Run("without config.json", func(t *testing.T) {
		f, _, _ := testutils.NewFactory(nil)

		cmd := newInitCmd(f)
		cmd.fileReader = func(path string) ([]byte, error) {
			return nil, os.ErrNotExist
		}

		err := cmd.runInitCmdLine()
		require.EqualError(t, err, "Failed to open config.json file")
	})

	t.Run("init.env not empty", func(t *testing.T) {
		f, _, _ := testutils.NewFactory(nil)

		cmd := newInitCmd(f)

		// Specified init.env file but it cannot be read correctly
		cmd.fileReader = func(path string) ([]byte, error) {
			return []byte(`{"init": {"cmd": "ls", "env": "./azion/init.env"}}`), nil
		}
		cmd.envLoader = func(path string) ([]string, error) {
			return nil, os.ErrNotExist
		}

		err := cmd.runInitCmdLine()
		require.ErrorIs(t, err, os.ErrNotExist)
	})

	t.Run("without specifing init.env", func(t *testing.T) {
		f, stdout, _ := testutils.NewFactory(nil)

		cmd := newInitCmd(f)
		cmd.fileReader = func(path string) ([]byte, error) {
			return []byte(`{"init": {"cmd": "ls"}}`), nil
		}
		cmd.envLoader = func(path string) ([]string, error) {
			return nil, nil
		}
		cmd.commandRunner = func(cmd string, env []string) (string, int, error) {
			if env != nil {
				return "", -1, errors.New("unexpected env")
			}
			return "my command output", 0, nil
		}

		err := cmd.runInitCmdLine()
		require.NoError(t, err)

		require.NoError(t, err)
		require.Contains(t, stdout.String(), "my command output")
	})

	t.Run("no init.cmd", func(t *testing.T) {
		f, stdout, _ := testutils.NewFactory(nil)

		cmd := newInitCmd(f)
		cmd.fileReader = func(path string) ([]byte, error) {
			return []byte(`{"init": {}}`), nil
		}

		err := cmd.runInitCmdLine()
		require.NoError(t, err)
		require.NotContains(t, stdout.String(), "Running init command")
	})

	t.Run("full", func(t *testing.T) {
		f, stdout, _ := testutils.NewFactory(nil)

		envs := []string{"UEBA=OBA", "FAZER=UM_PENSO"}
		cmd := newInitCmd(f)
		cmd.fileReader = func(path string) ([]byte, error) {
			return []byte(`{"init": {"cmd": "ls", "env": "./azion/init.env"}}`), nil
		}
		cmd.envLoader = func(path string) ([]string, error) {
			return envs, nil
		}
		cmd.commandRunner = func(cmd string, env []string) (string, int, error) {
			if !reflect.DeepEqual(envs, env) {
				return "", -1, errors.New("unexpected env")
			}
			return "my command output", 0, nil
		}

		err := cmd.runInitCmdLine()
		require.NoError(t, err)

		require.NoError(t, err)
		require.Contains(t, stdout.String(), "my command output")
		require.Contains(t, stdout.String(), "Running init step command")

	})
}
