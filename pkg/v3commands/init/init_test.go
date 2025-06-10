package init

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/fs"
	"net/http"
	"os"
	"testing"

	"github.com/AlecAivazis/survey/v2"
	"github.com/aziontech/azion-cli/pkg/cmd/deploy"
	"github.com/aziontech/azion-cli/pkg/cmd/dev"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/config"
	"github.com/aziontech/azion-cli/pkg/github"
	"github.com/aziontech/azion-cli/pkg/httpmock"
	"github.com/aziontech/azion-cli/pkg/iostreams"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/testutils"
	"github.com/go-git/go-git/v5"
	"github.com/spf13/cobra"
	"go.uber.org/zap/zapcore"
)

var infoJsonData = `{"preset":"astro"}`

func TestNewCmd(t *testing.T) {
	mock := &httpmock.Registry{}
	f, _, _ := testutils.NewFactory(mock)
	NewCmd(f)
}

var cloneOptions git.CloneOptions

func Test_initCmd_Run(t *testing.T) {
	logger.New(zapcore.DebugLevel)

	type fields struct {
		name                  string
		preset                string
		auto                  bool
		packageManager        string
		pathWorkingDir        string
		globalFlagAll         bool
		f                     *cmdutil.Factory
		git                   github.Github
		getWorkDir            func() (string, error)
		fileReader            func(path string) ([]byte, error)
		isDirEmpty            func(dirpath string) (bool, error)
		cleanDir              func(dirpath string) error
		writeFile             func(filename string, data []byte, perm fs.FileMode) error
		openFile              func(name string) (*os.File, error)
		removeAll             func(path string) error
		rename                func(oldpath string, newpath string) error
		envLoader             func(path string) ([]string, error)
		stat                  func(path string) (fs.FileInfo, error)
		mkdir                 func(path string, perm os.FileMode) error
		gitPlainClone         func(path string, isBare bool, o *git.CloneOptions) (*git.Repository, error)
		commandRunner         func(envVars []string, comm string) (string, int, error)
		commandRunnerOutput   func(f *cmdutil.Factory, comm string, envVars []string) (string, error)
		commandRunInteractive func(f *cmdutil.Factory, comm string) error
		deployCmd             func(f *cmdutil.Factory) *deploy.DeployCmd
		devCmd                func(f *cmdutil.Factory) *dev.DevCmd
		changeDir             func(dir string) error
		askOne                func(p survey.Prompt, response interface{}, opts ...survey.AskOpt) error
		load                  func(filenames ...string) (err error)
		dir                   func() config.DirPath
		mkdirTemp             func(dir, pattern string) (string, error)
		readAll               func(r io.Reader) ([]byte, error)
		get                   func(url string) (resp *http.Response, err error)
		marshalIndent         func(v any, prefix, indent string) ([]byte, error)
		unmarshal             func(data []byte, v any) error
	}

	type args struct {
		c   *cobra.Command
		in1 []string
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "success flow",
			fields: fields{
				auto: true,
				fileReader: func(filename string) ([]byte, error) {
					return []byte(infoJsonData), nil
				},
				f: &cmdutil.Factory{
					Flags: cmdutil.Flags{
						GlobalFlagAll: false,
						Format:        "",
						Out:           "",
						NoColor:       false,
					},
					IOStreams: iostreams.System(),
				},
				globalFlagAll: false,
				name:          "project-piece",
				preset:        "vite",
				getWorkDir: func() (string, error) {
					return "/path/full", nil
				},
				get: func(url string) (resp *http.Response, err error) {
					b, err := os.ReadFile("./.fixtures/project_samples.json")
					if err != nil {
						return nil, err
					}

					responseBody := io.NopCloser(bytes.NewReader(b))
					resp = &http.Response{
						StatusCode: 200,
						Body:       responseBody,
						Header:     make(http.Header),
					}
					return resp, nil
				},
				readAll:   io.ReadAll,
				unmarshal: json.Unmarshal,
				askOne: func(p survey.Prompt,
					response interface{},
					opts ...survey.AskOpt,
				) error {
					return nil
				},
				dir: config.Dir,
				mkdirTemp: func(dir, pattern string) (string, error) {
					return "", nil
				},
				removeAll: os.RemoveAll,
				rename: func(oldpath, newpath string) error {
					return nil
				},
				mkdir:         func(path string, perm os.FileMode) error { return nil },
				marshalIndent: json.MarshalIndent,
				writeFile: func(filename string, data []byte, perm fs.FileMode) error {
					return nil
				},
				changeDir: func(dir string) error { return nil },
				commandRunnerOutput: func(f *cmdutil.Factory, comm string, envVars []string) (string, error) {
					return "3.2.1", nil
				},
				commandRunInteractive: func(f *cmdutil.Factory, comm string) error {
					return nil
				},
				load: func(filenames ...string) (err error) {
					os.Setenv("preset", "vite")
					return nil
				},
				git: github.Github{
					Clone: func(cloneOptions *git.CloneOptions, url, path string) error {
						return nil
					},
				},
			},
		},
		{
			name: "error getWorkDir",
			fields: fields{
				fileReader: func(filename string) ([]byte, error) {
					return []byte(infoJsonData), nil
				},
				auto: true,
				f: &cmdutil.Factory{
					Flags: cmdutil.Flags{
						GlobalFlagAll: false,
						Format:        "",
						Out:           "",
						NoColor:       false,
					},
					IOStreams: iostreams.System(),
				},
				globalFlagAll: false,
				name:          "project-piece",
				preset:        "vite",
				getWorkDir: func() (string, error) {
					return "/path/full", errors.New("error getWorkDir")
				},
				get: func(url string) (resp *http.Response, err error) {
					b, err := os.ReadFile("./.fixtures/project_samples.json")
					if err != nil {
						return nil, err
					}

					responseBody := io.NopCloser(bytes.NewReader(b))
					resp = &http.Response{
						StatusCode: 200,
						Body:       responseBody,
						Header:     make(http.Header),
					}
					return resp, nil
				},
				readAll:   io.ReadAll,
				unmarshal: json.Unmarshal,
				askOne: func(p survey.Prompt,
					response interface{},
					opts ...survey.AskOpt,
				) error {
					return nil
				},
				dir: config.Dir,
				mkdirTemp: func(dir, pattern string) (string, error) {
					return "", nil
				},
				removeAll: os.RemoveAll,
				rename: func(oldpath, newpath string) error {
					return nil
				},
				mkdir:         func(path string, perm os.FileMode) error { return nil },
				marshalIndent: json.MarshalIndent,
				writeFile: func(filename string, data []byte, perm fs.FileMode) error {
					return nil
				},
				changeDir: func(dir string) error { return nil },
				commandRunnerOutput: func(f *cmdutil.Factory, comm string, envVars []string) (string, error) {
					return "3.2.1", nil
				},
				commandRunInteractive: func(f *cmdutil.Factory, comm string) error {
					return nil
				},
				load: func(filenames ...string) (err error) {
					os.Setenv("preset", "vite")
					return nil
				},
				git: github.Github{
					Clone: func(cloneOptions *git.CloneOptions, url, path string) error {
						return nil
					},
				},
			},
			wantErr: true,
		},
		{
			name: "error http.Get",
			fields: fields{
				fileReader: func(filename string) ([]byte, error) {
					return []byte(infoJsonData), nil
				},
				auto: true,
				f: &cmdutil.Factory{
					Flags: cmdutil.Flags{
						GlobalFlagAll: false,
						Format:        "",
						Out:           "",
						NoColor:       false,
					},
					IOStreams: iostreams.System(),
				},
				globalFlagAll: true,
				name:          "",
				preset:        "vite",
				getWorkDir: func() (string, error) {
					return "/path/full", nil
				},
				get: func(url string) (resp *http.Response, err error) {
					b, err := os.ReadFile("./.fixtures/project_samples.json")
					if err != nil {
						return nil, err
					}

					responseBody := io.NopCloser(bytes.NewReader(b))
					resp = &http.Response{
						StatusCode: 200,
						Body:       responseBody,
						Header:     make(http.Header),
					}
					return resp, errors.New("error http.Get")
				},
				readAll:   io.ReadAll,
				unmarshal: json.Unmarshal,
				askOne: func(p survey.Prompt,
					response interface{},
					opts ...survey.AskOpt,
				) error {
					return nil
				},
				dir: config.Dir,
				mkdirTemp: func(dir, pattern string) (string, error) {
					return "", nil
				},
				removeAll: os.RemoveAll,
				rename: func(oldpath, newpath string) error {
					return nil
				},
				mkdir:         func(path string, perm os.FileMode) error { return nil },
				marshalIndent: json.MarshalIndent,
				writeFile: func(filename string, data []byte, perm fs.FileMode) error {
					return nil
				},
				changeDir: func(dir string) error { return nil },
				commandRunnerOutput: func(f *cmdutil.Factory, comm string, envVars []string) (string, error) {
					return "3.2.1", nil
				},
				commandRunInteractive: func(f *cmdutil.Factory, comm string) error {
					return nil
				},
				load: func(filenames ...string) (err error) {
					os.Setenv("preset", "vite")
					return nil
				},
				git: github.Github{
					Clone: func(cloneOptions *git.CloneOptions, url, path string) error {
						return nil
					},
				},
			},
			wantErr: true,
		},
		{
			name: "error askForInput",
			fields: fields{
				fileReader: func(filename string) ([]byte, error) {
					return []byte(infoJsonData), nil
				},
				auto: true,
				f: &cmdutil.Factory{
					Flags: cmdutil.Flags{
						GlobalFlagAll: false,
						Format:        "",
						Out:           "",
						NoColor:       false,
					},
					IOStreams: iostreams.System(),
				},
				globalFlagAll: false,
				name:          "",
				preset:        "vite",
				getWorkDir: func() (string, error) {
					return "/path/full", nil
				},
				get: func(url string) (resp *http.Response, err error) {
					b, err := os.ReadFile("./.fixtures/project_samples.json")
					if err != nil {
						return nil, err
					}

					responseBody := io.NopCloser(bytes.NewReader(b))
					resp = &http.Response{
						StatusCode: 200,
						Body:       responseBody,
						Header:     make(http.Header),
					}
					return resp, nil
				},
				readAll:   io.ReadAll,
				unmarshal: json.Unmarshal,
				askOne: func(p survey.Prompt,
					response interface{},
					opts ...survey.AskOpt,
				) error {
					return errors.New("error askOne")
				},
				dir: config.Dir,
				mkdirTemp: func(dir, pattern string) (string, error) {
					return "", nil
				},
				removeAll: os.RemoveAll,
				rename: func(oldpath, newpath string) error {
					return nil
				},
				mkdir:         func(path string, perm os.FileMode) error { return nil },
				marshalIndent: json.MarshalIndent,
				writeFile: func(filename string, data []byte, perm fs.FileMode) error {
					return nil
				},
				changeDir: func(dir string) error { return nil },
				commandRunnerOutput: func(f *cmdutil.Factory, comm string, envVars []string) (string, error) {
					return "3.2.1", nil
				},
				commandRunInteractive: func(f *cmdutil.Factory, comm string) error {
					return nil
				},
				load: func(filenames ...string) (err error) {
					os.Setenv("preset", "vite")
					return nil
				},
				git: github.Github{
					Clone: func(cloneOptions *git.CloneOptions, url, path string) error {
						return nil
					},
				},
			},
			wantErr: true,
		},
		{
			name: "error expected status OK",
			fields: fields{
				fileReader: func(filename string) ([]byte, error) {
					return []byte(infoJsonData), nil
				},
				auto: true,
				f: &cmdutil.Factory{
					Flags: cmdutil.Flags{
						GlobalFlagAll: false,
						Format:        "",
						Out:           "",
						NoColor:       false,
					},
					IOStreams: iostreams.System(),
				},
				globalFlagAll: false,
				name:          "",
				preset:        "vite",
				getWorkDir: func() (string, error) {
					return "/path/full", nil
				},
				get: func(url string) (resp *http.Response, err error) {
					b, err := os.ReadFile("./.fixtures/project_samples.json")
					if err != nil {
						return nil, err
					}

					responseBody := io.NopCloser(bytes.NewReader(b))
					resp = &http.Response{
						StatusCode: 201,
						Body:       responseBody,
						Header:     make(http.Header),
					}
					return resp, nil
				},
				readAll:   io.ReadAll,
				unmarshal: json.Unmarshal,
				askOne: func(p survey.Prompt,
					response interface{},
					opts ...survey.AskOpt,
				) error {
					return nil
				},
				dir: config.Dir,
				mkdirTemp: func(dir, pattern string) (string, error) {
					return "", nil
				},
				removeAll: os.RemoveAll,
				rename: func(oldpath, newpath string) error {
					return nil
				},
				mkdir:         func(path string, perm os.FileMode) error { return nil },
				marshalIndent: json.MarshalIndent,
				writeFile: func(filename string, data []byte, perm fs.FileMode) error {
					return nil
				},
				changeDir: func(dir string) error { return nil },
				commandRunnerOutput: func(f *cmdutil.Factory, comm string, envVars []string) (string, error) {
					return "3.2.1", nil
				},
				commandRunInteractive: func(f *cmdutil.Factory, comm string) error {
					return nil
				},
				load: func(filenames ...string) (err error) {
					os.Setenv("preset", "vite")
					return nil
				},
				git: github.Github{
					Clone: func(cloneOptions *git.CloneOptions, url, path string) error {
						return nil
					},
				},
			},
			wantErr: true,
		},
		{
			name: "error readAll",
			fields: fields{
				fileReader: func(filename string) ([]byte, error) {
					return nil, nil
				},
				auto: true,
				f: &cmdutil.Factory{
					Flags: cmdutil.Flags{
						GlobalFlagAll: false,
						Format:        "",
						Out:           "",
						NoColor:       false,
					},
					IOStreams: iostreams.System(),
				},
				globalFlagAll: false,
				name:          "",
				preset:        "vite",
				getWorkDir: func() (string, error) {
					return "/path/full", nil
				},
				get: func(url string) (resp *http.Response, err error) {
					b, err := os.ReadFile("./.fixtures/project_samples.json")
					if err != nil {
						return nil, err
					}

					responseBody := io.NopCloser(bytes.NewReader(b))
					resp = &http.Response{
						StatusCode: 200,
						Body:       responseBody,
						Header:     make(http.Header),
					}
					return resp, nil
				},
				readAll: func(r io.Reader) ([]byte, error) {
					return []byte(""), errors.New("error readAll")
				},
				unmarshal: json.Unmarshal,
				askOne: func(p survey.Prompt,
					response interface{},
					opts ...survey.AskOpt,
				) error {
					return nil
				},
				dir: config.Dir,
				mkdirTemp: func(dir, pattern string) (string, error) {
					return "", nil
				},
				removeAll: os.RemoveAll,
				rename: func(oldpath, newpath string) error {
					return nil
				},
				mkdir:         func(path string, perm os.FileMode) error { return nil },
				marshalIndent: json.MarshalIndent,
				writeFile: func(filename string, data []byte, perm fs.FileMode) error {
					return nil
				},
				changeDir: func(dir string) error { return nil },
				commandRunnerOutput: func(f *cmdutil.Factory, comm string, envVars []string) (string, error) {
					return "3.2.1", nil
				},
				commandRunInteractive: func(f *cmdutil.Factory, comm string) error {
					return nil
				},
				load: func(filenames ...string) (err error) {
					os.Setenv("preset", "vite")
					return nil
				},
				git: github.Github{
					Clone: func(cloneOptions *git.CloneOptions, url, path string) error {
						return nil
					},
				},
			},
			wantErr: true,
		},
		{
			name: "error unmarshal",
			fields: fields{
				fileReader: func(filename string) ([]byte, error) {
					return nil, nil
				},
				auto: true,
				f: &cmdutil.Factory{
					Flags: cmdutil.Flags{
						GlobalFlagAll: false,
						Format:        "",
						Out:           "",
						NoColor:       false,
					},
					IOStreams: iostreams.System(),
				},
				globalFlagAll: false,
				name:          "project-piece",
				preset:        "vite",
				getWorkDir: func() (string, error) {
					return "/path/full", nil
				},
				get: func(url string) (resp *http.Response, err error) {
					b, err := os.ReadFile("./.fixtures/project_samples.json")
					if err != nil {
						return nil, err
					}

					responseBody := io.NopCloser(bytes.NewReader(b))
					resp = &http.Response{
						StatusCode: 200,
						Body:       responseBody,
						Header:     make(http.Header),
					}
					return resp, nil
				},
				readAll: io.ReadAll,
				unmarshal: func(data []byte, v any) error {
					return errors.New("error unmarshal")
				},
				askOne: func(p survey.Prompt,
					response interface{},
					opts ...survey.AskOpt,
				) error {
					return errors.New("")
				},
				dir: config.Dir,
				mkdirTemp: func(dir, pattern string) (string, error) {
					return "", nil
				},
				removeAll: os.RemoveAll,
				rename: func(oldpath, newpath string) error {
					return nil
				},
				mkdir:         func(path string, perm os.FileMode) error { return nil },
				marshalIndent: json.MarshalIndent,
				writeFile: func(filename string, data []byte, perm fs.FileMode) error {
					return nil
				},
				changeDir: func(dir string) error { return nil },
				commandRunnerOutput: func(f *cmdutil.Factory, comm string, envVars []string) (string, error) {
					return "3.2.1", nil
				},
				commandRunInteractive: func(f *cmdutil.Factory, comm string) error {
					return nil
				},
				load: func(filenames ...string) (err error) {
					os.Setenv("preset", "vite")
					return nil
				},
				git: github.Github{
					Clone: func(cloneOptions *git.CloneOptions, url, path string) error {
						return nil
					},
				},
			},
			wantErr: true,
		},
		{
			name: "error askOne",
			fields: fields{
				fileReader: func(filename string) ([]byte, error) {
					return nil, nil
				},
				auto: true,
				f: &cmdutil.Factory{
					Flags: cmdutil.Flags{
						GlobalFlagAll: false,
						Format:        "",
						Out:           "",
						NoColor:       false,
					},
					IOStreams: iostreams.System(),
				},
				globalFlagAll: false,
				name:          "",
				preset:        "vite",
				getWorkDir: func() (string, error) {
					return "/path/full", nil
				},
				get: func(url string) (resp *http.Response, err error) {
					b, err := os.ReadFile("./.fixtures/project_samples.json")
					if err != nil {
						return nil, err
					}

					responseBody := io.NopCloser(bytes.NewReader(b))
					resp = &http.Response{
						StatusCode: 200,
						Body:       responseBody,
						Header:     make(http.Header),
					}
					return resp, nil
				},
				readAll:   io.ReadAll,
				unmarshal: json.Unmarshal,
				askOne: func(p survey.Prompt,
					response interface{},
					opts ...survey.AskOpt,
				) error {
					return errors.New("error askOne")
				},
				dir: config.Dir,
				mkdirTemp: func(dir, pattern string) (string, error) {
					return "", nil
				},
				removeAll: os.RemoveAll,
				rename: func(oldpath, newpath string) error {
					return nil
				},
				mkdir:         func(path string, perm os.FileMode) error { return nil },
				marshalIndent: json.MarshalIndent,
				writeFile: func(filename string, data []byte, perm fs.FileMode) error {
					return nil
				},
				changeDir: func(dir string) error { return nil },
				commandRunnerOutput: func(f *cmdutil.Factory, comm string, envVars []string) (string, error) {
					return "3.2.1", nil
				},
				commandRunInteractive: func(f *cmdutil.Factory, comm string) error {
					return nil
				},
				load: func(filenames ...string) (err error) {
					os.Setenv("preset", "vite")
					return nil
				},
				git: github.Github{
					Clone: func(cloneOptions *git.CloneOptions, url, path string) error {
						return nil
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &initCmd{
				name:                  tt.fields.name,
				preset:                tt.fields.preset,
				auto:                  tt.fields.auto,
				packageManager:        tt.fields.packageManager,
				pathWorkingDir:        tt.fields.pathWorkingDir,
				globalFlagAll:         tt.fields.globalFlagAll,
				f:                     tt.fields.f,
				getWorkDir:            tt.fields.getWorkDir,
				fileReader:            tt.fields.fileReader,
				isDirEmpty:            tt.fields.isDirEmpty,
				cleanDir:              tt.fields.cleanDir,
				writeFile:             tt.fields.writeFile,
				openFile:              tt.fields.openFile,
				removeAll:             tt.fields.removeAll,
				rename:                tt.fields.rename,
				envLoader:             tt.fields.envLoader,
				stat:                  tt.fields.stat,
				mkdir:                 tt.fields.mkdir,
				gitPlainClone:         tt.fields.gitPlainClone,
				commandRunner:         tt.fields.commandRunner,
				commandRunnerOutput:   tt.fields.commandRunnerOutput,
				commandRunInteractive: tt.fields.commandRunInteractive,
				deployCmd:             tt.fields.deployCmd,
				devCmd:                tt.fields.devCmd,
				changeDir:             tt.fields.changeDir,
				askOne:                tt.fields.askOne,
				load:                  tt.fields.load,
				dir:                   tt.fields.dir,
				mkdirTemp:             tt.fields.mkdirTemp,
				readAll:               tt.fields.readAll,
				get:                   tt.fields.get,
				marshalIndent:         tt.fields.marshalIndent,
				unmarshal:             tt.fields.unmarshal,
				git:                   tt.fields.git,
			}
			if err := cmd.Run(tt.args.c, tt.args.in1); (err != nil) != tt.wantErr {
				t.Errorf("initCmd.Run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_initCmd_deps(t *testing.T) {
	logger.New(zapcore.DebugLevel)

	type fields struct {
		name                  string
		preset                string
		auto                  bool
		packageManager        string
		pathWorkingDir        string
		globalFlagAll         bool
		f                     *cmdutil.Factory
		git                   github.Github
		getWorkDir            func() (string, error)
		fileReader            func(path string) ([]byte, error)
		isDirEmpty            func(dirpath string) (bool, error)
		cleanDir              func(dirpath string) error
		writeFile             func(filename string, data []byte, perm fs.FileMode) error
		openFile              func(name string) (*os.File, error)
		removeAll             func(path string) error
		rename                func(oldpath string, newpath string) error
		envLoader             func(path string) ([]string, error)
		stat                  func(path string) (fs.FileInfo, error)
		mkdir                 func(path string, perm os.FileMode) error
		gitPlainClone         func(path string, isBare bool, o *git.CloneOptions) (*git.Repository, error)
		commandRunner         func(envVars []string, comm string) (string, int, error)
		commandRunnerOutput   func(f *cmdutil.Factory, comm string, envVars []string) (string, error)
		commandRunInteractive func(f *cmdutil.Factory, comm string) error
		deployCmd             func(f *cmdutil.Factory) *deploy.DeployCmd
		devCmd                func(f *cmdutil.Factory) *dev.DevCmd
		changeDir             func(dir string) error
		askOne                func(p survey.Prompt, response interface{}, opts ...survey.AskOpt) error
		load                  func(filenames ...string) (err error)
		dir                   func() config.DirPath
		mkdirTemp             func(dir, pattern string) (string, error)
		readAll               func(r io.Reader) ([]byte, error)
		get                   func(url string) (resp *http.Response, err error)
		marshalIndent         func(v any, prefix, indent string) ([]byte, error)
		unmarshal             func(data []byte, v any) error
		DetectPackageManager  func(pathWorkDir string) string
	}

	type args struct {
		c    *cobra.Command
		m    string
		msgs *[]string
	}

	mock := &httpmock.Registry{}
	f, _, _ := testutils.NewFactory(mock)

	cmd := NewCmd(f)

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "success flow",
			fields: fields{
				fileReader: func(filename string) ([]byte, error) {
					return nil, nil
				},
				f: &cmdutil.Factory{
					Flags: cmdutil.Flags{
						GlobalFlagAll: false,
						Format:        "",
						Out:           "",
						NoColor:       false,
					},
					IOStreams: iostreams.System(),
				},
				commandRunInteractive: func(f *cmdutil.Factory, comm string) error {
					return nil
				},
				getWorkDir: func() (string, error) {
					return "/path", nil
				},
				DetectPackageManager: func(pathWorkDir string) string {
					return "npm"
				},
			},
			args: args{
				c:    cmd,
				m:    "message",
				msgs: &[]string{},
			},
		},
		{
			name: "error depsInstall",
			fields: fields{
				fileReader: func(filename string) ([]byte, error) {
					return nil, nil
				},
				f: &cmdutil.Factory{
					Flags: cmdutil.Flags{
						GlobalFlagAll: false,
						Format:        "",
						Out:           "",
						NoColor:       false,
					},
					IOStreams: iostreams.System(),
				},
				commandRunInteractive: func(f *cmdutil.Factory, comm string) error {
					return errors.New("error depsInstall")
				},
				getWorkDir: func() (string, error) {
					return "/path", nil
				},
				DetectPackageManager: func(pathWorkDir string) string {
					return "npm"
				},
			},
			args: args{
				c:    cmd,
				m:    "message",
				msgs: &[]string{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &initCmd{
				name:                  tt.fields.name,
				preset:                tt.fields.preset,
				auto:                  tt.fields.auto,
				packageManager:        tt.fields.packageManager,
				pathWorkingDir:        tt.fields.pathWorkingDir,
				globalFlagAll:         tt.fields.globalFlagAll,
				f:                     tt.fields.f,
				git:                   tt.fields.git,
				getWorkDir:            tt.fields.getWorkDir,
				fileReader:            tt.fields.fileReader,
				isDirEmpty:            tt.fields.isDirEmpty,
				cleanDir:              tt.fields.cleanDir,
				writeFile:             tt.fields.writeFile,
				openFile:              tt.fields.openFile,
				removeAll:             tt.fields.removeAll,
				rename:                tt.fields.rename,
				envLoader:             tt.fields.envLoader,
				stat:                  tt.fields.stat,
				mkdir:                 tt.fields.mkdir,
				gitPlainClone:         tt.fields.gitPlainClone,
				commandRunner:         tt.fields.commandRunner,
				commandRunnerOutput:   tt.fields.commandRunnerOutput,
				commandRunInteractive: tt.fields.commandRunInteractive,
				deployCmd:             tt.fields.deployCmd,
				devCmd:                tt.fields.devCmd,
				changeDir:             tt.fields.changeDir,
				askOne:                tt.fields.askOne,
				load:                  tt.fields.load,
				dir:                   tt.fields.dir,
				mkdirTemp:             tt.fields.mkdirTemp,
				readAll:               tt.fields.readAll,
				get:                   tt.fields.get,
				marshalIndent:         tt.fields.marshalIndent,
				unmarshal:             tt.fields.unmarshal,
				DetectPackageManager:  tt.fields.DetectPackageManager,
			}
			if err := cmd.deps(tt.args.c, tt.args.m, tt.args.msgs); (err != nil) != tt.wantErr {
				t.Errorf("initCmd.deps() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
