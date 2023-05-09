package edge_applications

import (
	"io/fs"
	"log"
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

// TestNewCmd `go test -run TestNewCmd -v -cover`
func TestNewCmd(t *testing.T) {
	type args struct {
		init    func(f *cmdutil.Factory) *initcmd.InitCmd
		build   func(f *cmdutil.Factory) *buildcmd.BuildCmd
		publish func(f *cmdutil.Factory) *publishcmd.PublishCmd
	}
	tests := []struct {
		name string
		args args
	}{
		{
			// go test -run TestNewCmd/flow_success_full_init,_build_and_publish -v -cover
			name: "flow success full init, build and publish",
			args: args{
				init: func(f *cmdutil.Factory) *initcmd.InitCmd {
					return &initcmd.InitCmd{
						Io:         f.IOStreams,
						GetWorkDir: func() (string, error) { return "", nil },
						FileReader: func(path string) ([]byte, error) {
							return []byte(`{"init": {"cmd": "ls", "env": "./azion/init.env", "output-ctrl": "on-error"}, "type": "nextjs" , "dependencies": { "next": "12.2.5" }}`), nil
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
							return []byte(`{"build": {"cmd": "./azion/webdev.sh build", "env": "./azion/init.env", "output-ctrl": "on-error"}, "type": "nextjs" , "dependencies": { "next": "12.2.5" }}`), nil
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
						GetAzionJsonContent:   func() (*contracts.AzionApplicationOptions, error) { return &contracts.AzionApplicationOptions{}, nil },
						WriteAzionJsonContent: func(conf *contracts.AzionApplicationOptions) error { return nil },
						FileReader: func(path string) ([]byte, error) {
							return []byte(`{"publish": {"pre_cmd": "./azion/webdev.sh publish", "env": "./azion/init.env", "output-ctrl": "on-error"}, "type": "nextjs" , "dependencies": { "next": "12.2.5" }}`), nil
						},
						WriteFile: func(filename string, data []byte, perm fs.FileMode) error { return nil },
						GetAzionJsonCdn: func() (*contracts.AzionApplicationCdn, error) {
							return &contracts.AzionApplicationCdn{}, nil
						},
					}
				},
			},
		},
		{
			// go test -run TestNewCmd/flow_success_full_nextjs_init,_build_and_publih -v -cover
			name: "flow success full nextjs init, build and publih",
			args: args{
				init: func(f *cmdutil.Factory) *initcmd.InitCmd {
					return &initcmd.InitCmd{
						Io:         f.IOStreams,
						GetWorkDir: func() (string, error) { return "", nil },
						FileReader: func(path string) ([]byte, error) {
							return []byte(`{"init": {"cmd": "ls", "env": "./azion/init.env", "output-ctrl": "on-error"}, "type": "nextjs" , "dependencies": { "next": "12.2.5" }}`), nil
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
							return []byte(`{"build": {"cmd": "./azion/webdev.sh build", "env": "./azion/init.env", "output-ctrl": "on-error"}, "type": "nextjs" , "dependencies": { "next": "12.2.5" }}`), nil
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
						GetAzionJsonContent:   func() (*contracts.AzionApplicationOptions, error) { return &contracts.AzionApplicationOptions{}, nil },
						WriteAzionJsonContent: func(conf *contracts.AzionApplicationOptions) error { return nil },
						FileReader: func(path string) ([]byte, error) {
							return []byte(`{"publish": {"pre_cmd": "./azion/webdev.sh publish", "env": "./azion/init.env", "output-ctrl": "on-error"}, "type": "nextjs" , "dependencies": { "next": "12.2.5" }}`), nil
						},
						WriteFile: func(filename string, data []byte, perm fs.FileMode) error { return nil },
						GetAzionJsonCdn: func() (*contracts.AzionApplicationCdn, error) {
							return &contracts.AzionApplicationCdn{}, nil
						},
					}
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &httpmock.Registry{}
			Mock(mock)
			fMock, _, _ := testutils.NewFactory(mock)

			if cmd := NewCmd(fMock); cmd != nil {
				init := tt.args.init(fMock)
				errInit := initcmd.NewCobraCmd(init).Execute()
				if errInit != nil {
					log.Fatal(errInit)
					return
				}

				build := tt.args.build(fMock)
				errBuild := buildcmd.NewCobraCmd(build).Execute()
				if errBuild != nil {
					log.Fatal(errBuild)
					return
				}

				publish := tt.args.publish(fMock)
				publish.BuildCmd = func(f *cmdutil.Factory) *buildcmd.BuildCmd {
					return build
				}
				publish.F = fMock

				errPublish := publishcmd.NewCobraCmd(publish).Execute()
				if errPublish != nil {
					log.Fatal(errPublish)
					return
				}
			}
		})
	}
}
