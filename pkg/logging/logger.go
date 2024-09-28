package logging

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"time"
)

var l *slog.Logger

type Logger struct {
	*slog.Logger
}

type Config struct {
	SystemName string
	Env        string
	Level      string
	Output     io.Writer
}

const (
	LevelTrace = slog.Level(-8)
	LevelFatal = slog.Level(12)
	LevelPanic = slog.Level(13)
	LevelInfo  = slog.LevelInfo
	LevelDebug = slog.LevelDebug
	LevelWarn  = slog.LevelWarn
	LevelError = slog.LevelError
)

var LevelNames = map[slog.Leveler]string{
	LevelTrace: "TRACE",
	LevelFatal: "FATAL",
	LevelPanic: "PANIC",
	LevelInfo:  "INFO",
	LevelDebug: "DEBUG",
	LevelWarn:  "WARN",
	LevelError: "ERROR",
}

func InitLogging(cfg *Config) error {
	if cfg.SystemName == "" {
		return fmt.Errorf("system.name is required property for logging")
	}
	if cfg.Env == "" {
		return fmt.Errorf("env is required property for logging")
	}
	if cfg.Output == nil {
		cfg.Output = os.Stdout
	}
	if cfg.Level == "" {
		cfg.Level = "INFO"
	}

	logLevel, err := getLevelLogging(cfg.Level)
	if err != nil {
		return err
	}

	l = slog.New(slog.NewJSONHandler(cfg.Output, &slog.HandlerOptions{
		Level: logLevel,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.LevelKey {
				level := a.Value.Any().(slog.Level)
				levelLabel, exists := LevelNames[level]
				if !exists {
					levelLabel = level.String()
				}

				a.Value = slog.StringValue(levelLabel)
			}
			if a.Key == slog.TimeKey {

				a.Value = slog.StringValue(time.Now().Format("2006-01-02T15:04:05.000Z"))
			}

			return a
		}}))
	l = l.With("system.name", cfg.SystemName)
	l = l.With("env", cfg.Env)

	return nil
}

func getLevelLogging(level string) (slog.Level, error) {
	switch level {
	case "INFO":
		return LevelInfo, nil
	case "TRACE":
		return LevelTrace, nil
	case "FATAL":
		return LevelFatal, nil
	case "PANIC":
		return LevelPanic, nil
	case "DEBUG":
		return LevelDebug, nil
	case "WARN":
		return LevelWarn, nil
	case "ERROR":
		return LevelError, nil
	default:
		return 0, errors.New("invalid level logging")
	}
}

func Info(arg string) {
	l.Info(arg)
}

func Warn(arg string) {
	l.Warn(arg)
}

func Error(arg string) {
	l.Error(arg)
}

func Debug(arg string) {
	l.Debug(arg)
}

func Fatal(arg string) {
	l.Log(context.Background(), LevelFatal, arg)
	os.Exit(1)
}

func Trace(arg string) {
	l.Log(context.Background(), LevelTrace, arg)
}

func Panic(arg string) {
	l.Log(context.Background(), LevelPanic, arg)
	panic(arg)
}

func Infof(format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	l.Info(msg)
}

func Warnf(format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	l.Warn(msg)
}

func Errorf(format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	l.Error(msg)
}

func Debugf(format string, args ...any) {
	l.Debug(fmt.Sprintf(format, args...))
}

func Fatalf(format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	l.Log(context.Background(), LevelFatal, msg)
	os.Exit(1)
}

func Tracef(format string, args ...any) {
	l.Log(context.Background(), LevelTrace, fmt.Sprintf(format, args...))
}

func Panicf(format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	l.Log(context.Background(), LevelPanic, msg)
	panic(msg)
}
