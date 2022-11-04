package init

import (
	"bytes"
	"errors"
	"github.com/aziontech/azion-cli/pkg/httpmock"
	"io"
	"io/fs"
	"os"
	"strings"
	"testing"

	msg "github.com/aziontech/azion-cli/messages/webapp"
	"github.com/aziontech/azion-cli/pkg/testutils"
	"github.com/aziontech/azion-cli/utils"
	"github.com/stretchr/testify/require"
)

func TestCobraCmd(t *testing.T) {
	t.Run("without package.json", func(t *testing.T) {
		f, _, _ := testutils.NewFactory(nil)

		initCmd := newInitCmd(f)

		initCmd.FileReader = func(path string) ([]byte, error) {
			return nil, os.ErrNotExist
		}

		cmd := newCobraCmd(initCmd)

		cmd.SetArgs([]string{"--name", "SUUPA_DOOPA", "--type", "javascript"})

		err := cmd.Execute()

		require.ErrorIs(t, err, msg.ErrorPackageJsonNotFound)
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

		require.ErrorIs(t, err, msg.ErrorYesAndNoOptions)
	})

	t.Run("success with javascript", func(t *testing.T) {
		mock := &httpmock.Registry{}
		f, stdout, _ := testutils.NewFactory(mock)
		initCmd := newInitCmd(f)

		initCmd.CommandRunner = func(cmd string, envs []string) (string, int, error) {
			return "", 0, nil
		}
		initCmd.FileReader = func(path string) ([]byte, error) {
			return []byte(`{"init": {"cmd": "ls"},"type":"javascript"}`), nil
		}
		addGitignor = func(cmd *InitCmd, path string) error {
			return nil
		}
		initCmd.WriteFile = func(filename string, data []byte, perm fs.FileMode) error {
			return nil
		}
		initCmd.Rename = func(oldpath string, newpath string) error {
			return nil
		}
		initCmd.Stat = func(path string) (fs.FileInfo, error) {
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
		require.Contains(t, stdout.String(), `Template successfully fetched and configured`)
	})

	t.Run("success with javascript using flag -y", func(t *testing.T) {
		mock := &httpmock.Registry{}
		f, stdout, _ := testutils.NewFactory(mock)
		initCmd := newInitCmd(f)

		initCmd.CommandRunner = func(cmd string, envs []string) (string, int, error) {
			return "", 0, nil
		}
		initCmd.FileReader = func(path string) ([]byte, error) {
			return []byte(`{"init": {"cmd": "ls"},"type":"javascript"}`), nil
		}
		addGitignor = func(cmd *InitCmd, path string) error {
			return nil
		}
		initCmd.WriteFile = func(filename string, data []byte, perm fs.FileMode) error {
			return nil
		}
		initCmd.Rename = func(oldpath string, newpath string) error {
			return nil
		}
		initCmd.Stat = func(path string) (fs.FileInfo, error) {
			if !strings.HasSuffix(path, "package.json") {
				return nil, os.ErrNotExist
			}
			return nil, nil
		}

		cmd := newCobraCmd(initCmd)

		cmd.SetArgs([]string{"--name", "SUUPA_DOOPA", "--type", "javascript", "-y"})

		err := cmd.Execute()

		require.NoError(t, err)
		require.Contains(t, stdout.String(), `Template successfully fetched and configured`)
	})

	t.Run("does not overwrite contents", func(t *testing.T) {
		mock := &httpmock.Registry{}
		f, _, _ := testutils.NewFactory(mock)

		initCmd := newInitCmd(f)

		initCmd.CommandRunner = func(cmd string, envs []string) (string, int, error) {
			return "", 0, nil
		}
		initCmd.FileReader = func(path string) ([]byte, error) {
			return []byte(`{"init": {"cmd": "ls"},"type":"javascript"}`), nil
		}
		addGitignor = func(cmd *InitCmd, path string) error {
			return nil
		}
		initCmd.WriteFile = func(filename string, data []byte, perm fs.FileMode) error {
			return nil
		}
		initCmd.Rename = func(oldpath string, newpath string) error {
			return errors.New("unexpected rename")
		}
		initCmd.Stat = func(path string) (fs.FileInfo, error) {
			if !strings.HasSuffix(path, "package.json") {
				return nil, os.ErrNotExist
			}
			return nil, nil
		}
		initCmd.IsDirEmpty = func(dirpath string) (bool, error) {
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
		mock := &httpmock.Registry{}
		f, _, _ := testutils.NewFactory(mock)

		initCmd := newInitCmd(f)

		initCmd.CommandRunner = func(cmd string, envs []string) (string, int, error) {
			return "", 0, nil
		}
		initCmd.FileReader = func(path string) ([]byte, error) {
			return []byte(`{"init": {"cmd": "ls"},"type":"javascript"}`), nil
		}
		addGitignor = func(cmd *InitCmd, path string) error {
			return nil
		}
		initCmd.WriteFile = func(filename string, data []byte, perm fs.FileMode) error {
			return nil
		}
		initCmd.Rename = func(oldpath string, newpath string) error {
			return errors.New("unexpected rename")
		}
		initCmd.Stat = func(path string) (fs.FileInfo, error) {
			if !strings.HasSuffix(path, "package.json") {
				return nil, os.ErrNotExist
			}
			return nil, nil
		}
		initCmd.IsDirEmpty = func(dirpath string) (bool, error) {
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

		initCmd.Stat = func(path string) (fs.FileInfo, error) {
			if !strings.HasSuffix(path, "package.json") {
				return nil, os.ErrNotExist
			}
			return nil, nil
		}
		initCmd.IsDirEmpty = func(dirpath string) (bool, error) {
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
		cmd.FileReader = func(path string) ([]byte, error) {
			return nil, os.ErrNotExist
		}

		i := InitInfo{}
		err := cmd.runInitCmdLine(&i)
		require.EqualError(t, err, "Failed to open the config.json file. The file doesn't exist, is corrupted, or has an invalid JSON format. Verify if the file format is JSON or fix its content according to the JSON format specification at https://www.json.org/json-en.html")
	})

	t.Run("init.env not empty", func(t *testing.T) {
		f, _, _ := testutils.NewFactory(nil)
		cmd := newInitCmd(f)

		// Specified init.env file but it cannot be read correctly
		cmd.FileReader = func(path string) ([]byte, error) {
			return []byte(`{"init": {"cmd": "ls", "env": "./azion/init.env"},"type":"javascript"}`), nil
		}
		cmd.EnvLoader = func(path string) ([]string, error) {
			return nil, msg.ErrReadEnvFile
		}

		addGitignor = func(cmd *InitCmd, path string) error {
			return nil
		}
		cmd.WriteFile = func(filename string, data []byte, perm fs.FileMode) error {
			return nil
		}

		i := InitInfo{}
		err := cmd.runInitCmdLine(&i)
		require.ErrorIs(t, err, msg.ErrReadEnvFile)
	})

	t.Run("Failed to run the command specified", func(t *testing.T) {
		f, _, _ := testutils.NewFactory(nil)

		cmd := newInitCmd(f)
		cmd.FileReader = func(path string) ([]byte, error) {
			return []byte(`{"init": {"cmd": "ls"},"type":"javascript"}`), nil
		}
		addGitignor = func(cmd *InitCmd, path string) error {
			return nil
		}
		cmd.WriteFile = func(filename string, data []byte, perm fs.FileMode) error {
			return nil
		}
		cmd.EnvLoader = func(path string) ([]string, error) {
			return nil, nil
		}
		cmd.CommandRunner = func(cmd string, env []string) (string, int, error) {
			return "my command output", 0, nil
		}

		i := InitInfo{}
		err := cmd.runInitCmdLine(&i)
		require.ErrorIs(t, err, utils.ErrorRunningCommand)
	})

	t.Run("Success with Javacript", func(t *testing.T) {
		f, _, _ := testutils.NewFactory(nil)

		envs := []string{"UEBA=OBA", "FAZER=UM_PENSO"}
		cmd := newInitCmd(f)

		addGitignor = func(cmd *InitCmd, path string) error {
			return nil
		}

		cmd.WriteFile = func(filename string, data []byte, perm fs.FileMode) error {
			return nil
		}

		cmd.Rename = func(oldpath string, newpath string) error {
			return nil
		}

		cmd.Stat = func(path string) (fs.FileInfo, error) {
			return nil, nil
		}

		cmd.FileReader = func(path string) ([]byte, error) {
			return []byte(`{"init": {"cmd": "ls", "env": "./azion/init.env"}, "type": "javascript"}`), nil
		}

		cmd.EnvLoader = func(path string) ([]string, error) {
			return envs, nil
		}

		cmd.CommandRunner = func(cmd string, env []string) (string, int, error) {
			return "my command output", 0, nil
		}

		cmd.Mkdir = func(path string, perm os.FileMode) error {
			return nil
		}

		i := InitInfo{
			TypeLang: "javascript",
		}
		err := cmd.runInitCmdLine(&i)
		require.NoError(t, err)
	})

	t.Run("success with NextJS", func(t *testing.T) {
		f, _, _ := testutils.NewFactory(nil)

		envs := []string{"UEBA=OBA", "FAZER=UM_PENSO"}
		cmd := newInitCmd(f)

		addGitignor = func(cmd *InitCmd, path string) error {
			return nil
		}

		cmd.WriteFile = func(filename string, data []byte, perm fs.FileMode) error {
			return nil
		}

		cmd.Rename = func(oldpath string, newpath string) error {
			return nil
		}

		cmd.Stat = func(path string) (fs.FileInfo, error) {
			return nil, nil
		}

		cmd.FileReader = func(path string) ([]byte, error) {
			return []byte(`{"init": {"cmd": "ls", "env": "./azion/init.env"}, "type": "nextjs"}`), nil
		}

		cmd.EnvLoader = func(path string) ([]string, error) {
			return envs, nil
		}

		cmd.CommandRunner = func(cmd string, env []string) (string, int, error) {
			return "my command output", 0, nil
		}

		cmd.Mkdir = func(path string, perm os.FileMode) error {
			return nil
		}

		i := InitInfo{
			TypeLang: "nextjs",
		}
		err := cmd.runInitCmdLine(&i)
		require.NoError(t, err)
	})

	t.Run("success with Flareact", func(t *testing.T) {
		f, _, _ := testutils.NewFactory(nil)

		envs := []string{"UEBA=OBA", "FAZER=UM_PENSO"}
		cmd := newInitCmd(f)

		addGitignor = func(cmd *InitCmd, path string) error {
			return nil
		}

		cmd.WriteFile = func(filename string, data []byte, perm fs.FileMode) error {
			return nil
		}

		cmd.Rename = func(oldpath string, newpath string) error {
			return nil
		}

		cmd.Stat = func(path string) (fs.FileInfo, error) {
			return nil, nil
		}

		cmd.FileReader = func(path string) ([]byte, error) {
			return []byte(`{"init": {"cmd": "ls", "env": "./azion/init.env"}, "type": "flareact"}`), nil
		}

		cmd.EnvLoader = func(path string) ([]string, error) {
			return envs, nil
		}

		cmd.CommandRunner = func(cmd string, env []string) (string, int, error) {
			return "my command output", 0, nil
		}

		cmd.Mkdir = func(path string, perm os.FileMode) error {
			return nil
		}

		i := InitInfo{
			TypeLang: "flareact",
		}
		err := cmd.runInitCmdLine(&i)
		require.NoError(t, err)
	})

}
