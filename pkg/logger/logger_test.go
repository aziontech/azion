package logger

import (
	"bytes"
	"go.uber.org/zap/zapcore"
	"testing"
)

func TestFInfo(t *testing.T) {
	type args struct {
		message string
	}
	tests := []struct {
		name  string
		args  args
		wantW string
		level zapcore.Level
	}{
		{
			name: "level debug",
			args: args{
				"level debug",
			},
			wantW: "level debug",
			level: zapcore.DebugLevel,
		},
		{
			name: "level information",
			args: args{
				"level information",
			},
			wantW: "level information",
			level: zapcore.InfoLevel,
		},
		{
			name: "level error",
			args: args{
				"level error",
			},
			wantW: "",
			level: zapcore.ErrorLevel,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			New(tt.level)
			FInfo(w, tt.args.message)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("FInfo() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}
