package init

import (
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
	"github.com/aziontech/azion-cli/pkg/iostreams"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/utils"
	"github.com/go-git/go-git/v5"
	"go.uber.org/zap/zapcore"

	msg "github.com/aziontech/azion-cli/messages/init"
)

func Test_initCmd_createTemplateAzion(t *testing.T) {
	logger.New(zapcore.DebugLevel)

	type fields struct {
		name                  string
		preset                string
		auto                  bool
		mode                  string
		packageManager        string
		pathWorkingDir        string
		globalFlagAll         bool
		f                     *cmdutil.Factory
		io                    *iostreams.IOStreams
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
		dir                   func() (config.DirPath, error)
		mkdirTemp             func(dir, pattern string) (string, error)
		readAll               func(r io.Reader) ([]byte, error)
		get                   func(url string) (resp *http.Response, err error)
		marshalIndent         func(v any, prefix, indent string) ([]byte, error)
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "success flow",
			fields: fields{
				name:   "project_piece",
				preset: "vite",
				mode:   "deliver",
				mkdir: func(path string, perm os.FileMode) error {
					return nil
				},
				marshalIndent: json.MarshalIndent,
				writeFile: func(filename string, data []byte, perm fs.FileMode) error {
					return nil
				},
			},
			wantErr: false,
		},
		{
			name: "error mkdir",
			fields: fields{
				name:   "project_piece",
				preset: "vite",
				mode:   "deliver",
				mkdir: func(path string, perm os.FileMode) error {
					return msg.ErrorFailedCreatingAzionDirectory
				},
				marshalIndent: json.MarshalIndent,
				writeFile: func(filename string, data []byte, perm fs.FileMode) error {
					return nil
				},
			},
			wantErr: true,
		},
		{
			name: "error marshalIndent",
			fields: fields{
				name:   "project_piece",
				preset: "vite",
				mode:   "deliver",
				mkdir: func(path string, perm os.FileMode) error {
					return nil
				},
				marshalIndent: func(v any, prefix, indent string) ([]byte, error) {
					return []byte(""), errors.New("error marshalIndent")
				},
				writeFile: func(filename string, data []byte, perm fs.FileMode) error {
					return nil
				},
			},
			wantErr: true,
		},
		{
			name: "error writeFile",
			fields: fields{
				name:   "project_piece",
				preset: "vite",
				mode:   "deliver",
				mkdir: func(path string, perm os.FileMode) error {
					return nil
				},
				marshalIndent: json.MarshalIndent,
				writeFile: func(filename string, data []byte, perm fs.FileMode) error {
					return utils.ErrorInternalServerError
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
				mode:                  tt.fields.mode,
				packageManager:        tt.fields.packageManager,
				pathWorkingDir:        tt.fields.pathWorkingDir,
				globalFlagAll:         tt.fields.globalFlagAll,
				f:                     tt.fields.f,
				io:                    tt.fields.io,
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
			}
			if err := cmd.createTemplateAzion(); (err != nil) != tt.wantErr {
				t.Errorf("initCmd.createTemplateAzion() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
