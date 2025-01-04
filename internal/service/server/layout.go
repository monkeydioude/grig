package server

import (
	"fmt"
	"monkeydioude/grig/internal/errors"
	"monkeydioude/grig/internal/service/os"
	element "monkeydioude/grig/pkg/html/elements"
	"net/http"
	"sync"
)

type Layout struct {
	OS os.OS
	// AppsServices *fs.Dir[model.Service]
	// JosukeConfig   *model.Josuke
	// CapybaraConfig *model.Capybara
	ServerConfig ServerConfig
	Navigation   element.Nav
	mutex        sync.Mutex
}

// Handler our basic generic route handler
type Handler func(http.ResponseWriter, *http.Request) error

// Methods vector of available HTTP Methods
var Methods = [5]string{"GET", "POST", "PUT", "PATCH", "DELETE"}

// WithMethod is a geeneric wrapper around a generic handler, forcing the a HTTP verb
func (l *Layout) WithMethod(method string, handler Handler) func(http.ResponseWriter, *http.Request) {
	// #StephenCurrying
	return func(w http.ResponseWriter, req *http.Request) {
		resBuff := NewResponseWriterBuffer(w)
		for _, m := range Methods {
			// a method matches
			if m == method {
				if err := handler(resBuff, req); err != nil {
					errors.WriteError(err, resBuff)
				}
				_, err := resBuff.End()
				if err != nil {
					errors.WriteError(errors.InternalServerError(err), resBuff)
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
func (l *Layout) Get(handler Handler) func(http.ResponseWriter, *http.Request) {
	return l.WithMethod("GET", handler)
}

// Post is a wrapper around a generic handler, forcing the POST HTTP verb
func (l *Layout) Post(handler Handler) func(http.ResponseWriter, *http.Request) {
	return l.WithMethod("POST", handler)
}

// Put is a wrapper around a generic handler, forcing the PUT HTTP verb
func (l *Layout) Put(handler Handler) func(http.ResponseWriter, *http.Request) {
	return l.WithMethod("PUT", handler)
}

// Patch is a wrapper around a generic handler, forcing the PATCH HTTP verb
func (l *Layout) Patch(handler Handler) func(http.ResponseWriter, *http.Request) {
	return l.WithMethod("PATCH", handler)
}

// Delete is a wrapper around a generic handler, forcing the DELETE HTTP verb
func (l *Layout) Delete(handler Handler) func(http.ResponseWriter, *http.Request) {
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

func (l *Layout) LockMutex() UnlockMutexFn {
	l.mutex.Lock()
	return l.mutex.Unlock
}
