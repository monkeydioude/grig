package server

import (
	"fmt"
	"log/slog"
	element "monkeydioude/grig/pkg/html/elements"
	"monkeydioude/grig/pkg/os"
	"monkeydioude/grig/pkg/server/http_errors"
	"net/http"
	"sync"
)

const X_REQUEST_ID_LABEL = "X-Request-ID"
const NO_X_REQUEST_ID = "no_x_request_id"

type Layout[T any] struct {
	OS           os.OS
	ServerConfig T
	Navigation   element.Nav
	mutex        sync.Mutex
}

// Handler our basic generic route handler
type Handler func(http.ResponseWriter, *http.Request, *slog.Logger) error

// Methods vector of available HTTP Methods
var Methods = [5]string{"GET", "POST", "PUT", "PATCH", "DELETE"}

// WithMethod is a geeneric wrapper around a generic handler, forcing the a HTTP verb
func WithMethod(method string, handler Handler) func(http.ResponseWriter, *http.Request) {
	// #StephenCurrying
	return func(w http.ResponseWriter, req *http.Request) {
		resBuff := NewResponseWriterBuffer(w)
		if resBuff == nil {
			http_errors.WriteError(http_errors.InternalServerError(fmt.Errorf("Layout.WithMethod(): %w", ErrNilPointer)), resBuff)
			return
		}
		for _, m := range Methods {
			// a method matches
			if m == method {
				logger := slog.Default().With(X_REQUEST_ID_LABEL, req.Context().Value(X_REQUEST_ID_LABEL))
				if err := handler(resBuff, req, logger); err != nil {
					http_errors.WriteError(err, resBuff)
				}
				_, err := resBuff.End()
				if err != nil {
					http_errors.WriteError(http_errors.InternalServerError(err), resBuff)
				}
				return
			}
		}
		// no method matched the one provided over the array of available methods
		w.WriteHeader(405)
		w.Write([]byte(fmt.Sprintf("Method %s not allowd", req.Method)))
	}
}

// Get is a wrapper around a generic handler, forcing the GET HTTP verb
func (*Layout[any]) Get(handler Handler) func(http.ResponseWriter, *http.Request) {
	return WithMethod("GET", handler)
}

// Post is a wrapper around a generic handler, forcing the POST HTTP verb
func (*Layout[any]) Post(handler Handler) func(http.ResponseWriter, *http.Request) {
	return WithMethod("POST", handler)
}

// Put is a wrapper around a generic handler, forcing the PUT HTTP verb
func (*Layout[any]) Put(handler Handler) func(http.ResponseWriter, *http.Request) {
	return WithMethod("PUT", handler)
}

// Patch is a wrapper around a generic handler, forcing the PATCH HTTP verb
func (*Layout[any]) Patch(handler Handler) func(http.ResponseWriter, *http.Request) {
	return WithMethod("PATCH", handler)
}

// Delete is a wrapper around a generic handler, forcing the DELETE HTTP verb
func (*Layout[any]) Delete(handler Handler) func(http.ResponseWriter, *http.Request) {
	return WithMethod("DELETE", handler)
}

// func (l *Layout) SetParams(
// 	AppsServices *fs.Dir[model.Service],
// 	JosukeConfig *model.Josuke,
// 	CapybaraConfig *model.Capybara,
// ) {
// 	l.mutex.Lock()
// 	defer l.mutex.Unlock()
// 	// l.AppsServices = AppsServices
// 	// l.JosukeConfig = JosukeConfig
// 	// l.CapybaraConfig = CapybaraConfig
// }

type UnlockMutexFn = func()

func (l *Layout[any]) LockMutex() UnlockMutexFn {
	l.mutex.Lock()
	return l.mutex.Unlock
}
