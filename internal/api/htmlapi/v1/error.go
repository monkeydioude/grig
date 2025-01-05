package v1

import (
	"context"
	"monkeydioude/grig/internal/html/layouts"
	"monkeydioude/grig/internal/html/pages"
	"monkeydioude/grig/pkg/server/http_errors"
	"net/http"
)

func (h Handler) NotFound(w http.ResponseWriter, r *http.Request) error {
	layout := layouts.Main(h.Layout.Navigation, pages.Error(http_errors.ErrInternalServerError.Error()))
	layout.Render(context.Background(), w)
	return http_errors.NotFound()
}

func (h Handler) InternalServer(w http.ResponseWriter, r *http.Request) error {
	layout := layouts.Main(h.Layout.Navigation, pages.Error(http_errors.ErrInternalServerError.Error()))
	layout.Render(context.Background(), w)
	return http_errors.InternalServerError(http_errors.ErrInternalServerError)
}
