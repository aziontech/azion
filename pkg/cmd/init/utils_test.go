package init

import (
	"encoding/json"
	"errors"
	"os"
	"testing"

	"github.com/AlecAivazis/survey/v2"
	msg "github.com/aziontech/azion-cli/messages/init"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/logger"
	vulcanPkg "github.com/aziontech/azion-cli/pkg/vulcan"
	"go.uber.org/zap/zapcore"
)

func Test_initCmd_askForInput(t *testing.T) {
	logger.New(zapcore.DebugLevel)

	type fields struct {
		askOne func(p survey.Prompt, response interface{}, opts ...survey.AskOpt) error
	}
	type args struct {
		msg       string
		defaultIn string
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		want     string
		wantErr  bool
		readFile func(filename string) ([]byte, error)
	}{
		{
			name: "success flow",
			readFile: func(filename string) ([]byte, error) {
				return nil, nil
			},
			fields: fields{
				askOne: func(
					p survey.Prompt,
					response interface{},
					opts ...survey.AskOpt,
				) error {
					return nil
				},
			},
			args: args{
				msg:       "",
				defaultIn: "",
			},
			want:    "",
			wantErr: false,
		},
		{
			name: "error askOne",
			fields: fields{
				askOne: func(
					p survey.Prompt,
					response interface{},
					opts ...survey.AskOpt,
				) error {
					return errors.New("error askOne")
				},
			},
			args: args{
				msg:       "",
				defaultIn: "",
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &initCmd{
				askOne:     tt.fields.askOne,
				fileReader: tt.readFile,
			}
			got, err := cmd.askForInput(tt.args.msg, tt.args.defaultIn)
			if (err != nil) != tt.wantErr {
				t.Errorf("initCmd.askForInput() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("initCmd.askForInput() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_initCmd_selectVulcanTemplates(t *testing.T) {
	logger.New(zapcore.DebugLevel)

	type fields struct {
		preset                string
		auto                  bool
		packageManager        string
		pathWorkingDir        string
		commandRunner         func(envVars []string, comm string) (string, int, error)
		commandRunnerOutput   func(f *cmdutil.Factory, comm string, envVars []string) (string, error)
		commandRunInteractive func(f *cmdutil.Factory, comm string) error
		load                  func(filenames ...string) (err error)
		fileReader            func(path string) ([]byte, error)
		unmarshal             func(data []byte, v any) error
	}
	type args struct {
		vul *vulcanPkg.VulcanPkg
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
				unmarshal: json.Unmarshal,
				fileReader: func(filename string) ([]byte, error) {
					return []byte(infoJsonData), nil
				},
				commandRunnerOutput: func(
					f *cmdutil.Factory,
					comm string,
					envVars []string,
				) (string, error) {
					return "3.2.1", nil
				},
				commandRunInteractive: func(f *cmdutil.Factory, comm string) error {
					return nil
				},
				load: func(filenames ...string) (err error) {
					os.Setenv("preset", "vanilla")
					return nil
				},
			},
			args: args{
				vul: &vulcanPkg.VulcanPkg{
					CheckVulcanMajor: func(
						currentVersion string,
						f *cmdutil.Factory,
						vulcan *vulcanPkg.VulcanPkg,
					) error {
						return nil
					},
					Command: func(flags, params string, f *cmdutil.Factory) string {
						return "init"
					},
				},
			},
			wantErr: false,
		},
		{
			name: "success flow with the flags",
			fields: fields{
				unmarshal: json.Unmarshal,
				fileReader: func(filename string) ([]byte, error) {
					return []byte(infoJsonData), nil
				},
				preset:         "vanilla",
				pathWorkingDir: "./azion/pathmock",
				commandRunnerOutput: func(
					f *cmdutil.Factory,
					comm string,
					envVars []string,
				) (string, error) {
					return "3.2.1", nil
				},
				commandRunInteractive: func(f *cmdutil.Factory, comm string) error {
					return nil
				},
				load: func(filenames ...string) (err error) {
					os.Setenv("preset", "vite")
					return nil
				},
			},
			args: args{
				vul: &vulcanPkg.VulcanPkg{
					CheckVulcanMajor: func(
						currentVersion string,
						f *cmdutil.Factory,
						vulcan *vulcanPkg.VulcanPkg,
					) error {
						return nil
					},
					Command: func(flags, params string, f *cmdutil.Factory) string {
						return "init"
					},
				},
			},
			wantErr: false,
		},
		{
			name: "error command runner output",
			fields: fields{
				preset:         "vanilla",
				pathWorkingDir: "./azion/pathmock",
				commandRunnerOutput: func(
					f *cmdutil.Factory,
					comm string,
					envVars []string,
				) (string, error) {
					return "3.2.1", errors.New("error commandRunnerOutput")
				},
				commandRunInteractive: func(f *cmdutil.Factory, comm string) error {
					return nil
				},
				load: func(filenames ...string) (err error) {
					os.Setenv("preset", "vanilla")
					return nil
				},
			},
			args: args{
				vul: &vulcanPkg.VulcanPkg{
					CheckVulcanMajor: func(
						currentVersion string,
						f *cmdutil.Factory,
						vulcan *vulcanPkg.VulcanPkg,
					) error {
						return nil
					},
					Command: func(flags, params string, f *cmdutil.Factory) string {
						return "init"
					},
				},
			},
			wantErr: true,
		},
		{
			name: "error CheckVulcanMajor",
			fields: fields{
				preset:         "vanilla",
				pathWorkingDir: "./azion/pathmock",
				commandRunnerOutput: func(
					f *cmdutil.Factory,
					comm string,
					envVars []string,
				) (string, error) {
					return "3.2.1", nil
				},
				commandRunInteractive: func(f *cmdutil.Factory, comm string) error {
					return nil
				},
				load: func(filenames ...string) (err error) {
					os.Setenv("preset", "vanilla")
					return nil
				},
			},
			args: args{
				vul: &vulcanPkg.VulcanPkg{
					CheckVulcanMajor: func(
						currentVersion string,
						f *cmdutil.Factory,
						vulcan *vulcanPkg.VulcanPkg,
					) error {
						return errors.New("error CheckVulcanMajor")
					},
					Command: func(flags, params string, f *cmdutil.Factory) string {
						return "init"
					},
				},
			},
			wantErr: true,
		},
		{
			name: "error getVulcanEnvInfo",
			fields: fields{
				preset:         "vanilla",
				pathWorkingDir: "./azion/pathmock",
				commandRunnerOutput: func(
					f *cmdutil.Factory,
					comm string,
					envVars []string,
				) (string, error) {
					return "3.2.1", nil
				},
				commandRunInteractive: func(f *cmdutil.Factory, comm string) error {
					return nil
				},
				load: func(filenames ...string) (err error) {
					return errors.New("error load")
				},
				fileReader: func(path string) ([]byte, error) {
					return nil, errors.New("error reading info.json")
				},
			},
			args: args{
				vul: &vulcanPkg.VulcanPkg{
					CheckVulcanMajor: func(
						currentVersion string,
						f *cmdutil.Factory,
						vulcan *vulcanPkg.VulcanPkg,
					) error {
						return nil
					},
					Command: func(flags, params string, f *cmdutil.Factory) string {
						return "init"
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &initCmd{
				preset:                tt.fields.preset,
				auto:                  tt.fields.auto,
				packageManager:        tt.fields.packageManager,
				pathWorkingDir:        tt.fields.pathWorkingDir,
				commandRunner:         tt.fields.commandRunner,
				commandRunnerOutput:   tt.fields.commandRunnerOutput,
				commandRunInteractive: tt.fields.commandRunInteractive,
				load:                  tt.fields.load,
				fileReader:            tt.fields.fileReader,
				unmarshal:             tt.fields.unmarshal,
			}
			if err := cmd.selectVulcanTemplates(tt.args.vul); (err != nil) != tt.wantErr {
				t.Errorf("initCmd.selectVulcanTemplates() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_initCmd_depsInstall(t *testing.T) {
	logger.New(zapcore.DebugLevel)

	type fields struct {
		packageManager        string
		commandRunInteractive func(f *cmdutil.Factory, comm string) error
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "flow completed with success",
			fields: fields{
				packageManager: "npm",
				commandRunInteractive: func(f *cmdutil.Factory, comm string) error {
					return nil
				},
			},
			wantErr: false,
		},
		{
			name: "error run command interactive",
			fields: fields{
				packageManager: "npm",
				commandRunInteractive: func(f *cmdutil.Factory, comm string) error {
					return msg.ErrorDeps
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &initCmd{
				packageManager:        tt.fields.packageManager,
				commandRunInteractive: tt.fields.commandRunInteractive,
			}
			if err := cmd.depsInstall(); (err != nil) != tt.wantErr {
				t.Errorf("initCmd.depsInstall() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_initCmd_getVulcanInfo(t *testing.T) {
	logger.New(zapcore.DebugLevel)

	type fields struct {
		pathWorkingDir string
		unmarshal      func(data []byte, v interface{}) error
	}
	tests := []struct {
		name       string
		fields     fields
		mockFile   string
		wantPreset string
		wantErr    bool
		readFile   func(filename string) ([]byte, error)
	}{
		{
			name: "flow completed with success",
			fields: fields{
				pathWorkingDir: "/path/to/working/dir",
				unmarshal: func(data []byte, v interface{}) error {
					// Mocking unmarshalling process
					*(v.(*map[string]string)) = map[string]string{
						"preset": "astro",
					}
					return nil
				},
			},
			wantPreset: "astro",
			wantErr:    false,
			readFile: func(filename string) ([]byte, error) {
				return []byte(`{"preset": "astro"}`), nil
			},
		},
		{
			name: "error reading the file",
			fields: fields{
				pathWorkingDir: "/path/to/working/dir",
				unmarshal: func(data []byte, v interface{}) error {
					return errors.New("error unmarshalling json")
				},
			},
			wantPreset: "",
			wantErr:    true,
			readFile: func(filename string) ([]byte, error) {
				return nil, errors.New("error reading json")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &initCmd{
				pathWorkingDir: tt.fields.pathWorkingDir,
				unmarshal:      tt.fields.unmarshal,
				fileReader:     tt.readFile,
			}

			gotPreset, err := cmd.getVulcanInfo()
			if (err != nil) != tt.wantErr {
				t.Errorf("initCmd.getVulcanInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotPreset != tt.wantPreset {
				t.Errorf("initCmd.getVulcanInfo() gotPreset = %v, want %v", gotPreset, tt.wantPreset)
			}
		})
	}
}
