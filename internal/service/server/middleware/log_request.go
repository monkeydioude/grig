package middleware

import (
	"log"
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

func JsonApiLogRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[%s] >>> API call on %s", r.Header.Get(consts.X_REQUEST_ID_LABEL), r.URL)
		rec := &responseRecorder{rw: w, status: 200}
		handler.ServeHTTP(rec, r)
		if rec.status >= 400 {
			data := []byte{}
			if rec.data != nil {
				data = *rec.data
			}
			log.Printf("[%s] <<< %d on API %s, response body: %s", r.Header.Get(consts.X_REQUEST_ID_LABEL), rec.status, r.URL, string(data))
		} else {
			log.Printf("[%s] <<< %d on API %s", r.Header.Get(consts.X_REQUEST_ID_LABEL), rec.status, r.URL)
		}
	})
}
