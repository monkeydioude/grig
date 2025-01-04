package middleware

import (
	"log"
	"monkeydioude/grig/internal/errors"
	"net/http"
	"runtime/debug"
)

func PanicRecover(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("[ERR ] panic recovered: %v\n Stack trace:\n %s\n", r, debug.Stack())
				errors.WriteError(errors.UnknownInternalServerError(), w)
			}
		}()
		handler.ServeHTTP(w, r)
	})
}
