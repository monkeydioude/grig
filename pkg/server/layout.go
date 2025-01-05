package server

import (
	"fmt"
	element "monkeydioude/grig/pkg/html/elements"
	"monkeydioude/grig/pkg/os"
	"monkeydioude/grig/pkg/server/http_errors"
	"net/http"
	"sync"
)

type Layout[T any] struct {
	OS           os.OS
	ServerConfig T
	Navigation   element.Nav
	mutex        sync.Mutex
}

// Handler our basic generic route handler
type Handler func(http.ResponseWriter, *http.Request) error

// Methods vector of available HTTP Methods
var Methods = [5]string{"GET", "POST", "PUT", "PATCH", "DELETE"}

// WithMethod is a geeneric wrapper around a generic handler, forcing the a HTTP verb
func (l *Layout[any]) WithMethod(method string, handler Handler) func(http.ResponseWriter, *http.Request) {
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
				if err := handler(resBuff, req); err != nil {
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
func (l *Layout[any]) Get(handler Handler) func(http.ResponseWriter, *http.Request) {
	return l.WithMethod("GET", handler)
}

// Post is a wrapper around a generic handler, forcing the POST HTTP verb
func (l *Layout[any]) Post(handler Handler) func(http.ResponseWriter, *http.Request) {
	return l.WithMethod("POST", handler)
}

// Put is a wrapper around a generic handler, forcing the PUT HTTP verb
func (l *Layout[any]) Put(handler Handler) func(http.ResponseWriter, *http.Request) {
	return l.WithMethod("PUT", handler)
}

// Patch is a wrapper around a generic handler, forcing the PATCH HTTP verb
func (l *Layout[any]) Patch(handler Handler) func(http.ResponseWriter, *http.Request) {
	return l.WithMethod("PATCH", handler)
}

// Delete is a wrapper around a generic handler, forcing the DELETE HTTP verb
func (l *Layout[any]) Delete(handler Handler) func(http.ResponseWriter, *http.Request) {
	return l.WithMethod("DELETE", handler)
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
