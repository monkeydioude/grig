package v1

import (
	"context"
	"monkeydioude/grig/internal/html/layouts"
	"monkeydioude/grig/internal/html/pages"
	"monkeydioude/grig/pkg/html/elements"
	"net/http"
)

func (h Handler) CapybaraList(w http.ResponseWriter, r *http.Request, nav elements.Nav) error {
	layout := layouts.Main(nav, pages.CapybaraList(&h.Layout.ServerConfig))
	layout.Render(context.Background(), w)
	return nil
}
