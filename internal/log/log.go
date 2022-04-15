package log

import (
	"github.com/davecgh/go-spew/spew"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	DefaultLogger = NewLoggerX(zap.DebugLevel, zap.AddCallerSkip(1))
)

type Logger struct {
	*zap.SugaredLogger
}

func NewLoggerX(level zapcore.Level, opts ...zap.Option) *Logger {
	logger, err := NewLogger(level, opts...)
	if err != nil {
		panic(err)
	}
	return logger
}

func NewLogger(level zapcore.Level, opts ...zap.Option) (*Logger, error) {
	l, err := zap.Config{
		DisableStacktrace: true,
		Level:             zap.NewAtomicLevelAt(level),
		Development:       true,
		Encoding:          "console",
		EncoderConfig:     zap.NewDevelopmentEncoderConfig(),
		OutputPaths:       []string{"stderr"},
		ErrorOutputPaths:  []string{"stderr"},
	}.Build(opts...)

	if err != nil {
		return nil, err
	}

	return &Logger{SugaredLogger: l.WithOptions(zap.AddCallerSkip(1)).Sugar()}, nil
}

func (l *Logger) Debugf(a string, args ...interface{}) {
	l.SugaredLogger.Debugf(a, args...)
}

func Debugf(a string, args ...interface{}) {
	DefaultLogger.Debugf(a, args...)
}

func (l *Logger) Infof(a string, args ...interface{}) {
	l.SugaredLogger.Debugf(a, args...)
}

func (l *Logger) Printf(a string, args ...interface{}) {
	l.SugaredLogger.Debugf(a, args...)
}

func Printf(a string, args ...interface{}) {
	DefaultLogger.Printf(a, args...)
}

func Infof(a string, args ...interface{}) {
	DefaultLogger.Infof(a, args...)
}

func (l *Logger) Errorf(a string, args ...interface{}) {
	l.SugaredLogger.Errorf(a, args...)
}

func Errorf(a string, args ...interface{}) {
	DefaultLogger.Errorf(a, args...)
}

func (l *Logger) Debug(args ...interface{}) {
	l.SugaredLogger.Debug(args...)
}

func Debug(args ...interface{}) {
	DefaultLogger.Debug(args...)
}

func (l *Logger) Dump(args ...interface{}) {
	l.SugaredLogger.Debug(spew.Sdump(args...))
}

func Dump(args ...interface{}) {
	DefaultLogger.Debug(spew.Sdump(args...))
}

func (l *Logger) Info(args ...interface{}) {
	l.SugaredLogger.Info(args...)
}

func Info(args ...interface{}) {
	DefaultLogger.Info(args...)
}
