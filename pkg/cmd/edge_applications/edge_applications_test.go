package edge_applications

import (
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	buildcmd "github.com/aziontech/azion-cli/pkg/cmd/edge_applications/build"
	initcmd "github.com/aziontech/azion-cli/pkg/cmd/edge_applications/init"
	publishcmd "github.com/aziontech/azion-cli/pkg/cmd/edge_applications/publish"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/httpmock"

	"github.com/aziontech/azion-cli/pkg/testutils"
	"github.com/go-git/go-git/v5"
	"github.com/stretchr/testify/require"
)

func Mock(mock *httpmock.Registry) {
	mock.Register(
		httpmock.REST("POST", "edge_functions"),
		httpmock.JSONFromFile(".fixtures/edge_function.json"),
	)

	mock.Register(
		httpmock.REST("POST", "edge_applications"),
		httpmock.JSONFromFile(".fixtures/edge_application.json"),
	)

	mock.Register(
		httpmock.REST("PATCH", "edge_applications/777"),
		httpmock.JSONFromFile(".fixtures/edge_application.json"),
	)

	mock.Register(
		httpmock.REST("POST", "edge_applications/777/functions_instances"),
		httpmock.JSONFromFile(".fixtures/edge_application.json"),
	)

	mock.Register(
		httpmock.REST("GET", "edge_applications/777/rules_engine/request/rules"),
		httpmock.JSONFromFile(".fixtures/rules.json"),
	)

	mock.Register(
		httpmock.REST("PATCH", "edge_applications/777/rules_engine/request/rules/137056"),
		httpmock.JSONFromFile(".fixtures/rule.json"),
	)

	mock.Register(
		httpmock.REST("POST", "domains"),
		httpmock.JSONFromFile(".fixtures/domain.json"),
	)

	mock.Register(
		httpmock.REST("POST", ""),
		httpmock.JSONFromFile(".fixtures/domain.json"),
	)
}

