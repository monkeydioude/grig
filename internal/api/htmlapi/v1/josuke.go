package v1

import (
	"context"
	"monkeydioude/grig/internal/html/element"
	"monkeydioude/grig/internal/html/layouts"
	"monkeydioude/grig/internal/html/pages"
	"net/http"
)

func (h Handler) JosukeList(w http.ResponseWriter, r *http.Request, nav element.Nav) error {
	layout := layouts.Main(nav, pages.CapybaraList(&h.Layout.ServerConfig))
	layout.Render(context.Background(), w)
	return nil
}
