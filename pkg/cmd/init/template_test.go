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
		mode          string
		writeFile     func(filename string, data []byte, perm fs.FileMode) error
		mkdir         func(path string, perm os.FileMode) error
		marshalIndent func(v any, prefix, indent string) ([]byte, error)
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
				name:          tt.fields.name,
				preset:        tt.fields.preset,
				mode:          tt.fields.mode,
				writeFile:     tt.fields.writeFile,
				mkdir:         tt.fields.mkdir,
				marshalIndent: tt.fields.marshalIndent,
			}
			if err := cmd.createTemplateAzion(); (err != nil) != tt.wantErr {
				t.Errorf("initCmd.createTemplateAzion() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