func TestEdgeApplicationsCmd(t *testing.T) {
	t.Run("nextjs testing", func(t *testing.T) {
		mock := &httpmock.Registry{}
		Mock(mock)

		options := &contracts.AzionApplicationOptions{}

		f, _, _ := testutils.NewFactory(mock)
		envs := []string{"UEBA=OBA", "FAZER=UM_PENSO"}

		// INIT

		initCmd := initcmd.NewInitCmd(f)

		initCmd.LookPath = func(bin string) (string, error) {
			return "", nil
		}

		initCmd.LookPath = func(bin string) (string, error) {
			return "", nil
		}

		initCmd.WriteFile = func(filename string, data []byte, perm fs.FileMode) error {
			return nil
		}

		initCmd.Rename = func(oldpath string, newpath string) error {
			return nil
		}

		initCmd.Stat = func(path string) (fs.FileInfo, error) {
			return nil, nil
		}

		initCmd.FileReader = func(path string) ([]byte, error) {
			return []byte(`{"init": {"cmd": "ls", "env": "./azion/init.env", "output-ctrl": "on-error"}, "type": "nextjs"}`), nil
		}

		initCmd.EnvLoader = func(path string) ([]string, error) {
			return envs, nil
		}

		initCmd.CommandRunner = func(cmd string, env []string) (string, int, error) {
			return "my command output", 0, nil
		}

		initCmd.Mkdir = func(path string, perm os.FileMode) error {
			return nil
		}

		initCmd.GitPlainClone = func(path string, isBare bool, o *git.CloneOptions) (*git.Repository, error) {
			return &git.Repository{}, nil
		}

		cmdInit := initcmd.NewCobraCmd(initCmd)

		cmdInit.SetArgs([]string{"--name", "functional_testing_nextjs", "--type", "nextjs"})
		err := cmdInit.Execute()

		//BUILD

		buildCommand := buildcmd.NewBuildCmd(f)

		buildCommand.FileReader = func(path string) ([]byte, error) {
			return []byte(`{"build": {"cmd": "./azion/webdev.sh build", "env": "./azion/init.env", "output-ctrl": "on-error"}, "type": "nextjs"}`), nil
		}

		buildCommand.CommandRunner = func(cmd string, envs []string) (string, int, error) {
			return "Build completed", 0, nil
		}

		buildCommand.Stat = func(path string) (fs.FileInfo, error) {
			return nil, nil
		}

		buildCommand.EnvLoader = func(path string) ([]string, error) {
			return []string{}, nil
		}

		buildCommand.VersionId = func(dir string) (string, error) {
			return "123456789", nil
		}

		buildCommand.GetWorkDir = func() (string, error) {
			return "", nil
		}

		buildCommand.WriteFile = func(filename string, data []byte, perm fs.FileMode) error {
			return nil
		}

		cmdBuild := buildcmd.NewCobraCmd(buildCommand)

		errBuild := cmdBuild.Execute()

		//PUBLISH

		publishCmd := publishcmd.NewPublishCmd(f)

		publishCmd.Open = func(name string) (*os.File, error) {
			return nil, nil
		}

		publishCmd.FilepathWalk = func(root string, fn filepath.WalkFunc) error {
			return nil
		}

		publishCmd.FileReader = func(path string) ([]byte, error) {
			return []byte(`{"publish": {"pre_cmd": "./azion/webdev.sh publish", "env": "./azion/init.env", "output-ctrl": "on-error"}, "type": "nextjs"}`), nil
		}

		publishCmd.EnvLoader = func(path string) ([]string, error) {
			return []string{}, nil
		}

		publishCmd.CommandRunner = func(cmd string, env []string) (string, int, error) {
			return "", 0, nil
		}

		publishCmd.BuildCmd = func(f *cmdutil.Factory) *buildcmd.BuildCmd {
			return buildCommand
		}

		publishCmd.GetAzionJsonContent = func() (*contracts.AzionApplicationOptions, error) {
			return options, nil
		}

		publishCmd.WriteAzionJsonContent = func(conf *contracts.AzionApplicationOptions) error {
			return nil
		}

		cmdPublish := publishcmd.NewCobraCmd(publishCmd)

		errPublish := cmdPublish.Execute()

		require.NoError(t, err)

		require.NoError(t, errBuild)

		require.NoError(t, errPublish)

	})

	t.Run("flareact testing", func(t *testing.T) {
		mock := &httpmock.Registry{}
		Mock(mock)

		options := &contracts.AzionApplicationOptions{}

		f, _, _ := testutils.NewFactory(mock)
		envs := []string{"UEBA=OBA", "FAZER=UM_PENSO"}

		// INIT

		initCmd := initcmd.NewInitCmd(f)

		initCmd.LookPath = func(bin string) (string, error) {
			return "", nil
		}

		initCmd.WriteFile = func(filename string, data []byte, perm fs.FileMode) error {
			return nil
		}

		initCmd.Rename = func(oldpath string, newpath string) error {
			return nil
		}

		initCmd.Stat = func(path string) (fs.FileInfo, error) {
			return nil, nil
		}

		initCmd.FileReader = func(path string) ([]byte, error) {
			return []byte(`{"init": {"cmd": "ls", "env": "./azion/init.env", "output-ctrl": "on-error"}, "type": "flareact"}`), nil
		}

		initCmd.EnvLoader = func(path string) ([]string, error) {
			return envs, nil
		}

		initCmd.CommandRunner = func(cmd string, env []string) (string, int, error) {
			return "my command output", 0, nil
		}

		initCmd.Mkdir = func(path string, perm os.FileMode) error {
			return nil
		}

		initCmd.GitPlainClone = func(path string, isBare bool, o *git.CloneOptions) (*git.Repository, error) {
			return &git.Repository{}, nil
		}

		cmdInit := initcmd.NewCobraCmd(initCmd)

		cmdInit.SetArgs([]string{"--name", "functional_testing_flareact", "--type", "flareact"})
		err := cmdInit.Execute()

		//BUILD

		buildCommand := buildcmd.NewBuildCmd(f)

		buildCommand.FileReader = func(path string) ([]byte, error) {
			return []byte(`{"build": {"cmd": "./azion/webdev.sh build", "env": "./azion/init.env", "output-ctrl": "on-error"}, "type": "flareact"}`), nil
		}

		buildCommand.CommandRunner = func(cmd string, envs []string) (string, int, error) {
			return "Build completed", 0, nil
		}

		buildCommand.VersionId = func(dir string) (string, error) {
			return "123456789", nil
		}

		buildCommand.EnvLoader = func(path string) ([]string, error) {
			return []string{}, nil
		}

		buildCommand.Stat = func(path string) (fs.FileInfo, error) {
			return nil, nil
		}

		buildCommand.WriteFile = func(filename string, data []byte, perm fs.FileMode) error {
			return nil
		}

		buildCommand.GetWorkDir = func() (string, error) {
			return "", nil
		}

		cmdBuild := buildcmd.NewCobraCmd(buildCommand)

		errBuild := cmdBuild.Execute()

		//PUBLISH

		publishCmd := publishcmd.NewPublishCmd(f)

		publishCmd.Open = func(name string) (*os.File, error) {
			return nil, nil
		}

		publishCmd.FilepathWalk = func(root string, fn filepath.WalkFunc) error {
			return nil
		}

		publishCmd.FileReader = func(path string) ([]byte, error) {
			return []byte(`{"publish": {"pre_cmd": "./azion/webdev.sh publish", "env": "./azion/init.env", "output-ctrl": "on-error"}, "type": "flareact"}`), nil
		}

		publishCmd.EnvLoader = func(path string) ([]string, error) {
			return []string{}, nil
		}

		publishCmd.CommandRunner = func(cmd string, env []string) (string, int, error) {
			return "", 0, nil
		}

		publishCmd.BuildCmd = func(f *cmdutil.Factory) *buildcmd.BuildCmd {
			return buildCommand
		}

		publishCmd.GetAzionJsonContent = func() (*contracts.AzionApplicationOptions, error) {
			return options, nil
		}

		publishCmd.WriteAzionJsonContent = func(conf *contracts.AzionApplicationOptions) error {
			return nil
		}

		cmdPublish := publishcmd.NewCobraCmd(publishCmd)

		errPublish := cmdPublish.Execute()

		require.NoError(t, err)

		require.NoError(t, errBuild)

		require.NoError(t, errPublish)

	})
}

