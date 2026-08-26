package init

import (
	"encoding/json"
	"errors"
	"io/fs"
	"os"
	"testing"

	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/utils"
	"go.uber.org/zap/zapcore"

	msg "github.com/aziontech/azion-cli/messages/init"
)

func Test_initCmd_createTemplateAzion(t *testing.T) {
	logger.New(zapcore.DebugLevel)

	type fields struct {
		name          string
		preset        string
		writeFile     func(filename string, data []byte, perm fs.FileMode) error
		mkdir         func(path string, perm os.FileMode) error
		marshalIndent func(v any, prefix, indent string) ([]byte, error)
		fileReader    func(path string) ([]byte, error)
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "success flow",
			fields: fields{
				fileReader: func(filename string) ([]byte, error) {
					return nil, nil
				},
				name:   "project_piece",
				preset: "vite",
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
				fileReader: func(filename string) ([]byte, error) {
					return nil, nil
				},
				name:   "project_piece",
				preset: "vite",
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
				fileReader: func(filename string) ([]byte, error) {
					return nil, nil
				},
				name:   "project_piece",
				preset: "vite",
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
				fileReader: func(filename string) ([]byte, error) {
					return nil, nil
				},
				name:   "project_piece",
				preset: "vite",
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
				name:          tt.fields.name,
				preset:        tt.fields.preset,
				writeFile:     tt.fields.writeFile,
				mkdir:         tt.fields.mkdir,
				marshalIndent: tt.fields.marshalIndent,
				fileReader:    tt.fields.fileReader,
			}
			if err := cmd.createTemplateAzion(); (err != nil) != tt.wantErr {
				t.Errorf("initCmd.createTemplateAzion() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_initCmd_createTemplateAzion_ConfigDir(t *testing.T) {
	logger.New(zapcore.DebugLevel)

	tests := []struct {
		name        string
		projectPath string
		expectedDir string
	}{
		{name: "default config dir", projectPath: "azion", expectedDir: "myproject/azion"},
		{name: "custom config dir", projectPath: "myconf", expectedDir: "myproject/myconf"},
		{name: "nested config dir", projectPath: "conf/azion", expectedDir: "myproject/conf/azion"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var createdDir, writtenFile string
			cmd := &initCmd{
				name:           "project_piece",
				preset:         "vite",
				pathWorkingDir: "myproject",
				projectPath:    tt.projectPath,
				marshalIndent:  json.MarshalIndent,
				fileReader:     func(string) ([]byte, error) { return nil, nil },
				mkdir: func(path string, perm os.FileMode) error {
					createdDir = path
					return nil
				},
				writeFile: func(filename string, data []byte, perm fs.FileMode) error {
					writtenFile = filename
					return nil
				},
			}

			if err := cmd.createTemplateAzion(); err != nil {
				t.Fatalf("createTemplateAzion() error = %v", err)
			}
			if createdDir != tt.expectedDir {
				t.Errorf("created directory = %q, want %q", createdDir, tt.expectedDir)
			}
			expectedFile := tt.expectedDir + "/azion.json"
			if writtenFile != expectedFile {
				t.Errorf("wrote azion.json to %q, want %q", writtenFile, expectedFile)
			}
		})
	}
}
