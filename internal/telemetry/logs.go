package telemetry

import (
	"context"
	"log/slog"
	"os"

	"go.opentelemetry.io/contrib/bridges/otelslog"

	sdkLog "go.opentelemetry.io/otel/sdk/log"
)

type Logger struct {
	*slog.Logger
}

type FanoutHandler struct {
	handlers []slog.Handler
}

func NewFanoutHandler(handlers ...slog.Handler) *FanoutHandler {
	return &FanoutHandler{handlers: handlers}
}

func (h *FanoutHandler) Enabled(ctx context.Context, level slog.Level) bool {
	for _, handler := range h.handlers {
		if handler.Enabled(ctx, level) {
			return true
		}
	}
	return false
}

func (h *FanoutHandler) Handle(ctx context.Context, record slog.Record) error {
	for _, handler := range h.handlers {
		_ = handler.Handle(ctx, record)
	}
	return nil
}

func (h *FanoutHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	handlers := make([]slog.Handler, 0, len(h.handlers))
	for _, handler := range h.handlers {
		handlers = append(handlers, handler.WithAttrs(attrs))
	}
	return &FanoutHandler{handlers: handlers}
}

func (h *FanoutHandler) WithGroup(name string) slog.Handler {
	handlers := make([]slog.Handler, 0, len(h.handlers))
	for _, handler := range h.handlers {
		handlers = append(handlers, handler.WithGroup(name))
	}
	return &FanoutHandler{handlers: handlers}
}

func NewLogger(name string, loggerProvider *sdkLog.LoggerProvider) *Logger {
	stdout := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})

	if loggerProvider == nil {
		return &Logger{slog.New(stdout)}
	}

	otel := otelslog.NewHandler(name, otelslog.WithLoggerProvider(loggerProvider))
	fanout := NewFanoutHandler(stdout, otel)

	return &Logger{slog.New(fanout)}
}
