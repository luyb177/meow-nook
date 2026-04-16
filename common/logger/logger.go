// Package logger wraps go.uber.org/zap and provides a single, project-wide
// logger that is initialised once at start-up via Init and then accessible
// through the module-level helpers (Info, Error, …).
package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var global *zap.Logger

// Level constants that can be used in configuration.
const (
	LevelDebug = "debug"
	LevelInfo  = "info"
	LevelWarn  = "warn"
	LevelError = "error"
)

// Config holds logger configuration.
type Config struct {
	Level      string // debug | info | warn | error
	Encoding   string // json | console
	OutputPath string // "stdout" or a file path
}

// Init initialises the global logger.  It is safe to call multiple times; each
// call replaces the previous global logger.
func Init(cfg Config) error {
	level := zap.InfoLevel
	if err := level.UnmarshalText([]byte(cfg.Level)); err != nil {
		level = zap.InfoLevel
	}

	encoding := cfg.Encoding
	if encoding == "" {
		encoding = "json"
	}

	outputPath := cfg.OutputPath
	if outputPath == "" {
		outputPath = "stdout"
	}

	encCfg := zapcore.EncoderConfig{
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
	}

	var enc zapcore.Encoder
	if encoding == "console" {
		enc = zapcore.NewConsoleEncoder(encCfg)
	} else {
		enc = zapcore.NewJSONEncoder(encCfg)
	}

	var sink zapcore.WriteSyncer
	if outputPath == "stdout" {
		sink = zapcore.Lock(os.Stdout)
	} else {
		f, err := os.OpenFile(outputPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
		if err != nil {
			return err
		}
		sink = zapcore.Lock(f)
	}

	core := zapcore.NewCore(enc, sink, zap.NewAtomicLevelAt(level))
	global = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	return nil
}

// mustInit initialises a no-op logger so that the helpers below never panic
// before Init has been called.
func init() {
	global = zap.NewNop()
}

// ──────────────────────────────────────────────
// Module-level helpers
// ──────────────────────────────────────────────

// With returns a child logger with extra fields attached.
func With(fields ...zap.Field) *zap.Logger {
	return global.With(fields...)
}

// Debug logs a message at DEBUG level.
func Debug(msg string, fields ...zap.Field) {
	global.Debug(msg, fields...)
}

// Info logs a message at INFO level.
func Info(msg string, fields ...zap.Field) {
	global.Info(msg, fields...)
}

// Warn logs a message at WARN level.
func Warn(msg string, fields ...zap.Field) {
	global.Warn(msg, fields...)
}

// Error logs a message at ERROR level.
func Error(msg string, fields ...zap.Field) {
	global.Error(msg, fields...)
}

// Fatal logs a message at FATAL level and then calls os.Exit(1).
func Fatal(msg string, fields ...zap.Field) {
	global.Fatal(msg, fields...)
}

// Sync flushes buffered log entries.  Call this before process exit.
func Sync() {
	_ = global.Sync()
}
