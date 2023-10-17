package logger

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"kits/api/src/common/configs"
	"os"
	"time"
)

const callerSkip = 2

type Logger interface {
	Info(ctx context.Context, msg string, v ...interface{})
	Debug(ctx context.Context, msg string, v ...interface{})
	Warn(ctx context.Context, err error, msg string, v ...interface{})
	Error(ctx context.Context, err error, msg string, v ...interface{})
	Fatal(ctx context.Context, err error, msg string, v ...interface{})
}

type noopLogger struct{}

func (noopLogger) Info(context.Context, string, ...interface{})         {}
func (noopLogger) Debug(context.Context, string, ...interface{})        {}
func (noopLogger) Warn(context.Context, error, string, ...interface{})  {}
func (noopLogger) Fatal(context.Context, error, string, ...interface{}) {}
func (noopLogger) Error(context.Context, error, string, ...interface{}) {}

type StandardLogger struct {
	l *zap.SugaredLogger
}

func (s StandardLogger) getScopedLogger(ctx context.Context) []interface{} {
	out := make([]interface{}, 0)
	ctxMeta := contextValue(ctx)
	if len(ctxMeta) == 0 {
		return out
	}
	for k, v := range ctxMeta {
		if v != nil && v != "" {
			out = append(out, k, v)
		}
	}
	return out
}

func (s StandardLogger) Info(ctx context.Context, msg string, v ...interface{}) {
	scoped := s.getScopedLogger(ctx)
	if len(scoped) > 0 {
		s.l.Infow(fmt.Sprintf(msg, v...), scoped)
	}
	s.l.Infow(fmt.Sprintf(msg, v...))

}

func (s StandardLogger) Debug(ctx context.Context, msg string, v ...interface{}) {
	scoped := s.getScopedLogger(ctx)
	if len(scoped) > 0 {
		s.l.Debugw(fmt.Sprintf(msg, v...), scoped)
	}
	s.l.Debugw(fmt.Sprintf(msg, v...))
}

func (s StandardLogger) Warn(ctx context.Context, err error, msg string, v ...interface{}) {
	scoped := s.getScopedLogger(ctx)
	scoped = append(scoped, "err", err.Error())
	if len(scoped) > 0 {
		s.l.Warnw(fmt.Sprintf(msg, v...), scoped)
	}
	s.l.Warnw(fmt.Sprintf(msg, v...))
}

func (s StandardLogger) Error(ctx context.Context, err error, msg string, v ...interface{}) {
	scoped := s.getScopedLogger(ctx)
	scoped = append(scoped, "err", err.Error())
	if len(scoped) > 0 {
		s.l.Errorw(fmt.Sprintf(msg, v...), scoped)
	}
	s.l.Errorw(fmt.Sprintf(msg, v...))
}

func (s StandardLogger) Fatal(ctx context.Context, err error, msg string, v ...interface{}) {
	scoped := s.getScopedLogger(ctx)
	scoped = append(scoped, "err", err.Error())
	if len(scoped) > 0 {
		s.l.Fatalw(fmt.Sprintf(msg, v...), scoped)
	}
	s.l.Fatalw(fmt.Sprintf(msg, v...))
}

func SyslogTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05"))
}

func CustomLevelEncoder(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString("[" + level.CapitalString() + "]")
}

func NewLogger(config *configs.Config) Logger {
	encoderConfig := zapcore.EncoderConfig{
		MessageKey:   "message",
		LevelKey:     "level",
		TimeKey:      "time",
		CallerKey:    "caller",
		EncodeCaller: zapcore.FullCallerEncoder,
		EncodeTime:   SyslogTimeEncoder,
		EncodeLevel:  CustomLevelEncoder,
	}

	var encoder zapcore.Encoder
	var level zapcore.Level
	if config.Mode.IsProd() {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
		level = zap.InfoLevel
	} else {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
		level = zap.DebugLevel
	}
	core := zapcore.NewCore(encoder, zapcore.AddSync(os.Stderr), level)
	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel), zap.AddCallerSkip(callerSkip)).Sugar()
	return StandardLogger{l: logger}
}

func NewFxLogger() *zap.SugaredLogger {
	encoderConfig := zapcore.EncoderConfig{
		MessageKey:   "message",
		LevelKey:     "level",
		TimeKey:      "time",
		CallerKey:    "caller",
		EncodeCaller: zapcore.FullCallerEncoder,
		EncodeTime:   SyslogTimeEncoder,
		EncodeLevel:  CustomLevelEncoder,
	}

	var encoder zapcore.Encoder
	var level zapcore.Level
	encoder = zapcore.NewJSONEncoder(encoderConfig)
	level = zap.InfoLevel
	core := zapcore.NewCore(encoder, zapcore.AddSync(os.Stderr), level)
	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel), zap.AddCallerSkip(callerSkip)).Sugar()
	return logger
}
