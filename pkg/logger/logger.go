package logger

import (
	"context"
	"log/slog"
	"os"
)

type Logger struct {
	*slog.Logger
}

func New() *Logger {
	opts := &slog.HandlerOptions{
		Level:     slog.LevelInfo,
		AddSource: true,
	}

	handler := slog.NewTextHandler(os.Stdout, opts)
	logger := slog.New(handler)

	return &Logger{Logger: logger}
}

func (l *Logger) Info(ctx context.Context, msg string, args ...any) {
	args = append([]any{"level", "INFO"}, args...)
	l.Logger.InfoContext(ctx, msg, args...)
}

func (l *Logger) Warn(ctx context.Context, msg string, args ...any) {
	args = append([]any{"level", "WARN"}, args...)
	l.Logger.WarnContext(ctx, msg, args...)
}

func (l *Logger) Error(ctx context.Context, msg string, args ...any) {
	args = append([]any{"level", "ERROR"}, args...)
	l.Logger.ErrorContext(ctx, msg, args...)
}

func (l *Logger) Debug(ctx context.Context, msg string, args ...any) {
	args = append([]any{"level", "DEBUG"}, args...)
	l.Logger.DebugContext(ctx, msg, args...)
}

func (l *Logger) Critical(ctx context.Context, msg string, args ...any) {
	args = append([]any{"level", "CRITICAL"}, args...)
	l.Logger.ErrorContext(ctx, msg, args...)
}
func (l *Logger) Panic(ctx context.Context, msg string, args ...any) {
	args = append([]any{"level", "PANIC"}, args...)
	l.Logger.ErrorContext(ctx, msg, args...)
	panic(msg)
}
