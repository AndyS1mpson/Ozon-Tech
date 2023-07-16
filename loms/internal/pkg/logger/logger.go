package logger

import (
	"io"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	global *zap.SugaredLogger
	defaultLevel = zap.NewAtomicLevelAt(zap.InfoLevel)
)

// Init logger singleton
func init() {
	global = New(defaultLevel, os.Stdout)
}

// Create a logger
func New(level zap.AtomicLevel, sink io.Writer, opts ...zap.Option) *zap.SugaredLogger {
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(zapcore.EncoderConfig{
			TimeKey:        "ts",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			MessageKey:     "message",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		}),
		zapcore.AddSync(sink),
		level,
	)

	return zap.New(core, opts...).Sugar()
}

// Set logger environment
func SetLoggerByEnvironment(environment string) {
	if environment == "PRODUCTION" {
		global = New(
			zap.NewAtomicLevelAt(zap.ErrorLevel),
			os.Stdout,
			zap.WithCaller(false),
			zap.AddStacktrace(zap.NewAtomicLevelAt(zap.PanicLevel)))
	}
}

// Log info message
func Info(args ...interface{}) {
	global.Info(args...)
}

// Log error
func Error(args ...interface{}) {
	global.Error(args...)
}

// func Errorf(args ...interface{}) {
// 	global.Errorf(args...)
// }

// Log fatal error
func Fatalf(template string, args ...interface{}) {
	global.Fatalf(template, args...)
}


// func withTraceID(ctx context.Context) *zap.SugaredLogger {
// 	span := opentracing.SpanFromContext(ctx)
// 	if span == nil {
// 		return global
// 	}

// 	if sc, ok := span.Context().(jaeger.SpanContext); ok {
// 		return global.Desugar().With(
// 			zap.Stringer("trace_id", sc.TraceID()),
// 			zap.Stringer("span_id", sc.SpanID()),
// 		).Sugar()
// 	}

// 	return global
// }

