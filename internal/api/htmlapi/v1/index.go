package v1

import (
	"context"
	"monkeydioude/grig/internal/html/layouts"
	"monkeydioude/grig/internal/html/pages"
	"net/http"
)

func (h Handler) Index(w http.ResponseWriter, r *http.Request) error {
	if r.URL.Path != "/" {
		return h.NotFound(w, r)
	}
	layout := layouts.Main(h.Layout.Navigation, pages.Index(&h.Layout.ServerConfig))
	layout.Render(context.Background(), w)
	return nil
}
