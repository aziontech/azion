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

	"github.com/go-git/go-git/v5/plumbing"
	"github.com/stretchr/testify/mock"
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

	t.Run("does not overwrite contents", func(t *testing.T) {
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
		initCmd.IsDirEmpty = func(dirpath string) (bool, error) {
			return false, nil
		}
		initCmd.Mkdir = func(path string, perm os.FileMode) error {
			return nil
		}

		cmd := NewCobraCmd(initCmd)

		cmd.SetArgs([]string{"--name", "SUUPA_DOOPA", "--template", "static"})

		in := bytes.NewBuffer(nil)
		in.WriteString("no\n")
		f.IOStreams.In = io.NopCloser(in)

		err := cmd.Execute()

		require.NoError(t, err)
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

func TestInitCmd(t *testing.T) {
	logger.New(zapcore.DebugLevel)

}

func Test_fetchTemplates(t *testing.T) {
	logger.New(zapcore.DebugLevel)

	t.Run("tests flow full template", func(t *testing.T) {
		f, _, _ := testutils.NewFactory(nil)
		cmd := NewInitCmd(f)

		cmd.CreateTempDir = func(dir string, pattern string) (string, error) {
			return "", nil
		}

		cmd.GitPlainClone = func(path string, isBare bool, o *git.CloneOptions) (*git.Repository, error) {
			return nil, nil
		}

		cmd.Rename = func(oldpath string, newpath string) error {
			return nil
		}
		cmd.Mkdir = func(path string, perm os.FileMode) error {
			return nil
		}

		err := cmd.fetchTemplates(&InitInfo{})
		require.NoError(t, err)
	})
}

type mockReferenceIter struct {
	mock.Mock
}

func (m *mockReferenceIter) ForEach(f func(*plumbing.Reference) error) error {
	args := m.Called(f)
	return args.Error(0)
}

func Test_SortTag(t *testing.T) {
	logger.New(zapcore.DebugLevel)

	mockIter := new(mockReferenceIter)

	// defines the mock action
	mockIter.On("ForEach", mock.Anything).Run(func(args mock.Arguments) {
		f := args.Get(0).(func(*plumbing.Reference) error)
		refs := []*plumbing.Reference{
			plumbing.NewReferenceFromStrings("refs/tags/v0.1.0", "beefdead"),
			plumbing.NewReferenceFromStrings("refs/tags/v0.2.0", "deadbeef"),
			plumbing.NewReferenceFromStrings("refs/tags/v0.10.5", "deafbeef"),
			plumbing.NewReferenceFromStrings("refs/tags/v0.5.0", "deafbeef"),
			plumbing.NewReferenceFromStrings("refs/tags/v0.10.0", "deafbeef"),
			plumbing.NewReferenceFromStrings("refs/tags/v1.10.0", "deafbeef"),
			plumbing.NewReferenceFromStrings("refs/tags/v0.10.1", "deafbeef"),
			plumbing.NewReferenceFromStrings("refs/tags/v0.10.2", "deafbeef"),
			plumbing.NewReferenceFromStrings("refs/tags/v0.11.0-dev.13", "deafbeef"),
		}

		for _, ref := range refs {
			if err := f(ref); err != nil {
				panic(err)
			}
		}
	}).Return(nil)

	result, err := sortTag(mockIter, TemplateMajor)
	require.NoError(t, err)

	expected := "refs/tags/v0.10.5"
	require.Equal(t, expected, result)
}
