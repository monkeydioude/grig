package v1

import (
	"log/slog"
	cErrors "monkeydioude/grig/internal/errors"
	"monkeydioude/grig/internal/html/blocks"
	"monkeydioude/grig/internal/html/layouts"
	"monkeydioude/grig/internal/html/pages"
	"monkeydioude/grig/internal/model"
	"monkeydioude/grig/internal/service/services"
	"monkeydioude/grig/pkg/dt"
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

func (h Handler) AddServiceByFilepath(
	w http.ResponseWriter,
	r *http.Request,
	logger *slog.Logger,
	p *services.Filepath,
) error {
	if p == nil {
		return http_errors.InternalServerError(pkgErrors.Wrap(cErrors.ErrNilPointer, "AddServiceByFilename: *services.ServiceFilename"))
	}

	srv, err := p.TryLoadAndParse()
	if err != nil {
		logger.Error("AddServiceByFilename: *services.TryLoadAndParse", "error", err)
		if os.IsNotExist(err) {
			return http_errors.BadRequest(cErrors.ErrServicesInvalidFilepath)
		}
		return http_errors.InternalServerError(cErrors.ErrServicesUnableFileParsing)
	}
	// we don't want to add an already existing file
	slice, res := dt.AppendUnique(h.Layout.ServerConfig.AppsServicesPaths, srv.Path)
	if res == false {
		return http_errors.BadRequest(cErrors.ErrServicesFilepathExists)
	}
	h.Layout.ServerConfig.AppsServicesPaths = slice
	if err := h.Layout.ServerConfig.Save(); err != nil {
		return http_errors.InternalServerError(cErrors.ErrWritingFile)
	}
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
	return blocks.ServicesEnvironmentBlock(index, "", model.Service{
		Name: r.URL.Query().Get("parent_name"),
	}).Render(r.Context(), w)
}
