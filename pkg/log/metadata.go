package log

import (
	"context"
	"log/slog"
)

type contextKey string

const contextKeyMetadata contextKey = "log/metadata"

type Metadata map[string]any

type metadataHandler struct {
	next slog.Handler
}

func (h *metadataHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.next.Enabled(ctx, level)
}

func (h *metadataHandler) Handle(ctx context.Context, record slog.Record) error {
	if value, ok := ctx.Value(contextKeyMetadata).(Metadata); ok && len(value) > 0 {
		record = record.Clone()
		record.AddAttrs(slog.Any(attrKeyMetadata, value))
	}
	return h.next.Handle(ctx, record)
}

func (h *metadataHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &metadataHandler{next: h.next.WithAttrs(attrs)}
}

func (h *metadataHandler) WithGroup(name string) slog.Handler {
	return &metadataHandler{next: h.next.WithGroup(name)}
}

func withMetadataHandler(handler slog.Handler) slog.Handler {
	return &metadataHandler{next: handler}
}

func AddMetadataToContext(ctx context.Context, metadata Metadata) context.Context {
	value := metadata
	if v, ok := ctx.Value(contextKeyMetadata).(Metadata); ok && len(v) > 0 {
		value = Metadata{}
		for k, v := range v {
			value[k] = v
		}
		for k, v := range metadata {
			value[k] = v
		}
	}
	return context.WithValue(ctx, contextKeyMetadata, value)
}
