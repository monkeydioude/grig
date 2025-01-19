package v1

import (
	"log/slog"
	customErrors "monkeydioude/grig/internal/errors"
	"monkeydioude/grig/internal/model"
	"net/http"
)

func (h Handler) ServicesSave(
	w http.ResponseWriter,
	r *http.Request,
	logger *slog.Logger,
	srvcs *map[string]model.Service,
) error {
	if r == nil || srvcs == nil {
		logger.Error("api.ServicesSave", "error", customErrors.ErrNilPointer)
		return customErrors.ErrInvalidProvidedParameters
	}

	appSrvcs, err := h.Layout.ServerConfig.AppsServicesPaths.WithUpdate(srvcs, logger)
	if err != nil {
		logger.Error("api.ServicesSave::WithUpdate", "error", err)
		return customErrors.ErrServicesServicesUpdateFail
	}

	h.Layout.ServerConfig.AppsServicesPaths = appSrvcs
	if err := h.Layout.ServerConfig.Save(); err != nil {
		logger.Error("api.ServicesSave::Save", "error", err)
		return customErrors.ErrConfigSaveFail
	}
	return nil
}
