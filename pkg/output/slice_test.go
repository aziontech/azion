package output

import (
	"bytes"
	"testing"

	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/logger"
	"go.uber.org/zap/zapcore"
)

func TestSliceOutput_Format(t *testing.T) {
	logger.New(zapcore.DebugLevel)
	outBuffer := &bytes.Buffer{}

	type fields struct {
		Messages      []string
		GeneralOutput GeneralOutput
	}
	tests := []struct {
		name    string
		fields  fields
		want    bool
		wantErr bool
	}{
		{
			name: "format with flags",
			fields: fields{
				Messages: []string{"success"},
				GeneralOutput: GeneralOutput{
					Out:   outBuffer,
					Flags: cmdutil.Flags{Out: "", Format: "json"},
				},
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "format no flags",
			fields: fields{
				Messages: []string{"success"},
				GeneralOutput: GeneralOutput{
					Out:   outBuffer,
					Flags: cmdutil.Flags{Out: "", Format: ""},
				},
			},
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &SliceOutput{
				Messages:      tt.fields.Messages,
				GeneralOutput: tt.fields.GeneralOutput,
			}
			got, err := i.Format()
			if (err != nil) != tt.wantErr {
				t.Errorf("SliceOutput.Format() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("SliceOutput.Format() = %v, want %v", got, tt.want)
			}
		})
	}
}
