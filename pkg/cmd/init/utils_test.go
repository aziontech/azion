package init

import (
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
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "success flow",
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
				askOne: tt.fields.askOne,
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
		mode                  string
		packageManager        string
		pathWorkingDir        string
		commandRunner         func(envVars []string, comm string) (string, int, error)
		commandRunnerOutput   func(f *cmdutil.Factory, comm string, envVars []string) (string, error)
		commandRunInteractive func(f *cmdutil.Factory, comm string) error
		load                  func(filenames ...string) (err error)
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
			name: "success flow ",
			fields: fields{
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
					os.Setenv("mode", "compute")
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
				preset:         "vanilla",
				mode:           "compute",
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
					os.Setenv("mode", "deliver")
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
				mode:           "compute",
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
					os.Setenv("mode", "compute")
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
				mode:           "compute",
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
					os.Setenv("mode", "compute")
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
			name: "error commandRunInteractive",
			fields: fields{
				preset:         "vanilla",
				mode:           "compute",
				pathWorkingDir: "./azion/pathmock",
				commandRunnerOutput: func(
					f *cmdutil.Factory,
					comm string,
					envVars []string,
				) (string, error) {
					return "3.2.1", nil
				},
				commandRunInteractive: func(f *cmdutil.Factory, comm string) error {
					return errors.New("error commandRunInteractive")
				},
				load: func(filenames ...string) (err error) {
					os.Setenv("preset", "vanilla")
					os.Setenv("mode", "compute")
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
			name: "error getVulcanEnvInfo",
			fields: fields{
				preset:         "vanilla",
				mode:           "compute",
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
				mode:                  tt.fields.mode,
				packageManager:        tt.fields.packageManager,
				pathWorkingDir:        tt.fields.pathWorkingDir,
				commandRunner:         tt.fields.commandRunner,
				commandRunnerOutput:   tt.fields.commandRunnerOutput,
				commandRunInteractive: tt.fields.commandRunInteractive,
				load:                  tt.fields.load,
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

func Test_initCmd_getVulcanEnvInfo(t *testing.T) {
	logger.New(zapcore.DebugLevel)

	type fields struct {
		load func(filenames ...string) (err error)
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		want1   string
		wantErr bool
	}{
		{
			name: "flow completed with success",
			fields: fields{
				load: func(filenames ...string) (err error) {
					os.Setenv("preset", "vanilla")
					os.Setenv("mode", "compute")
					return nil
				},
			},
			want:    "vanilla",
			want1:   "compute",
			wantErr: false,
		},
		{
			name: "error load envirements on .vulcan",
			fields: fields{
				load: func(filenames ...string) (err error) {
					return errors.New("error loading .vulcan file")
				},
			},
			want:    "",
			want1:   "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &initCmd{
				load: tt.fields.load,
			}
			got, got1, err := cmd.getVulcanInfo()
			if (err != nil) != tt.wantErr {
				t.Errorf("initCmd.getVulcanEnvInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("initCmd.getVulcanEnvInfo() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("initCmd.getVulcanEnvInfo() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
