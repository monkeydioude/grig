package v1

import (
	"context"
	"monkeydioude/grig/internal/html/blocks"
	"monkeydioude/grig/internal/html/layouts"
	"monkeydioude/grig/internal/html/pages"
	"monkeydioude/grig/internal/model"
	"monkeydioude/grig/pkg/html/elements"
	"net/http"
	"strconv"

	"github.com/pkg/errors"
)

func (h Handler) JosukeList(w http.ResponseWriter, r *http.Request, nav elements.Nav) error {
	layout := layouts.Main(nav, pages.JosukeList(&h.Layout.ServerConfig))
	return layout.Render(context.Background(), w)
}

func (h Handler) JosukeHookBlock(w http.ResponseWriter, r *http.Request) error {
	indexStr := r.URL.Query().Get("index")
	index := 0
	if indexStr != "" {
		it, err := strconv.Atoi(indexStr)
		if err != nil {
			return errors.Wrap(err, "JosukeHookBlock")
		}
		index = it
	}

	return blocks.JosukeHook(index, pages.GetHookInputName, model.Hook{}).Render(context.Background(), w)
}
