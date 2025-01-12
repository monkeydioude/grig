package middleware

import (
	"log/slog"
	"monkeydioude/grig/internal/consts"
	"net/http"
)

type responseRecorder struct {
	rw     http.ResponseWriter
	status int
	data   *[]byte
}

func (r *responseRecorder) Header() http.Header {
	return r.rw.Header()
}

func (r *responseRecorder) Write(data []byte) (int, error) {
	r.data = &data
	return r.rw.Write(data)
}

func (r *responseRecorder) WriteHeader(code int) {
	r.status = code
	r.rw.WriteHeader(code)
}

func (rec *responseRecorder) HandleResponse(r *http.Request) {
	if rec.status < 400 {
		slog.Info("<<<", slog.String(consts.X_REQUEST_ID_LABEL, r.Header.Get(X_REQUEST_ID_LABEL)), slog.Int("status", rec.status), slog.String("method", r.Method), slog.String("url", r.URL.String()))
		return
	}
	// 400+
	logMethod := slog.Warn
	data := []byte{}
	if rec.data != nil {
		data = *rec.data
	}
	if rec.status >= 500 {
		logMethod = slog.Error
	}
	logMethod("<<<", slog.String(consts.X_REQUEST_ID_LABEL, r.Header.Get(X_REQUEST_ID_LABEL)), slog.Int("status", rec.status), slog.String("method", r.Method), slog.String("url", r.URL.String()), slog.String("response_body", string(data)))
}

func JsonApiLogRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.Info(">>>", slog.String(consts.X_REQUEST_ID_LABEL, r.Header.Get(X_REQUEST_ID_LABEL)), slog.String("method", r.Method), slog.String("url", r.URL.String()))
		rec := &responseRecorder{rw: w, status: 200}
		handler.ServeHTTP(rec, r)
		rec.HandleResponse(r)
	})
}
