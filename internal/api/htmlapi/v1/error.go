package v1

import (
	"context"
	"monkeydioude/grig/internal/errors"
	"monkeydioude/grig/internal/html/layouts"
	"monkeydioude/grig/internal/html/pages"
	"net/http"
)

func (h Handler) NotFound(w http.ResponseWriter, r *http.Request) error {
	layout := layouts.Main(h.Layout.Navigation, pages.Error(errors.ErrInternalServerError.Error()))
	layout.Render(context.Background(), w)
	return errors.NotFound()
}

func (h Handler) InternalServer(w http.ResponseWriter, r *http.Request) error {
	layout := layouts.Main(h.Layout.Navigation, pages.Error(errors.ErrInternalServerError.Error()))
	layout.Render(context.Background(), w)
	return errors.InternalServerError(errors.ErrInternalServerError)
}
