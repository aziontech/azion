package build

import (
	"fmt"
	"io"
	"io/fs"
	"testing"

	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/iostreams"
	"github.com/aziontech/azion-cli/pkg/logger"
	vulcanPkg "github.com/aziontech/azion-cli/pkg/vulcan"
	"go.uber.org/zap/zapcore"
)

func TestBuildCmd_vulcan(t *testing.T) {
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
		vul          *vulcanPkg.VulcanPkg
		conf         *contracts.AzionApplicationOptions
		vulcanParams string
		fields       *contracts.BuildInfo
		msgs         *[]string
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
					Command: func(
						flags, params string,
						f *cmdutil.Factory,
					) string {
						installEdgeFunctions := "npx --yes %s edge-functions%s %s"
						versionVulcan := "@3.2.1"
						return fmt.Sprintf(
							installEdgeFunctions,
							flags,
							versionVulcan,
							params,
						)
					},
				},
				conf:         &contracts.AzionApplicationOptions{},
				vulcanParams: "",
				fields:       &contracts.BuildInfo{},
				msgs:         &[]string{},
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
			if err := b.vulcan(tt.args.vul,
				tt.args.conf,
				tt.args.vulcanParams,
				tt.args.fields,
				tt.args.msgs,
			); (err != nil) != tt.wantErr {
				t.Errorf("BuildCmd.vulcan() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
