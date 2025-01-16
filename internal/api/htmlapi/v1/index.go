package v1

import (
	"log/slog"
	"monkeydioude/grig/internal/html/layouts"
	"monkeydioude/grig/internal/html/pages"
	"monkeydioude/grig/pkg/html/elements"
	"net/http"
)

func (h Handler) Index(w http.ResponseWriter, r *http.Request, _ *slog.Logger, nav elements.Nav) error {
	if r.URL.Path != "/" {
		return h.NotFound(w, r)
	}
	layout := layouts.Main(nav, pages.Index(&h.Layout.ServerConfig))
	layout.Render(r.Context(), w)
	return nil
}
