package logger

import (
	"bytes"
	"testing"

	"go.uber.org/zap/zapcore"
)

func TestFInfo(t *testing.T) {
	type args struct {
		message string
	}
	tests := []struct {
		name   string
		args   args
		wantW  string
		level  zapcore.Level
		debug  bool
		silent bool
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
			wantW: "level error",
			level: zapcore.ErrorLevel,
		},
		{
			name: "silent mode",
			args: args{
				"silent mode",
			},
			wantW:  "",
			level:  zapcore.InfoLevel,
			silent: true,
		},
		{
			name: "debug mode",
			args: args{
				"debug mode",
			},
			wantW: "debug mode",
			level: zapcore.InfoLevel,
			debug: true,
		},
		{
			name: "level warning",
			args: args{
				"level warning",
			},
			wantW: "level warning",
			level: zapcore.WarnLevel,
		},
		{
			name: "level fatal",
			args: args{
				"level fatal",
			},
			wantW: "level fatal",
			level: zapcore.FatalLevel,
		},
		{
			name: "level panic",
			args: args{
				"level panic",
			},
			wantW: "level panic",
			level: zapcore.PanicLevel,
		},
		{
			name: "level info with debug mode",
			args: args{
				"level info with debug mode",
			},
			wantW: "level info with debug mode",
			level: zapcore.InfoLevel,
			debug: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			log := Logger{Debug: tt.debug, Silent: tt.silent}
			LogLevel(log)
			FInfo(w, tt.args.message)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("FInfo() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func TestError(t *testing.T) {
	type args struct {
		message string
		fields  []zapcore.Field
	}
	tests := []struct {
		name   string
		args   args
		wantW  string
		debug  bool
		silent bool
	}{
		{
			name: "error message in silent mode",
			args: args{
				message: "error message in silent mode",
				fields:  nil,
			},
			wantW:  "",
			silent: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			log := Logger{Debug: tt.debug, Silent: tt.silent}
			LogLevel(log)
			Error(tt.args.message, tt.args.fields...)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("Error() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func TestDebug(t *testing.T) {
	type args struct {
		message string
		fields  []zapcore.Field
	}
	tests := []struct {
		name   string
		args   args
		wantW  string
		debug  bool
		silent bool
	}{
		{
			name: "debug message in silent mode",
			args: args{
				message: "debug message in silent mode",
				fields:  nil,
			},
			wantW:  "",
			silent: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			log := Logger{Debug: tt.debug, Silent: tt.silent}
			LogLevel(log)
			Debug(tt.args.message, tt.args.fields...)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("Debug() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func TestFInfoFlags(t *testing.T) {
	type args struct {
		message string
		format  string
		out     string
	}
	tests := []struct {
		name   string
		args   args
		wantW  string
		debug  bool
		silent bool
	}{
		{
			name: "info flags with format and output",
			args: args{
				message: "info flags with format and output",
				format:  "yaml",
				out:     "output",
			},
			wantW: "",
		},
		{
			name: "info flags in silent mode",
			args: args{
				message: "info flags in silent mode",
				format:  "",
				out:     "",
			},
			wantW:  "",
			silent: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			log := Logger{Debug: tt.debug, Silent: tt.silent}
			LogLevel(log)
			FInfoFlags(w, tt.args.message, tt.args.format, tt.args.out)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("FInfoFlags() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}
