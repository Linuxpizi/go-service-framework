package logger

import (
	"os"

	"golang.org/x/xerrors"

	zap "go.uber.org/zap"
	zapcore "go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	DebugLevel   = "debug"
	InfoLevel    = "info"
	WarningLevel = "warning"
	ErrorLevel   = "error"
)

var myLogger *zap.Logger

// https://github.com/uber-go/zap/blob/master/FAQ.md
func Init(level, logFile string) error {
	var zapLevel zapcore.Level
	switch level {
	case DebugLevel:
		zapLevel = zap.DebugLevel
	case InfoLevel:
		zapLevel = zap.InfoLevel
	case WarningLevel:
		zapLevel = zap.WarnLevel
	case ErrorLevel:
		zapLevel = zap.ErrorLevel
	default:
		return xerrors.Errorf("unknow log level %s", level)
	}

	fileLog := zapcore.AddSync(&lumberjack.Logger{
		Filename:   logFile,
		MaxSize:    100,
		MaxBackups: 10,
		MaxAge:     15,
		Compress:   true,
	})
	consoleLog := zapcore.Lock(os.Stdout)

	encode := zapcore.NewJSONEncoder(zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	})
	core := zapcore.NewTee(
		zapcore.NewCore(encode, fileLog, zapLevel),
		zapcore.NewCore(encode, consoleLog, zapLevel),
	)
	myLogger = zap.New(core).WithOptions(zap.AddCaller())

	return nil
}

func Sugar() *zap.SugaredLogger {
	return myLogger.Sugar()
}

func Sync() error {
	return myLogger.Sync()
}
