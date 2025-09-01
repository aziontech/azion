package output

import (
	"bytes"
	"testing"

	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/logger"
	"go.uber.org/zap/zapcore"
)

func TestListOutput_Format(t *testing.T) {
	logger.New(zapcore.DebugLevel)
	outBuffer := &bytes.Buffer{}

	type fields struct {
		GeneralOutput GeneralOutput
		Columns       []string
		Lines         [][]string
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
				Columns: []string{
					"ID",
					"NAME",
					"ACTIVE",
				},
				Lines: [][]string{
					{
						"1665081615",
						"Update Application",
						"true",
					},
					{
						"1665081616",
						"testazion2",
						"true",
					},
					{
						"1694024475",
						"New Application",
						"true",
					},
				},
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
				Columns: []string{
					"ID",
					"NAME",
					"ACTIVE",
				},
				Lines: [][]string{
					{
						"1665081615",
						"Update Application",
						"true",
					},
					{
						"1665081616",
						"testazion2",
						"true",
					},
					{
						"1694024475",
						"New Application",
						"true",
					},
				},
			},
			want:    false,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &ListOutput{
				GeneralOutput: tt.fields.GeneralOutput,
				Columns:       tt.fields.Columns,
				Lines:         tt.fields.Lines,
			}
			got, err := l.Format()
			if (err != nil) != tt.wantErr {
				t.Errorf("ListOutput.Format() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ListOutput.Format() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestListOutput_Output(t *testing.T) {
	logger.New(zapcore.DebugLevel)
	outBuffer := &bytes.Buffer{}

	type fields struct {
		GeneralOutput GeneralOutput
		Columns       []string
		Lines         [][]string
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "success flow with flags",
			fields: fields{
				GeneralOutput: GeneralOutput{
					Out:   outBuffer,
					Flags: cmdutil.Flags{Out: "", Format: "json"},
				},
				Columns: []string{
					"ID",
					"NAME",
					"ACTIVE",
				},
				Lines: [][]string{
					{
						"1665081615",
						"Update Application",
						"true",
					},
					{
						"1665081616",
						"testazion2",
						"true",
					},
					{
						"1694024475",
						"New Application",
						"true",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &ListOutput{
				GeneralOutput: tt.fields.GeneralOutput,
				Columns:       tt.fields.Columns,
				Lines:         tt.fields.Lines,
			}
			c.Output()
		})
	}
}
