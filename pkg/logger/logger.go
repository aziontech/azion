package logger

import (
	"fmt"
	"io"

	"github.com/fatih/color"
	"github.com/aziontech/tablecli"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	log *zap.Logger
)

type Logger struct {
	Debug    bool
	Silent   bool
	LogLevel string
}

func LogLevel(logger Logger) {
	switch {
	case logger.Debug:
		New(zapcore.DebugLevel)
	case logger.Silent:
		New(zapcore.ErrorLevel)
	case "debug" == logger.LogLevel:
		New(zapcore.DebugLevel)
	case "error" == logger.LogLevel:
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
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.EpochTimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	encoderConfig.TimeKey = "timestamp"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.StacktraceKey = ""

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

// FInfo I need to check if the debug is false because the error comes in the debug also as true
func FInfo(w io.Writer, message string) {
	if !(log.Core().Enabled(zapcore.ErrorLevel) && !log.Core().Enabled(zapcore.DebugLevel)) ||
		!(log.Core().Enabled(zapcore.ErrorLevel) && !log.Core().Enabled(zapcore.InfoLevel)) {
		fmt.Fprintf(w, message) // nolint:all
	}
}

func PrintHeader(table tablecli.Table, format string) {
	if !(log.Core().Enabled(zapcore.ErrorLevel) && !log.Core().Enabled(zapcore.DebugLevel)) ||
		!(log.Core().Enabled(zapcore.ErrorLevel) && !log.Core().Enabled(zapcore.InfoLevel)) {
		table.PrintHeader(format)
	}
}

func PrintRow(table tablecli.Table, format string, row []string) {
	if !(log.Core().Enabled(zapcore.ErrorLevel) && !log.Core().Enabled(zapcore.DebugLevel)) ||
		!(log.Core().Enabled(zapcore.ErrorLevel) && !log.Core().Enabled(zapcore.InfoLevel)) {
		table.PrintRow(format, row)
	}
}

func Info(message string, fields ...zap.Field) {
	log.Info(message, fields...)
}

func Fatal(message string, fields ...zap.Field) {
	log.Fatal(message, fields...)
}

func Debug(message string, fields ...zap.Field) {
	log.Debug(message, fields...)
}

func Error(message string, fields ...zap.Field) {
	log.Error(message, fields...)
}

func LogSuccess(w io.Writer, message string) {
	formatSuccess := color.New(color.FgGreen).SprintfFunc()
	FInfo(w, formatSuccess("🚀 %s\n", message))
}

func LogSuccessBad(w io.Writer, message string) {
	formatSuccess := color.New(color.FgBlue).SprintfFunc()
	FInfo(w, formatSuccess("👋 %s\n", message))
}

func LogWarning(w io.Writer, message string) {
	formatWarning := color.New(color.FgYellow).SprintfFunc()
	FInfo(w, formatWarning("⚠️ %s\n", message))
}

func LogError(w io.Writer, message string) {
	formatError := color.New(color.FgRed).SprintfFunc()
	FInfo(w, formatError("🐛 %s\n", message))
}
