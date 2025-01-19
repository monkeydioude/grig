package middleware

import (
	"context"
	"monkeydioude/grig/pkg/server"
	"net/http"

	"github.com/google/uuid"
)

func JsonApiXRequestID(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		xRequestID := server.NO_X_REQUEST_ID
		tmpXReqID := r.Header.Get(server.X_REQUEST_ID_LABEL)
		if tmpXReqID != "" {
			xRequestID = tmpXReqID
		} else {
			xRequestID = uuid.NewString()
			r.Header.Add(server.X_REQUEST_ID_LABEL, xRequestID)
		}
		ctx := context.WithValue(r.Context(), server.X_REQUEST_ID_LABEL, xRequestID)
		handler.ServeHTTP(w, r.WithContext(ctx))
		w.Header().Add(server.X_REQUEST_ID_LABEL, xRequestID)
	})
}
