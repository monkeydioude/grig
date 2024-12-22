package v1

import (
	"context"
	"monkeydioude/grig/internal/html/layouts"
	"monkeydioude/grig/internal/html/pages"
	"net/http"
)

func (h Handler) CapybaraList(w http.ResponseWriter, r *http.Request) error {
	layout := layouts.Main(h.Layout.Navigation, pages.CapybaraList(&h.Layout.ServerConfig))
	layout.Render(context.Background(), w)
	return nil
}
