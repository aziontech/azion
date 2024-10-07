package output

import (
	"bytes"
	"errors"
	"testing"

	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/logger"
	"go.uber.org/zap/zapcore"
)

func TestErrorOutput_Format(t *testing.T) {
	logger.New(zapcore.DebugLevel)
	outBuffer := &bytes.Buffer{}

	type fields struct {
		GeneralOutput GeneralOutput
		Err           error
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
				GeneralOutput: GeneralOutput{
					Out:   outBuffer,
					Flags: cmdutil.Flags{Out: "", Format: "json"},
				},
				Err: errors.New("error"),
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "format no flags",
			fields: fields{
				GeneralOutput: GeneralOutput{
					Out:   outBuffer,
					Flags: cmdutil.Flags{Out: "", Format: ""},
				},
				Err: errors.New("error"),
			},
			want:    false,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &ErrorOutput{
				GeneralOutput: tt.fields.GeneralOutput,
				Err:           tt.fields.Err,
			}
			got, err := e.Format()
			if (err != nil) != tt.wantErr {
				t.Errorf("ErrorOutput.Format() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ErrorOutput.Format() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestErrorOutput_Output(t *testing.T) {
	logger.New(zapcore.DebugLevel)
	outBuffer := &bytes.Buffer{}

	type fields struct {
		GeneralOutput GeneralOutput
		Err           error
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "success with flags",
			fields: fields{
				GeneralOutput: GeneralOutput{
					Out:   outBuffer,
					Flags: cmdutil.Flags{Out: "", Format: "json"},
				},
				Err: errors.New("error"),
			},
		},
		{
			name: "success no flags",
			fields: fields{
				GeneralOutput: GeneralOutput{
					Out:   outBuffer,
					Flags: cmdutil.Flags{Out: "", Format: ""},
				},
				Err: errors.New("error"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &ErrorOutput{
				GeneralOutput: tt.fields.GeneralOutput,
				Err:           tt.fields.Err,
			}
			e.Output()
		})
	}
}
