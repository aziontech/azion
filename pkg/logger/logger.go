package logger

import (
	"fmt"
	"io"

	"github.com/aziontech/tablecli"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	log *zap.Logger
)

type Logger struct {
	Debug    bool   `json:"-" yaml:"-" toml:"-"`
	Silent   bool   `json:"-" yaml:"-" toml:"-"`
	LogLevel string `json:"-" yaml:"-" toml:"-"`
}

func LogLevel(logger Logger) {
	switch {
	case logger.Debug:
		New(zapcore.DebugLevel)
	case logger.Silent:
		New(zapcore.ErrorLevel)
	case logger.LogLevel == "debug":
		New(zapcore.DebugLevel)
	case logger.LogLevel == "error":
		New(zapcore.ErrorLevel)
	default:
		New(zapcore.InfoLevel)
	}
}

func New(level zapcore.Level) {
	var err error

	config := zap.NewProductionConfig()

	logLevel := zap.NewAtomicLevel()
	logLevel.SetLevel(level)
	config.Level = logLevel

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	config.EncoderConfig = encoderConfig
	config.Encoding = "console"
	config.OutputPaths = []string{"stdout"}
	config.ErrorOutputPaths = []string{"stderr"}
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder

	log, err = config.Build(zap.AddCallerSkip(1))

	if err != nil {
		panic(err)
	}
}

func shouldPrint() bool {
	core := log.Core()
	errorEnabled := core.Enabled(zapcore.ErrorLevel)
	debugEnabled := core.Enabled(zapcore.DebugLevel)
	infoEnabled := core.Enabled(zapcore.InfoLevel)

	return !(errorEnabled && !debugEnabled && !infoEnabled)
}

// FInfo I need to check if the debug is false because the error comes in the debug also as true
func FInfo(w io.Writer, message string) {
	if shouldPrint() {
		fmt.Fprintf(w, "%s", message) // nolint:all
	}
}

func FInfoFlags(w io.Writer, message, format, out string) {
	if len(format) > 0 || len(out) > 0 {
		return
	}

	if shouldPrint() {
		fmt.Fprintf(w, "%s", message) // nolint:all
	}
}

func PrintHeader(table tablecli.Table, format string) {
	if shouldPrint() {
		table.PrintHeader(format)
	}
}

func PrintRow(table tablecli.Table, format string, row []string) {
	if shouldPrint() {
		table.PrintRow(format, row)
	}
}

func Debug(message string, fields ...zap.Field) {
	log.Debug(message, fields...)
}

func Error(message string, fields ...zap.Field) {
	log.Error(message, fields...)
}
