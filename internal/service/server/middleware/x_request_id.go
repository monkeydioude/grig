package middleware

import (
	"monkeydioude/grig/internal/consts"
	"net/http"

	"github.com/google/uuid"
)

func JsonApiXRequestID(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		xRequestID := consts.NO_X_REQUEST_ID
		tmpXReqID := r.Header.Get(consts.X_REQUEST_ID_LABEL)
		if tmpXReqID != "" {
			xRequestID = tmpXReqID
		} else {
			xRequestID = uuid.NewString()
			r.Header.Add(consts.X_REQUEST_ID_LABEL, xRequestID)
		}
		handler.ServeHTTP(w, r)
		w.Header().Add(consts.X_REQUEST_ID_LABEL, xRequestID)
	})
}
