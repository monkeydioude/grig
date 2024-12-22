package v1

import (
	"context"
	"monkeydioude/grig/internal/html/page_data"
	"monkeydioude/grig/internal/templ/layouts"
	"monkeydioude/grig/internal/templ/pages"
	"net/http"
)

func (h Handler) CapybaraList(w http.ResponseWriter, r *http.Request) error {
	data := page_data.CapybaraList(&h.Layout.ServerConfig)
	layout := layouts.Main(h.Layout.Navigation, data, pages.CapybaraList(data))
	layout.Render(context.Background(), w)
	return nil
}
