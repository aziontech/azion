package build

import (
	"io"
	"io/fs"
	"testing"

	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/iostreams"
	"github.com/aziontech/azion-cli/pkg/logger"
	"go.uber.org/zap/zapcore"
)

func TestNewCmd(t *testing.T) {
	f := &cmdutil.Factory{
		Flags: cmdutil.Flags{
			GlobalFlagAll: false,
			Format:        "",
			Out:           "",
			NoColor:       false,
		},
		IOStreams: iostreams.System(),
	}
	NewCmd(f)
}

func TestBuildCmd_ExternalRun(t *testing.T) {
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
		fields   *contracts.BuildInfoV3
		confPath string
		msgs     *[]string
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
				fields: &contracts.BuildInfoV3{
					ProjectPath: "",
				},
				confPath: "",
				msgs:     &[]string{},
			},
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
			if err := b.ExternalRun(tt.args.fields, tt.args.confPath, tt.args.msgs); (err != nil) != tt.wantErr {
				t.Errorf("BuildCmd.ExternalRun() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewBuildCmd(t *testing.T) {
	f := &cmdutil.Factory{
		Flags: cmdutil.Flags{
			GlobalFlagAll: false,
			Format:        "",
			Out:           "",
			NoColor:       false,
		},
		IOStreams: iostreams.System(),
	}
	NewBuildCmd(f)
}
