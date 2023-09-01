package link

import (
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"os"
	"testing"

	"github.com/aziontech/azion-cli/pkg/logger"
	"go.uber.org/zap/zapcore"

	msg "github.com/aziontech/azion-cli/messages/link"
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
		linkCmd := NewLinkCmd(f)

		linkCmd.LookPath = func(bin string) (string, error) {
			return "", nil
		}

		linkCmd.FileReader = func(path string) ([]byte, error) {
			return []byte(`{"init": {"cmd": "ls", "output-ctrl": "on-error"}, "type":"static"}`), nil
		}
		linkCmd.WriteFile = func(filename string, data []byte, perm fs.FileMode) error {
			return nil
		}
		linkCmd.Rename = func(oldpath string, newpath string) error {
			return nil
		}
		linkCmd.Mkdir = func(path string, perm os.FileMode) error {
			return nil
		}
		linkCmd.GitPlainClone = func(path string, isBare bool, o *git.CloneOptions) (*git.Repository, error) {
			return &git.Repository{}, nil
		}
		linkCmd.ShouldConfigure = func(info *LinkInfo) (bool, error) {
			return true, nil
		}

		cmd := NewCobraCmd(linkCmd, f)

		cmd.SetArgs([]string{"--name", "SUUPA_DOOPA", "--preset", "simple"})

		in := bytes.NewBuffer(nil)
		in.WriteString("yes\n")
		f.IOStreams.In = io.NopCloser(in)

		err := cmd.Execute()

		require.NoError(t, err)
		require.Contains(t, stdout.String(), fmt.Sprintf(msg.EdgeApplicationsLinkSuccessful, "SUUPA_DOOPA"))
	})

	t.Run("success with static", func(t *testing.T) {
		mock := &httpmock.Registry{}
		f, stdout, _ := testutils.NewFactory(mock)
		linkCmd := NewLinkCmd(f)

		linkCmd.LookPath = func(bin string) (string, error) {
			return "", nil
		}
		linkCmd.FileReader = func(path string) ([]byte, error) {
			return []byte(`{"init": {"cmd": "ls", "output-ctrl": "on-error"}, "type":"static"}`), nil
		}
		linkCmd.WriteFile = func(filename string, data []byte, perm fs.FileMode) error {
			return nil
		}
		linkCmd.Rename = func(oldpath string, newpath string) error {
			return nil
		}
		linkCmd.Mkdir = func(path string, perm os.FileMode) error {
			return nil
		}
		linkCmd.GitPlainClone = func(path string, isBare bool, o *git.CloneOptions) (*git.Repository, error) {
			return &git.Repository{}, nil
		}
		linkCmd.ShouldConfigure = func(info *LinkInfo) (bool, error) {
			return true, nil
		}

		cmd := NewCobraCmd(linkCmd, f)

		cmd.SetArgs([]string{"--name", "SUUPA_DOOPA", "--preset", "static"})

		in := bytes.NewBufferString("yes\n")
		f.IOStreams.In = io.NopCloser(in)

		err := cmd.Execute()

		require.NoError(t, err)
		require.Contains(t, stdout.String(), fmt.Sprintf(msg.EdgeApplicationsLinkSuccessful+"\n", "SUUPA_DOOPA"))
	})
}
