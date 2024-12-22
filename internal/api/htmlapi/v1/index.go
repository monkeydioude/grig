package v1

import (
	"context"
	"monkeydioude/grig/internal/html/page_data"
	"monkeydioude/grig/internal/templ/layouts"
	"monkeydioude/grig/internal/templ/pages"
	"net/http"
)

func (h Handler) Index(w http.ResponseWriter, r *http.Request) error {
	if r.URL.Path != "/" {
		return h.NotFound(w, r)
	}
	page := page_data.Index(&h.Layout.ServerConfig)
	layout := layouts.Main(h.Layout.Navigation, page, pages.Index(page))
	layout.Render(context.Background(), w)
	return nil
}
