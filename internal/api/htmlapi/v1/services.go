package v1

import (
	"fmt"
	"log/slog"
	cErrors "monkeydioude/grig/internal/errors"
	"monkeydioude/grig/internal/html/blocks"
	"monkeydioude/grig/internal/html/layouts"
	"monkeydioude/grig/internal/html/pages"
	"monkeydioude/grig/internal/model"
	"monkeydioude/grig/internal/service/services"
	"monkeydioude/grig/pkg/dt"
	pkgErrors "monkeydioude/grig/pkg/errors"
	"monkeydioude/grig/pkg/html/elements"
	"monkeydioude/grig/pkg/server/http_errors"
	"monkeydioude/grig/pkg/trans_types"
	"net/http"
	"os"
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
	index, err := trans_types.AtoiOr0(r.URL.Query().Get("index"))
	if err != nil {
		logger.Warn("ServicesEnvironmentBlock AtoiOr0", "error", err)
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
	return blocks.ServicesService(index, srv).Render(r.Context(), w)
}

func (h Handler) ServicesEnvironmentBlock(w http.ResponseWriter, r *http.Request, logger *slog.Logger) error {
	index, err := trans_types.AtoiOr0(r.URL.Query().Get("index"))
	if err != nil {
		logger.Warn("ServicesEnvironmentBlock AtoiOr0", "error", err)
	}
	return blocks.ServicesEnvironmentBlock(r.URL.Query().Get("parent_name"), index, "", model.Service{
		Name: r.URL.Query().Get("parent_name"),
	}).Render(r.Context(), w)
}

func (h Handler) ServicesServiceBlock(w http.ResponseWriter, r *http.Request, logger *slog.Logger) error {
	index, err := trans_types.AtoiOr0(r.URL.Query().Get("index"))
	if err != nil {
		logger.Warn("ServicesEnvironmentBlock AtoiOr0", "error", err)
	}
	fmt.Println(index, r.URL.Query().Get("parent_name"))
	return blocks.ServicesService(index, model.Service{}).Render(r.Context(), w)
}
