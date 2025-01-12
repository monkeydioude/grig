package middleware

import (
	"net/http"

	"github.com/google/uuid"
)

func JsonApiXRequestID(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		xRequestID := NO_X_REQUEST_ID
		tmpXReqID := r.Header.Get(X_REQUEST_ID_LABEL)
		if tmpXReqID != "" {
			xRequestID = tmpXReqID
		} else {
			xRequestID = uuid.NewString()
			r.Header.Add(X_REQUEST_ID_LABEL, xRequestID)
		}
		// slog.SetDefault(slog.With(X_REQUEST_ID_LABEL, xRequestID))
		handler.ServeHTTP(w, r)
		w.Header().Add(X_REQUEST_ID_LABEL, xRequestID)
	})
}
