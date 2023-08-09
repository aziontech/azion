package init

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"strings"
	"testing"

	"github.com/aziontech/azion-cli/pkg/logger"
	"go.uber.org/zap/zapcore"

	msg "github.com/aziontech/azion-cli/messages/init"
	"github.com/aziontech/azion-cli/pkg/httpmock"
	"github.com/aziontech/azion-cli/pkg/testutils"
	"github.com/go-git/go-git/v5"

	"github.com/stretchr/testify/require"
)

func TestCobraCmd(t *testing.T) {
	logger.New(zapcore.InfoLevel)
	t.Run("success with CDN", func(t *testing.T) {
		mock := &httpmock.Registry{}
		f, stdout, _ := testutils.NewFactory(mock)
		initCmd := NewInitCmd(f)

		initCmd.LookPath = func(bin string) (string, error) {
			return "", nil
		}

		initCmd.FileReader = func(path string) ([]byte, error) {
			return []byte(`{"init": {"cmd": "ls", "output-ctrl": "on-error"}, "type":"static"}`), nil
		}
		initCmd.WriteFile = func(filename string, data []byte, perm fs.FileMode) error {
			return nil
		}
		initCmd.Rename = func(oldpath string, newpath string) error {
			return nil
		}
		initCmd.Mkdir = func(path string, perm os.FileMode) error {
			return nil
		}
		initCmd.GitPlainClone = func(path string, isBare bool, o *git.CloneOptions) (*git.Repository, error) {
			return &git.Repository{}, nil
		}

		cmd := NewCobraCmd(initCmd)

		cmd.SetArgs([]string{"--name", "SUUPA_DOOPA", "--template", "simple"})

		in := bytes.NewBuffer(nil)
		in.WriteString("yes\n")
		f.IOStreams.In = io.NopCloser(in)

		err := cmd.Execute()

		require.NoError(t, err)
		require.Contains(t, stdout.String(), fmt.Sprintf(msg.EdgeApplicationsInitSuccessful, "SUUPA_DOOPA"))
	})

	t.Run("with -y and -n flags", func(t *testing.T) {
		f, _, _ := testutils.NewFactory(nil)

		initCmd := NewInitCmd(f)

		cmd := NewCobraCmd(initCmd)

		cmd.SetArgs([]string{"--name", "BLEBLEBLE", "--template", "demeuamor", "-y", "-n"})

		err := cmd.Execute()

		require.ErrorIs(t, err, msg.ErrorYesAndNoOptions)
	})

	t.Run("success with static", func(t *testing.T) {
		mock := &httpmock.Registry{}
		f, stdout, _ := testutils.NewFactory(mock)
		initCmd := NewInitCmd(f)

		initCmd.LookPath = func(bin string) (string, error) {
			return "", nil
		}
		initCmd.FileReader = func(path string) ([]byte, error) {
			return []byte(`{"init": {"cmd": "ls", "output-ctrl": "on-error"}, "type":"static"}`), nil
		}
		initCmd.WriteFile = func(filename string, data []byte, perm fs.FileMode) error {
			return nil
		}
		initCmd.Rename = func(oldpath string, newpath string) error {
			return nil
		}
		initCmd.Mkdir = func(path string, perm os.FileMode) error {
			return nil
		}
		initCmd.GitPlainClone = func(path string, isBare bool, o *git.CloneOptions) (*git.Repository, error) {
			return &git.Repository{}, nil
		}

		cmd := NewCobraCmd(initCmd)

		cmd.SetArgs([]string{"--name", "SUUPA_DOOPA", "--template", "static"})

		in := bytes.NewBuffer(nil)
		in.WriteString("yes\n")
		f.IOStreams.In = io.NopCloser(in)

		err := cmd.Execute()

		require.NoError(t, err)
		require.Contains(t, stdout.String(), fmt.Sprintf(msg.EdgeApplicationsInitSuccessful+"\n", "SUUPA_DOOPA"))
	})

	t.Run("does not overwrite contents using flag -n", func(t *testing.T) {
		mock := &httpmock.Registry{}
		f, _, _ := testutils.NewFactory(mock)

		initCmd := NewInitCmd(f)

		initCmd.LookPath = func(bin string) (string, error) {
			return "", nil
		}

		initCmd.FileReader = func(path string) ([]byte, error) {
			return []byte(`{"init": {"cmd": "ls", "output-ctrl": "on-error"},"type":"static"}`), nil
		}
		initCmd.WriteFile = func(filename string, data []byte, perm fs.FileMode) error {
			return nil
		}
		initCmd.Rename = func(oldpath string, newpath string) error {
			return errors.New("unexpected rename")
		}
		initCmd.GitPlainClone = func(path string, isBare bool, o *git.CloneOptions) (*git.Repository, error) {
			return &git.Repository{}, nil
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
		initCmd.Mkdir = func(path string, perm os.FileMode) error {
			return nil
		}

		cmd := NewCobraCmd(initCmd)

		cmd.SetArgs([]string{"--name", "SUUPA_DOOPA", "--template", "static", "-n"})

		err := cmd.Execute()

		require.NoError(t, err)
	})
}
