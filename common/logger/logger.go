package logger

import (
	"context"
	"errors"
	"os"
	"strings"
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	global *zap.Logger
	once   sync.Once
	inited bool
)

// Level constants.
const (
	LevelDebug = "debug"
	LevelInfo  = "info"
	LevelWarn  = "warn"
	LevelError = "error"
)

const (
	OutputStdout = "stdout"
)

type ctxKey string

const (
	ctxKeyRequestID ctxKey = "request_id"
	ctxKeyTraceID   ctxKey = "trace_id"
)

// Config holds logger configuration.
type Config struct {
	Level      string // debug | info | warn | error
	OutputPath string // "stdout" or a file path

	// If OutputPath is a file, also write to stdout/stderr.
	AlsoToStdout bool

	// Encodings:
	// - ConsoleEncoding: console/json (dev usually console)
	// - FileEncoding: json/console (file usually json)
	ConsoleEncoding string
	FileEncoding    string

	Development      bool
	EnableStacktrace bool
	Service          string
	Env              string
	InitialFields    map[string]string
}

func init() {
	// Make helpers safe before Init.
	global = zap.NewNop()
}

// Init must be called ONCE during startup.
// If called more than once, it returns an error.
func Init(cfg Config) error {
	called := false
	var initErr error

	once.Do(func() {
		called = true
		initErr = initOnce(cfg)
		if initErr == nil {
			inited = true
		}
	})

	if !called {
		return errors.New("logger: Init already called (this logger supports Init once per process)")
	}
	return initErr
}

func initOnce(cfg Config) error {
	// ---- level
	level := zap.InfoLevel
	if strings.TrimSpace(cfg.Level) != "" {
		if err := level.UnmarshalText([]byte(cfg.Level)); err != nil {
			return err
		}
	}

	// ---- defaults
	outputPath := strings.TrimSpace(cfg.OutputPath)
	if outputPath == "" {
		outputPath = OutputStdout
	}

	consoleEncName := strings.TrimSpace(cfg.ConsoleEncoding)
	fileEncName := strings.TrimSpace(cfg.FileEncoding)

	if consoleEncName == "" {
		if cfg.Development {
			consoleEncName = "console"
		} else {
			consoleEncName = "json"
		}
	}
	if fileEncName == "" {
		fileEncName = "json"
	}

	// ---- cores
	var cores []zapcore.Core

	if outputPath == OutputStdout {
		consoleEnc := newEncoder(consoleEncName, cfg.Development)
		stdoutCore, stderrCore := buildStdCores(consoleEnc, level)
		cores = append(cores, stdoutCore, stderrCore)
	} else {
		// file core
		fileEnc := newEncoder(fileEncName, false)
		cores = append(cores, buildFileCore(fileEnc, level, outputPath))

		if cfg.AlsoToStdout {
			consoleEnc := newEncoder(consoleEncName, cfg.Development)
			stdoutCore, stderrCore := buildStdCores(consoleEnc, level)
			cores = append(cores, stdoutCore, stderrCore)
		}
	}

	if len(cores) == 0 {
		return errors.New("logger: no output cores configured")
	}

	core := zapcore.NewTee(cores...)

	// ---- options
	opts := []zap.Option{
		zap.AddCaller(),
		zap.AddCallerSkip(1),
	}

	if cfg.EnableStacktrace {
		opts = append(opts, zap.AddStacktrace(zapcore.ErrorLevel))
	}

	// fixed fields
	fixed := make([]zap.Field, 0, 2+len(cfg.InitialFields))
	if cfg.Service != "" {
		fixed = append(fixed, zap.String("service", cfg.Service))
	}
	if cfg.Env != "" {
		fixed = append(fixed, zap.String("env", cfg.Env))
	}
	for k, v := range cfg.InitialFields {
		k = strings.TrimSpace(k)
		if k == "" {
			continue
		}
		fixed = append(fixed, zap.String(k, v))
	}
	if len(fixed) > 0 {
		opts = append(opts, zap.Fields(fixed...))
	}

	global = zap.New(core, opts...)
	return nil
}

func newEncoder(encoding string, development bool) zapcore.Encoder {
	encCfg := zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
		EncodeDuration: zapcore.MillisDurationEncoder,
	}

	if encoding == "console" && development {
		encCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
	} else if encoding == "console" {
		encCfg.EncodeLevel = zapcore.CapitalLevelEncoder
	} else {
		encCfg.EncodeLevel = zapcore.LowercaseLevelEncoder
	}

	if encoding == "console" {
		return zapcore.NewConsoleEncoder(encCfg)
	}
	return zapcore.NewJSONEncoder(encCfg)
}

func buildStdCores(enc zapcore.Encoder, min zapcore.Level) (zapcore.Core, zapcore.Core) {
	stdoutCore := zapcore.NewCore(
		enc,
		zapcore.Lock(os.Stdout),
		zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return lvl < zapcore.ErrorLevel && lvl >= min
		}),
	)
	stderrCore := zapcore.NewCore(
		enc,
		zapcore.Lock(os.Stderr),
		zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return lvl >= zapcore.ErrorLevel && lvl >= min
		}),
	)
	return stdoutCore, stderrCore
}

func buildFileCore(enc zapcore.Encoder, min zapcore.Level, path string) zapcore.Core {
	fileSink := zapcore.AddSync(&lumberjack.Logger{
		Filename:   path,
		MaxSize:    100,
		MaxBackups: 10,
		MaxAge:     7,
		Compress:   true,
	})
	return zapcore.NewCore(enc, fileSink, zap.NewAtomicLevelAt(min))
}

// ──────────────────────────────────────────────
// Context helpers
// ──────────────────────────────────────────────

func WithRequestID(ctx context.Context, requestID string) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	return context.WithValue(ctx, ctxKeyRequestID, strings.TrimSpace(requestID))
}

func RequestIDFromContext(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if v, ok := ctx.Value(ctxKeyRequestID).(string); ok {
		return v
	}
	return ""
}

func WithTraceID(ctx context.Context, traceID string) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	return context.WithValue(ctx, ctxKeyTraceID, strings.TrimSpace(traceID))
}

func TraceIDFromContext(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if v, ok := ctx.Value(ctxKeyTraceID).(string); ok {
		return v
	}
	return ""
}

func FromContext(ctx context.Context) *zap.Logger {
	l := global
	if rid := RequestIDFromContext(ctx); rid != "" {
		l = l.With(zap.String("request_id", rid))
	}
	if tid := TraceIDFromContext(ctx); tid != "" {
		l = l.With(zap.String("trace_id", tid))
	}
	return l
}

// ──────────────────────────────────────────────
// Module-level helpers
// ──────────────────────────────────────────────

func With(fields ...zap.Field) *zap.Logger { return global.With(fields...) }

func Debug(msg string, fields ...zap.Field) { global.Debug(msg, fields...) }
func Info(msg string, fields ...zap.Field)  { global.Info(msg, fields...) }
func Warn(msg string, fields ...zap.Field)  { global.Warn(msg, fields...) }
func Error(msg string, fields ...zap.Field) { global.Error(msg, fields...) }
func Fatal(msg string, fields ...zap.Field) { global.Fatal(msg, fields...) }

func Sync() {
	// before Init it's Nop, safe.
	_ = global.Sync()
}

// Useful helper
func SinceMS(start time.Time) int64 { return time.Since(start).Milliseconds() }

// Optional: if you want to check initialization status.
func Inited() bool { return inited }
