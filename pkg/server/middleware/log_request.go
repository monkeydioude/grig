package middleware

import (
	"log/slog"
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
		slog.InfoContext(r.Context(), "<<<", slog.Int("status", rec.status), slog.String("method", r.Method), slog.String("url", r.URL.String()))
		return
	}
	// 400+
	logMethod := slog.WarnContext
	data := []byte{}
	if rec.data != nil {
		data = *rec.data
	}
	if rec.status >= 500 {
		logMethod = slog.ErrorContext
	}
	logMethod(r.Context(), "<<<", slog.Int("status", rec.status), slog.String("method", r.Method), slog.String("url", r.URL.String()), slog.String("response_body", string(data)))
}

func JsonApiLogRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.InfoContext(r.Context(), ">>>", slog.String("method", r.Method), slog.String("url", r.URL.String()))
		rec := &responseRecorder{rw: w, status: 200}
		handler.ServeHTTP(rec, r)
		rec.HandleResponse(r)
	})
}
