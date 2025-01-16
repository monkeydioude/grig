package api

import (
	"log/slog"
	"net/http"
)

func Healthcheck(w http.ResponseWriter, req *http.Request, logger *slog.Logger) error {
	logger.Info("healtcheck")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("{\"health\": \"OK\"}"))
	return nil
}
