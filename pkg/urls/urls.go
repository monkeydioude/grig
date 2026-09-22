// Package urls builds links for an app served under a sub-path, such as
// /grig, by a reverse proxy that forwards the prefix instead of stripping
// it. Routes stay declared relative to the app root; the prefix is added
// when a link is rendered.
package urls

import (
	"context"
	"net/http"
	"strings"
)

// Prefixer prepends the app's base path to app-relative paths. Its zero
// value prepends nothing, so an app served at the root needs no special
// casing.
type Prefixer struct {
	base string
}

// New normalizes a configured base path: "grig", "/grig" and "/grig/"
// all become "/grig". An empty base path serves the app at the root.
func New(basePath string) Prefixer {
	trimmed := strings.Trim(strings.TrimSpace(basePath), "/")
	if trimmed == "" {
		return Prefixer{}
	}
	return Prefixer{base: "/" + trimmed}
}

// Base is the normalized base path, empty when serving at the root.
func (p Prefixer) Base() string {
	return p.base
}

// Path prefixes an app-relative path. Anything that is not a rooted path
// (absolute URLs, protocol-relative URLs, anchors, query-only links) is
// left alone.
func (p Prefixer) Path(rel string) string {
	if p.base == "" || !strings.HasPrefix(rel, "/") || strings.HasPrefix(rel, "//") {
		return rel
	}
	if rel == "/" {
		return p.base + "/"
	}
	return p.base + rel
}

type ctxKey struct{}

// NewContext carries p to the templates rendered for a request.
func NewContext(ctx context.Context, p Prefixer) context.Context {
	return context.WithValue(ctx, ctxKey{}, p)
}

// FromContext returns the carried Prefixer, or a zero one (no prefix)
// when none was set.
func FromContext(ctx context.Context) Prefixer {
	p, _ := ctx.Value(ctxKey{}).(Prefixer)
	return p
}

// Path is the shorthand templates use: urls.Path(ctx, "/josuke").
func Path(ctx context.Context, rel string) string {
	return FromContext(ctx).Path(rel)
}

// Middleware makes p available to every request's templates.
func Middleware(p Prefixer) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r.WithContext(NewContext(r.Context(), p)))
		})
	}
}
