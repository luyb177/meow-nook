// Package logger wraps go.uber.org/zap and provides a single, project-wide
// logger that is initialised once at start-up via Init and then accessible
// through the module-level helpers (Info, Error, …).
package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// todo trace 支持
// todo context 支持

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
	if global != nil {
		_ = global.Sync()
	}

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
		TimeKey:      "ts",
		LevelKey:     "level",
		CallerKey:    "caller",
		MessageKey:   "msg",
		EncodeLevel:  zapcore.LowercaseLevelEncoder,
		EncodeTime:   zapcore.ISO8601TimeEncoder,
		EncodeCaller: zapcore.ShortCallerEncoder,
	}

	var enc zapcore.Encoder
	if encoding == "console" {
		enc = zapcore.NewConsoleEncoder(encCfg)
	} else {
		enc = zapcore.NewJSONEncoder(encCfg)
	}

	var core zapcore.Core

	if cfg.OutputPath == "stdout" {
		stdoutCore := zapcore.NewCore(
			enc,
			zapcore.Lock(os.Stdout),
			zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
				return lvl < zapcore.ErrorLevel && lvl >= level
			}),
		)

		stderrCore := zapcore.NewCore(
			enc,
			zapcore.Lock(os.Stderr),
			zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
				return lvl >= zapcore.ErrorLevel
			}),
		)

		core = zapcore.NewTee(stdoutCore, stderrCore)

	} else {
		fileSink := zapcore.AddSync(&lumberjack.Logger{
			Filename:   cfg.OutputPath,
			MaxSize:    100,
			MaxBackups: 10,
			MaxAge:     7,
			Compress:   true,
		})

		core = zapcore.NewCore(enc, fileSink, zap.NewAtomicLevelAt(level))
	}

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

// Sugared returns a sugared logger that provides a more ergonomic, but slightly slower, API.
// For example, instead of `logger.Info("Failed to fetch URL.", zap.String("url", url), zap.Int("attempt", 3))`,
// you can write `logger.Sugared().Infow("Failed to fetch URL.", "url", url, "attempt", 3)`.
func Sugared() *zap.SugaredLogger {
	return global.Sugar()
}
