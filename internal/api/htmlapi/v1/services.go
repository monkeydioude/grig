package v1

import (
	"log/slog"
	cErrors "monkeydioude/grig/internal/errors"
	"monkeydioude/grig/internal/html/blocks"
	"monkeydioude/grig/internal/html/layouts"
	"monkeydioude/grig/internal/html/pages"
	"monkeydioude/grig/internal/model"
	"monkeydioude/grig/internal/service/services"
	"monkeydioude/grig/pkg/errors"
	pkgErrors "monkeydioude/grig/pkg/errors"
	"monkeydioude/grig/pkg/html/elements"
	"monkeydioude/grig/pkg/server/http_errors"
	"net/http"
	"os"
	"strconv"
)

func (h Handler) ServicesList(w http.ResponseWriter, r *http.Request, logger *slog.Logger, nav elements.Nav) error {
	layout := layouts.Main(nav, pages.ServicesList(&h.Layout.ServerConfig, logger))
	return layout.Render(r.Context(), w)
}

func (h Handler) AddServiceByFilename(
	w http.ResponseWriter,
	r *http.Request,
	logger *slog.Logger,
	p *services.ServiceFilename,
) error {
	if p == nil {
		return http_errors.InternalServerError(pkgErrors.Wrap(cErrors.ErrNilPointer, "AddServiceByFilename: *services.ServiceFilename"))
	}

	srv, err := p.TryLoadAndParse()
	if err != nil {
		if os.IsNotExist(err) {
			return http_errors.BadRequest(pkgErrors.Wrap(err, "AddServiceByFilename: *services.TryLoadAndParse"))
		}
		return http_errors.InternalServerError(pkgErrors.Wrap(err, "AddServiceByFilename: *services.TryLoadAndParse"))
	}
	h.Layout.ServerConfig.AppsServicesPaths = append(h.Layout.ServerConfig.AppsServicesPaths, srv.Path)
	return blocks.ServicesService(srv).Render(r.Context(), w)
}

func (h Handler) ServicesEnvironmentBlock(w http.ResponseWriter, r *http.Request, _ *slog.Logger) error {
	indexStr := r.URL.Query().Get("index")
	index := 0
	if indexStr != "" {
		it, err := strconv.Atoi(indexStr)
		if err != nil {
			return errors.Wrap(err, "JosukeBranchAction")
		}
		index = it
	}
	return blocks.ServicesEnvironmentBlock(index, "", model.Service{}).Render(r.Context(), w)
}
