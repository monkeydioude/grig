package middleware

import (
	"log"
	"monkeydioude/grig/internal/consts"
	"net/http"
)

const (
	Reset  = "\033[0m"
	Blue   = "\033[1;34m"
	Purple = "\033[1;35m"
	Red    = "\033[1;31m"
	Green  = "\033[1;32m"
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
	if rec.status < 300 {
		log.Printf("%s[%s] <<< %d on API %s %s%s", Green, r.Header.Get(consts.X_REQUEST_ID_LABEL), rec.status, r.Method, r.URL, Reset)
		return
	}
	// 400+
	color := Purple
	data := []byte{}
	if rec.data != nil {
		data = *rec.data
	}
	if rec.status >= 500 {
		color = Red
	}
	log.Printf("%s[%s] <<< %d on API %s %s, response body: %+v %s", color, r.Header.Get(consts.X_REQUEST_ID_LABEL), rec.status, r.Method, r.URL, string(data), Reset)
}

func JsonApiLogRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s[%s] >>> API call on %s %s%s", Blue, r.Header.Get(consts.X_REQUEST_ID_LABEL), r.Method, r.URL, Reset)
		rec := &responseRecorder{rw: w, status: 200}
		handler.ServeHTTP(rec, r)
		rec.HandleResponse(r)
	})
}
