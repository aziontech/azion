package build

import (
	"io"
	"io/fs"
	"testing"

	msg "github.com/aziontech/azion-cli/messages/build"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/iostreams"
	"github.com/aziontech/azion-cli/pkg/logger"
	"go.uber.org/zap/zapcore"
)

func TestBuildCmd_run(t *testing.T) {
	logger.New(zapcore.DebugLevel)

	type fields struct {
		Io                    *iostreams.IOStreams
		WriteFile             func(filename string, data []byte, perm fs.FileMode) error
		CommandRunnerStream   func(out io.Writer, cmd string, envvars []string) error
		CommandRunInteractive func(f *cmdutil.Factory, comm string) error
		CommandRunner         func(f *cmdutil.Factory, comm string, envVars []string) (string, error)
		FileReader            func(path string) ([]byte, error)
		GetAzionJsonContent   func(pathConf string) (*contracts.AzionApplicationOptions, error)
		WriteAzionJsonContent func(conf *contracts.AzionApplicationOptions, confPath string) error
		EnvLoader             func(path string) ([]string, error)
		Stat                  func(path string) (fs.FileInfo, error)
		GetWorkDir            func() (string, error)
		f                     *cmdutil.Factory
	}
	type args struct {
		fields *contracts.BuildInfo
		msgs   *[]string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "flow completed with success",
			fields: fields{
				GetAzionJsonContent: func(pathConf string) (*contracts.AzionApplicationOptions, error) {
					return &contracts.AzionApplicationOptions{}, nil
				},
				Io: iostreams.System(),
				CommandRunner: func(
					f *cmdutil.Factory,
					comm string,
					envVars []string,
				) (string, error) {
					return "", nil
				},
				CommandRunInteractive: func(f *cmdutil.Factory, comm string) error {
					return nil
				},
				WriteAzionJsonContent: func(
					conf *contracts.AzionApplicationOptions,
					confPath string,
				) error {
					return nil
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
			},
			args: args{
				fields: &contracts.BuildInfo{
					ProjectPath: "",
				},
				msgs: &[]string{},
			},
		},
		{
			name: "flow completed with success, fields full values",
			fields: fields{
				GetAzionJsonContent: func(pathConf string) (*contracts.AzionApplicationOptions, error) {
					return &contracts.AzionApplicationOptions{}, nil
				},
				Io: iostreams.System(),
				CommandRunner: func(
					f *cmdutil.Factory,
					comm string,
					envVars []string,
				) (string, error) {
					return "", nil
				},
				CommandRunInteractive: func(f *cmdutil.Factory, comm string) error {
					return nil
				},
				WriteAzionJsonContent: func(
					conf *contracts.AzionApplicationOptions,
					confPath string,
				) error {
					return nil
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
			},
			args: args{
				fields: &contracts.BuildInfo{
					ProjectPath:   "",
					Preset:        "vanilla",
					Entry:         "no",
					NodePolyfills: "true",
					OwnWorker:     "true",
					IsFirewall:    true,
				},
				msgs: &[]string{},
			},
		},
		{
			name: "Error Get Azion json",
			fields: fields{
				GetAzionJsonContent: func(pathConf string) (*contracts.AzionApplicationOptions, error) {
					return &contracts.AzionApplicationOptions{}, msg.ErrorBuilding
				},
				Io: iostreams.System(),
				CommandRunner: func(
					f *cmdutil.Factory,
					comm string,
					envVars []string,
				) (string, error) {
					return "", nil
				},
				CommandRunInteractive: func(f *cmdutil.Factory, comm string) error {
					return nil
				},
				WriteAzionJsonContent: func(
					conf *contracts.AzionApplicationOptions,
					confPath string,
				) error {
					return nil
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
			},
			args: args{
				fields: &contracts.BuildInfo{
					ProjectPath:   "",
					Preset:        "vanilla",
					Entry:         "no",
					NodePolyfills: "true",
					OwnWorker:     "true",
					IsFirewall:    true,
				},
				msgs: &[]string{},
			},
			wantErr: true,
		},
		{
			name: "Error parse NodePolyfills",
			fields: fields{
				GetAzionJsonContent: func(pathConf string) (*contracts.AzionApplicationOptions, error) {
					return &contracts.AzionApplicationOptions{}, nil
				},
				Io: iostreams.System(),
				CommandRunner: func(
					f *cmdutil.Factory,
					comm string,
					envVars []string,
				) (string, error) {
					return "", nil
				},
				CommandRunInteractive: func(f *cmdutil.Factory, comm string) error {
					return nil
				},
				WriteAzionJsonContent: func(
					conf *contracts.AzionApplicationOptions,
					confPath string,
				) error {
					return nil
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
			},
			args: args{
				fields: &contracts.BuildInfo{
					ProjectPath:   "",
					Preset:        "vanilla",
					Entry:         "no",
					NodePolyfills: "adf",
					OwnWorker:     "true",
					IsFirewall:    true,
				},
				msgs: &[]string{},
			},
			wantErr: true,
		},
		{
			name: "Error parse OwnWorker",
			fields: fields{
				GetAzionJsonContent: func(pathConf string) (*contracts.AzionApplicationOptions, error) {
					return &contracts.AzionApplicationOptions{}, nil
				},
				Io: iostreams.System(),
				CommandRunner: func(
					f *cmdutil.Factory,
					comm string,
					envVars []string,
				) (string, error) {
					return "", nil
				},
				CommandRunInteractive: func(f *cmdutil.Factory, comm string) error {
					return nil
				},
				WriteAzionJsonContent: func(
					conf *contracts.AzionApplicationOptions,
					confPath string,
				) error {
					return nil
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
			},
			args: args{
				fields: &contracts.BuildInfo{
					ProjectPath:   "",
					Preset:        "vanilla",
					Entry:         "no",
					NodePolyfills: "true",
					OwnWorker:     "adf",
					IsFirewall:    true,
				},
				msgs: &[]string{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &BuildCmd{
				Io:                    tt.fields.Io,
				WriteFile:             tt.fields.WriteFile,
				CommandRunnerStream:   tt.fields.CommandRunnerStream,
				CommandRunInteractive: tt.fields.CommandRunInteractive,
				CommandRunner:         tt.fields.CommandRunner,
				FileReader:            tt.fields.FileReader,
				GetAzionJsonContent:   tt.fields.GetAzionJsonContent,
				WriteAzionJsonContent: tt.fields.WriteAzionJsonContent,
				EnvLoader:             tt.fields.EnvLoader,
				Stat:                  tt.fields.Stat,
				GetWorkDir:            tt.fields.GetWorkDir,
				f:                     tt.fields.f,
			}
			if err := b.run(tt.args.fields, tt.args.msgs); (err != nil) != tt.wantErr {
				t.Errorf("BuildCmd.run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
