package v1

import (
	"context"
	"log/slog"
	customErrors "monkeydioude/grig/internal/errors"
	"monkeydioude/grig/internal/model"
	"monkeydioude/grig/pkg/os"
	"net/http"
	"time"
)

type ServicesPayload struct {
	Services []model.Service
}

func (h Handler) ServicesSave(
	w http.ResponseWriter,
	r *http.Request,
	logger *slog.Logger,
	srvcs *ServicesPayload,
) error {
	if r == nil || srvcs == nil {
		logger.Error("api.ServicesSave", "error", customErrors.ErrNilPointer)
		return customErrors.ErrInvalidProvidedParameters
	}

	appSrvcs, err := h.Layout.ServerConfig.AppsServicesPaths.WithUpdate((*srvcs).Services, logger)
	if err != nil {
		logger.Error("api.ServicesSave::WithUpdate", "error", err)
		return customErrors.ErrServicesServicesUpdateFail
	}

	h.Layout.ServerConfig.AppsServicesPaths = appSrvcs
	if err := h.Layout.ServerConfig.Save(); err != nil {
		logger.Error("api.ServicesSave::Save", "error", err)
		return customErrors.ErrConfigSaveFail
	}

	ctx, cancelFn := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFn()
	os.DaemonReload(ctx, logger)
	return nil
}
