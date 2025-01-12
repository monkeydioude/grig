package middleware

import (
	"log/slog"
	"monkeydioude/grig/pkg/server/http_errors"
	"net/http"
	"runtime/debug"
)

func PanicRecover(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				slog.Error("panic recovered: %v\n Stack trace:\n %s\n", r, debug.Stack())
				http_errors.WriteError(http_errors.UnknownInternalServerError(), w)
			}
		}()
		handler.ServeHTTP(w, r)
	})
}