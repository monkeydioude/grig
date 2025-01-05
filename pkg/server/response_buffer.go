package server

import (
	"bytes"
	"fmt"
	"net/http"
	"sync"
)

// ResponseWriterBuffer stores the status code
// and the response data, and only calls `WriteHeader`
// and `Write“ though the End() method.
// This allows for some better response handling.
// Current case is, some pkg calls `Write`
// before we have the chance to call `WriteHeader`.
// By definition, any calls to `Write` without `WriteHeader`
// being called beforehand will force a `WriteHeader(http.StatusCodeOk)`.
type ResponseWriterBuffer struct {
	rw           http.ResponseWriter
	status       int
	responseData *bytes.Buffer
	mutex        sync.Mutex
}

func (r *ResponseWriterBuffer) Header() http.Header {
	return r.rw.Header()
}

func (r *ResponseWriterBuffer) Write(data []byte) (int, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	if r.responseData == nil {
		return 0, fmt.Errorf("ResponseWriterBuffer.Write(): %w", ErrNilPointer)
	}
	return r.responseData.Write(data)
}

func (r *ResponseWriterBuffer) WriteHeader(code int) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.status = code
}

func (r *ResponseWriterBuffer) End() (int, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	if r.responseData == nil || r.rw == nil {
		return 0, fmt.Errorf("ResponseWriterBuffer.End(): %w", ErrNilPointer)
	}
	r.rw.WriteHeader(r.status)
	return r.rw.Write(r.responseData.Bytes())
}

func NewResponseWriterBuffer(w http.ResponseWriter) *ResponseWriterBuffer {
	return &ResponseWriterBuffer{
		status:       http.StatusOK,
		rw:           w,
		responseData: bytes.NewBuffer([]byte{}),
	}
}