// TestNewCmd `go test -run TestNewCmd -v -cover`
func TestNewCmd(t *testing.T) {
	var build *buildcmd.BuildCmd

	mock := &httpmock.Registry{}
	Mock(mock)
	fMock, _, _ := testutils.NewFactory(mock)

	type args struct {
		factory *cmdutil.Factory
		init    func(f *cmdutil.Factory) *initcmd.InitCmd
		build   func(f *cmdutil.Factory) *buildcmd.BuildCmd
		publish func(f *cmdutil.Factory) *publishcmd.PublishCmd
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "flow success full, init, build, publish",
			args: args{
				factory: fMock,
				init: func(f *cmdutil.Factory) *initcmd.InitCmd {
					return &initcmd.InitCmd{
						Io:         f.IOStreams,
						GetWorkDir: func() (string, error) { return "", nil },
						FileReader: func(path string) ([]byte, error) {
							return []byte(`{"init": {"cmd": "ls", "env": "./azion/init.env", "output-ctrl": "on-error"}, "type": "nextjs"}`), nil
						},
						CommandRunner: func(cmd string, envvars []string) (string, int, error) { return "", 0, nil },
						LookPath:      func(bin string) (string, error) { return "", nil },
						IsDirEmpty:    func(dirpath string) (bool, error) { return true, nil },
						CleanDir:      func(dirpath string) error { return nil },
						WriteFile:     func(filename string, data []byte, perm fs.FileMode) error { return nil },
						OpenFile:      func(name string) (*os.File, error) { return nil, nil },
						RemoveAll:     func(path string) error { return nil },
						Rename:        func(oldpath, newpath string) error { return nil },
						CreateTempDir: func(dir, pattern string) (string, error) { return "", nil },
						EnvLoader:     func(path string) ([]string, error) { return []string{}, nil },
						Stat:          func(path string) (fs.FileInfo, error) { return nil, nil },
						Mkdir:         func(path string, perm os.FileMode) error { return nil },
						GitPlainClone: func(path string, isBare bool, o *git.CloneOptions) (*git.Repository, error) {
							return &git.Repository{}, nil
						},
					}
				},
				build: func(f *cmdutil.Factory) *buildcmd.BuildCmd {
					return &buildcmd.BuildCmd{
						Io: f.IOStreams,
						FileReader: func(path string) ([]byte, error) {
							return []byte(`{"build": {"cmd": "./azion/webdev.sh build", "env": "./azion/init.env", "output-ctrl": "on-error"}, "type": "nextjs"}`), nil
						},
						CommandRunner:      func(cmd string, envs []string) (string, int, error) { return "Build completed", 0, nil },
						ConfigRelativePath: "/azion/config.json",
						GetWorkDir:         func() (string, error) { return "", nil },
						EnvLoader:          func(path string) ([]string, error) { return []string{}, nil },
						WriteFile:          func(filename string, data []byte, perm fs.FileMode) error { return nil },
						Stat:               func(path string) (fs.FileInfo, error) { return nil, nil },
						VersionId:          func(dir string) (string, error) { return "123456789", nil },
					}
				},
				publish: func(f *cmdutil.Factory) *publishcmd.PublishCmd {
					return &publishcmd.PublishCmd{
						Io:                    f.IOStreams,
						Open:                  func(name string) (*os.File, error) { return nil, nil },
						GetWorkDir:            func() (string, error) { return "", nil },
						FilepathWalk:          func(root string, fn filepath.WalkFunc) error { return nil },
						EnvLoader:             func(path string) ([]string, error) { return []string{}, nil },
						CommandRunner:         func(cmd string, env []string) (string, int, error) { return "", 0, nil },
						BuildCmd:              func(f *cmdutil.Factory) *buildcmd.BuildCmd { return build },
						GetAzionJsonContent:   func() (*contracts.AzionApplicationOptions, error) { return &contracts.AzionApplicationOptions{}, nil },
						WriteAzionJsonContent: func(conf *contracts.AzionApplicationOptions) error { return nil },
						FileReader: func(path string) ([]byte, error) {
							return []byte(`{"publish": {"pre_cmd": "./azion/webdev.sh publish", "env": "./azion/init.env", "output-ctrl": "on-error"}, "type": "nextjs"}`), nil
						},
						WriteFile: func(filename string, data []byte, perm fs.FileMode) error { return nil },
						GetAzionJsonCdn: func() (*contracts.AzionApplicationCdn, error) {
							return &contracts.AzionApplicationCdn{}, nil
						},
						F: fMock,
					}
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if cmd := NewCmd(tt.args.factory); cmd != nil {
				init := tt.args.init(tt.args.factory)
				require.NoError(t, initcmd.NewCobraCmd(init).Execute())
				build = tt.args.build(tt.args.factory)
				require.NoError(t, buildcmd.NewCobraCmd(build).Execute())
                publish := tt.args.publish(tt.args.factory)
				require.NoError(t, publishcmd.NewCobraCmd(publish).Execute())
			}
		})
	}
}
