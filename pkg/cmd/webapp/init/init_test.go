package init

import (
	"bytes"
	"errors"
	"io"
	"io/fs"
	"os"
	"strings"
	"testing"

	msg "github.com/aziontech/azion-cli/messages/webapp"
	"github.com/aziontech/azion-cli/pkg/httpmock"
	"github.com/aziontech/azion-cli/pkg/testutils"
	"github.com/aziontech/azion-cli/utils"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/storer"
	"github.com/go-git/go-git/v5/storage/memory"
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

		cmd.SetArgs([]string{"--name", "SUUPA_DOOPA", "--type", "javascript"})

		err := cmd.Execute()

		require.ErrorIs(t, err, msg.ErrorPackageJsonNotFound)
	})

	t.Run("with unsupported type", func(t *testing.T) {
		f, _, _ := testutils.NewFactory(nil)

		initCmd := NewInitCmd(f)

		cmd := NewCobraCmd(initCmd)

		cmd.SetArgs([]string{"--name", "BLEBLEBLE", "--type", "demeuamor"})

		err := cmd.Execute()

		require.ErrorIs(t, err, utils.ErrorUnsupportedType)
	})

	t.Run("with -y and -n flags", func(t *testing.T) {
		f, _, _ := testutils.NewFactory(nil)

		initCmd := NewInitCmd(f)

		cmd := NewCobraCmd(initCmd)

		cmd.SetArgs([]string{"--name", "BLEBLEBLE", "--type", "demeuamor", "-y", "-n"})

		err := cmd.Execute()

		require.ErrorIs(t, err, msg.ErrorYesAndNoOptions)
	})

	t.Run("success with javascript", func(t *testing.T) {
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
			return []byte(`{"init": {"cmd": "ls", "output-ctrl": "on-error"}, "type":"javascript"}`), nil
		}
		initCmd.WriteFile = func(filename string, data []byte, perm fs.FileMode) error {
			return nil
		}
		initCmd.Rename = func(oldpath string, newpath string) error {
			return nil
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

		cmd := NewCobraCmd(initCmd)

		cmd.SetArgs([]string{"--name", "SUUPA_DOOPA", "--type", "javascript"})

		in := bytes.NewBuffer(nil)
		in.WriteString("yes\n")
		f.IOStreams.In = io.NopCloser(in)

		err := cmd.Execute()
		// fmt.Println(err.Error())

		require.NoError(t, err)
		require.Contains(t, stdout.String(), `Template successfully fetched and configured`)
	})

	t.Run("success with javascript using flag -y", func(t *testing.T) {
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
			return []byte(`{"init": {"cmd": "ls", "output-ctrl": "on-error"},"type":"javascript"}`), nil
		}
		initCmd.WriteFile = func(filename string, data []byte, perm fs.FileMode) error {
			return nil
		}
		initCmd.Rename = func(oldpath string, newpath string) error {
			return nil
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

		cmd := NewCobraCmd(initCmd)

		cmd.SetArgs([]string{"--name", "SUUPA_DOOPA", "--type", "javascript", "-y"})

		err := cmd.Execute()

		require.NoError(t, err)
		require.Contains(t, stdout.String(), `Template successfully fetched and configured`)
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
			return []byte(`{"init": {"cmd": "ls", "output-ctrl": "on-error"},"type":"javascript"}`), nil
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

		cmd := NewCobraCmd(initCmd)

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

		initCmd := NewInitCmd(f)

		initCmd.LookPath = func(bin string) (string, error) {
			return "", nil
		}

		initCmd.CommandRunner = func(cmd string, envs []string) (string, int, error) {
			return "", 0, nil
		}
		initCmd.FileReader = func(path string) ([]byte, error) {
			return []byte(`{"init": {"cmd": "ls", "output-ctrl": "on-error"},"type":"javascript"}`), nil
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

		cmd := NewCobraCmd(initCmd)

		cmd.SetArgs([]string{"--name", "SUUPA_DOOPA", "--type", "javascript", "-n"})

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
			return []byte(`{"init": {"cmd": "notacmd", "output-ctrl": "disable"},"type":"javascript"}`), nil
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

		i := InitInfo{TypeLang: "javascript", PathWorkingDir: "."}
		err := cmd.runInitCmdLine(&i)
		require.ErrorIs(t, err, msg.ErrFailedToRunInitCommand)
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

	t.Run("success with Flareact", func(t *testing.T) {
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
			return []byte(`{"init": {"cmd": "ls", "env": "./azion/init.env", "output-ctrl": "on-error"}, "type": "flareact"}`), nil
		}

		cmd.EnvLoader = func(path string) ([]string, error) {
			return envs, nil
		}

		cmd.CommandRunner = func(cmd string, env []string) (string, int, error) {
			return "my command output", 0, nil
		}

		cmd.GitPlainClone = func(path string, isBare bool, o *git.CloneOptions) (*git.Repository, error) {
			return &git.Repository{}, nil
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

func Test_sortTag(t *testing.T) {
	r, _ := git.Clone(memory.NewStorage(), nil, &git.CloneOptions{URL: "https://github.com/MaxwelMazur/action-testing.git"})
	tags, _ := r.Tags()

	type args struct {
		tags   storer.ReferenceIter
		major  string
		branch string
	}
	tests := []struct {
		name    string
		args    args
		wantTag string
		wantErr bool
	}{
		{
			name: "branch main with major 0",
			args: args{
				tags:   tags,
				major:  "0",
				branch: "main",
			},
			wantTag: "refs/tags/v0.5.0",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotTag, err := sortTag(tt.args.tags, tt.args.major, tt.args.branch)
			if (err != nil) != tt.wantErr {
				t.Errorf("sortTag() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotTag != tt.wantTag {
				t.Errorf("sortTag() gotTag = %v, want %v", gotTag, tt.wantTag)
			}
		})
	}
}
