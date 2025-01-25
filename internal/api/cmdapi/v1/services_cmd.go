package v1

import (
	"log/slog"
	"monkeydioude/grig/internal/errors"
	"monkeydioude/grig/internal/service/parser"
	"monkeydioude/grig/pkg/os"
	"monkeydioude/grig/pkg/server/http_errors"
	"net/http"
)

func (h *Handler) CmdServiceRestart(w http.ResponseWriter, r *http.Request, logger *slog.Logger) error {
	if err := os.ServiceRestart(r.Context(), r.PathValue("service"), logger); err != nil {
		logger.Error("CmdServiceRestart", "error", err, slog.String("func", "monkeydioude/grig/pkg/os.ServiceRestart"))
		return http_errors.InternalServerError(errors.ErrServicesServicesRestartFail)
	}
	return nil
}

func (h *Handler) CmdServiceRestartAll(w http.ResponseWriter, r *http.Request, logger *slog.Logger) error {
	services, err := parser.IniServicesParser(h.Layout.ServerConfig.AppsServicesPaths)
	if err != nil {
		logger.Error("CmdServiceRestartAll", "error", err, slog.String("func", "monkeydioude/grig/internal/service/parser.IniServicesParser"))
		return http_errors.InternalServerError(errors.ErrServicesUnableFileParsing)
	}
	for _, s := range services {
		if err := os.ServiceRestart(r.Context(), s.Name, logger); err != nil {
			return http_errors.InternalServerError(errors.ErrServicesUnableFileParsing)
		}
	}
	return nil
}
