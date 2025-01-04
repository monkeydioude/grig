package v1

import (
	"context"
	"monkeydioude/grig/internal/html/elements"
	"monkeydioude/grig/internal/html/layouts"
	"monkeydioude/grig/internal/html/pages"
	"net/http"
)

func (h Handler) JosukeList(w http.ResponseWriter, r *http.Request, nav elements.Nav) error {
	layout := layouts.Main(nav, pages.JosukeList(&h.Layout.ServerConfig))
	return layout.Render(context.Background(), w)
}
