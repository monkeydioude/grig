package v1

import (
	"context"
	"log/slog"
	"monkeydioude/grig/internal/html/blocks"
	"monkeydioude/grig/internal/html/layouts"
	"monkeydioude/grig/internal/html/pages"
	"monkeydioude/grig/internal/model"
	"monkeydioude/grig/pkg/errors"
	"monkeydioude/grig/pkg/html/elements"
	"net/http"
	"strconv"
)

func (h Handler) CapybaraList(w http.ResponseWriter, r *http.Request, _ *slog.Logger, nav elements.Nav) error {
	layout := layouts.Main(nav, pages.CapybaraList(&h.Layout.ServerConfig))
	layout.Render(context.Background(), w)
	return nil
}

func (h Handler) CapybaraServiceBlock(w http.ResponseWriter, r *http.Request, _ *slog.Logger) error {
	indexStr := r.URL.Query().Get("index")
	index := 0
	if indexStr != "" {
		it, err := strconv.Atoi(indexStr)
		if err != nil {
			return errors.Wrap(err, "CapybaraServiceBlock")
		}
		index = it
	}
	return blocks.CapybaraService(pages.GetServiceInputName, index, model.ServiceDefinition{}).Render(context.Background(), w)
}
