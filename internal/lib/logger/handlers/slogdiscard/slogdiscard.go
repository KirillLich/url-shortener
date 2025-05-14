package slogdiscard

import (
	"context"
	"log/slog"
)

type DiscardHandler struct{}

func NewDiscardLogger() *slog.Logger {
	return slog.New(NewDiscardHandler())
}

func NewDiscardHandler() slog.Handler {
	return &DiscardHandler{}
}

func (_ DiscardHandler) Enabled(_ context.Context, _ slog.Level) bool {
	return false
}

func (_ DiscardHandler) Handle(_ context.Context, _ slog.Record) error {
	return nil
}

func (h *DiscardHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return h
}

func (h *DiscardHandler) WithGroup(name string) slog.Handler {
	return h
}
