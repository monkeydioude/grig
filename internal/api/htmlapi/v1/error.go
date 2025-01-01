package v1

import (
	"context"
	"monkeydioude/grig/internal/html/layouts"
	"monkeydioude/grig/internal/html/pages"
	"net/http"
)

func (h Handler) NotFound(w http.ResponseWriter, r *http.Request) error {
	layout := layouts.Main(h.Layout.Navigation, pages.Error("Not found"))
	w.WriteHeader(http.StatusNotFound)
	layout.Render(context.Background(), w)
	return nil
}
