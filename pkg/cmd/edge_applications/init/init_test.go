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

	msg "github.com/aziontech/azion-cli/messages/edge_applications"
	"github.com/aziontech/azion-cli/pkg/httpmock"
	"github.com/aziontech/azion-cli/pkg/testutils"
	"github.com/aziontech/azion-cli/utils"
	"github.com/go-git/go-git/v5"

	"github.com/go-git/go-git/v5/plumbing"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCobraCmd(t *testing.T) {
	t.Run("without package.json", func(t *testing.T) {
		f, _, _ := testutils.NewFactory(nil)

		initCmd := NewInitCmd(f)

		initCmd.FileReader = func(path string) ([]byte, error) {
			return nil, os.ErrNotExist
		}

		cmd := NewCobraCmd(initCmd)

		cmd.SetArgs([]string{"--name", "SUUPA_DOOPA", "--type", "static", "-y"})
		require.NoError(t, cmd.Execute())
	})

	t.Run("with -y and -n flags", func(t *testing.T) {
		f, _, _ := testutils.NewFactory(nil)

		initCmd := NewInitCmd(f)

		cmd := NewCobraCmd(initCmd)

		cmd.SetArgs([]string{"--name", "BLEBLEBLE", "--type", "demeuamor", "-y", "-n"})

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

		initCmd.CommandRunner = func(cmd string, envs []string) (string, int, error) {
			return "", 0, nil
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

		cmd.SetArgs([]string{"--name", "SUUPA_DOOPA", "--type", "static"})

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

		initCmd.CommandRunner = func(cmd string, envs []string) (string, int, error) {
			return "", 0, nil
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

		cmd.SetArgs([]string{"--name", "SUUPA_DOOPA", "--type", "static"})

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

		initCmd.CommandRunner = func(cmd string, envs []string) (string, int, error) {
			return "", 0, nil
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

		cmd.SetArgs([]string{"--name", "SUUPA_DOOPA", "--type", "static", "-n"})

		err := cmd.Execute()

		require.NoError(t, err)
	})

	t.Run("invalid option", func(t *testing.T) {
		f, _, _ := testutils.NewFactory(nil)

		initCmd := NewInitCmd(f)

		cmd := NewCobraCmd(initCmd)

		initCmd.LookPath = func(bin string) (string, error) {
			return "", nil
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

		initCmd.FileReader = func(path string) ([]byte, error) {
			return []byte("{\n  \"name\": \"vanillajs-app\",\n  \"version\": \"1.0.0\",\n  \"main\": \"index.js\",\n  \"scripts\": {\n    \"test\": \"echo \\\"Error: no test specified\\\" && exit 1\",\n    \"build\": \"azioncli edge_applications build\",\n    \"deploy\": \"azioncli edge_applications publish\"\n  },\n  \"repository\": {\n    \"type\": \"git\",\n    \"url\": \"git+https://github.com/aziontech/azioncli-template.git\"\n  },\n  \"author\": \"\",\n  \"license\": \"ISC\",\n  \"bugs\": {\n    \"url\": \"https://github.com/aziontech/azioncli-template/issues\"\n  },\n  \"homepage\": \"https://github.com/aziontech/azioncli-template#readme\",\n  \"description\": \"\",\n  \"devDependencies\": {\n    \"clean-webpack-plugin\": \"^4.0.0\",\n    \"webpack-cli\": \"^4.9.2\",\n  \"next\": \"12.2.5\"\n  }\n}\n"), nil
		}

		initCmd.IsDirEmpty = func(dirpath string) (bool, error) {
			return false, nil
		}
		initCmd.Mkdir = func(path string, perm os.FileMode) error {
			return nil
		}

		cmd.SetArgs([]string{"--name", "SUUPA_DOOPA", "--type", "static"})

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

		cmd := NewInitCmd(f)

		cmd.LookPath = func(bin string) (string, error) {
			return "", nil
		}

		cmd.FileReader = func(path string) ([]byte, error) {
			return nil, os.ErrNotExist
		}
		cmd.GitPlainClone = func(path string, isBare bool, o *git.CloneOptions) (*git.Repository, error) {
			return &git.Repository{}, nil
		}
		cmd.Mkdir = func(path string, perm os.FileMode) error {
			return nil
		}

		i := InitInfo{}
		err := cmd.runInitCmdLine(&i)
		require.EqualError(t, err, "Failed to open the config.json file. The file doesn't exist, is corrupted, or has an invalid JSON format. Verify if the file format is JSON or fix its content according to the JSON format specification at https://www.json.org/json-en.html")
	})

	t.Run("init.env not empty", func(t *testing.T) {
		f, _, _ := testutils.NewFactory(nil)
		cmd := NewInitCmd(f)

		cmd.LookPath = func(bin string) (string, error) {
			return "", nil
		}

		// Specified init.env file but it cannot be read correctly
		cmd.FileReader = func(path string) ([]byte, error) {
			return []byte(`{"init": {"cmd": "ls", "env": "./azion/init.env", "output-ctrl": "on-error"},"type":"javascript"}`), nil
		}
		cmd.EnvLoader = func(path string) ([]string, error) {
			return nil, msg.ErrReadEnvFile
		}

		cmd.WriteFile = func(filename string, data []byte, perm fs.FileMode) error {
			return nil
		}
		cmd.GitPlainClone = func(path string, isBare bool, o *git.CloneOptions) (*git.Repository, error) {
			return &git.Repository{}, nil
		}
		cmd.Mkdir = func(path string, perm os.FileMode) error {
			return nil
		}

		i := InitInfo{}
		err := cmd.runInitCmdLine(&i)
		require.ErrorIs(t, err, msg.ErrReadEnvFile)
	})

	t.Run("Failed to run the command specified", func(t *testing.T) {
		f, _, _ := testutils.NewFactory(nil)

		cmd := NewInitCmd(f)

		cmd.LookPath = func(bin string) (string, error) {
			return "", nil
		}

		cmd.FileReader = func(path string) ([]byte, error) {
			return []byte(`{"init": {"cmd": "notacmd", "output-ctrl": "disable"},"type":"static"}`), nil
		}
		cmd.WriteFile = func(filename string, data []byte, perm fs.FileMode) error {
			return nil
		}
		cmd.EnvLoader = func(path string) ([]string, error) {
			return nil, nil
		}
		cmd.CommandRunner = func(cmd string, env []string) (string, int, error) {
			return "cmd", 0, errors.New("Failed to run the command specified in the template (config.json)")
		}
		cmd.Mkdir = func(path string, perm os.FileMode) error {
			return nil
		}

		i := InitInfo{TypeLang: "static", PathWorkingDir: "."}
		err := cmd.runInitCmdLine(&i)
		require.ErrorIs(t, err, utils.ErrorUnsupportedType)
	})

	t.Run("success with NextJS", func(t *testing.T) {
		f, _, _ := testutils.NewFactory(nil)

		envs := []string{"UEBA=OBA", "FAZER=UM_PENSO"}
		cmd := NewInitCmd(f)

		cmd.LookPath = func(bin string) (string, error) {
			return "", nil
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
			return []byte(`{"init": {"cmd": "ls", "env": "./azion/init.env", "output-ctrl": "on-error"}, "type": "nextjs"}`), nil
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

		cmd.GitPlainClone = func(path string, isBare bool, o *git.CloneOptions) (*git.Repository, error) {
			return &git.Repository{}, nil
		}

		i := InitInfo{
			TypeLang: "nextjs",
		}
		err := cmd.runInitCmdLine(&i)
		require.NoError(t, err)
	})
}

func Test_fetchTemplates(t *testing.T) {
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

func Test_formatTag(t *testing.T) {
	type args struct {
		tag string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "case branch dev",
			args: args{
				tag: "refs/tags/v0.1.0-dev.2",
			},
			want: "0102",
		},
		{
			name: "case branch main",
			args: args{
				tag: "refs/tags/v0.1.0",
			},
			want: "010",
		},
		{
			name: "case major not exist",
			args: args{
				tag: "refs/tags/v0.1.0",
			},
			want: "010",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := formatTag(tt.args.tag); got != tt.want {
				t.Errorf("formatTag() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_checkBranch(t *testing.T) {
	type args struct {
		num    string
		branch string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "branch dev, with tag dev",
			args: args{
				num:    "1234",
				branch: "dev",
			},
			want: "1234",
		},
		{
			name: "branch any, with tag dev",
			args: args{
				num:    "1234",
				branch: "any",
			},
			want: "",
		},
		{
			name: "branch any, with tag any",
			args: args{
				num:    "123",
				branch: "any",
			},
			want: "123",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := checkBranch(tt.args.num, tt.args.branch); got != tt.want {
				t.Errorf("checkBranch() = %v, want %v", got, tt.want)
			}
		})
	}
}

type mockReferenceIter struct {
	mock.Mock
}

func (m *mockReferenceIter) ForEach(f func(*plumbing.Reference) error) error {
	args := m.Called(f)
	return args.Error(0)
}

func Test_SortTag(t *testing.T) {
	mockIter := new(mockReferenceIter)

	// defines the mock action
	mockIter.On("ForEach", mock.Anything).Run(func(args mock.Arguments) {
		f := args.Get(0).(func(*plumbing.Reference) error)
		refs := []*plumbing.Reference{
			plumbing.NewReferenceFromStrings("refs/tags/v2.0.0", "beefdead"),
			plumbing.NewReferenceFromStrings("refs/tags/v1.0.0", "deadbeef"),
			plumbing.NewReferenceFromStrings("refs/tags/v3.0.0", "deafbeef"),
		}

		for _, ref := range refs {
			if err := f(ref); err != nil {
				panic(err)
			}
		}
	}).Return(nil)

	// call func sortTag with the mock
	result, err := sortTag(mockIter, "2", "main")
	require.NoError(t, err)

	expected := "refs/tags/v2.0.0"
	require.Equal(t, expected, result)
}
