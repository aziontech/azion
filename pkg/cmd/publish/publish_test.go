package publish

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

		publishCmd := newPublishCmd(f)

		publishCmd.fileReader = func(path string) ([]byte, error) {
			return nil, os.ErrNotExist
		}

		cmd := newCobraCmd(publishCmd)

		cmd.SetArgs([]string{"--name", "SUUPA_DOOPA", "--type", "javascript"})

		err := cmd.Execute()

		require.ErrorIs(t, err, ErrorPackageJsonNotFound)
	})

	t.Run("with unsupported type", func(t *testing.T) {
		f, _, _ := testutils.NewFactory(nil)

		publishCmd := newPublishCmd(f)

		cmd := newCobraCmd(publishCmd)

		cmd.SetArgs([]string{"--name", "BLEBLEBLE", "--type", "demeuamor"})

		err := cmd.Execute()

		require.ErrorIs(t, err, utils.ErrorUnsupportedType)
	})

	t.Run("with -y and -n flags", func(t *testing.T) {
		f, _, _ := testutils.NewFactory(nil)

		publishCmd := newPublishCmd(f)

		cmd := newCobraCmd(publishCmd)

		cmd.SetArgs([]string{"--name", "BLEBLEBLE", "--type", "demeuamor", "-y", "-n"})

		err := cmd.Execute()

		require.ErrorIs(t, err, ErrorYesAndNoOptions)
	})

	t.Run("success with javascript", func(t *testing.T) {
		mock := &httpmock.Registry{}
		f, stdout, _ := testutils.NewFactory(mock)
		publishCmd := newPublishCmd(f)

		publishCmd.commandRunner = func(cmd string, envs []string) (string, int, error) {
			if !strings.HasPrefix(cmd, GIT) && !strings.HasPrefix(cmd, "ls") {
				return "", -1, errors.New("unexpected command")
			}
			return "", 0, nil
		}
		publishCmd.fileReader = func(path string) ([]byte, error) {
			return []byte(`{"publish": {"cmd": "ls"}}`), nil
		}
		publishCmd.writeFile = func(filename string, data []byte, perm fs.FileMode) error {
			return nil
		}
		publishCmd.rename = func(oldpath string, newpath string) error {
			return nil
		}
		publishCmd.stat = func(path string) (fs.FileInfo, error) {
			if !strings.HasSuffix(path, "package.json") {
				return nil, os.ErrNotExist
			}
			return nil, nil
		}

		cmd := newCobraCmd(publishCmd)

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

		publishCmd := newPublishCmd(f)

		publishCmd.commandRunner = func(cmd string, envs []string) (string, int, error) {
			if !strings.HasPrefix(cmd, GIT) && !strings.HasPrefix(cmd, "ls") {
				return "", -1, errors.New("unexpected command")
			}
			return "", 0, nil
		}
		publishCmd.fileReader = func(path string) ([]byte, error) {
			return []byte(`{"publish": {"cmd": "ls"}}`), nil
		}
		publishCmd.writeFile = func(filename string, data []byte, perm fs.FileMode) error {
			return nil
		}
		publishCmd.rename = func(oldpath string, newpath string) error {
			return nil
		}
		publishCmd.stat = func(path string) (fs.FileInfo, error) {
			if !strings.HasSuffix(path, "package.json") {
				return nil, os.ErrNotExist
			}
			return nil, nil
		}

		cmd := newCobraCmd(publishCmd)

		cmd.SetArgs([]string{"--name", "SUUPA_DOOPA", "--type", "javascript", "-y"})

		err := cmd.Execute()

		require.NoError(t, err)
		require.Contains(t, stdout.String(), `Template successfully fetched and configured
`)
	})

	t.Run("does not overwrite contents", func(t *testing.T) {
		mock := &httpmock.Registry{}
		f, _, _ := testutils.NewFactory(mock)

		publishCmd := newPublishCmd(f)

		publishCmd.commandRunner = func(cmd string, envs []string) (string, int, error) {
			if !strings.HasPrefix(cmd, GIT) && !strings.HasPrefix(cmd, "ls") {
				return "", -1, errors.New("unexpected command")
			}
			return "", 0, nil
		}
		publishCmd.fileReader = func(path string) ([]byte, error) {
			return []byte(`{"publish": {"cmd": "ls"}}`), nil
		}
		publishCmd.writeFile = func(filename string, data []byte, perm fs.FileMode) error {
			return nil
		}
		publishCmd.rename = func(oldpath string, newpath string) error {
			return errors.New("unexpected rename")
		}
		publishCmd.stat = func(path string) (fs.FileInfo, error) {
			if !strings.HasSuffix(path, "package.json") {
				return nil, os.ErrNotExist
			}
			return nil, nil
		}
		publishCmd.isDirEmpty = func(dirpath string) (bool, error) {
			return false, nil
		}

		cmd := newCobraCmd(publishCmd)

		cmd.SetArgs([]string{"--name", "SUUPA_DOOPA", "--type", "javascript"})

		in := bytes.NewBuffer(nil)
		in.WriteString("no\n")
		f.IOStreams.In = io.NopCloser(in)

		err := cmd.Execute()

		require.NoError(t, err)
	})

	t.Run("does not overwrite contents using flag -n", func(t *testing.T) {
		f, _, _ := testutils.NewFactory(nil)

		publishCmd := newPublishCmd(f)

		publishCmd.commandRunner = func(cmd string, envs []string) (string, int, error) {
			if !strings.HasPrefix(cmd, GIT) && !strings.HasPrefix(cmd, "ls") {
				return "", -1, errors.New("unexpected command")
			}
			return "", 0, nil
		}
		publishCmd.fileReader = func(path string) ([]byte, error) {
			return []byte(`{"publish": {"cmd": "ls"}}`), nil
		}
		publishCmd.writeFile = func(filename string, data []byte, perm fs.FileMode) error {
			return errors.New("unexpected write")
		}
		publishCmd.rename = func(oldpath string, newpath string) error {
			return errors.New("unexpected rename")
		}
		publishCmd.stat = func(path string) (fs.FileInfo, error) {
			if !strings.HasSuffix(path, "package.json") {
				return nil, os.ErrNotExist
			}
			return nil, nil
		}
		publishCmd.isDirEmpty = func(dirpath string) (bool, error) {
			return false, nil
		}

		cmd := newCobraCmd(publishCmd)

		cmd.SetArgs([]string{"--name", "SUUPA_DOOPA", "--type", "javascript", "-n"})

		err := cmd.Execute()

		require.NoError(t, err)
	})

	t.Run("invalid option", func(t *testing.T) {
		f, _, _ := testutils.NewFactory(nil)

		publishCmd := newPublishCmd(f)

		cmd := newCobraCmd(publishCmd)

		publishCmd.stat = func(path string) (fs.FileInfo, error) {
			if !strings.HasSuffix(path, "package.json") {
				return nil, os.ErrNotExist
			}
			return nil, nil
		}
		publishCmd.isDirEmpty = func(dirpath string) (bool, error) {
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

func TestPublishCmd(t *testing.T) {
	t.Run("without config.json", func(t *testing.T) {
		f, _, _ := testutils.NewFactory(nil)

		cmd := newPublishCmd(f)
		cmd.fileReader = func(path string) ([]byte, error) {
			return nil, os.ErrNotExist
		}

		err := cmd.runPublishCmdLine()
		require.EqualError(t, err, "Failed to open config.json file")
	})

	t.Run("publish.env not empty", func(t *testing.T) {
		f, _, _ := testutils.NewFactory(nil)

		cmd := newPublishCmd(f)

		// Specified publish.env file but it cannot be read correctly
		cmd.fileReader = func(path string) ([]byte, error) {
			return []byte(`{"publish": {"cmd": "ls", "env": "./azion/publish.env"}}`), nil
		}
		cmd.envLoader = func(path string) ([]string, error) {
			return nil, os.ErrNotExist
		}

		err := cmd.runPublishCmdLine()
		require.ErrorIs(t, err, os.ErrNotExist)
	})

	t.Run("without specifing publish.env", func(t *testing.T) {
		f, stdout, _ := testutils.NewFactory(nil)

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
			return "my command output", 0, nil
		}

		err := cmd.runPublishCmdLine()
		require.NoError(t, err)

		require.NoError(t, err)
		require.Contains(t, stdout.String(), "my command output")
	})

	t.Run("no publish.cmd", func(t *testing.T) {
		f, stdout, _ := testutils.NewFactory(nil)

		cmd := newPublishCmd(f)
		cmd.fileReader = func(path string) ([]byte, error) {
			return []byte(`{"publish": {}}`), nil
		}

		err := cmd.runPublishCmdLine()
		require.NoError(t, err)
		require.NotContains(t, stdout.String(), "Running publish command")
	})

	t.Run("full", func(t *testing.T) {
		f, stdout, _ := testutils.NewFactory(nil)

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
			return "my command output", 0, nil
		}

		err := cmd.runPublishCmdLine()
		require.NoError(t, err)

		require.NoError(t, err)
		require.Contains(t, stdout.String(), "my command output")
		require.Contains(t, stdout.String(), "Running publish step command")

	})
}
