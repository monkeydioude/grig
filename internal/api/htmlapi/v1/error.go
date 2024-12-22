package v1

import (
	"context"
	"monkeydioude/grig/internal/html/page_data"
	"monkeydioude/grig/internal/templ/layouts"
	"monkeydioude/grig/internal/templ/pages"
	"net/http"
)

func (h Handler) NotFound(w http.ResponseWriter, r *http.Request) error {
	layout := layouts.Main(h.Layout.Navigation, page_data.Error("Not found"), pages.NotFound())
	layout.Render(context.Background(), w)
	return nil
}
