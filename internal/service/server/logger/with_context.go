package logger

import (
	"context"
	"log/slog"
)

type slogWithContext struct {
	handler slog.Handler
	ctxKeys []string
}

func (swc *slogWithContext) Enabled(ctx context.Context, lvl slog.Level) bool {
	return swc.handler.Enabled(ctx, lvl)
}
func (swc *slogWithContext) Handle(ctx context.Context, rec slog.Record) error {
	attrs := make([]slog.Attr, len(swc.ctxKeys))
	for _, key := range swc.ctxKeys {
		val := ctx.Value(key)
		if val == nil {
			continue
		}
		attrs = append(attrs, slog.Any(key, val))
	}
	return swc.handler.WithAttrs(attrs).Handle(ctx, rec)
}

func (swc *slogWithContext) WithAttrs(attrs []slog.Attr) slog.Handler {
	return swc.handler.WithAttrs(attrs)
}

func (swc *slogWithContext) WithGroup(name string) slog.Handler {
	return swc.handler.WithGroup(name)
}

func SlogWithContext(handler slog.Handler, ctxKeys ...string) *slogWithContext {
	return &slogWithContext{
		handler: handler,
		ctxKeys: ctxKeys,
	}
}
