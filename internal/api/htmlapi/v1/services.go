package v1

import (
	"context"
	"log/slog"
	"monkeydioude/grig/internal/html/layouts"
	"monkeydioude/grig/internal/html/pages"
	"monkeydioude/grig/pkg/html/elements"
	"net/http"
)

func (h Handler) ServicesList(w http.ResponseWriter, r *http.Request, logger *slog.Logger, nav elements.Nav) error {
	layout := layouts.Main(nav, pages.ServicesList(&h.Layout.ServerConfig, logger))
	return layout.Render(context.Background(), w)
}
