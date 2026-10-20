package v1

import (
	"errors"
	"fmt"
	"log/slog"
	"monkeydioude/grig/internal/service/restart"
	"monkeydioude/grig/pkg/server/http_errors"
	"net/http"
)

// wantsRestart tells whether "Save & restart" was used rather than "Save".
func wantsRestart(r *http.Request) bool {
	return r.URL.Query().Get("restart") == "1"
}

// restartError reports a restart that failed after the config was saved.
func restartError(err error, logger *slog.Logger) error {
	logger.Error("restart after save", "error", err)
	switch {
	case errors.Is(err, restart.ErrJosukeBusy):
		return http_errors.Conflict(fmt.Errorf("saved, not restarted: %w", err))
	case errors.Is(err, restart.ErrNoUnit), errors.Is(err, restart.ErrNoJobsAPI):
		return http_errors.BadRequest(fmt.Errorf("saved, not restarted: %w", err))
	}
	// system errors may hold quotes, which break the JSON error body: logged only
	return http_errors.InternalServerError(errors.New("saved, but the restart failed: see grig logs"))
}
